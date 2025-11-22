# Upcoming Features

## Planned

### Caching
- In-memory cache with 5-10 minute TTL
- Reduce API calls to ThespaceDevs

### Additional Endpoints
- `/launches/recent` - Recently completed launches
- `/launches/:id` - Specific launch details

### CORS Support
- Enable frontend access from different domains

### Rate Limiting
- Protect against abuse
- Stay within API limits

## Completed

### ✅ Filter European Launches
- Filter launches from European spaceports (Kourou, Esrange, SaxaVord, etc.)
- Configurable location IDs and country codes
- Returns European launches when enabled in config

### ✅ API Key Authentication
- Secure API with authentication middleware
- DEBUG mode for development (skips auth)
- Environment-based configuration
