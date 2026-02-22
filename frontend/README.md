# VisionCloud Frontend

A modern React frontend for the VisionCloud AI Image Upscaler service.

## Features

- 🎨 **Modern UI** - Clean, dark-themed interface with gradient accents
- 📤 **Drag & Drop Upload** - Easy image upload with visual feedback
- 📊 **Real-time Progress** - Upload and processing progress indicator
- 💓 **Health Monitoring** - Live backend connection status
- 🖼️ **Image Gallery** - Categorized view of processed images
- 📱 **Responsive Design** - Works on desktop and mobile
- 🔔 **Toast Notifications** - User-friendly feedback messages

## Tech Stack

- **React 19** - UI library
- **Vite** - Build tool and dev server
- **Axios** - HTTP client
- **React Dropzone** - Drag and drop file uploads
- **React Toastify** - Toast notifications
- **CSS Modules** - Component-scoped styling

## Project Structure

```
frontend/
├── public/              # Static assets
├── src/
│   ├── components/     # React components
│   │   ├── Header.jsx
│   │   ├── HealthStatus.jsx
│   │   ├── UploadSection.jsx
│   │   ├── ImageGallery.jsx
│   │   └── ImageCard.jsx
│   ├── hooks/          # Custom React hooks
│   │   ├── useHealthCheck.js
│   │   └── useImageUpload.js
│   ├── services/       # API services
│   │   ├── api.js
│   │   └── imageService.js
│   ├── App.jsx         # Main app component
│   ├── main.jsx        # Entry point
│   └── index.css       # Global styles
├── index.html          # HTML template
├── package.json        # Dependencies
├── vite.config.js      # Vite configuration
└── Dockerfile          # Docker build
```

## Getting Started

### Prerequisites

- Node.js 18+
- Backend running on http://localhost:8080

### Installation

```bash
cd frontend
npm install
```

### Development

```bash
npm run dev
```

The app will be available at http://localhost:3000

### Build for Production

```bash
npm run build
```

Output will be in the `dist/` directory.

### Docker

```bash
# Build and run with docker-compose
docker-compose up frontend
```

## API Integration

The frontend connects to the VisionCloud backend API:

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/api/health` | GET | Health check |
| `/api/images/upload` | POST | Upload image |
| `/api/images/{folder}/{filename}` | GET | Get image info |
| `/api/images/list/{folder}` | GET | List images |

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `VITE_API_URL` | Backend API URL | `http://localhost:8080` |

## Features in Detail

### Upload Section
- Drag and drop file upload
- File type validation (JPG, PNG, WebP)
- File size limit (100MB)
- Upload progress indicator
- Processing result display

### Image Gallery
- Three categories: Good Quality, Upscaled, Couldn't Upscale
- Quality score visualization
- Processing metadata display
- S3 link to view processed images

### Health Status
- Real-time backend connection monitoring
- Visual status indicator
- Auto-retry on connection loss

## Browser Support

- Chrome 90+
- Firefox 88+
- Safari 14+
- Edge 90+

## License

MIT
