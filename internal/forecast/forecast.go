// Package forecast defines the provider-agnostic weather domain model.
//
// The structs in this file represent the internal weather schema used by the
// CLI and printer. The schema intentionally mirrors the original DarkSky API
// because the output/formatting layer depends on it.
//
// Weather providers (Open-Meteo, future NOAA, etc.) must translate their API
// responses into these structs.

package forecast

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// Forecast contains the information returned from the server
// when requesting the forecast.
type Forecast struct {
	Alerts    []Alert       `json:"alerts"`
	Currently Weather       `json:"currently"`
	Code      int           `json:"code"`
	Daily     TimeDelimited `json:"daily"`
	Error     string        `json:"error"`
	Flags     Flags         `json:"flags"`
	Hourly    TimeDelimited `json:"hourly"`
	Latitude  float64       `json:"latitude"`
	Longitude float64       `json:"longitude"`
	Offset    float64       `json:"offset"`
	Timezone  string        `json:"timezone"`
}

// Alert contains any weather alerts happening at the location.
type Alert struct {
	Description string `json:"description"`
	Expires     int64  `json:"expires"`
	Time        int64  `json:"time"`
	Title       string `json:"title"`
	URI         string `json:"uri"`
}

// Flags describes the flags on a forecast.
type Flags struct {
	Units string `json:"units"`
}

// Weather describes details about the weather for the location.
type Weather struct {
	ApparentTemperature        float64 `json:"apparentTemperature"`
	ApparentTemperatureMax     float64 `json:"apparentTemperatureMax"`
	ApparentTemperatureMaxTime int64   `json:"apparentTemperatureMaxTime"`
	ApparentTemperatureMin     float64 `json:"apparentTemperatureMin"`
	ApparentTemperatureMinTime int64   `json:"apparentTemperatureMinTime"`
	CloudCover                 float64 `json:"cloudCover"`
	DewPoint                   float64 `json:"dewPoint"`
	Humidity                   float64 `json:"humidity"`
	Icon                       string  `json:"icon"`
	NearestStormDistance       float64 `json:"nearestStormDistance"`
	NearestStormBearing        float64 `json:"nearestStormBearing"`
	Ozone                      float64 `json:"ozone"`
	PrecipIntensity            float64 `json:"precipIntensity"`
	PrecipIntensityMax         float64 `json:"precipIntensityMax"`
	PrecipIntensityMaxTime     int64   `json:"precipIntensityMaxTime"`
	PrecipProbability          float64 `json:"precipProbability"`
	PrecipType                 string  `json:"precipType"`
	Pressure                   float64 `json:"pressure"`
	Summary                    string  `json:"summary"`
	SunriseTime                int64   `json:"sunriseTime"`
	SunsetTime                 int64   `json:"sunsetTime"`
	Temperature                float64 `json:"temperature"`
	TemperatureMax             float64 `json:"temperatureMax"`
	TemperatureMaxTime         int64   `json:"temperatureMaxTime"`
	TemperatureMin             float64 `json:"temperatureMin"`
	TemperatureMinTime         int64   `json:"temperatureMinTime"`
	Time                       int64   `json:"time"`
	Visibility                 float64 `json:"visibility"`
	WindBearing                float64 `json:"windBearing"`
	WindSpeed                  float64 `json:"windSpeed"`
}

// TimeDelimited describes the data for the time series.
type TimeDelimited struct {
	Data    []Weather `json:"data"`
	Icon    string    `json:"icon"`
	Summary string    `json:"summary"`
}

// Request describes the request posted to the forecast api.
type Request struct {
	Latitude  float64  `json:"lat"`
	Longitude float64  `json:"lng"`
	Units     string   `json:"units"`
	Exclude   []string `json:"exclude"`
}

func GetFromServer(uri string, data Request) (Forecast, error) {
	jsonByte, err := json.Marshal(data)
	if err != nil {
		return Forecast{}, fmt.Errorf("marshaling forecast request failed: %v", err)
	}

	resp, err := http.Post(uri, "application/json", bytes.NewBuffer(jsonByte))
	if err != nil {
		return Forecast{}, fmt.Errorf("request to weather server failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Forecast{}, fmt.Errorf("weather server returned status: %v", resp.Status)
	}

	var fc Forecast
	if err := json.NewDecoder(resp.Body).Decode(&fc); err != nil {
		return Forecast{}, fmt.Errorf("decoding forecast response failed: %v", err)
	}

	return fc, nil
}
