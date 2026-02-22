# VisionCloud Image Upscaling Pipeline

Complete image upscaling pipeline with Go backend and Python deep learning upscaler.

## Architecture

### Flow
1. **Image Upload** → Go backend receives image via HTTP POST  
2. **Quality Assessment** → Evaluates image resolution against threshold
3. **Routing**:
   - **Good Quality** (above threshold) → Upload to `good_quality/` folder
   - **Needs Upscaling** (below threshold) → Send to PyTorch model
4. **Upscaling** → PyTorch directly upscales image via subprocess
5. **Result Handling**:
   - **Success** → Upload upscaled image to `upscaled/` folder
   - **Failure** → Upload original to `couldn't_upscale/` folder

### S3 Folder Structure
```
s3://visioncloud-bucket/
├── good_quality/          # Images already at good quality
├── upscaled/              # Successfully upscaled images
└── couldn't_upscale/      # Images that failed upscaling
```

## Setup

### Prerequisites
- Go 1.20+
- Python 3.9+
- AWS Account with S3 access
- CUDA (optional, for GPU acceleration)

### Backend Setup (Go)

```bash
cd backend
go mod download
```

### Python Upscaler Setup

```bash
cd python/upscaler
pip install -r requirements.txt
```

## Configuration

Create a `.env` file in the project root:

```env
# Backend
PORT=8080
AWS_REGION=us-east-1
S3_BUCKET=your-bucket-name
QUALITY_THRESHOLD=0.5      # 0-1, images below this need upscaling
UPSCALE_SCALE=2            # 2x or 4x upscaling
UPSCALER_SERVICE_URL=http://localhost:5000/upscale

# AWS Credentials (use AWS SDK default chain or set here)
AWS_ACCESS_KEY_ID=your_key
AWS_SECRET_ACCESS_KEY=your_secret

# Python Upscaler
MODEL_PATH=./models/upscaler.pth
DEVICE=cuda                # or 'cpu'
```

## Running the Services

### Start Go Backend
```bash
cd backend
go run main.go
```

The backend will be available at `http://localhost:8080`

**Note:** The Python upscaler will be called automatically when needed. No separate service needs to be running.

## API Endpoints

### Upload Image
**POST** `/api/images/upload`

Upload an image for processing through the pipeline.

```bash
curl -X POST http://localhost:8080/api/images/upload \
  -F "image=@image.jpg"
```

Response:
```json
{
  "success": true,
  "message": "Image processed: success",
  "result": {
    "original_key": "image.jpg",
    "status": "success",
    "folder": "upscaled",
    "s3_url": "https://bucket.s3.amazonaws.com/upscaled/image.jpg",
    "quality_score": 0.35,
    "upscale_scale": 2,
    "processed_at": "2024-01-15T10:30:45Z"
  }
}
```

### Health Check
**GET** `/api/health`

```bash
curl http://localhost:8080/api/health
```

Response:
```json
{
  "status": "healthy"
}
```

## Quality Threshold Explanation

The quality score is calculated based on image resolution:
- Score = (width × height) / (3840 × 2160) [4K resolution]
- Examples:
  - 1920×1080 (Full HD) → score ≈ 0.25
  - 2560×1440 (2K) → score ≈ 0.45
  - 3840×2160 (4K) → score = 1.0

If `QUALITY_THRESHOLD=0.5`:
- Images < 0.5 score → upscaled
- Images ≥ 0.5 score → stored as good_quality

## Error Handling

| Status | Folder | Meaning |
|--------|--------|---------|
| success | good_quality | Image already high quality |
| success | upscaled | Successfully upscaled low-quality image |
| error | couldn't_upscale | Failed to upscale, original stored for review |

## Deployment

### Docker Images

Create `docker-compose.yml`:

```yaml
version: '3.8'
services:
  backend:
    build:
      context: ./backend
      dockerfile: Dockerfile
    ports:
      - "8080:8080"
    environment:
      - PORT=8080
      - AWS_REGION=us-east-1
      - S3_BUCKET=visioncloud-bucket
      - QUALITY_THRESHOLD=0.5
      - UPSCALER_SERVICE_URL=http://upscaler:5000/upscale
    depends_on:
      - upscaler

  upscaler:
    build:
      context: ./python/upscaler
      dockerfile: Dockerfile
    ports:
      - "5000:5000"
    environment:
      - DEVICE=cuda
      - MODEL_PATH=/model/upscaler.pth
    volumes:
      - ./models:/model
```

### AWS Deployment

1. **Backend**: Deploy Go binary to EC2, ECS, or Lambda
2. **Python Service**: Deploy to EC2, Docker on ECS, or SageMaker
3. **Configure S3**: Set up IAM roles with S3 access

## Performance Tuning

- **Batch Processing**: Process multiple images concurrently
- **Model Caching**: Load model once, reuse for multiple requests
- **GPU Acceleration**: Use CUDA for faster upscaling
- **S3 Configuration**: Use S3 Transfer Acceleration for large files

## Monitoring & Logging

Monitor processing metrics:
- Images per folder
- Quality score distribution
- Upscaling success rate
- Processing time per image

## Troubleshooting

### Images not uploading to S3
- Check AWS credentials and IAM permissions
- Verify bucket name and region
- Ensure S3 bucket exists

### Upscaler service errors
- Verify Python service is running: `curl http://localhost:5000/health`
- Check model file exists at path specified in `MODEL_PATH`
- Review Python service logs for errors

### Quality assessment issues
- Ensure image file is valid (JPEG, PNG, WebP)
- Check quality threshold value (0-1)
- Verify image dimensions can be read

## Future Enhancements

- [ ] Batch image processing API
- [ ] Multiple upscaling models selection
- [ ] Custom quality metrics beyond resolution
- [ ] Image metadata preservation (EXIF)
- [ ] WebUI for monitoring
- [ ] Webhook notifications
- [ ] Rate limiting and authentication
