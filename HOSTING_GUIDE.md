# VisionCloud Hosting Guide

## Option 1: Host on Your Computer (Local Network)
If you want others to access it from your local network:

```bash
# Find your computer's IP address
ipconfig  # Windows
# Look for "IPv4 Address" under your active connection

# Edit docker-compose.yml, change ports to:
ports:
  - "0.0.0.0:3000:3000"  # Frontend
  - "0.0.0.0:8080:8080"  # Backend

# Restart containers
docker-compose down
docker-compose up -d
```

Others can now access via: `http://YOUR_IP:3000`

---

## Option 2: Railway (Easiest Cloud Hosting - FREE)
Railway offers free hosting with automatic deployments.

### Steps:
1. Push your code to GitHub
2. Go to https://railway.app and sign up
3. Click "New Project" → "Deploy from GitHub repo"
4. Select your VisionCloud repo
5. Railway will auto-detect your docker-compose.yml
6. Add environment variables in Railway dashboard:
   - `AWS_ACCESS_KEY_ID`
   - `AWS_SECRET_ACCESS_KEY`
   - `AWS_REGION=ca-central-1`
   - `S3_BUCKET=visionindex-achebe`
7. Deploy! You'll get a public URL

---

## Option 3: Render (Free Tier Available)
Another easy cloud option.

### Steps:
1. Go to https://render.com
2. Sign up with GitHub
3. Click "New" → "Blueprint"
4. Connect your GitHub repo
5. Render will use your docker-compose.yml
6. Add environment variables
7. Deploy

---

## Option 4: AWS (Most Professional)
Best for production use.

### Option 4A: AWS ECS (Easiest AWS option)
```bash
# Install AWS Copilot
brew install aws/tap/copilot-cli  # Mac
# or download from AWS website for Windows

# Initialize
copilot init --app visioncloud --name backend --type 'Load Balanced Web Service' --dockerfile './backend/Dockerfile'

# Deploy
copilot deploy
```

### Option 4B: AWS EC2 (More Control)
1. Launch EC2 instance (t3.medium recommended)
2. Install Docker:
   ```bash
   sudo yum update -y
   sudo yum install docker -y
   sudo service docker start
   sudo usermod -a -G docker ec2-user
   ```
3. Clone your repo and run:
   ```bash
   docker-compose up -d
   ```

---

## Option 5: Vercel + Render (Hybrid)
- **Frontend**: Deploy to Vercel (free, fast CDN)
- **Backend**: Deploy to Render or Railway

### Frontend on Vercel:
```bash
cd frontend
npm install -g vercel
vercel --prod
```

### Backend on Render:
- Follow Option 3 above for just the backend

---

## Quick Start: Railway (Recommended)

**Why Railway?**
- ✅ Free tier available
- ✅ Automatic HTTPS
- ✅ Easy GitHub integration
- ✅ Auto-deploy on push
- ✅ Built-in environment variables

### Step-by-Step:

1. **Push to GitHub** (if not already):
   ```bash
   git add .
   git commit -m "Containerized VisionCloud"
   git push origin main
   ```

2. **Go to Railway**:
   - Visit https://railway.app
   - Sign up with GitHub
   - Click "New Project"

3. **Deploy**:
   - Select "Deploy from GitHub repo"
   - Choose your VisionCloud repository
   - Railway auto-detects Docker setup

4. **Add Environment Variables**:
   - In Railway dashboard, go to Variables
   - Add:
     - `AWS_ACCESS_KEY_ID` = your_key
     - `AWS_SECRET_ACCESS_KEY` = your_secret
     - `AWS_REGION` = ca-central-1
     - `S3_BUCKET` = visionindex-achebe

5. **Get Public URL**:
   - Railway provides a `.up.railway.app` URL
   - Custom domains available in settings

---

## Domain & HTTPS

Once hosted, add a custom domain:

1. Buy domain (Namecheap, GoDaddy, etc.)
2. Add DNS records pointing to your host
3. Enable HTTPS (usually automatic on cloud platforms)

---

## Monitoring

After hosting, monitor your app:
- **Railway**: Built-in metrics dashboard
- **AWS**: CloudWatch
- **Render**: Built-in logs and metrics

---

## Need Help?

Choose your path:
- **Fastest**: Railway (5 minutes setup)
- **Free forever**: Render or Vercel+Render combo
- **Production scale**: AWS ECS or EKS
- **Learning**: AWS EC2 (full control)
