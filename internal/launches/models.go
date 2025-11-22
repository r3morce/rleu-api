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
