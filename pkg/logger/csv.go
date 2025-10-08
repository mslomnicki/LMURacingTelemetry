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

type CSVLogger struct {
	filename    string
	driverStats map[string]*models.DriverStats
	session     *models.SessionData
}

func NewCSVLogger(session *models.SessionData) (*CSVLogger, error) {
	if session == nil {
		return nil, fmt.Errorf("session data is required")
	}

	now := time.Now()
	trackName := strings.ReplaceAll(session.TrackName, " ", "_")
	sessionName := strings.ReplaceAll(session.Session, " ", "_")
	filename := fmt.Sprintf("%s_%s_%s_telemetry.csv",
		now.Format("2006-01-02_15-04-05"),
		trackName,
		sessionName)

	return &CSVLogger{
		filename:    filename,
		driverStats: make(map[string]*models.DriverStats),
		session:     session,
	}, nil
}

func (l *CSVLogger) UpdateDriver(stats *models.DriverStats) {
	key := stats.DriverName
	l.driverStats[key] = stats
	if err := l.WriteCurrentState(); err != nil {
		fmt.Printf("Error writing CSV file: %v\n", err)
	}
}

func (l *CSVLogger) WriteCurrentState() error {
	file, err := os.Create(l.filename)
	if err != nil {
		return fmt.Errorf("failed to create CSV file: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	writer.Comma = ';'
	defer writer.Flush()

	header := []string{
		"Position",
		"SteamID", "DriverName", "CarClass", "VehicleModel", "VehicleName",
		"LapsCompleted", "MaxSpeed", "BestLapTime",
		"BestSector1", "BestSector2", "BestSector3",
		"MaxSpeedOnBestLap", "BestLapTimeCalculated", "BestSector1Calculated", "BestSector2Calculated", "BestSector3Calculated", "MaxSpeedOnBestLapCalc",
	}
	if err := writer.Write(header); err != nil {
		return fmt.Errorf("failed to write CSV header: %w", err)
	}

	var drivers []*models.DriverStats
	for _, stats := range l.driverStats {
		drivers = append(drivers, stats)
	}

	sort.Slice(drivers, func(i, j int) bool {
		return drivers[i].Position < drivers[j].Position
	})

	for _, stats := range drivers {
		record := []string{
			fmt.Sprintf("%d", stats.Position),
			fmt.Sprintf("%d", stats.SteamID),
			stats.DriverName,
			stats.CarClass,
			stats.VehicleModel,
			stats.VehicleName,
			fmt.Sprintf("%d", stats.LapsCompleted),
			fmt.Sprintf("%.1f", stats.MaxSpeed),
			formatTime(stats.BestLapTime),
			formatTime(stats.BestSector1),
			formatTime(stats.BestSector2),
			formatTime(stats.BestSector3),
			fmt.Sprintf("%.1f", stats.MaxSpeedOnBestLap),
			formatTime(stats.BestLapTimeCalculated),
			formatTime(stats.BestSector1Calculated),
			formatTime(stats.BestSector2Calculated),
			formatTime(stats.BestSector3Calculated),
			fmt.Sprintf("%.1f", stats.MaxSpeedOnBestLapCalc),
		}

		if err := writer.Write(record); err != nil {
			return fmt.Errorf("failed to write CSV record: %w", err)
		}
	}

	return nil
}

func (l *CSVLogger) Close() error {
	return l.WriteCurrentState()
}

func formatTime(seconds float64) string {
	if seconds <= 0 {
		return "N/A"
	}
	minutes := int(seconds) / 60
	secs := seconds - float64(minutes*60)
	return fmt.Sprintf("%d:%06.3f", minutes, secs)
}
