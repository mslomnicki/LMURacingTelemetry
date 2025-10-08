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

type Display struct {
	app        *tview.Application
	sessionBox *tview.TextView
	driversBox *tview.TextView
	statsBox   *tview.TextView
	versionBox *tview.TextView
}

func NewDisplay() *Display {
	return &Display{
		app: tview.NewApplication(),
	}
}

func (d *Display) Setup() {
	d.sessionBox = tview.NewTextView()
	d.sessionBox.SetBorder(true).SetTitle(" Session Info ").SetTitleAlign(tview.AlignLeft)
	d.sessionBox.SetDynamicColors(true)

	d.driversBox = tview.NewTextView()
	d.driversBox.SetBorder(true).SetTitle(" All Drivers - Live Data ").SetTitleAlign(tview.AlignLeft)
	d.driversBox.SetDynamicColors(true)

	d.statsBox = tview.NewTextView()
	d.statsBox.SetBorder(true).SetTitle(" Driver Statistics & Records ").SetTitleAlign(tview.AlignLeft)
	d.statsBox.SetDynamicColors(true)

	grid := tview.NewGrid().
		SetRows(4, 0, 0).
		SetBorders(true)

	grid.AddItem(d.sessionBox, 0, 0, 1, 1, 0, 0, false).
		AddItem(d.driversBox, 1, 0, 1, 1, 0, 0, true).
		AddItem(d.statsBox, 2, 0, 1, 1, 0, 0, true)

	frame := tview.NewFrame(grid).
		SetBorders(0, 0, 0, 0, 0, 0).
		AddText(fmt.Sprintf("LMU Racing Telemetry %s", Version), true, tview.AlignCenter, tcell.ColorBlue)

	d.app.SetRoot(frame, true).EnableMouse(true)

	d.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyCtrlC || event.Rune() == 'q' || event.Rune() == 'Q' {
			d.app.Stop()
		}
		return event
	})
}

func (d *Display) UpdateSession(session *models.SessionData) {
	if session == nil {
		return
	}

	sessionText := fmt.Sprintf(
		"[yellow]Track:[-] %s  [green]Session:[-] %s  [cyan]Phase:[-] %d\n"+
			"[white]Event Time:[-] %s  [orange]Cars:[-] %d/%d  [red]Track:[-] %.1f°C  [blue]Air:[-] %.1f°C  [gray]Rain:[-] %.1f%%",
		session.TrackName,
		session.Session,
		session.GamePhase,
		formatTime(session.CurrentEventTime),
		session.NumberOfVehicles,
		session.MaxPlayers,
		session.TrackTemp,
		session.AmbientTemp,
		session.Raining*100,
	)
	d.sessionBox.SetText(sessionText)
}

func (d *Display) UpdateDrivers(drivers map[string]*models.StandingsData) {
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

	maxDriverName := 6
	maxClassName := 5
	maxVehicleName := 7
	maxStatus := 6

	for _, driver := range driverList {
		if len(driver.DriverName) > maxDriverName {
			maxDriverName = len(driver.DriverName)
		}
		if len(driver.CarClass) > maxClassName {
			maxClassName = len(driver.CarClass)
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

	if maxDriverName > 40 {
		maxDriverName = 40
	}
	if maxClassName > 20 {
		maxClassName = 20
	}
	if maxVehicleName > 40 {
		maxVehicleName = 40
	}
	if maxStatus > 20 {
		maxStatus = 20
	}

	headerFormat := fmt.Sprintf("%%-%ds %%-%ds %%-%ds %%-%ds %%-%ds %%-%ds %%-%ds %%-%ds %%-%ds\n",
		3, maxDriverName, maxClassName, maxVehicleName, 4, 8, 8, 6, maxStatus)
	dataFormat := fmt.Sprintf("%%-%dd %%-%ds %%-%ds %%-%ds %%-%dd %%-%ds %%-%ds %%-%d.0f %%-%ds\n",
		3, maxDriverName, maxClassName, maxVehicleName, 4, 8, 8, 6, maxStatus)

	var driversText strings.Builder

	driversText.WriteString(fmt.Sprintf(headerFormat,
		"Pos", "Driver", "Class", "Vehicle", "Laps", "CurLap", "BestLap", "Speed", "Status"))
	totalWidth := 3 + 1 + maxDriverName + 1 + maxClassName + 1 + maxVehicleName + 1 + 4 + 1 + 8 + 1 + 8 + 1 + 6 + 1 + maxStatus
	driversText.WriteString(strings.Repeat("-", totalWidth) + "\n")

	for _, driver := range driverList {
		status := driver.PitState
		if driver.Flag != "" && driver.Flag != "green" {
			status = driver.Flag
		}

		line := fmt.Sprintf(dataFormat,
			driver.Position,
			truncate(driver.DriverName, maxDriverName),
			truncate(driver.CarClass, maxClassName),
			truncate(driver.VehicleName, maxVehicleName),
			driver.LapsCompleted,
			formatTime(driver.TimeIntoLap),
			formatTime(driver.BestLapTime),
			driver.CarVelocity.Velocity*3.6,
			truncate(status, maxStatus),
		)
		driversText.WriteString(line)
	}

	d.driversBox.SetText(driversText.String())
}

func (d *Display) UpdateStats(stats map[string]*models.DriverStats) {
	statsList := make([]*models.DriverStats, 0, len(stats))
	for _, stat := range stats {
		statsList = append(statsList, stat)
	}

	sort.Slice(statsList, func(i int, j int) bool {
		if statsList[i].BestLapTime <= 0 && statsList[j].BestLapTime <= 0 {
			return statsList[i].DriverName < statsList[j].DriverName
		}
		return statsList[i].BestLapTime > statsList[j].BestLapTime
	})

	if len(statsList) == 0 {
		d.statsBox.SetText("No driver statistics available...")
		return
	}

	maxDriverName := 6
	maxVehicleName := 7
	maxClassName := 5

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

	if maxDriverName > 25 {
		maxDriverName = 25
	}
	if maxVehicleName > 20 {
		maxVehicleName = 20
	}
	if maxClassName > 15 {
		maxClassName = 15
	}

	headerFormat := fmt.Sprintf("%%-%ds %%-%ds %%-%ds %%6s %%8s %%8s %%8s %%8s %%7s %%8s %%8s %%8s %%8s %%6s\n",
		maxDriverName, maxClassName, maxVehicleName)
	dataFormat := fmt.Sprintf("%%-%ds %%-%ds %%-%ds %%6.1f %%8s %%8s %%8s %%8s %%7.1f %%8s %%8s %%8s %%8s %%6.1f\n",
		maxDriverName, maxClassName, maxVehicleName)

	var statsText strings.Builder

	statsText.WriteString(fmt.Sprintf(headerFormat,
		"Driver", "Class", "Vehicle", "MaxSpd", "BestLap", "BestS1", "BestS2", "BestS3", "MaxSpdC", "BestLapC", "BestS1C", "BestS2C", "BestS3C", "MaxSpdBC"))
	totalWidth := maxDriverName + 1 + maxClassName + 1 + maxVehicleName + 1 + 6 + 1 + 8 + 1 + 8 + 1 + 8 + 1 + 8 + 1 + 7 + 1 + 8 + 1 + 8 + 1 + 8 + 1 + 8 + 1 + 6
	statsText.WriteString(strings.Repeat("-", totalWidth) + "\n")

	for _, stat := range statsList {
		line := fmt.Sprintf(dataFormat,
			truncate(stat.DriverName, maxDriverName),
			truncate(stat.CarClass, maxClassName),
			truncate(stat.VehicleName, maxVehicleName),
			stat.MaxSpeed,
			formatTime(stat.BestLapTime),
			formatTime(stat.BestSector1),
			formatTime(stat.BestSector2),
			formatTime(stat.BestSector3),
			stat.MaxSpeedOnBestLapCalc,
			formatTime(stat.BestLapTimeCalculated),
			formatTime(stat.BestSector1Calculated),
			formatTime(stat.BestSector2Calculated),
			formatTime(stat.BestSector3Calculated),
			stat.MaxSpeedOnBestLapCalc,
		)
		statsText.WriteString(line)
	}

	d.statsBox.SetText(statsText.String())
}

func (d *Display) Draw() {
	d.app.Draw()
}

func (d *Display) Run() error {
	return d.app.Run()
}

func (d *Display) Stop() {
	d.app.Stop()
}

func truncate(s string, length int) string {
	if len(s) <= length {
		return s
	}
	return s[:length-3] + "..."
}

func formatTime(seconds float64) string {
	if seconds <= 0 {
		return "N/A"
	}
	minutes := int(seconds) / 60
	secs := seconds - float64(minutes*60)
	return fmt.Sprintf("%d:%06.3f", minutes, secs)
}

func floatDiffers(a, b float64) bool {
	const eps = 0.001
	if a == 0 && b == 0 {
		return false
	}
	return (a > b && a-b > eps) || (b > a && b-a > eps)
}

func padColorString(s string, width int) string {
	visible := 0
	result := ""
	inTag := false
	for i := 0; i < len(s); i++ {
		if s[i] == '[' {
			inTag = true
		}
		if !inTag {
			if visible < width {
				result += string(s[i])
				visible++
			}
		} else {
			result += string(s[i])
		}
		if s[i] == ']' {
			inTag = false
		}
	}
	for visible < width {
		result += " "
		visible++
	}
	return result
}
