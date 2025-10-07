package ui

import (
	"fmt"
	"sort"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/mslomnicki/LMURacingTelemetry/pkg/models"
	"github.com/rivo/tview"
)

var Version = "dev"

// Display handles the terminal user interface
type Display struct {
	app        *tview.Application
	sessionBox *tview.TextView
	driversBox *tview.TextView
	statsBox   *tview.TextView
	versionBox *tview.TextView
}

// NewDisplay creates a new UI display
func NewDisplay() *Display {
	return &Display{
		app: tview.NewApplication(),
	}
}

// Setup initializes the UI layout
func (d *Display) Setup() {
	// Create text views for different data sections
	d.sessionBox = tview.NewTextView()
	d.sessionBox.SetBorder(true).SetTitle(" Session Info ").SetTitleAlign(tview.AlignLeft)
	d.sessionBox.SetDynamicColors(true)

	d.driversBox = tview.NewTextView()
	d.driversBox.SetBorder(true).SetTitle(" All Drivers - Live Data ").SetTitleAlign(tview.AlignLeft)
	d.driversBox.SetDynamicColors(true)

	d.statsBox = tview.NewTextView()
	d.statsBox.SetBorder(true).SetTitle(" Driver Statistics & Records ").SetTitleAlign(tview.AlignLeft)
	d.statsBox.SetDynamicColors(true)

	// Create a grid layout with 4 rows and 2 columns for bottom row
	grid := tview.NewGrid().
		SetRows(4, 0, 0).
		SetColumns(0, 60).
		SetBorders(true)

	// Add components to grid
	grid.AddItem(d.sessionBox, 0, 0, 1, 1, 0, 0, false).
		AddItem(d.driversBox, 1, 0, 1, 1, 0, 0, true).
		AddItem(d.statsBox, 2, 0, 1, 1, 0, 0, true)

	// Wrap grid in a Frame with app name and version as the title
	frame := tview.NewFrame(grid).
		SetBorders(0, 0, 0, 0, 0, 0).
		AddText(fmt.Sprintf("LMU Racing Telemetry %s", Version), true, tview.AlignCenter, tcell.ColorBlue)

	d.app.SetRoot(frame, true).EnableMouse(true)

	// Handle Ctrl+C and 'q' to exit
	d.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyCtrlC || event.Rune() == 'q' || event.Rune() == 'Q' {
			d.app.Stop()
		}
		return event
	})
}

// UpdateSession updates the session information display
func (d *Display) UpdateSession(session *models.SessionData) {
	if session == nil {
		return
	}

	sessionText := fmt.Sprintf(
		"[yellow]Track:[-] %s  [green]Session:[-] %s  [cyan]Phase:[-] %d\n"+
			"[white]Event Time:[-] %.1fs  [orange]Cars:[-] %d/%d  [red]Track:[-] %.1f°C  [blue]Air:[-] %.1f°C  [gray]Rain:[-] %.1f%%",
		session.TrackName,
		session.Session,
		session.GamePhase,
		session.CurrentEventTime,
		session.NumberOfVehicles,
		session.MaxPlayers,
		session.TrackTemp,
		session.AmbientTemp,
		session.Raining*100,
	)
	d.sessionBox.SetText(sessionText)
}

// UpdateDrivers updates the live drivers display
func (d *Display) UpdateDrivers(drivers map[string]*models.StandingsData) {
	// Convert map to slice and sort by position
	driverList := make([]*models.StandingsData, 0, len(drivers))
	for _, driver := range drivers {
		driverList = append(driverList, driver)
	}

	sort.Slice(driverList, func(i, j int) bool {
		return driverList[i].Position < driverList[j].Position
	})

	if len(driverList) == 0 {
		d.driversBox.SetText("No drivers connected...")
		return
	}

	// Calculate dynamic column widths based on content
	maxDriverName := 6  // "Driver"
	maxVehicleName := 7 // "Vehicle"
	maxStatus := 6      // "Status"

	for _, driver := range driverList {
		if len(driver.DriverName) > maxDriverName {
			maxDriverName = len(driver.DriverName)
		}
		if len(driver.VehicleName) > maxVehicleName {
			maxVehicleName = len(driver.VehicleName)
		}

		status := driver.PitState
		if driver.Flag != "" && driver.Flag != "green" {
			status = driver.Flag
		}
		if len(status) > maxStatus {
			maxStatus = len(status)
		}
	}

	// Limit maximum column widths to keep table readable
	if maxDriverName > 40 {
		maxDriverName = 40
	}
	if maxVehicleName > 40 {
		maxVehicleName = 40
	}
	if maxStatus > 20 {
		maxStatus = 20
	}

	// Create format strings with dynamic widths
	headerFormat := fmt.Sprintf("%%-%ds %%-%ds %%-%ds %%-%ds %%-%ds %%-%ds %%-%ds %%-%ds %%-%ds\n",
		3, maxDriverName, maxVehicleName, 4, 8, 8, 6, 5, maxStatus)
	dataFormat := fmt.Sprintf("%%-%dd %%-%ds %%-%ds %%-%dd %%-%ds %%-%ds %%-%d.0f %%-%d.1f %%-%ds\n",
		3, maxDriverName, maxVehicleName, 4, 8, 8, 6, 5, maxStatus)

	var driversText strings.Builder

	// Header
	driversText.WriteString(fmt.Sprintf(headerFormat,
		"Pos", "Driver", "Vehicle", "Laps", "CurLap", "BestLap", "Speed", "Fuel", "Status"))

	totalWidth := 3 + 1 + maxDriverName + 1 + maxVehicleName + 1 + 4 + 1 + 8 + 1 + 8 + 1 + 6 + 1 + 5 + 1 + maxStatus
	driversText.WriteString(strings.Repeat("-", totalWidth) + "\n")

	// Driver data
	for _, driver := range driverList {
		status := driver.PitState
		if driver.Flag != "" && driver.Flag != "green" {
			status = driver.Flag
		}

		line := fmt.Sprintf(dataFormat,
			driver.Position,
			truncate(driver.DriverName, maxDriverName),
			truncate(driver.VehicleName, maxVehicleName),
			driver.LapsCompleted,
			formatTime(driver.TimeIntoLap),
			formatTime(driver.BestLapTime),
			driver.CarVelocity.Velocity*3.6,
			driver.FuelFraction*100,
			truncate(status, maxStatus),
		)
		driversText.WriteString(line)
	}

	d.driversBox.SetText(driversText.String())
}

// UpdateStats updates the driver statistics display
func (d *Display) UpdateStats(stats map[string]*models.DriverStats) {
	// Convert map to slice and sort by best lap time
	statsList := make([]*models.DriverStats, 0, len(stats))
	for _, stat := range stats {
		statsList = append(statsList, stat)
	}

	sort.Slice(statsList, func(i, j int) bool {
		if statsList[i].BestLapTime == 0 && statsList[j].BestLapTime == 0 {
			return statsList[i].DriverName < statsList[j].DriverName
		}
		if statsList[i].BestLapTime == 0 {
			return false
		}
		if statsList[j].BestLapTime == 0 {
			return true
		}
		return statsList[i].BestLapTime < statsList[j].BestLapTime
	})

	if len(statsList) == 0 {
		d.statsBox.SetText("No driver statistics available...")
		return
	}

	// Calculate dynamic column widths based on content
	maxDriverName := 6  // "Driver"
	maxVehicleName := 7 // "Vehicle"
	maxClassName := 5   // "Class"

	for _, stat := range statsList {
		if len(stat.DriverName) > maxDriverName {
			maxDriverName = len(stat.DriverName)
		}
		if len(stat.VehicleName) > maxVehicleName {
			maxVehicleName = len(stat.VehicleName)
		}
		if len(stat.CarClass) > maxClassName {
			maxClassName = len(stat.CarClass)
		}
	}

	// Limit maximum column widths to keep table readable
	if maxDriverName > 25 {
		maxDriverName = 25
	}
	if maxVehicleName > 20 {
		maxVehicleName = 20
	}
	if maxClassName > 15 {
		maxClassName = 15
	}

	// Create format strings with dynamic widths
	headerFormat := fmt.Sprintf("%%-%ds %%-%ds %%-%ds %%-%ds %%-%ds %%-%ds %%-%ds %%-%ds\n",
		maxDriverName, maxVehicleName, maxClassName, 6, 8, 8, 8, 8)
	dataFormat := fmt.Sprintf("%%-%ds %%-%ds %%-%ds %%-%d.0f %%-%ds %%-%ds %%-%ds %%-%ds\n",
		maxDriverName, maxVehicleName, maxClassName, 6, 8, 8, 8, 8)

	var statsText strings.Builder

	// Header
	statsText.WriteString(fmt.Sprintf(headerFormat,
		"Driver", "Vehicle", "Class", "MaxSpd", "BestLap", "BestS1", "BestS2", "BestS3"))

	totalWidth := maxDriverName + 1 + maxVehicleName + 1 + maxClassName + 1 + 6 + 1 + 8 + 1 + 8 + 1 + 8 + 1 + 8
	statsText.WriteString(strings.Repeat("-", totalWidth) + "\n")

	// Stats data
	for _, stat := range statsList {
		line := fmt.Sprintf(dataFormat,
			truncate(stat.DriverName, maxDriverName),
			truncate(stat.VehicleName, maxVehicleName),
			truncate(stat.CarClass, maxClassName),
			stat.MaxSpeed,
			formatTime(stat.BestLapTime),
			formatTime(stat.BestSector1),
			formatTime(stat.BestSector2),
			formatTime(stat.BestSector3),
		)
		statsText.WriteString(line)
	}

	d.statsBox.SetText(statsText.String())
}

// Draw refreshes the display
func (d *Display) Draw() {
	d.app.Draw()
}

// Run starts the UI application
func (d *Display) Run() error {
	return d.app.Run()
}

// Stop stops the UI application
func (d *Display) Stop() {
	d.app.Stop()
}

// truncate shortens strings to fit display columns
func truncate(s string, length int) string {
	if len(s) <= length {
		return s
	}
	return s[:length-3] + "..."
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
