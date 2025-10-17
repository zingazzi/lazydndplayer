#!/bin/bash
# RUN_ME.sh - Quick start script for LazyDnDPlayer

echo "════════════════════════════════════════════════════════"
echo "  LazyDnDPlayer v2.0 - D&D Character Manager"
echo "════════════════════════════════════════════════════════"
echo ""
echo "Starting the application with sample character..."
echo ""
echo "New Features in v2.0:"
echo "  ✓ Tab-based navigation (no more sidebar)"
echo "  ✓ Dice roller always visible (bottom right)"
echo "  ✓ Actions panel always visible (bottom left)"
echo "  ✓ Quick dice rolls: d1-d7 (d6 = roll d20)"
echo "  ✓ Shift+R for long rest"
echo ""
echo "Quick Tips:"
echo "  • Press ? for help"
echo "  • Press 1-5 to switch tabs"
echo "  • Press d6 to roll a d20"
echo "  • Press Tab to cycle through tabs"
echo "  • Press q to quit"
echo ""
echo "════════════════════════════════════════════════════════"
echo ""

# Run with sample character
./lazydndplayer -file ./data/sample_character.json
