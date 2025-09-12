// Package apq provides Automatic Persisted Queries (APQ) domain models and logic.
package apq

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
)

// APQCache defines the interface for APQ cache operations.
type APQCache interface {
	Get(ctx context.Context, hash string) (string, bool, error)
	Set(ctx context.Context, hash, query string, ttl time.Duration) error
	Delete(ctx context.Context, hash string) error
	Clear(ctx context.Context) error
}

// APQManager manages Automatic Persisted Queries.
type APQManager struct {
	cache APQCache
	ttl   time.Duration
}

// NewAPQManager creates a new APQManager.
func NewAPQManager(cache APQCache, ttl time.Duration) *APQManager {
	if ttl == 0 {
		ttl = 24 * time.Hour // Default TTL
	}

	return &APQManager{
		cache: cache,
		ttl:   ttl,
	}
}

// APQRequest represents an APQ request.
type APQRequest struct {
	Query      string                 `json:"query,omitempty"`
	Variables  map[string]interface{} `json:"variables,omitempty"`
	Extensions *APQExtensions         `json:"extensions,omitempty"`
}

// APQExtensions represents APQ extensions.
type APQExtensions struct {
	PersistedQuery *PersistedQuery `json:"persistedQuery,omitempty"`
}

// PersistedQuery represents a persisted query.
type PersistedQuery struct {
	Version int    `json:"version"`
	Hash    string `json:"sha256Hash"`
}

// APQResponse represents an APQ response.
type APQResponse struct {
	Query     string                 `json:"query,omitempty"`
	Variables map[string]interface{} `json:"variables,omitempty"`
	Errors    []APQError             `json:"errors,omitempty"`
}

// APQError represents an APQ error.
type APQError struct {
	Message   string                 `json:"message"`
	Code      string                 `json:"code,omitempty"`
	Extension map[string]interface{} `json:"extension,omitempty"`
}

// ProcessRequest processes an APQ request.
func (m *APQManager) ProcessRequest(ctx context.Context, req *APQRequest) (*APQResponse, error) {
	// If query is provided, try to persist it
	if req.Query != "" {
		hash := m.hashQuery(req.Query)

		// Store the query in cache
		if err := m.cache.Set(ctx, hash, req.Query, m.ttl); err != nil {
			return nil, fmt.Errorf("failed to persist query: %w", err)
		}

		return &APQResponse{
			Query:     req.Query,
			Variables: req.Variables,
		}, nil
	}

	// If only hash is provided, try to retrieve the query
	if req.Extensions != nil && req.Extensions.PersistedQuery != nil {
		hash := req.Extensions.PersistedQuery.Hash

		query, found, err := m.cache.Get(ctx, hash)
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve persisted query: %w", err)
		}

		if !found {
			return &APQResponse{
				Errors: []APQError{
					{
						Message: "PersistedQueryNotFound",
						Code:    "PERSISTED_QUERY_NOT_FOUND",
						Extension: map[string]interface{}{
							"code":    "PERSISTED_QUERY_NOT_FOUND",
							"message": "The persisted query was not found",
						},
					},
				},
			}, nil
		}

		return &APQResponse{
			Query:     query,
			Variables: req.Variables,
		}, nil
	}

	// Neither query nor hash provided
	return &APQResponse{
		Errors: []APQError{
			{
				Message: "Invalid request: either query or persisted query hash must be provided",
				Code:    "INVALID_REQUEST",
			},
		},
	}, nil
}

// hashQuery generates a SHA256 hash of the query.
func (m *APQManager) hashQuery(query string) string {
	hash := sha256.Sum256([]byte(query))
	return hex.EncodeToString(hash[:])
}

// GetQueryHash returns the hash for a given query.
func (m *APQManager) GetQueryHash(query string) string {
	return m.hashQuery(query)
}

// IsPersistedQuery checks if a request is a persisted query request.
func (req *APQRequest) IsPersistedQuery() bool {
	return req.Extensions != nil &&
		req.Extensions.PersistedQuery != nil &&
		req.Query == ""
}

// IsNewQuery checks if a request is a new query to be persisted.
func (req *APQRequest) IsNewQuery() bool {
	return req.Query != "" &&
		(req.Extensions == nil ||
			req.Extensions.PersistedQuery == nil)
}

