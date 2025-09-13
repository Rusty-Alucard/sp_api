// Package graphql provides GraphQL domain models and logic.
package graphql

import (
	"context"
	"fmt"
	"strconv"

	"github.com/vektah/gqlparser/v2/ast"
	"github.com/vektah/gqlparser/v2/parser"
)

// QueryAnalyzer analyzes GraphQL queries for complexity and cost.
type QueryAnalyzer struct {
	complexityRules   map[string]ComplexityRule
	defaultComplexity int
}

// ComplexityRule defines complexity calculation rules for fields.
type ComplexityRule struct {
	BaseComplexity int
	Multiplier     string // Field name for multiplier (e.g., "limit", "first")
}

// NewQueryAnalyzer creates a new QueryAnalyzer.
func NewQueryAnalyzer() *QueryAnalyzer {
	return &QueryAnalyzer{
		complexityRules:   make(map[string]ComplexityRule),
		defaultComplexity: 1,
	}
}

// AddComplexityRule adds a complexity rule for a field.
func (qa *QueryAnalyzer) AddComplexityRule(fieldPath string, rule ComplexityRule) {
	qa.complexityRules[fieldPath] = rule
}

// AnalyzeQuery analyzes a GraphQL query and returns its complexity.
func (qa *QueryAnalyzer) AnalyzeQuery(_ context.Context, query string, variables map[string]interface{}) (
	*QueryAnalysis, error) {
	// Parse the query
	doc, err := parser.ParseQuery(&ast.Source{Input: query})
	if err != nil {
		return nil, fmt.Errorf("failed to parse query: %w", err)
	}

	// Analyze the query
	analysis := &QueryAnalysis{
		Complexity: 0,
		Depth:      0,
		FieldCount: 0,
		Fields:     make(map[string]int),
	}

	for _, operation := range doc.Operations {
		complexity, depth, fieldCount, fields := qa.analyzeOperation(operation, variables)
		analysis.Complexity += complexity
		if depth > analysis.Depth {
			analysis.Depth = depth
		}
		analysis.FieldCount += fieldCount

		// Merge field complexities
		for field, count := range fields {
			analysis.Fields[field] += count
		}
	}

	return analysis, nil
}

// QueryAnalysis represents the analysis result of a GraphQL query.
type QueryAnalysis struct {
	Complexity int
	Depth      int
	FieldCount int
	Fields     map[string]int // Field path -> complexity
}

// IsWithinLimit checks if the query is within the specified complexity limit.
func (qa *QueryAnalysis) IsWithinLimit(limit int) bool {
	return qa.Complexity <= limit
}

// GetComplexityBreakdown returns a breakdown of complexity by field.
func (qa *QueryAnalysis) GetComplexityBreakdown() map[string]int {
	return qa.Fields
}

// analyzeOperation analyzes a GraphQL operation.
func (qa *QueryAnalyzer) analyzeOperation(
	operation *ast.OperationDefinition,
	variables map[string]interface{},
) (int, int, int, map[string]int) {
	complexity := 0
	depth := 0
	fieldCount := 0
	fields := make(map[string]int)

	for _, selection := range operation.SelectionSet {
		c, d, fc, f := qa.analyzeSelection(selection, "", 0, variables)
		complexity += c
		if d > depth {
			depth = d
		}
		fieldCount += fc

		for field, count := range f {
			fields[field] += count
		}
	}

	return complexity, depth, fieldCount, fields
}

// analyzeSelection analyzes a GraphQL selection.
func (qa *QueryAnalyzer) analyzeSelection(
	selection ast.Selection,
	path string,
	depth int,
	variables map[string]interface{},
) (int, int, int, map[string]int) {
	switch sel := selection.(type) {
	case *ast.Field:
		return qa.analyzeField(sel, path, depth, variables)
	case *ast.InlineFragment:
		return qa.analyzeInlineFragment(sel, path, depth, variables)
	case *ast.FragmentSpread:
		// Fragment spreads are handled by the parser
		return 0, depth, 0, make(map[string]int)
	default:
		return 0, depth, 0, make(map[string]int)
	}
}

// analyzeField analyzes a GraphQL field.
func (qa *QueryAnalyzer) analyzeField(
	field *ast.Field,
	path string,
	depth int,
	variables map[string]interface{},
) (int, int, int, map[string]int) {
	fieldPath := qa.buildFieldPath(path, field.Name)

	// Get complexity rule for this field
	rule, exists := qa.complexityRules[fieldPath]
	if !exists {
		rule = ComplexityRule{BaseComplexity: qa.defaultComplexity}
	}

	// Calculate base complexity
	complexity := rule.BaseComplexity

	// Apply multiplier if specified
	if rule.Multiplier != "" {
		multiplier := qa.getMultiplier(field, rule.Multiplier, variables)
		complexity *= multiplier
	}

	// Analyze nested selections
	nestedComplexity := 0
	nestedDepth := depth + 1
	nestedFieldCount := 1
	fields := map[string]int{fieldPath: complexity}

	if field.SelectionSet != nil {
		for _, selection := range field.SelectionSet {
			c, d, fc, f := qa.analyzeSelection(selection, fieldPath, nestedDepth, variables)
			nestedComplexity += c
			if d > nestedDepth {
				nestedDepth = d
			}
			nestedFieldCount += fc

			for field, count := range f {
				fields[field] += count
			}
		}
	}

	return complexity + nestedComplexity, nestedDepth, nestedFieldCount, fields
}

// analyzeInlineFragment analyzes an inline fragment.
func (qa *QueryAnalyzer) analyzeInlineFragment(
	fragment *ast.InlineFragment,
	path string,
	depth int,
	variables map[string]interface{},
) (int, int, int, map[string]int) {
	complexity := 0
	maxDepth := depth
	fieldCount := 0
	fields := make(map[string]int)

	for _, selection := range fragment.SelectionSet {
		c, d, fc, f := qa.analyzeSelection(selection, path, depth, variables)
		complexity += c
		if d > maxDepth {
			maxDepth = d
		}
		fieldCount += fc

		for field, count := range f {
			fields[field] += count
		}
	}

	return complexity, maxDepth, fieldCount, fields
}

// buildFieldPath builds a field path from the current path and field name.
func (qa *QueryAnalyzer) buildFieldPath(currentPath, fieldName string) string {
	if currentPath == "" {
		return fieldName
	}
	return currentPath + "." + fieldName
}

// getMultiplier extracts the multiplier value from field arguments or variables.
func (qa *QueryAnalyzer) getMultiplier(field *ast.Field, multiplierName string, variables map[string]interface{}) int {
	// Check field arguments first
	for _, arg := range field.Arguments {
		if arg.Name == multiplierName {
			if arg.Value.Kind == ast.IntValue {
				value, _ := strconv.Atoi(arg.Value.Raw)
				return value
			}
			if arg.Value.Kind == ast.Variable {
				if variables != nil {
					if value, ok := variables[arg.Value.Raw]; ok {
						if intValue, ok := value.(int); ok {
							return intValue
						}
						if floatValue, ok := value.(float64); ok {
							return int(floatValue)
						}
					}
				}
			}
		}
	}

	// Default multiplier
	return 1
}

// QueryComplexityError represents a query complexity exceeded error.
type QueryComplexityError struct {
	Message        string
	Complexity     int
	Limit          int
	FieldBreakdown map[string]int
}

func (e *QueryComplexityError) Error() string {
	return e.Message
}

// NewQueryComplexityError creates a new query complexity error.
func NewQueryComplexityError(complexity, limit int, breakdown map[string]int) *QueryComplexityError {
	return &QueryComplexityError{
		Message:        fmt.Sprintf("query complexity %d exceeds limit %d", complexity, limit),
		Complexity:     complexity,
		Limit:          limit,
		FieldBreakdown: breakdown,
	}
}
