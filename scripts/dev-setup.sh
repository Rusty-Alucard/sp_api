#!/bin/bash

# Development setup script for SP API project
# This script sets up the development environment with all necessary tools

set -e

echo "🚀 Setting up SP API development environment..."

# Check if CircleCI CLI is installed
if ! command -v circleci &> /dev/null; then
    echo "📦 Installing CircleCI CLI..."
    if [[ "$OSTYPE" == "darwin"* ]]; then
        # macOS
        if command -v brew &> /dev/null; then
            brew install circleci
        else
            curl -fLSs https://circle.ci/cli | bash
        fi
    else
        # Linux
        curl -fLSs https://circle.ci/cli | bash
    fi
else
    echo "✅ CircleCI CLI already installed"
fi

# Check if pre-commit is installed
if ! command -v pre-commit &> /dev/null; then
    echo "📦 Installing pre-commit..."
    if command -v pip3 &> /dev/null; then
        pip3 install pre-commit
    elif command -v pip &> /dev/null; then
        pip install pre-commit
    else
        echo "❌ pip not found. Please install Python and pip first."
        exit 1
    fi
else
    echo "✅ pre-commit already installed"
fi

# Install pre-commit hooks
echo "🔧 Installing pre-commit hooks..."
pre-commit install

# Validate CircleCI config
echo "🔍 Validating CircleCI configuration..."
make validate-ci

# Run initial lint check
echo "🔍 Running initial lint check..."
make lint

echo "✅ Development environment setup complete!"
echo ""
echo "📋 Available commands:"
echo "  make validate-ci    - Validate CircleCI configuration"
echo "  make validate-all   - Run all validations (CI + lint)"
echo "  make lint          - Run golangci-lint"
echo "  make test          - Run tests"
echo ""
echo "🎯 Pre-commit hooks are now active and will run on every commit."
