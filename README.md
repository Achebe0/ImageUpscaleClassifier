# CloudSnap IQ

CloudSnap IQ is an automated image quality evaluation system that routes uploads to different S3 buckets based on visual correctness. It uses a Go backend for processing and a React frontend for user interaction.

## Architecture

### Backend (Go + Python + PyTorch)
- **Framework**: `net/http` with `gorilla/mux`
- **Architecture**: Clean layered architecture (Handler -> Service -> Storage)
- **ML Upscaling**: Python script using **PyTorch (SRCNN)** for super-resolution.
- **Storage**:
  - AWS S3 (Approved vs Rejected buckets)
  - AWS DynamoDB (Metadata storage)

### Frontend (React)
- **Build Tool**: Vite
- **Styling**: TailwindCSS
- **Features**: Image upload, quality visualization, gallery view

## Setup

### Prerequisites
- Go 1.21+
- Python 3.8+
- Node.js 18+
- AWS Account with credentials configured

### Backend Setup
1. Navigate to `backend/`
2. Copy `.env.example` to `.env` and fill in your AWS details.
3. Install Python dependencies:
   ```bash
   pip install -r requirements.txt
   ```
   *(This installs torch, torchvision, opencv-python, etc.)*
4. Run `go mod tidy`
5. Start server: `go run cmd/server/main.go`

### Frontend Setup
1. Navigate to `frontend/` (Create this using Vite)
2. `npm install`
3. `npm run dev`

## API Endpoints

- `POST /api/upload`: Upload an image for evaluation
- `GET /api/images`: List processed images with quality scores
- `GET /api/health`: Health check

## Quality Evaluation Logic

The system evaluates images based on:
1. **Upscaling**: Images are upscaled using a PyTorch SRCNN model (Super-Resolution Convolutional Neural Network).
2. **Sharpness Check**: Compares Laplacian variance of original vs upscaled.
3. **Score**: Normalized score (0-100).

Score > 50 is considered **High Quality**.
