package launches

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"rleu-api/config"
)

// PadCache stores pad/location information
type PadCache struct {
	mu      sync.RWMutex
	pads    map[int]PadInfo
	lastFetch time.Time
}

// PadInfo contains cached pad and location data
type PadInfo struct {
	PadName      string
	LocationName string
	CountryCode  string
}

var (
	padCache     *PadCache
	cacheOnce    sync.Once
)

// GetPadCache returns the singleton pad cache instance
func GetPadCache() *PadCache {
	cacheOnce.Do(func() {
		padCache = &PadCache{
			pads: make(map[int]PadInfo),
		}
	})
	return padCache
}

// FetchAndCachePads fetches pad data from the API and caches it
func (pc *PadCache) FetchAndCachePads() error {
	pc.mu.Lock()
	defer pc.mu.Unlock()

	cfg := config.Get()

	// Build URL with detailed mode to get pad information
	url := fmt.Sprintf("%s?limit=%d&mode=detailed&format=%s",
		cfg.API.BaseURL,
		cfg.API.Limit,
		cfg.API.Format,
	)

	// Add European location filtering if enabled
	if cfg.EuropeanPads.Enabled && len(cfg.EuropeanPads.LocationIDs) > 0 {
		url += "&location__ids=" + joinInts(cfg.EuropeanPads.LocationIDs, ",")
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make API request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API returned status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	var detailedResponse LaunchResponseDetailed
	err = json.Unmarshal(body, &detailedResponse)
	if err != nil {
		return fmt.Errorf("failed to parse JSON response: %w", err)
	}

	// Cache pad information
	for _, launch := range detailedResponse.Results {
		if launch.Pad.ID > 0 {
			info := PadInfo{
				PadName:      launch.Pad.Name,
				LocationName: launch.Pad.Location.Name,
			}

			// Get country code from API
			if launch.Pad.Country != nil && launch.Pad.Country.Alpha2Code != "" {
				info.CountryCode = launch.Pad.Country.Alpha2Code
			} else if launch.Pad.Location.Country != nil && launch.Pad.Location.Country.Alpha2Code != "" {
				info.CountryCode = launch.Pad.Location.Country.Alpha2Code
			}

			pc.pads[launch.Pad.ID] = info
		}
	}

	pc.lastFetch = time.Now()

	if os.Getenv("DEBUG") == "true" {
		log.Printf("DEBUG: Cached %d pads", len(pc.pads))
	}

	return nil
}

// GetPadInfo retrieves cached pad information
func (pc *PadCache) GetPadInfo(padID int) (PadInfo, bool) {
	pc.mu.RLock()
	defer pc.mu.RUnlock()

	info, ok := pc.pads[padID]
	return info, ok
}

// IsStale checks if cache needs refreshing (older than 24 hours)
func (pc *PadCache) IsStale() bool {
	pc.mu.RLock()
	defer pc.mu.RUnlock()

	if pc.lastFetch.IsZero() {
		return true
	}

	return time.Since(pc.lastFetch) > 24*time.Hour
}
