package forecast

// weatherCodeToSummary converts Open-Meteo weather codes into
// DarkSky-style summaries used by the CLI output layer.
func WeatherCodeToSummary(code int) string {
	switch code {
	case 0:
		return "Clear"
	case 1, 2:
		return "Partly Cloudy"
	case 3:
		return "Cloudy"
	case 45, 48:
		return "Fog"
	case 51, 53, 55:
		return "Drizzle"
	case 61, 63, 65:
		return "Rain"
	case 71, 73, 75:
		return "Snow"
	case 95, 96, 99:
		return "Thunderstorm"
	default:
		return "Cloudy"
	}
}

// summaryToIcon converts a summary into a DarkSky-compatible icon string.
func SummaryToIcon(summary string) string {
	switch summary {
	case "Clear":
		return "clear-day"
	case "Partly Cloudy":
		return "partly-cloudy-day"
	case "Cloudy":
		return "cloudy"
	case "Rain":
		return "rain"
	case "Snow":
		return "snow"
	case "Fog":
		return "fog"
	case "Thunderstorm":
		return "thunderstorm"
	default:
		return "cloudy"
	}
}
