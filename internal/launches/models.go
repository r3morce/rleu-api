package launches

import (
	"strings"

	"rleu-api/config"
)

// Data structures matching the ThespaceDevs API response

// LaunchResponse is the top-level API response
type LaunchResponse struct {
	Count    int      `json:"count"`
	Next     *string  `json:"next"`
	Previous *string  `json:"previous"`
	Results  []Launch `json:"results"`
}

// Launch represents a single rocket launch
type Launch struct {
	ID               string       `json:"id"`
	URL              string       `json:"url"`
	Name             string       `json:"name"`
	ResponseMode     string       `json:"response_mode"`
	Slug             string       `json:"slug"`
	LaunchDesignator *string      `json:"launch_designator"`
	Status           Status       `json:"status"`
	LastUpdated      string       `json:"last_updated"`
	NET              string       `json:"net"` // "No Earlier Than" - planned launch time
	NetPrecision     NetPrecision `json:"net_precision"`
	WindowEnd        string       `json:"window_end"`
	WindowStart      string       `json:"window_start"`
	Image            *Image       `json:"image"`
	Infographic      *string      `json:"infographic"`
}

// Status represents launch status
type Status struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Abbrev      string `json:"abbrev"`
	Description string `json:"description"`
}

// NetPrecision indicates launch time accuracy
type NetPrecision struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Abbrev      string `json:"abbrev"`
	Description string `json:"description"`
}

// Image contains launch image data
type Image struct {
	ID           int      `json:"id"`
	Name         string   `json:"name"`
	ImageURL     string   `json:"image_url"`
	ThumbnailURL string   `json:"thumbnail_url"`
	Credit       *string  `json:"credit"`
	License      License  `json:"license"`
	SingleUse    bool     `json:"single_use"`
	Variants     []string `json:"variants"`
}

// License contains image licensing info
type License struct {
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	Priority int     `json:"priority"`
	Link     *string `json:"link"`
}

// Detailed API response structures
type Country struct {
	Alpha2Code string `json:"alpha_2_code"`
	Alpha3Code string `json:"alpha_3_code"`
	Name       string `json:"name"`
}

type Pad struct {
	ID       int      `json:"id"`
	Name     string   `json:"name"`
	Country  *Country `json:"country"`
	Location Location `json:"location"`
}

type Location struct {
	ID      int      `json:"id"`
	Name    string   `json:"name"`
	Country *Country `json:"country"`
}

type MissionPatch struct {
	Name      string `json:"name"`
	ImageURL  string `json:"image_url"`
	Priority  int    `json:"priority"`
	Agency    *int   `json:"agency"`
}

// SimplePad is the minimal pad structure in normal/list mode
type SimplePad struct {
	ID int `json:"id"`
}

// LaunchDetailed extends Launch with detailed fields
type LaunchDetailed struct {
	Name           string          `json:"name"`
	NET            string          `json:"net"`
	WindowStart    string          `json:"window_start"`
	WindowEnd      string          `json:"window_end"`
	Pad            Pad             `json:"pad"`
	MissionPatches []MissionPatch  `json:"mission_patches"`
	Image          *Image          `json:"image"`
}

// LaunchNormal represents a launch in normal/list mode (minimal data)
type LaunchNormal struct {
	Name           string         `json:"name"`
	NET            string         `json:"net"`
	WindowStart    string         `json:"window_start"`
	WindowEnd      string         `json:"window_end"`
	Pad            SimplePad      `json:"pad"`
	MissionPatches []MissionPatch `json:"mission_patches"`
	Image          *Image         `json:"image"`
}

// LaunchResponseNormal is the response for normal/list mode
type LaunchResponseNormal struct {
	Count    int            `json:"count"`
	Next     *string        `json:"next"`
	Previous *string        `json:"previous"`
	Results  []LaunchNormal `json:"results"`
}

// CompactLaunch is the simplified response for the API
type CompactLaunch struct {
	Name           string   `json:"name"`
	NET            string   `json:"net"`
	WindowStart    string   `json:"window_start"`
	WindowEnd      string   `json:"window_end"`
	Location       string   `json:"location"`
	Country        string   `json:"country"`
	MissionPatches []string `json:"mission_patches,omitempty"`
	ImageURL       string   `json:"image_url,omitempty"`
}

// CompactLaunchResponse is the compact API response
type CompactLaunchResponse struct {
	Count   int             `json:"count"`
	Results []CompactLaunch `json:"results"`
}

// transformImageURL replaces -dev with -prod in image URLs if configured
func transformImageURL(url string) string {
	cfg := config.Get()
	if cfg.API.UseProductionImage {
		return strings.ReplaceAll(url, "-dev", "-prod")
	}
	return url
}

// ToCompact converts a LaunchDetailed to CompactLaunch
func (ld *LaunchDetailed) ToCompact() CompactLaunch {
	locationName := ""
	if ld.Pad.Location.Name != "" {
		locationName = ld.Pad.Location.Name
	} else {
		locationName = ld.Pad.Name
	}

	compact := CompactLaunch{
		Name:        ld.Name,
		NET:         ld.NET,
		WindowStart: ld.WindowStart,
		WindowEnd:   ld.WindowEnd,
		Location:    locationName,
		Country:     "",
	}

	// Get country code from API response
	// Priority: Pad.Country > Location.Country > fallback to mapping
	if ld.Pad.Country != nil && ld.Pad.Country.Alpha2Code != "" {
		compact.Country = ld.Pad.Country.Alpha2Code
	} else if ld.Pad.Location.Country != nil && ld.Pad.Location.Country.Alpha2Code != "" {
		compact.Country = ld.Pad.Location.Country.Alpha2Code
	} else {
		// Fall back to location name mapping for older data
		compact.Country = GetCountryCode(ld.Pad.Location.Name)
	}

	// Add mission patch URLs
	for _, patch := range ld.MissionPatches {
		if patch.ImageURL != "" {
			compact.MissionPatches = append(compact.MissionPatches, transformImageURL(patch.ImageURL))
		}
	}

	// Add image URL if available
	if ld.Image != nil {
		compact.ImageURL = transformImageURL(ld.Image.ImageURL)
	}

	return compact
}

// ToCompact converts a LaunchNormal to CompactLaunch using cached pad data
func (ln *LaunchNormal) ToCompact() CompactLaunch {
	compact := CompactLaunch{
		Name:        ln.Name,
		NET:         ln.NET,
		WindowStart: ln.WindowStart,
		WindowEnd:   ln.WindowEnd,
		Location:    "",
		Country:     "",
	}

	// Try to get pad info from cache
	cache := GetPadCache()
	if info, ok := cache.GetPadInfo(ln.Pad.ID); ok {
		compact.Location = info.LocationName
		compact.Country = info.CountryCode
	}

	// Add mission patch URLs
	for _, patch := range ln.MissionPatches {
		if patch.ImageURL != "" {
			compact.MissionPatches = append(compact.MissionPatches, transformImageURL(patch.ImageURL))
		}
	}

	// Add image URL if available
	if ln.Image != nil {
		compact.ImageURL = transformImageURL(ln.Image.ImageURL)
	}

	return compact
}

