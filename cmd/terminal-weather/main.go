package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/gitscosh/terminal-weather/internal/forecast"
	"github.com/gitscosh/terminal-weather/internal/geo"
	"github.com/gitscosh/terminal-weather/internal/openmeteo"
	"github.com/gitscosh/terminal-weather/internal/output"
	"github.com/mitchellh/colorstring"
)

var (
	location     string
	units        string
	days         int
	ignoreAlerts bool
	hideIcon     bool
	noForecast   bool
	jsonOut      bool
)

//go:generate go run ../../internal/icons/generate.go

func main() {
	// --- FLAGS ---
	flag.StringVar(&location, "location", "", "Location (city or zip). Example: weather -location \"San Diego\"")
	flag.StringVar(&location, "l", "", "Location shorthand")

	flag.StringVar(&units, "units", "auto", "Units (auto, us, si, ca, uk2)")
	flag.StringVar(&units, "u", "auto", "Units shorthand")

	flag.IntVar(&days, "days", 0, "Number of forecast days")
	flag.IntVar(&days, "d", 0, "Days shorthand")

	flag.BoolVar(&ignoreAlerts, "ignore-alerts", false, "Ignore alerts")
	flag.BoolVar(&hideIcon, "hide-icon", false, "Hide icon output")
	flag.BoolVar(&noForecast, "no-forecast", false, "Hide hourly forecast")
	flag.BoolVar(&jsonOut, "json", false, "Output raw JSON")

	flag.Parse()

	// Allow: weather 00000  OR  weather "City Name"
	if location == "" && flag.NArg() > 0 {
		location = flag.Arg(0)
	}

	var fc forecast.Forecast

	// --- LOCATION RESOLUTION ---
	if location == "" {
		coords, err := geo.DetectLocationByIP()
		if err != nil {
			log.Fatal(err)
		}

		fc, err = openmeteo.GetOpenMeteoForecast(coords.Latitude, coords.Longitude, units)
		if err != nil {
			log.Fatal(err)
		}

		location = fmt.Sprintf("%s, %s", coords.City, coords.Region)

	} else {
		coords, err := geo.GeocodeLocation(location)
		if err != nil {
			printError(err)
		}

		fc, err = openmeteo.GetOpenMeteoForecast(coords.Latitude, coords.Longitude, units)
		if err != nil {
			printError(err)
		}

		location = fmt.Sprintf("%s, %s", coords.City, coords.Region)
	}

	// --- JSON MODE ---
	if jsonOut {
		jsn, err := json.MarshalIndent(&fc, "", "  ")
		if err != nil {
			printError(err)
		}
		fmt.Println(string(jsn))
		return
	}

	// --- PRINT OUTPUT ---
	if err := output.PrintCurrent(fc, location, ignoreAlerts, hideIcon); err != nil {
		printError(err)
	}

	if days > 0 {
		output.PrintDaily(fc, days)
	}
}

func printError(err error) {
	fmt.Println(colorstring.Color("[red]" + err.Error()))
	os.Exit(1)
}
