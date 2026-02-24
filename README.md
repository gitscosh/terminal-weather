# terminal-weather

Weather in your terminal — modernized.

This project is a **modern rewrite inspired by**
https://github.com/genuinetools/weather.

The original project depended on the DarkSky API and a companion server.
DarkSky was shut down in 2023, which made the original architecture unusable.

This repository keeps the spirit and terminal UX of the original tool,
while rebuilding the internals around a simple, self-contained CLI.

---
## Credits & Attribution

This project is heavily inspired by:

https://github.com/genuinetools/weather

The original project and terminal UI concepts were created by the
Genuinetools authors. This repository preserves the MIT license and
credit while modernizing the implementation for current APIs.

This project is a personal rewrite and is **not an official continuation**
of the original project.

---
## What Changed

### Removed
The original architecture required:

• DarkSky API (shut down)  
• Google Geocoding API  
• A hosted weather server  
• API keys and configuration  

All of this has been removed.

---
### Added / Modernized

This project is now a **pure CLI** with no server and no API keys.

It now uses:

| Feature | Provider |
|---|---|
| Weather data | Open-Meteo |
| Geocoding | Open-Meteo |
| Auto-location | IP based lookup |

Key improvements:

• No API keys required  
• No server required  
• Faster startup and simpler architecture  
• Works entirely as a standalone CLI  
• Modern Go project layout (`cmd/ + internal/`)  
• HTTP timeouts and improved reliability  
• Cleaner output formatting  
• Maintained as a personal project

---

## Installation

### Build from source

Requires Go 1.22+

```bash
git clone https://github.com/gitscosh/terminal-weather
cd terminal-weather
go build ./cmd/weather
```
### Run:

```bash
./weather
```
---
### Usage

```bash
weather [location]
```
If no location is provided, the tool attempts IP-based auto-location.

#### Flags

| Flag              | Description                            |
| ----------------- | -------------------------------------- |
| `-l, --location`  | City, state, or ZIP code               |
| `-u, --units`     | Units: `auto`, `us`, `si`, `ca`, `uk2` |
| `-d, --days`      | Number of forecast days                |
| `--hide-icon`     | Hide ASCII art icon                    |
| `--ignore-alerts` | Hide weather alerts                    |
| `--json`          | Output raw JSON                        |

---
#### Examples

Get local weather:

```bash
weather
```

#### Specify location:

```bash
weather "San Diego"
weather -l "Paris, France"
```

#### Metric units:

```bash
weather -u si "Berlin"
```

#### Forecast:

```bash
weather -d 3 "New York"
```

---
## Project Status

This is a hobby / personal project.

It is stable and usable, but not intended to be a full weather platform.

Future improvements may include:

- Better unit auto-detection
- Weather alerts support
- Improved forecast handling
- Output polish and CLI refinements

---
## License

MIT License.

See LICENSE for details.