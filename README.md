# Rocket Launch EU API

Go REST API for upcoming European rocket launch data with authentication and location filtering.

## Quick Start

```bash
# 1. Clone and setup
git clone <repo-url>
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

## Troubleshooting

### Port already in use
```bash
# Find process using port 8080
lsof -i :8080
# Kill it or change port in config/config.yaml
```

### Authentication errors
```bash
# Check your .env file exists
cat .env

# Verify API key is set
grep RLEU_API_KEY .env

# Test without auth in debug mode
DEBUG=true go run ./cmd/server
```

### Build errors
```bash
# Verify Go version
go version  # Should be 1.16+

# Clean and rebuild
go clean
go build -o rleu-api ./cmd/server
```

### API returns no data
- Check internet connection
- Verify ThespaceDevs API is up
- Try disabling European filtering in `config/config.yaml`

## Configuration Reference

### Environment Variables (.env)
```bash
RLEU_API_KEY=<uuid>          # Required: Your API authentication key
DEBUG=false                   # Optional: Enable debug mode (true/false)
SPACEDEVS_API_KEY=           # Optional: ThespaceDevs API key (higher limits)
```

### Config File (config/config.yaml)
```yaml
api:
  base_url: "https://lldev.thespacedevs.com/2.3.0/launches/upcoming/"
  limit: 5
  mode: "list"
  format: "json"

server:
  port: 8080
  host: "localhost"

european_pads:
  enabled: true
  location_ids: [159, 178, 157, 205]  # Sutherland, Esrange, SaxaVord
  country_codes: ["FRA", "GBR", "SWE", "NOR", ...]
```

## Contributing

Contributions welcome! See [UPCOMING.md](UPCOMING.md) for planned features.

## License

Personal learning project - use freely!

---

**Data:** [ThespaceDevs Launch Library](https://thespacedevs.com/)
