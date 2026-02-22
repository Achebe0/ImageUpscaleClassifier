@echo off
REM VisionCloud Quick Start Script for Windows

echo 🚀 Starting VisionCloud Pipeline...
echo.

REM Check if backend binary exists
if not exist "backend\main.exe" (
    echo 📦 Building Go backend...
    cd backend
    go build -o main.exe
    cd ..
)

REM Check Python requirements
python -c "import cv2, torch, numpy" >nul 2>&1
if %errorlevel% neq 0 (
    echo 📦 Installing Python dependencies...
    pip install -r python/upscaler/requirements.txt
)

echo.
echo ✅ Setup complete!
echo.
echo To run the backend:
echo   cd backend
echo   main.exe  (or: go run main.go)
echo.
echo The server will start on http://localhost:8080
echo.
echo Configuration:
echo   - Copy .env.example to .env and update AWS credentials
echo   - Set QUALITY_THRESHOLD (0-1) for quality assessment
echo   - Set UPSCALE_SCALE (2 or 4) for upscaling factor
echo.
echo Features:
echo   ✓ Go backend handles HTTP requests ^& S3 routing
echo   ✓ PyTorch upscaler runs via subprocess (no Flask overhead)
echo   ✓ Quality assessment based on image resolution
echo   ✓ Automatic S3 categorization
echo.
