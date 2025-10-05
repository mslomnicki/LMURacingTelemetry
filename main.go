package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/mslomnicki/LMURacingTelemetry/pkg/telemetry"
)

func main() {
	// Parse CLI arguments
	host := flag.String("host", "localhost", "WebSocket server hostname or IP address")
	port := flag.String("port", "6398", "WebSocket server port")
	flag.Parse()

	websocketURL := fmt.Sprintf("ws://%s:%s/websocket/controlpanel", *host, *port)

	// Create and configure telemetry monitor
	monitor := telemetry.NewMonitor(websocketURL)

	fmt.Printf("Starting LMU Racing Telemetry Monitor...\n")
	fmt.Printf("Connecting to %s\n", websocketURL)
	fmt.Printf("Press Ctrl+C to exit\n\n")

	// Run the telemetry monitor
	if err := monitor.Run(); err != nil {
		log.Fatalf("Error running telemetry monitor: %v", err)
	}
}
