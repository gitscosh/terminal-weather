package geo

import (
	"encoding/json"
	"fmt"

	"github.com/gitscosh/terminal-weather/internal/httpx"
)

type ipAPIResponse struct {
	City      string  `json:"city"`
	Region    string  `json:"region"`
	Country   string  `json:"country_name"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type Coordinates struct {
	Latitude  float64
	Longitude float64
	City      string
	Region    string
	Country   string
}

// DetectLocationByIP resolves approximate user coordinates using ipapi.
func DetectLocationByIP() (Coordinates, error) {
	resp, err := httpx.HTTPClient.Get("https://ipapi.co/json/")
	if err != nil {
		return Coordinates{}, err
	}
	defer resp.Body.Close()

	var ip ipAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&ip); err != nil {
		return Coordinates{}, err
	}

	if ip.Latitude == 0 || ip.Longitude == 0 {
		return Coordinates{}, fmt.Errorf("could not determine location from IP")
	}

	return Coordinates{
		Latitude:  ip.Latitude,
		Longitude: ip.Longitude,
		City:      ip.City,
		Region:    ip.Region,
		Country:   ip.Country,
	}, nil
}
