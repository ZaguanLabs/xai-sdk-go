#!/bin/bash
# Test script for models listing

if [ -z "$XAI_API_KEY" ]; then
    echo "Error: XAI_API_KEY environment variable is not set"
    echo "Usage: export XAI_API_KEY='your-key' && ./test_models.sh"
    exit 1
fi

echo "Testing models listing with API key: ${XAI_API_KEY:0:10}..."
go run examples/models/list.go
