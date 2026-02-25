#!/bin/bash

# VisionCloud Quick Hosting Script
# Usage: ./host.sh [railway|render|local]

set -e

HOSTING_TYPE=${1:-railway}

echo "🚀 VisionCloud Hosting - $HOSTING_TYPE mode"

case $HOSTING_TYPE in
  railway)
    echo "☁️  Deploying to Railway..."
    
    # Check if railway CLI is installed
    if ! command -v railway &> /dev/null; then
      echo "📦 Installing Railway CLI..."
      npm install -g @railway/cli
    fi
    
    # Login to Railway
    echo "🔑 Please login to Railway:"
    railway login
    
    # Link project
    echo "🔗 Linking to Railway project..."
    railway link
    
    # Set environment variables
    echo "⚙️  Setting environment variables..."
    echo "Please enter your AWS credentials:"
    read -p "AWS Access Key ID: " aws_key
    read -p "AWS Secret Access Key: " aws_secret
    
    railway variables set AWS_ACCESS_KEY_ID="$aws_key"
    railway variables set AWS_SECRET_ACCESS_KEY="$aws_secret"
    railway variables set AWS_REGION="ca-central-1"
    railway variables set S3_BUCKET="visionindex-achebe"
    
    # Deploy
    echo "🚀 Deploying to Railway..."
    railway up
    
    echo ""
    echo "✅ Deployment complete!"
    echo "🌐 Your app will be available at the URL shown above"
    echo "📊 View logs: railway logs"
    ;;
    
  render)
    echo "☁️  Deploying to Render..."
    echo ""
    echo "Render uses your render.yaml configuration."
    echo "Steps:"
    echo "1. Push code to GitHub"
    echo "2. Go to https://render.com"
    echo "3. Click 'New' → 'Blueprint'"
    echo "4. Connect your GitHub repo"
    echo "5. Render will auto-deploy using render.yaml"
    echo ""
    echo "⚙️  Make sure to set these environment variables in Render dashboard:"
    echo "   - AWS_ACCESS_KEY_ID"
    echo "   - AWS_SECRET_ACCESS_KEY"
    ;;
    
  local)
    echo "🏠 Hosting on local network..."
    
    # Get IP address
    IP=$(ipconfig | grep "IPv4" | head -1 | grep -oE '[0-9]+\.[0-9]+\.[0-9]+\.[0-9]+' || echo "localhost")
    echo "Your IP address: $IP"
    
    # Update docker-compose to bind to all interfaces
    sed -i 's/"3000:3000"/"0.0.0.0:3000:3000"/g' docker-compose.yml
    sed -i 's/"8080:8080"/"0.0.0.0:8080:8080"/g' docker-compose.yml
    
    # Start Docker
    echo "🐳 Starting Docker containers..."
    docker-compose down 2>/dev/null || true
    docker-compose up -d
    
    echo ""
    echo "✅ Local hosting complete!"
    echo "🌐 Access from this computer: http://localhost:3000"
    echo "🌐 Access from other devices: http://$IP:3000"
    echo ""
    echo "⚠️  Make sure your firewall allows port 3000 and 8080"
    ;;
    
  *)
    echo "Usage: ./host.sh [railway|render|local]"
    echo ""
    echo "Options:"
    echo "  railway - Deploy to Railway (easiest, free tier)"
    echo "  render  - Deploy to Render (free tier available)"
    echo "  local   - Host on your local network"
    exit 1
    ;;
esac
