#!/bin/bash
# rebuild.sh - Script to clean, build, and test the app

set -e

echo "=== Cleaning old builds ==="
rm -f lazydndplayer
go clean -cache

echo ""
echo "=== Building lazydndplayer ==="
go build -v -o lazydndplayer . 2>&1

echo ""
echo "✅ Build complete!"
echo ""
echo "=== Build Info ==="
ls -lh lazydndplayer
echo ""

echo "=== IMPORTANT: Remove old character file ==="
echo "Run: rm -f ~/.local/share/lazydndplayer/character.json"
echo "Or your character storage location"
echo ""

echo "=== To run the app ==="
echo "./lazydndplayer"
