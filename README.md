# VisionCloud

VisionCloud is an AI-powered image upscaling and quality assessment system. It automatically evaluates image quality and upscales low-quality images using deep learning models. The system features a modern React frontend and a Go backend with Python ML integration.

## Architecture

### Backend (Go + Python + PyTorch)
- **Framework**: Go `net/http` with clean architecture
- **Architecture**: Handler -> Service -> Storage layers
- **ML Upscaling**: Python script using **PyTorch (SRCNN)** for super-resolution
- **Quality Assessment**: Laplacian variance and sharpness detection
- **Storage**: AWS S3 with categorized folders:
  - `good_quality/` - High quality images (no upscaling needed)
  - `upscaled/` - Successfully upscaled images
  - `couldnt_upscale/` - Failed processing attempts

### Frontend (React + Vite)
- **Build Tool**: Vite 7.x for fast development and building
- **UI Library**: React 19 with hooks
- **Styling**: Custom CSS with CSS variables, dark theme
- **Key Features**:
  - 🎨 Modern dark UI with gradient accents
  - 📤 Drag & drop image upload
  - 📊 Real-time upload progress
  - 💓 Live backend health monitoring
  - 🖼️ Categorized image gallery
  - 📱 Fully responsive design
  - 🔔 Toast notifications

## Quick Start

### Prerequisites
- Go 1.21+
- Python 3.8+
- Node.js 18+
- AWS Account with credentials configured

### Option 1: Docker (Recommended)

```bash
# Start all services
docker-compose up --build

# Access points:
# - Frontend: http://localhost:3000
# - Backend API: http://localhost:8080
```

### Option 2: Manual Development

**Backend:**
```bash
cd backend
go mod tidy
go run main.go
```

**Frontend:**
```bash
cd frontend
npm install
npm run dev
```

## Project Structure

```
VisionCloud/
├── backend/              # Go backend
│   ├── handlers/        # HTTP handlers
│   ├── services/        # Business logic
│   ├── models/          # Data types
│   └── main.go          # Entry point
├── frontend/            # React frontend
│   ├── src/
│   │   ├── components/  # React components
│   │   ├── hooks/       # Custom hooks
│   │   └── services/    # API services
│   └── Dockerfile
├── python/upscaler/     # Python ML scripts
│   ├── upscale.py       # Main upscaling script
│   ├── model.py         # SRCNN model
│   └── inference.py     # Inference logic
└── docker-compose.yml   # Docker orchestration
```

## API Endpoints

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/api/health` | GET | Health check |
| `/api/images/upload` | POST | Upload image for processing |
| `/api/images/{folder}/{filename}` | GET | Get image info |
| `/api/images/list/{folder}` | GET | List images in folder |

### Upload Response

```json
{
  "success": true,
  "message": "Image processed: success",
  "result": {
    "original_key": "image.jpg",
    "status": "success",
    "folder": "upscaled",
    "quality_score": 0.45,
    "upscale_scale": 2,
    "s3_url": "https://bucket.s3.amazonaws.com/upscaled/image.jpg",
    "processed_at": "2024-01-15T10:30:00Z"
  }
}
```

## Quality Evaluation Logic

The system evaluates images through this pipeline:

1. **Quality Assessment**: Calculate Laplacian variance (sharpness metric)
2. **Decision**:
   - Score ≥ 0.5: Route to `good_quality/` folder
   - Score < 0.5: Attempt upscaling
3. **Upscaling**: Use PyTorch SRCNN model to 2x upscale
4. **Result**:
   - Success: Store in `upscaled/` folder
   - Failure: Store original in `couldnt_upscale/` folder

## Environment Variables

### Backend
| Variable | Description | Default |
|----------|-------------|---------|
| `PORT` | Server port | `8080` |
| `AWS_REGION` | AWS region | `us-east-1` |
| `S3_BUCKET` | S3 bucket name | `visioncloud-bucket` |
| `QUALITY_THRESHOLD` | Quality cutoff (0-1) | `0.5` |
| `UPSCALE_SCALE` | Upscale factor | `2` |

### Frontend
| Variable | Description | Default |
|----------|-------------|---------|
| `VITE_API_URL` | Backend API URL | `http://localhost:8080` |

## Tech Stack Details

### Backend Dependencies
- AWS SDK v2 for Go
- Standard library `net/http`

### Frontend Dependencies
- React 19
- Axios for HTTP requests
- React Dropzone for file uploads
- React Toastify for notifications

### Python Dependencies
- PyTorch
- OpenCV (cv2)
- NumPy
- Pillow

## Browser Support

- Chrome 90+
- Firefox 88+
- Safari 14+
- Edge 90+

## License

MIT
