package launches

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"rleu-api/config"
)

// ThrottleStatus represents API throttle information
type ThrottleStatus struct {
	Limit     int `json:"limit"`
	Remaining int `json:"remaining"`
}

// FetchUpcomingLaunches calls ThespaceDevs API and returns upcoming launches
func FetchUpcomingLaunches() (*LaunchResponse, error) {
	cfg := config.Get()
	url := cfg.GetAPIURL()

	// Add European location filtering if enabled
	if cfg.EuropeanPads.Enabled && len(cfg.EuropeanPads.LocationIDs) > 0 {
		url += "&location__ids=" + joinInts(cfg.EuropeanPads.LocationIDs, ",")
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make API request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var launchResponse LaunchResponse
	err = json.Unmarshal(body, &launchResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON response: %w", err)
	}

	// Check API throttle status in debug mode
	if os.Getenv("DEBUG") == "true" {
		checkThrottleStatus()
	}

	return &launchResponse, nil
}

// joinInts converts integer slice to comma-separated string
func joinInts(ints []int, sep string) string {
	if len(ints) == 0 {
		return ""
	}

	result := make([]string, len(ints))
	for i, num := range ints {
		result[i] = fmt.Sprintf("%d", num)
	}

	return strings.Join(result, sep)
}

// checkThrottleStatus checks remaining API calls (debug mode only)
func checkThrottleStatus() {
	throttleURL := "https://lldev.thespacedevs.com/2.3.0/api-throttle/"

	req, err := http.NewRequest("GET", throttleURL, nil)
	if err != nil {
		log.Printf("DEBUG: Failed to create throttle request: %v", err)
		return
	}

	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("DEBUG: Failed to check throttle status: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("DEBUG: Throttle API returned status code: %d", resp.StatusCode)
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("DEBUG: Failed to read throttle response: %v", err)
		return
	}

	var throttle ThrottleStatus
	err = json.Unmarshal(body, &throttle)
	if err != nil {
		log.Printf("DEBUG: Failed to parse throttle response: %v", err)
		return
	}

	log.Printf("DEBUG: API Throttle - Limit: %d, Remaining: %d", throttle.Limit, throttle.Remaining)
}
