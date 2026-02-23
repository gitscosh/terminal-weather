// Open-Meteo provider adapter.
//
// This file replaces the deprecated DarkSky API used by the original
// genuinetools/weather project.
//
// The CLI and printing code expect a DarkSky-compatible Forecast schema.
// Instead of rewriting the entire application, we translate Open-Meteo’s
// API response into the existing Forecast/Weather structs.
//
// Open-Meteo docs:
// https://open-meteo.com/en/docs

package openmeteo

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/gitscosh/terminal-weather/internal/forecast"
	"github.com/gitscosh/terminal-weather/internal/httpx"
)

type openMeteoResponse struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Timezone  string  `json:"timezone"`

	Current struct {
		Time                string  `json:"time"`
		Temperature         float64 `json:"temperature_2m"`
		ApparentTemperature float64 `json:"apparent_temperature"`
		WeatherCode         int     `json:"weathercode"`
		WindSpeed           float64 `json:"wind_speed_10m"`
		WindDirection       float64 `json:"wind_direction_10m"`
		Visibility          float64 `json:"visibility"`
		Pressure            float64 `json:"pressure_msl"`
		CloudCover          float64 `json:"cloud_cover"`
		DewPoint            float64 `json:"dew_point_2m"`
	} `json:"current"`

	Hourly struct {
		Time              []string  `json:"time"`
		Temperature       []float64 `json:"temperature_2m"`
		Humidity          []float64 `json:"relative_humidity_2m"`
		PrecipProbability []float64 `json:"precipitation_probability"`
	} `json:"hourly"`

	Daily struct {
		Time    []string  `json:"time"`
		TempMax []float64 `json:"temperature_2m_max"`
		TempMin []float64 `json:"temperature_2m_min"`
	} `json:"daily"`
}

func GetOpenMeteoForecast(lat, lon float64, units string) (forecast.Forecast, error) {
	tempUnit := "fahrenheit"
	windUnit := "mph"
	precipUnit := "inch"

	switch units {
	case "si", "ca", "uk2":
		tempUnit = "celsius"
		windUnit = "kmh"
		precipUnit = "mm"
	}
	url := fmt.Sprintf(
		"https://api.open-meteo.com/v1/forecast?latitude=%f&longitude=%f&current=temperature_2m,apparent_temperature,weathercode,wind_speed_10m,wind_direction_10m,visibility,pressure_msl,cloud_cover,dew_point_2m&hourly=temperature_2m,relative_humidity_2m,precipitation_probability&daily=temperature_2m_max,temperature_2m_min&temperature_unit=%s&wind_speed_unit=%s&precipitation_unit=%s",
		lat, lon, tempUnit, windUnit, precipUnit,
	)

	resp, err := httpx.HTTPClient.Get(url)
	if err != nil {
		return forecast.Forecast{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return forecast.Forecast{}, err
	}

	var om openMeteoResponse
	if err := json.Unmarshal(body, &om); err != nil {
		return forecast.Forecast{}, err
	}

	return translateOpenMeteoToForecast(om, units), nil
}

func translateOpenMeteoToForecast(om openMeteoResponse, units string) forecast.Forecast {
	fc := forecast.Forecast{
		Latitude:  om.Latitude,
		Longitude: om.Longitude,
		Timezone:  om.Timezone,
	}
	fc.Flags.Units = units

	// --- CURRENT CONDITIONS ---
	fc.Currently.Time = time.Now().Unix()
	fc.Currently.Temperature = om.Current.Temperature
	fc.Currently.ApparentTemperature = om.Current.ApparentTemperature
	fc.Currently.WindSpeed = om.Current.WindSpeed
	fc.Currently.WindBearing = om.Current.WindDirection
	fc.Currently.CloudCover = om.Current.CloudCover / 100
	fc.Currently.DewPoint = om.Current.DewPoint
	fc.Currently.Pressure = om.Current.Pressure

	// meters → miles, DarkSky caps at 10mi
	visMiles := om.Current.Visibility / 1609.34
	if visMiles > 10 {
		visMiles = 10
	}
	fc.Currently.Visibility = visMiles

	// weather code → DarkSky style summary/icon
	fc.Currently.Summary = forecast.WeatherCodeToSummary(om.Current.WeatherCode)
	fc.Currently.Icon = forecast.SummaryToIcon(fc.Currently.Summary)

	// first hourly datapoint used as "current" humidity & precip
	if len(om.Hourly.Humidity) > 0 {
		fc.Currently.Humidity = om.Hourly.Humidity[0] / 100
	}
	if len(om.Hourly.PrecipProbability) > 0 {
		fc.Currently.PrecipProbability = om.Hourly.PrecipProbability[0] / 100
	}

	// today's high/low (DarkSky compatibility)
	if len(om.Daily.TempMax) > 0 {
		fc.Currently.TemperatureMax = om.Daily.TempMax[0]
	}
	if len(om.Daily.TempMin) > 0 {
		fc.Currently.TemperatureMin = om.Daily.TempMin[0]
	}

	// --- HOURLY ---
	fc.Hourly.Summary = fc.Currently.Summary
	fc.Hourly.Icon = fc.Currently.Icon

	hours := len(om.Hourly.Time)
	if hours > 24 {
		hours = 24
	}

	for i := 0; i < hours; i++ {
		fc.Hourly.Data = append(fc.Hourly.Data, forecast.Weather{
			Time:              time.Now().Unix() + int64(i*3600),
			Temperature:       om.Hourly.Temperature[i],
			Humidity:          om.Hourly.Humidity[i] / 100,
			WindSpeed:         om.Current.WindSpeed,
			Summary:           fc.Currently.Summary,
			Icon:              fc.Currently.Icon,
			PrecipProbability: om.Hourly.PrecipProbability[i] / 100,
		})
	}

	// --- DAILY ---
	days := len(om.Daily.Time)

	for i := 0; i < days; i++ {
		fc.Daily.Data = append(fc.Daily.Data, forecast.Weather{
			Time:           time.Now().Unix() + int64(i*86400),
			Temperature:    (om.Daily.TempMax[i] + om.Daily.TempMin[i]) / 2,
			TemperatureMax: om.Daily.TempMax[i],
			TemperatureMin: om.Daily.TempMin[i],
			Summary:        fc.Currently.Summary,
			Icon:           fc.Currently.Icon,
		})
	}

	fc.Daily.Summary = "Daily forecast"
	fc.Daily.Icon = fc.Currently.Icon

	return fc
}
