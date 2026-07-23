#!/bin/bash
# KAVACH CSS Build Script
# Builds Tailwind CSS from source to static output
#
# Prerequisites:
#   npm install
#
# Usage:
#   ./build-css.sh        # One-time build (minified)
#   ./build-css.sh watch  # Watch mode for development

set -e

echo "🎨 KAVACH CSS Builder"
echo "━━━━━━━━━━━━━━━━━━━━"

if [ "$1" == "watch" ]; then
    echo "👀 Watch mode — rebuilding on changes..."
    npx tailwindcss -i ./static/css/input.css -o ./static/css/tailwind.css --watch
else
    echo "📦 Building minified CSS..."
    npx tailwindcss -i ./static/css/input.css -o ./static/css/tailwind.css --minify
    echo "✅ Built: static/css/tailwind.css ($(wc -c < ./static/css/tailwind.css) bytes)"
fi
