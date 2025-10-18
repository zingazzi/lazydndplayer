#!/bin/bash
# Quick test script to verify Makefile commands work

set -e

echo "=== Testing Makefile Commands ==="
echo ""

# Test help
echo "1. Testing 'make help'..."
make help
echo ""

# Test deps
echo "2. Testing 'make deps'..."
make deps
echo ""

# Test fmt
echo "3. Testing 'make fmt'..."
make fmt
echo ""

# Test vet
echo "4. Testing 'make vet'..."
make vet
echo ""

# Test build
echo "5. Testing 'make build'..."
make build
echo ""

# Test test
echo "6. Testing 'make test'..."
make test
echo ""

echo "=== All Makefile commands work! ==="

