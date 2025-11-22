package launches

// locationToCountry maps European location names to ISO 3166-1 alpha-2 country codes
// This is used as a fallback when the API doesn't provide country_code directly
var locationToCountry = map[string]string{
	"Esrange Space Center":    "SE", // Sweden
	"SaxaVord Spaceport":      "GB", // United Kingdom
	"Sutherland Spaceport":    "GB", // United Kingdom
	"Andøya Spaceport":        "NO", // Norway
	"El Arenosillo Test Centre": "ES", // Spain
	"Guiana Space Centre, French Guiana": "FR", // France (French Guiana)
	"Interarmy Special Vehicles Test Centre, French Algeria": "FR", // France
	"Broglio Space Center, Kenya": "IT", // Italian facility
}

// GetCountryCode returns the country code for a given location name
// Returns empty string if location is not found in the mapping
func GetCountryCode(locationName string) string {
	if code, ok := locationToCountry[locationName]; ok {
		return code
	}
	return ""
}
