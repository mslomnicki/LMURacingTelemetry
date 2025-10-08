package telemetry

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
	"github.com/mslomnicki/LMURacingTelemetry/pkg/logger"
	"github.com/mslomnicki/LMURacingTelemetry/pkg/models"
	"github.com/mslomnicki/LMURacingTelemetry/pkg/ui"
)

type DriverLapState struct {
	currentLapMaxSpeed   float64
	lastCompletedLaps    int
	lastValidTimeIntoLap float64
}

type Monitor struct {
	conn         *websocket.Conn
	display      *ui.Display
	csvLogger    *logger.CSVLogger
	drivers      map[string]*models.StandingsData
	driverStats  map[string]*models.DriverStats
	lapStates    map[string]*DriverLapState // Dodatkowy stan dla każdego kierowcy
	session      *models.SessionData
	websocketURL string
	reconnecting bool
	stopChan     chan struct{}
}

func NewMonitor(websocketURL string) *Monitor {
	return &Monitor{
		display:      ui.NewDisplay(),
		drivers:      make(map[string]*models.StandingsData),
		driverStats:  make(map[string]*models.DriverStats),
		lapStates:    make(map[string]*DriverLapState),
		websocketURL: websocketURL,
		stopChan:     make(chan struct{}),
	}
}

func (m *Monitor) Connect() error {
	var err error
	m.conn, _, err = websocket.DefaultDialer.Dial(m.websocketURL, nil)
	if err != nil {
		return fmt.Errorf("failed to connect to WebSocket: %w", err)
	}
	log.Printf("Connected to WebSocket: %s", m.websocketURL)
	return nil
}

func (m *Monitor) connectWithRetry() {
	backoff := time.Second
	maxBackoff := 30 * time.Second

	for {
		select {
		case <-m.stopChan:
			return
		default:
		}

		if !m.reconnecting {
			log.Printf("Attempting to connect to WebSocket...")
		} else {
			log.Printf("Attempting to reconnect to WebSocket...")
		}

		if err := m.Connect(); err != nil {
			if !m.reconnecting {
				log.Printf("Initial connection failed: %v. Retrying in %v...", err, backoff)
			} else {
				log.Printf("Reconnection failed: %v. Retrying in %v...", err, backoff)
			}

			select {
			case <-time.After(backoff):
				backoff = time.Duration(float64(backoff) * 1.5)
				if backoff > maxBackoff {
					backoff = maxBackoff
				}
			case <-m.stopChan:
				return
			}
			continue
		}

		if m.reconnecting {
			log.Printf("Reconnected successfully!")
			m.reconnecting = false
		}

		m.listenForMessages()

		if m.conn != nil {
			m.conn.Close()
		}

		select {
		case <-m.stopChan:
			return
		default:
			m.reconnecting = true
			backoff = time.Second
		}
	}
}

func (m *Monitor) Run() error {
	m.display.Setup()

	go m.connectWithRetry()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	go func() {
		<-interrupt
		close(m.stopChan)
		m.cleanup()
		m.display.Stop()
	}()

	return m.display.Run()
}

func (m *Monitor) listenForMessages() {
	defer func() {
		if m.conn != nil {
			m.conn.Close()
		}
	}()

	for {
		select {
		case <-m.stopChan:
			return
		default:
		}

		_, message, err := m.conn.ReadMessage()
		if err != nil {
			log.Printf("WebSocket read error: %v", err)
			return
		}

		var wsMsg models.WSMessage
		if err := json.Unmarshal(message, &wsMsg); err != nil {
			log.Printf("Error unmarshaling WebSocket message: %v", err)
			continue
		}

		bodyBytes, err := json.Marshal(wsMsg.Body)
		if err != nil {
			log.Printf("Error marshaling message body: %v", err)
			continue
		}

		m.handleMessage(wsMsg.Type, bodyBytes)
	}
}

func (m *Monitor) handleMessage(msgType string, body json.RawMessage) {
	switch msgType {
	case "standings":
		m.handleStandings(body)
	case "sessionInfo":
		m.handleSessionInfo(body)
	case "standingsHistory":
		// Currently not processed
	default:
		log.Printf("Unsupported message type: %s, body: %s", msgType, string(body))
	}

	m.updateDisplay()
}

func (m *Monitor) handleStandings(body json.RawMessage) {
	var standings []models.StandingsData
	if err := json.Unmarshal(body, &standings); err != nil {
		log.Printf("Error unmarshaling standings: %v", err)
		return
	}

	for _, driver := range standings {
		key := driver.DriverName
		m.drivers[key] = &driver
		m.updateDriverStats(&driver)
		m.logDriverData(&driver)
	}
}

func (m *Monitor) handleSessionInfo(body json.RawMessage) {
	var session models.SessionData
	if err := json.Unmarshal(body, &session); err != nil {
		log.Printf("Error unmarshaling session info: %v", err)
		return
	}

	m.session = &session

	if m.csvLogger == nil && m.session != nil {
		var err error
		m.csvLogger, err = logger.NewCSVLogger(m.session)
		if err != nil {
			log.Printf("Error initializing CSV logger: %v", err)
		} else {
			log.Printf("CSV logging initialized for %s - %s", m.session.TrackName, m.session.Session)
		}
	}
}

func (m *Monitor) updateDriverStats(driver *models.StandingsData) {
	key := driver.DriverName

	lapState, lapStateExists := m.lapStates[key]
	if !lapStateExists {
		lapState = &DriverLapState{
			lastCompletedLaps: driver.LapsCompleted,
		}
		m.lapStates[key] = lapState
	}

	const maxReasonableTimeIntoLap = 600.0 // 10 minutes in seconds
	if driver.TimeIntoLap < 0 || driver.TimeIntoLap > maxReasonableTimeIntoLap {
		if lapState.lastValidTimeIntoLap >= 0 {
			driver.TimeIntoLap = lapState.lastValidTimeIntoLap
		} else {
			driver.TimeIntoLap = 0
		}
	}
	lapState.lastValidTimeIntoLap = driver.TimeIntoLap

	stats, exists := m.driverStats[key]
	if !exists {
		stats = &models.DriverStats{
			DriverName:  driver.DriverName,
			VehicleName: driver.VehicleName,
			CarClass:    driver.CarClass,
		}
		m.driverStats[key] = stats
	}

	stats.DriverName = driver.DriverName
	stats.VehicleName = driver.VehicleName
	stats.CarClass = driver.CarClass
	stats.Position = driver.Position
	stats.LapsCompleted = driver.LapsCompleted
	stats.LastUpdate = time.Now()

	currentSpeed := driver.CarVelocity.Velocity * 3.6 // Convert to km/h

	if currentSpeed > stats.MaxSpeed {
		stats.MaxSpeed = currentSpeed
	}

	if currentSpeed > lapState.currentLapMaxSpeed {
		lapState.currentLapMaxSpeed = currentSpeed
	}

	if driver.LapsCompleted > lapState.lastCompletedLaps {
		if driver.LastLapTime > 0 && (stats.BestLapTimeCalculated == 0 || driver.LastLapTime < stats.BestLapTimeCalculated) {
			stats.MaxSpeedOnBestLapCalc = lapState.currentLapMaxSpeed
			stats.BestLapTimeCalculated = driver.LastLapTime

			stats.BestSector1Calculated = driver.LastSectorTime1
			stats.BestSector2Calculated = driver.LastSectorTime2 - driver.LastSectorTime1
			stats.BestSector3Calculated = driver.LastLapTime - driver.LastSectorTime2
		}
		lapState.currentLapMaxSpeed = currentSpeed
		lapState.lastCompletedLaps = driver.LapsCompleted
	}

	stats.BestLapTime = driver.BestLapTime
	stats.BestSector1 = driver.BestSectorTime1
	stats.BestSector2 = driver.BestSectorTime2 - driver.BestSectorTime1
	stats.BestSector3 = driver.BestLapTime - driver.BestSectorTime2
}

// logDriverData updates driver data in CSV logger (most recent state only)
func (m *Monitor) logDriverData(driver *models.StandingsData) {
	if m.csvLogger == nil {
		return
	}

	// Get stats for this driver
	key := driver.DriverName
	if stats, exists := m.driverStats[key]; exists {
		m.csvLogger.UpdateDriver(driver, stats)
	}
}

// updateDisplay refreshes all UI components
func (m *Monitor) updateDisplay() {
	m.display.UpdateSession(m.session)
	m.display.UpdateDrivers(m.drivers)
	m.display.UpdateStats(m.driverStats)
	m.display.Draw()
}

// cleanup performs cleanup operations on shutdown
func (m *Monitor) cleanup() {
	log.Println("Shutting down...")

	if m.csvLogger != nil {
		if err := m.csvLogger.Close(); err != nil {
			log.Printf("Error closing CSV logger: %v", err)
		} else {
			log.Println("CSV logging stopped")
		}
	}

	if m.conn != nil {
		m.conn.Close()
	}
}
