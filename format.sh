#!/bin/bash

# Auto-fix Go code formatting and imports
goimports -w .
gofumpt -w .
golines -w .

# Run linter (read-only)
golangci-lint run