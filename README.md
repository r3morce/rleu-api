# Rocket Launch EU API

[![Go Version](https://img.shields.io/badge/Go-1.25-00ADD8?style=flat&logo=go)](https://golang.org)
[![License](https://img.shields.io/badge/license-MIT-green.svg?style=flat)](LICENSE)
[![API](https://img.shields.io/badge/API-ThespaceDevs-orange?style=flat&logo=rocket)](https://thespacedevs.com/)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg?style=flat)](CONTRIBUTING.md)

Go REST API for upcoming European rocket launch data with authentication and location filtering.

> **⚠️ Development API Notice**
> This project currently uses the ThespaceDevs **development API** (`lldev.thespacedevs.com`).
> **Not intended for production use.** For production deployments, switch to the production API endpoint and obtain a proper API key from [ThespaceDevs](https://thespacedevs.com/).

## Quick Start

```bash
# 1. Clone and setup
git clone git@github.com:r3morce/rleu-api.git
cd rleu-api
cp .env.example .env

# 2. Generate API key
uuidgen  # Copy the output

# 3. Edit .env and add your key
RLEU_API_KEY=your-key-here

# 4. Run
go run ./cmd/server

# 5. Test (use your actual key)
curl -H "Authorization: Bearer YOUR_API_KEY" http://localhost:8080/launches
```

## Setup

### Requirements
- Go 1.16+ ([download](https://go.dev/dl/))

### Configuration

Create `.env` file:
```bash
RLEU_API_KEY=your-generated-uuid-here  # Required for auth
DEBUG=false                             # Set to true to skip auth
```

Edit `config/config.yaml` for API settings:
```yaml
api:
  limit: 5                    # Number of launches to fetch

server:
  port: 8080                  # Server port

european_pads:
  enabled: true               # Filter European launches only
```

## Build & Run

```bash
# Development (with hot reload)
go run ./cmd/server

# Production build
go build -o rleu-api ./cmd/server
./rleu-api

# Debug mode (no authentication)
DEBUG=true go run ./cmd/server
```

## API Usage

### Authentication

Include API key in Authorization header:
```bash
curl -H "Authorization: Bearer YOUR_API_KEY" http://localhost:8080/launches
```

### Endpoints

**GET /launches** - Get upcoming launches
- Returns 5 upcoming launches
- European filtering if enabled in config
- Requires authentication (unless DEBUG=true)

Response:
```json
{
  "count": 1,
  "results": [{
    "id": "...",
    "name": "Falcon 9 | Starlink",
    "net": "2025-11-22T23:41:00Z",
    "status": { "name": "Go for Launch" }
  }]
}
```

## Development

### Debug Mode
```bash
DEBUG=true ./rleu-api
curl http://localhost:8080/launches  # No auth needed
```

Debug features:
- ✓ Authentication disabled
- ✓ API throttle monitoring
- ✓ Detailed logging

**⚠️ Never use DEBUG=true in production**

### Project Structure
```
rleu-api/
├── cmd/server/main.go       # HTTP server & auth
├── internal/launches/       # API client & models
├── config/                  # Configuration
└── .env                     # Your secrets (git-ignored)
```

## Deployment

### Fly.io (Recommended)

Deploy to Fly.io with automatic HTTPS and European hosting:

```bash
# 1. Install Fly CLI
curl -L https://fly.io/install.sh | sh

# 2. Login
fly auth login

# 3. Launch app (follow prompts)
fly launch

# 4. Set your API key as a secret
fly secrets set RLEU_API_KEY=your-generated-uuid-here

# 5. Deploy
fly deploy

# 6. Check status
fly status

# Your API will be available at: https://rleu-api.fly.dev
```

**Configuration:**
- Region: Amsterdam (`ams`) - European data center
- Auto HTTPS with Let's Encrypt certificates
- Auto-scaling: scales to zero when idle (free tier friendly)
- Health checks on `/launches` endpoint

**Update deployment:**
```bash
fly deploy
```

**View logs:**
```bash
fly logs
```

### Docker

**With Docker Compose (recommended for local development):**

```bash
# Start (reads .env file automatically)
docker-compose up

# Start in background
docker-compose up -d

# View logs
docker-compose logs -f

# Stop
docker-compose down

# Rebuild and start
docker-compose up --build
```

**With Docker directly:**

```bash
# Build
docker build -t rleu-api .

# Run
docker run -p 8080:8080 -e RLEU_API_KEY=your-key-here rleu-api

# Run in debug mode
docker run -p 8080:8080 -e RLEU_API_KEY=your-key-here -e DEBUG=true rleu-api
```

**Test:**
```bash
curl -H "Authorization: Bearer your-key-here" http://localhost:8080/launches
```

---

**Data:** [ThespaceDevs Launch Library](https://thespacedevs.com/)
