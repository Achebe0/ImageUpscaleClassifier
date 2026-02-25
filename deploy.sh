#!/bin/bash

# VisionCloud Deployment Script
# Usage: ./deploy.sh [local|aws|k8s]

set -e

DEPLOYMENT_TYPE=${1:-local}

echo "🚀 VisionCloud Deployment - $DEPLOYMENT_TYPE mode"

case $DEPLOYMENT_TYPE in
  local)
    echo "📦 Building and starting local Docker containers..."
    
    # Check if Docker is running
    if ! docker info > /dev/null 2>&1; then
      echo "❌ Docker is not running. Please start Docker Desktop first."
      exit 1
    fi
    
    # Check for .env file
    if [ ! -f .env ]; then
      echo "⚠️  .env file not found. Creating template..."
      cat > .env << EOF
AWS_ACCESS_KEY_ID=your_access_key_here
AWS_SECRET_ACCESS_KEY=your_secret_key_here
AWS_REGION=ca-central-1
S3_BUCKET=visionindex-achebe
QUALITY_THRESHOLD=0.5
UPSCALE_SCALE=2
EOF
      echo "📝 Please edit .env file with your AWS credentials before deploying."
      exit 1
    fi
    
    # Build and deploy
    docker-compose down 2>/dev/null || true
    docker-compose build --no-cache
    docker-compose up -d
    
    echo ""
    echo "✅ Deployment complete!"
    echo "🌐 Frontend: http://localhost:3000"
    echo "🔌 Backend API: http://localhost:8080"
    echo "💚 Health Check: http://localhost:8080/api/health"
    echo ""
    echo "📊 View logs: docker-compose logs -f"
    echo "🛑 Stop: docker-compose down"
    ;;
    
  aws)
    echo "☁️  AWS EC2 Deployment"
    echo "Please ensure you have:"
    echo "  - AWS CLI configured"
    echo "  - EC2 instance running"
    echo "  - SSH key configured"
    echo ""
    echo "Run these commands on your EC2 instance:"
    echo "  1. git clone <your-repo>"
    echo "  2. cd VisionCloud"
    echo "  3. ./deploy.sh local"
    ;;
    
  k8s)
    echo "☸️  Kubernetes Deployment"
    
    # Check for kubectl
    if ! command -v kubectl &> /dev/null; then
      echo "❌ kubectl not found. Please install kubectl first."
      exit 1
    fi
    
    # Apply manifests
    kubectl apply -f k8s/namespace.yaml
    kubectl apply -f k8s/configmap.yaml
    kubectl apply -f k8s/secret.yaml
    kubectl apply -f k8s/backend-deployment.yaml
    kubectl apply -f k8s/backend-service.yaml
    kubectl apply -f k8s/frontend-deployment.yaml
    kubectl apply -f k8s/frontend-service.yaml
    kubectl apply -f k8s/ingress.yaml
    
    echo "✅ Kubernetes deployment applied!"
    echo "📊 Check status: kubectl get pods -n visioncloud"
    ;;
    
  *)
    echo "Usage: ./deploy.sh [local|aws|k8s]"
    echo ""
    echo "Options:"
    echo "  local - Deploy locally with Docker Compose"
    echo "  aws   - Deploy to AWS EC2 (manual steps)"
    echo "  k8s   - Deploy to Kubernetes"
    exit 1
    ;;
esac
