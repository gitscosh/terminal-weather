package geo

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// Open-Meteo geocoding response
type openMeteoGeoResponse struct {
	Results []struct {
		Name      string  `json:"name"`
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
		Country   string  `json:"country"`
		Admin1    string  `json:"admin1"`
	} `json:"results"`
}

// GeocodeResult is the minimal location info the CLI needs.
type GeocodeResult struct {
	City      string
	Region    string
	Country   string
	Latitude  float64
	Longitude float64
}

// GeocodeLocation converts a city/zip string → lat/lon using Open-Meteo.
func GeocodeLocation(query string) (GeocodeResult, error) {
	endpoint := fmt.Sprintf(
		"https://geocoding-api.open-meteo.com/v1/search?name=%s&count=1",
		url.QueryEscape(query),
	)

	resp, err := http.Get(endpoint)
	if err != nil {
		return GeocodeResult{}, err
	}
	defer resp.Body.Close()

	var data openMeteoGeoResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return GeocodeResult{}, err
	}

	if len(data.Results) == 0 {
		return GeocodeResult{}, fmt.Errorf("location not found: %s", query)
	}

	r := data.Results[0]

	return GeocodeResult{
		City:      r.Name,
		Region:    r.Admin1,
		Country:   r.Country,
		Latitude:  r.Latitude,
		Longitude: r.Longitude,
	}, nil
}
