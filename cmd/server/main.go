package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"rleu-api/config"
	"rleu-api/internal/launches"
)

func main() {
	config.LoadEnv()
	cfg := config.Get()

	// Print debug status
	if os.Getenv("DEBUG") == "true" {
		log.Println("DEBUG MODE ENABLED - API throttle status will be checked after each request")
		log.Printf("Config loaded: API limit=%d, Server port=%d", cfg.API.Limit, cfg.Server.Port)
	}

	http.HandleFunc("/launches", authMiddleware(handleLaunches))

	serverAddr := cfg.GetServerAddress()
	fmt.Printf("🚀 Rocket Launch EU API is running on http://localhost%s\n", serverAddr)
	fmt.Printf("📡 Try visiting: http://localhost%s/launches\n", serverAddr)

	log.Fatal(http.ListenAndServe(serverAddr, nil))
}

// authMiddleware validates the API key from the Authorization header
func authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Skip authentication in DEBUG mode
		if os.Getenv("DEBUG") == "true" {
			log.Println("DEBUG MODE: Skipping authentication")
			next(w, r)
			return
		}

		cfg := config.Get()

		// Check if API key is configured
		if cfg.API.Key == "" {
			log.Println("Unauthorized: API key not configured on server")
			http.Error(w, "Unauthorized: Server not configured", http.StatusInternalServerError)
			return
		}

		// Get the Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			log.Println("Unauthorized: Missing Authorization header")
			http.Error(w, "Unauthorized: Missing API key", http.StatusUnauthorized)
			return
		}

		// Expected format: "Bearer <api_key>" or "Token <api_key>"
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 {
			log.Println("Unauthorized: Invalid Authorization header format")
			http.Error(w, "Unauthorized: Invalid API key format", http.StatusUnauthorized)
			return
		}

		providedKey := parts[1]
		expectedKey := cfg.API.Key

		// Compare the API keys
		if providedKey != expectedKey {
			log.Println("Unauthorized: Invalid API key")
			http.Error(w, "Unauthorized: Invalid API key", http.StatusUnauthorized)
			return
		}

		// Valid API key, proceed to the handler
		next(w, r)
	}
}

// handleLaunches fetches launch data from ThespaceDevs and returns it as JSON
func handleLaunches(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	log.Println("Received request for /launches")

	launchData, err := launches.FetchUpcomingLaunches()
	if err != nil {
		log.Printf("Error fetching launches: %v", err)
		http.Error(w, "Failed to fetch launches", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(launchData)
	if err != nil {
		log.Printf("Error encoding JSON: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}

	log.Printf("Successfully returned %d launches", len(launchData.Results))
}
