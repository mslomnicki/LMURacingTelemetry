# LMU Racing Telemetry Monitor

A real-time telemetry monitoring system for Le Mans Ultimate with live data visualization and CSV logging.

## Features

- **Real-time Multi-Driver Monitoring**: Display all drivers simultaneously with live telemetry data including position, lap times, speed, and status
- **Fullscreen Display Modes**: Toggle fullscreen view for drivers or statistics panels for better visibility
- **Historical Statistics Tracking**: Track best lap times, sector times, and maximum speeds for each driver
- **CSV Data Logging**: Automatic logging of all telemetry data with organized filenames
- **Session Information**: Live display of track name, session type, weather conditions, and temperatures
- **Clean Terminal Interface**: Multi-panel interface optimized for terminal viewing

## Installation

### Download Pre-built Binaries

Download the latest release for your platform from the [Releases page](https://github.com/mslomnicki/LMURacingTelemetry/releases):

- **Windows AMD64**: `lmu-racing-telemetry-windows-amd64.exe`
- **Linux AMD64**: `lmu-racing-telemetry-linux-amd64`
- **macOS AMD64**: `lmu-racing-telemetry-macos-amd64`
- **macOS ARM64**: `lmu-racing-telemetry-macos-arm64`

## How to Run

### Basic Usage

```bash
# Connect to localhost (default)
./lmu-racing-telemetry

# Connect to specific IP address
./lmu-racing-telemetry -host 192.168.0.121

# Connect to custom port
./lmu-racing-telemetry -host 192.168.0.121 -port 8080
```

### Keyboard Controls

- **Ctrl+C** or **Q** - Quit the application
- **F** - Toggle fullscreen view for drivers panel
- **S** - Toggle fullscreen view for statistics panel

## Display Panels

The interface is divided into three main sections:

1. **Session Info Panel** (Top)
   - Track name and session type
   - Event time
   - Number of connected vehicles
   - Track and air temperature
   - Rain percentage

2. **All Drivers - Live Data Panel** (Middle)
   - Current position
   - Driver name and vehicle details
   - Laps completed
   - Current lap time
   - Best lap time
   - Current speed
   - Status (pit, flags, etc.)

3. **Driver Statistics & Records Panel** (Bottom)
   - Best lap times (official and calculated)
   - Best sector times (S1, S2, S3)
   - Maximum speeds
   - Per-driver historical records

## CSV Output

CSV files are automatically created with the format:
```
YYYY-MM-DD_HH-MM-SS_TrackName_SessionType_telemetry.csv
```

Example: `2025-10-10_16-40-39_Bahrain_International_Circuit_PRACTICE1_telemetry.csv`

The CSV file contains semicolon-delimited data with fields for driver name, vehicle, car class, laps, speeds, and all timing information.

## Requirements

- Le Mans Ultimate or compatible racing simulator with WebSocket telemetry enabled
- Terminal with color support for best viewing experience

## License

See LICENSE file for details.

## Author

Copyright (C) 2025 Marek Słomnicki <marek@slomnicki.net>
