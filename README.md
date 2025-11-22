# Rocket Launch EU (RLEU) API

A Go backend API service that fetches and displays upcoming European rocket launches using the [ThespaceDevs Launch Library API](https://thespacedevs.com/).

## What This Does

This backend provides a REST API endpoint that returns the next 5 upcoming rocket launches worldwide. It's a great beginner project for learning:
- Building HTTP servers in Go
- Making API requests to external services
- Working with JSON data
- Managing environment variables securely

## Project Structure

```
rleu-api/
├── cmd/
│   └── server/
│       └── main.go          # HTTP server and request handlers
├── internal/
│   └── launches/
│       ├── api.go           # ThespaceDevs API client
│       └── models.go        # Launch data structures
├── config/
│   └── env.go               # Environment variable loader
├── .env                     # Your API key (git-ignored)
├── .env.example             # Template for .env
├── .gitignore               # Git ignore rules
├── go.mod                   # Go module definition
├── README.md                # This file
└── UPCOMING.md              # Planned features
```

## Configuration

The application uses `config/config.yaml` for settings and environment variables for sensitive data.

### config/config.yaml

```yaml
api:
  base_url: "https://lldev.thespacedevs.com/2.3.0/launches/upcoming/"
  limit: 5          # Number of launches to fetch
  mode: "list"      # API response mode
  format: "json"    # Output format

server:
  port: 8080        # Server port

european_pads:
  enabled: true     # Enable European filtering
  location_ids: [159, 178, 157, 205]
  country_codes: ["FRA", "GBR", "SWE", ...]
```

### API Authentication

**Important:** This API requires authentication via API key (except in DEBUG mode).

#### Setting Up Authentication

1. **Generate an API Key**: Create a secure random key
   ```bash
   # On macOS/Linux:
   uuidgen
   # Example output: 016fdcae-62e4-4cb2-b64e-fcc412d33242
   ```

2. **Add to .env file**: Set the `RLEU_API_KEY` in your `.env` file:
   ```bash
   RLEU_API_KEY=016fdcae-62e4-4cb2-b64e-fcc412d33242
   ```

3. **Restart the server**: The API will now require this key for all requests

#### Making Authenticated Requests

Clients must include the API key in the `Authorization` header:

```bash
# Using Bearer format
curl -H "Authorization: Bearer 016fdcae-62e4-4cb2-b64e-fcc412d33242" \
  http://localhost:8080/launches

# Or using Token format
curl -H "Authorization: Token 016fdcae-62e4-4cb2-b64e-fcc412d33242" \
  http://localhost:8080/launches
```

**Without authentication**, requests will return `401 Unauthorized`.

### Debug Mode

Enable debug mode by setting `DEBUG=true`:
```bash
DEBUG=true go run ./cmd/server
# or
make run-debug
```

**In debug mode:**
- ⚠️ **Authentication is disabled** - All requests work without an API key
- API throttle status is checked after each request
- Detailed request logging is enabled

**WARNING:** Only use DEBUG mode for local development and testing. Never deploy with DEBUG=true in production.

## Getting Started

### Prerequisites

- Go 1.16 or higher installed ([Download Go](https://go.dev/dl/))
- Basic knowledge of Go from boot.dev or similar course

### Installation

1. **Clone or navigate to this repository**
   ```bash
   cd rleu-api
   ```

2. **Set up environment variables**
   ```bash
   cp .env.example .env
   ```

   Edit `.env` and set your `RLEU_API_KEY`:
   ```bash
   # Generate a key
   uuidgen

   # Add it to .env
   RLEU_API_KEY=your-generated-key-here
   ```

   See the [API Authentication](#api-authentication) section for details.

3. **Build the project**
   ```bash
   go build -o rleu-api ./cmd/server
   ```

4. **Run the server**
   ```bash
   ./rleu-api
   ```
   
   Or run directly without building:
   ```bash
   go run ./cmd/server
   ```

5. **Test it!**

   Use curl with your API key:
   ```bash
   curl -H "Authorization: Bearer YOUR_API_KEY" http://localhost:8080/launches
   ```

   Or test in DEBUG mode (no auth required):
   ```bash
   DEBUG=true ./rleu-api
   curl http://localhost:8080/launches
   ```

## API Endpoints

### GET /launches

Returns the next 5 upcoming rocket launches. If European filtering is enabled in config, only returns launches from European spaceports.

**European Filtering:**
- Set `european_pads.enabled: true` in `config/config.yaml`
- Uses location IDs: 159 (Sutherland), 178 (Esrange), 157/205 (SaxaVord)
- Country codes: FRA, GBR, SWE, etc.

**Example Response:**
```json
{
  "count": 105,
  "next": "https://lldev.thespacedevs.com/2.3.0/launches/upcoming/?format=json&limit=5&mode=list&offset=5",
  "previous": null,
  "results": [
    {
      "id": "a81fcb21-ed57-47ef-ab57-33ddd9c6ffcd",
      "name": "Falcon 9 Block 5 | Starlink Group 6-79",
      "net": "2025-11-22T23:41:00Z",
      "status": {
        "name": "Go for Launch",
        "abbrev": "Go"
      },
      "image": {
        "image_url": "https://...",
        "thumbnail_url": "https://..."
      }
      // ... more fields
    }
    // ... 4 more launches
  ]
}
```

**European Filtered Example:**
```json
{
  "count": 2,
  "results": [
    {
      "name": "Ariane 5 | Starlink Group 6-79",
      "location": "Kourou, French Guiana"
    },
    {
      "name": "Vega | Earth Observation Satellite",
      "location": "Kourou, French Guiana"
    }
  ]
}
```

## How It Works

1. **Client Request**: A client (browser, app, etc.) sends a GET request to `/launches`
2. **API Call**: The server calls the ThespaceDevs API to fetch upcoming launches
3. **Data Processing**: The JSON response is parsed into Go structs
4. **Response**: The launch data is returned to the client as JSON

```
Client → GET /launches → Handler → ThespaceDevs API → Parse JSON → Return to Client
```

## Learning Resources

### Understanding the Code

- **cmd/server/main.go**: Entry point - HTTP server and request handlers
- **internal/launches/models.go**: Go structs mapping to JSON data
- **internal/launches/api.go**: HTTP client for ThespaceDevs API
- **config/env.go**: Environment variable management

### Go Concepts Used

- HTTP servers with `net/http`
- JSON encoding/decoding
- Error handling with `error` types
- Structs and JSON tags
- File I/O for environment variables
- Pointers (for nullable fields)

## Upcoming Features

See [UPCOMING.md](UPCOMING.md) for planned features, including:
- Filtering for European launches only
- Caching to improve performance
- Additional endpoints

## Troubleshooting

### Port Already in Use
If you see "address already in use", another program is using port 8080. Either:
- Stop the other program
- Change the port in `cmd/server/main.go` (line 21)

### API Request Failed
- Check your internet connection
- The ThespaceDevs API might be temporarily down
- You might have hit rate limits (get an API key to increase limits)

### Build Errors
Make sure you're in the project directory and Go is properly installed:
```bash
go version  # Should show Go 1.16 or higher
```

## Contributing

This is a learning project! Feel free to:
- Add new features from [UPCOMING.md](UPCOMING.md)
- Improve error handling
- Add tests
- Refactor code for better organization

## License

This is a personal learning project. Use it however you like!

## Acknowledgments

- Launch data provided by [ThespaceDevs Launch Library](https://thespacedevs.com/)
- Built as a practice project after completing the boot.dev Go course
