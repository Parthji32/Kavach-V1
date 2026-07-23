@echo off
REM KAVACH CSS Build Script (Windows)
REM Prerequisites: npm install

echo KAVACH CSS Builder
echo ==================

if "%1"=="watch" (
    echo Watch mode - rebuilding on changes...
    npx tailwindcss -i ./static/css/input.css -o ./static/css/tailwind.css --watch
) else (
    echo Building minified CSS...
    npx tailwindcss -i ./static/css/input.css -o ./static/css/tailwind.css --minify
    echo Done: static/css/tailwind.css
)
