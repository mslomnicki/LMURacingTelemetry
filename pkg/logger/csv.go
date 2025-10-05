package logger

import (
	"encoding/csv"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/mslomnicki/LMURacingTelemetry/pkg/models"
)

// CSVLogger handles CSV file logging for telemetry data
type CSVLogger struct {
	filename    string
	driverData  map[string]*models.StandingsData
	driverStats map[string]*models.DriverStats
	session     *models.SessionData
}

// NewCSVLogger creates a new CSV logger with filename based on session data
func NewCSVLogger(session *models.SessionData) (*CSVLogger, error) {
	if session == nil {
		return nil, fmt.Errorf("session data is required")
	}

	// Create filename with date, track, and session
	now := time.Now()
	trackName := strings.ReplaceAll(session.TrackName, " ", "_")
	sessionName := strings.ReplaceAll(session.Session, " ", "_")
	filename := fmt.Sprintf("%s_%s_%s_telemetry.csv",
		now.Format("2006-01-02_15-04-05"),
		trackName,
		sessionName)

	return &CSVLogger{
		filename:    filename,
		driverData:  make(map[string]*models.StandingsData),
		driverStats: make(map[string]*models.DriverStats),
		session:     session,
	}, nil
}

// UpdateDriver updates the driver data (replacing previous data for this driver)
func (l *CSVLogger) UpdateDriver(driver *models.StandingsData, stats *models.DriverStats) {
	key := fmt.Sprintf("%d", driver.SlotID)
	l.driverData[key] = driver
	l.driverStats[key] = stats

	// Write current state to CSV file immediately
	if err := l.WriteCurrentState(); err != nil {
		// Log error but don't stop the program
		fmt.Printf("Error writing CSV file: %v\n", err)
	}
}

// WriteCurrentState writes the current state of all drivers to CSV file
func (l *CSVLogger) WriteCurrentState() error {
	file, err := os.Create(l.filename)
	if err != nil {
		return fmt.Errorf("failed to create CSV file: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	writer.Comma = ';' // Set delimiter to semicolon
	defer writer.Flush()

	// Write CSV header
	header := []string{
		"SteamID", "DriverName", "VehicleName", "CarClass",
		"LapsCompleted", "MaxSpeed", "BestLapTime",
		"BestSector1", "BestSector2", "BestSector3",
	}
	if err := writer.Write(header); err != nil {
		return fmt.Errorf("failed to write CSV header: %w", err)
	}

	// Sort drivers by position
	type driverPair struct {
		driver *models.StandingsData
		stats  *models.DriverStats
	}

	var drivers []driverPair
	for key, driver := range l.driverData {
		if stats, exists := l.driverStats[key]; exists {
			drivers = append(drivers, driverPair{driver, stats})
		}
	}

	sort.Slice(drivers, func(i, j int) bool {
		return drivers[i].driver.Position < drivers[j].driver.Position
	})

	// Write driver data
	for _, pair := range drivers {
		driver := pair.driver
		stats := pair.stats

		record := []string{
			fmt.Sprintf("%d", driver.SteamID),
			driver.DriverName,
			driver.VehicleName,
			driver.CarClass,
			fmt.Sprintf("%d", driver.LapsCompleted),
			fmt.Sprintf("%.1f", stats.MaxSpeed),
			formatTime(driver.BestLapTime),
			formatTime(stats.BestSector1),
			formatTime(stats.BestSector2),
			formatTime(stats.BestSector3),
		}

		if err := writer.Write(record); err != nil {
			return fmt.Errorf("failed to write CSV record: %w", err)
		}
	}

	return nil
}

// Close performs final write and cleanup
func (l *CSVLogger) Close() error {
	return l.WriteCurrentState()
}

// formatTime converts seconds to MM:SS.sss format
func formatTime(seconds float64) string {
	if seconds <= 0 {
		return "N/A"
	}
	minutes := int(seconds) / 60
	secs := seconds - float64(minutes*60)
	return fmt.Sprintf("%d:%06.3f", minutes, secs)
}
