# LMU Racing Telemetry Monitor

A comprehensive telemetry monitoring system for racing simulators with real-time data visualization and CSV logging.

## Features

- **Real-time Multi-Driver Monitoring**: Display all drivers simultaneously with live telemetry data
- **Historical Statistics Tracking**: Track best lap times, sector times, and maximum speeds per driver
- **CSV Data Logging**: Automatic logging with date/track/session-based filenames
- **Flexible Connectivity**: Connect to any racing simulator WebSocket endpoint
- **Clean Terminal UI**: Multi-panel interface with session info, live data, and statistics

## Installation

### Download Pre-built Binaries

Download the latest release for your platform from the [Releases page](https://github.com/mslomnicki/LMURacingTelemetry/releases):

- **Windows AMD64**: `lmu-racing-telemetry-windows-amd64.exe`
- **Linux AMD64**: `lmu-racing-telemetry-linux-amd64`
- **macOS AMD64**: `lmu-racing-telemetry-macos-amd64`
- **macOS ARM64**: `lmu-racing-telemetry-macos-arm64`

### Build from Source

```bash
# Clone the repository
git clone https://github.com/mslomnicki/LMURacingTelemetry.git
cd LMURacingTelemetry

# Build the application
go build -o lmu-racing-telemetry .
```

## Usage

```bash
# Connect to localhost (default)
./lmu-racing-telemetry

# Connect to specific IP address
./lmu-racing-telemetry -host 192.168.0.121

# Connect to custom port
./lmu-racing-telemetry -host 192.168.0.121 -port 8080

# Show help
./lmu-racing-telemetry -help
```

## Project Structure

```
├── main.go                     # Application entry point
├── pkg/
│   ├── models/
│   │   └── types.go           # Data structures for WebSocket messages
│   ├── telemetry/
│   │   └── monitor.go         # Core telemetry monitoring logic
│   ├── ui/
│   │   └── display.go         # Terminal user interface components
│   └── logger/
│       └── csv.go             # CSV data logging functionality
└── go.mod
```

## Package Overview

### `pkg/models`
Contains all data structures for WebSocket messages, driver data, and session information.

### `pkg/telemetry`
Core telemetry monitoring system that:
- Manages WebSocket connections
- Processes incoming telemetry data
- Maintains driver statistics
- Orchestrates UI updates and CSV logging

### `pkg/ui`
Terminal user interface using tview library:
- Session information panel
- Live driver data table
- Historical statistics display
- Keyboard controls (Ctrl+C to exit)

### `pkg/logger`
CSV logging functionality:
- Automatic filename generation based on date/track/session
- Real-time data logging for all drivers
- Proper file handling and cleanup

## CSV Output

CSV files are automatically created with the format:
`YYYY-MM-DD_HH-MM-SS_TrackName_SessionType_telemetry.csv`

Example: `2025-10-05_19-00-23_Silverstone_Practice_telemetry.csv`

**CSV Format:**
- **Delimiter**: Semicolon (`;`)
- **Fields**: DriverName, VehicleName, CarClass, LapsCompleted, MaxSpeed, BestLapTime, BestSector1, BestSector2, BestSector3
- **Timing Format**: MM:SS.sss (e.g., `1:23.456`)

## Dependencies

- `github.com/gorilla/websocket` - WebSocket client
- `github.com/rivo/tview` - Terminal UI framework
- `github.com/gdamore/tcell/v2` - Terminal control library

## WebSocket Protocol

The application expects WebSocket messages in the format:
```json
{
  "type": "standings|sessionInfo",
  "body": { ... }
}
```

Supports message types:
- `standings`: Driver telemetry data
- `sessionInfo`: Track and session information