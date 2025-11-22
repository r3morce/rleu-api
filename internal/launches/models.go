package launches

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
type Pad struct {
	Name     string   `json:"name"`
	Location Location `json:"location"`
}

type Location struct {
	Name        string  `json:"name"`
	CountryCode *string `json:"country_code"`
}

type MissionPatch struct {
	Name      string `json:"name"`
	ImageURL  string `json:"image_url"`
	Priority  int    `json:"priority"`
	Agency    *int   `json:"agency"`
}

// LaunchDetailed extends Launch with detailed fields
type LaunchDetailed struct {
	Name           string          `json:"name"`
	NET            string          `json:"net"`
	Pad            Pad             `json:"pad"`
	MissionPatches []MissionPatch  `json:"mission_patches"`
	Image          *Image          `json:"image"`
}

// CompactLaunch is the simplified response for the API
type CompactLaunch struct {
	Name           string   `json:"name"`
	StartTime      string   `json:"start_time"`
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

// ToCompact converts a LaunchDetailed to CompactLaunch
func (ld *LaunchDetailed) ToCompact() CompactLaunch {
	compact := CompactLaunch{
		Name:      ld.Name,
		StartTime: ld.NET,
		Location:  ld.Pad.Location.Name,
		Country:   "",
	}

	// Add country if available from API
	if ld.Pad.Location.CountryCode != nil {
		compact.Country = *ld.Pad.Location.CountryCode
	} else {
		// Fall back to location name mapping
		compact.Country = GetCountryCode(ld.Pad.Location.Name)
	}

	// Add mission patch URLs
	for _, patch := range ld.MissionPatches {
		if patch.ImageURL != "" {
			compact.MissionPatches = append(compact.MissionPatches, patch.ImageURL)
		}
	}

	// Add image URL if available
	if ld.Image != nil {
		compact.ImageURL = ld.Image.ImageURL
	}

	return compact
}
