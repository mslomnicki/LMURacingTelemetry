package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/mslomnicki/LMURacingTelemetry/pkg/telemetry"
	"github.com/mslomnicki/LMURacingTelemetry/pkg/ui"
)

func main() {
	logFile, err := os.OpenFile("LMURacingTelemetry.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	log.SetOutput(logFile)

	host := flag.String("host", "localhost", "WebSocket server hostname or IP address")
	wsPort := flag.String("ws-port", "6398", "WebSocket server port")
	restPort := flag.String("rest-port", "6397", "REST API server port")
	flag.Parse()

	monitor := telemetry.NewMonitor(*host, *wsPort, *restPort)

	fmt.Printf("Starting LMU Racing Telemetry Monitor %s...\n", ui.Version)
	fmt.Printf("Connecting to ws://%s:%s and REST http://%s:%s\n", *host, *wsPort, *host, *restPort)
	fmt.Printf("Press Ctrl+C to exit\n\n")

	if err := monitor.Run(); err != nil {
		log.Fatalf("Error running telemetry monitor: %v", err)
	}
}
