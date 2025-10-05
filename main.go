package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/mslomnicki/LMURacingTelemetry/pkg/telemetry"
	"github.com/mslomnicki/LMURacingTelemetry/pkg/ui"
)

func main() {
	// Parse CLI arguments
	host := flag.String("host", "localhost", "WebSocket server hostname or IP address")
	port := flag.String("port", "6398", "WebSocket server port")
	flag.Parse()

	websocketURL := fmt.Sprintf("ws://%s:%s/websocket/controlpanel", *host, *port)

	// Create and configure telemetry monitor
	monitor := telemetry.NewMonitor(websocketURL)

	fmt.Printf("Starting LMU Racing Telemetry Monitor %s...\n", ui.Version)
	fmt.Printf("Connecting to %s\n", websocketURL)
	fmt.Printf("Press Ctrl+C to exit\n\n")

	if err := monitor.Run(); err != nil {
		log.Fatalf("Error running telemetry monitor: %v", err)
	}
}
