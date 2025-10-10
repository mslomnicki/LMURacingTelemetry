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
var Year = "2025"

type Display struct {
	app               *tview.Application
	sessionBox        *tview.TextView
	driversBox        *tview.TextView
	statsBox          *tview.TextView
	versionBox        *tview.TextView
	grid              *tview.Grid
	frame             *tview.Frame
	fullscreenDrivers bool
	fullscreenStats   bool
	prevFocus         tview.Primitive
}

func NewDisplay() *Display {
	return &Display{
		app: tview.NewApplication(),
	}
}

func (d *Display) Setup() {
	d.sessionBox = tview.NewTextView()
	d.sessionBox.SetBorder(true).SetTitle(" [::b]Session Info[::-] ").SetTitleAlign(tview.AlignLeft)
	d.sessionBox.SetDynamicColors(true)

	d.driversBox = tview.NewTextView()
	d.driversBox.SetBorder(true).SetTitle(" [::b]All Drivers - Live Data[::-] ").SetTitleAlign(tview.AlignLeft)
	d.driversBox.SetDynamicColors(true)

	d.statsBox = tview.NewTextView()
	d.statsBox.SetBorder(true).SetTitle(" [::b]Driver Statistics & Records[::-] ").SetTitleAlign(tview.AlignLeft)
	d.statsBox.SetDynamicColors(true)

	d.grid = tview.NewGrid().
		SetRows(3, 0, 0).
		SetBorders(true)

	d.grid.AddItem(d.sessionBox, 0, 0, 1, 1, 0, 0, false).
		AddItem(d.driversBox, 1, 0, 1, 1, 0, 0, true).
		AddItem(d.statsBox, 2, 0, 1, 1, 0, 0, true)

	d.frame = tview.NewFrame(d.grid).
		SetBorders(0, 0, 0, 0, 0, 0).
		AddText(fmt.Sprintf("[::b]LMU Racing Telemetry[::-] %s", Version), true, tview.AlignCenter, tcell.ColorLightBlue).
		AddText(fmt.Sprintf("Copyright (C) %s Marek Słomnicki <marek@slomnicki.net>", Year), false, tview.AlignCenter, tcell.ColorBlue).
		AddText("Press Ctrl+C or Q to quit | F - fullscreen drivers | S - fullscreen stats", false, tview.AlignCenter, tcell.ColorGray)

	d.app.SetRoot(d.frame, true).EnableMouse(true)

	d.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyCtrlC || event.Rune() == 'q' || event.Rune() == 'Q' {
			d.app.Stop()
			return nil
		}
		if event.Rune() == 'f' || event.Rune() == 'F' {
			d.toggleDriversFullscreen()
			return nil
		}
		if event.Rune() == 's' || event.Rune() == 'S' {
			d.toggleStatsFullscreen()
			return nil
		}
		return event
	})
}

func (d *Display) UpdateSession(session *models.SessionData) {
	if session == nil {
		return
	}

	sessionText := fmt.Sprintf(
		"[yellow]Track:[-] %s  [green]Session:[-] %s  "+
			"[white]Event Time:[-] %s  [orange]Cars:[-] %d/%d  [red]Track:[-] %.1f°C  [blue]Air:[-] %.1f°C  [gray]Rain:[-] %.1f%%",
		session.TrackName,
		session.Session,
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
	maxVehicleModel := 7
	maxStatus := 6
	maxVehicleNumber := 4
	for _, driver := range driverList {
		if len(driver.DriverName) > maxDriverName {
			maxDriverName = len(driver.DriverName)
		}
		if len(driver.CarClass) > maxClassName {
			maxClassName = len(driver.CarClass)
		}
		if len(driver.VehicleModel) > maxVehicleModel {
			maxVehicleModel = len(driver.VehicleModel)
		}
		if len(driver.VehicleNumber) > maxVehicleNumber {
			maxVehicleNumber = len(driver.VehicleNumber)
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
	if maxVehicleModel > 40 {
		maxVehicleModel = 40
	}
	if maxVehicleNumber > 4 {
		maxVehicleNumber = 4
	}
	if maxStatus > 20 {
		maxStatus = 20
	}

	headerFormat := fmt.Sprintf("[yellow][::b]%%-%ds %%-%ds %%-%ds %%-%ds %%-%ds %%-%ds %%-%ds %%-%ds %%-%ds %%-%ds[::-][-]\n",
		3, maxDriverName, maxClassName, maxVehicleNumber, maxVehicleModel, 4, 8, 8, 6, maxStatus)
	dataFormat := fmt.Sprintf("%%-%dd %%-%ds %%-%ds %%%ds %%-%ds %%-%dd %%-%ds %%-%ds %%-%d.0f %%-%ds\n",
		3, maxDriverName, maxClassName, maxVehicleNumber, maxVehicleModel, 4, 8, 8, 6, maxStatus)

	var driversText strings.Builder

	driversText.WriteString(fmt.Sprintf(headerFormat,
		"Pos", "Driver", "Class", "No.", "Vehicle", "Laps", "CurLap", "BestLap", "Speed", "Status"))
	totalWidth := 3 + 1 + maxDriverName + 1 + maxClassName + 1 + maxVehicleNumber + 1 + maxVehicleModel + 1 + 4 + 1 + 8 + 1 + 8 + 1 + 6 + 1 + maxStatus
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
			truncate(driver.VehicleNumber, maxVehicleNumber),
			truncate(driver.VehicleModel, maxVehicleModel),
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
		timeA := statsList[i].BestLapTime
		if timeA <= 0 {
			timeA = 1e9
		}
		timeB := statsList[j].BestLapTime
		if timeB <= 0 {
			timeB = 1e9
		}
		return timeA < timeB
	})

	if len(statsList) == 0 {
		d.statsBox.SetText("No driver statistics available...")
		return
	}

	maxDriverName := 6
	maxVehicleModel := 7
	maxClassName := 5
	maxVehicleNumber := 4

	for _, stat := range statsList {
		if len(stat.DriverName) > maxDriverName {
			maxDriverName = len(stat.DriverName)
		}
		if len(stat.VehicleModel) > maxVehicleModel {
			maxVehicleModel = len(stat.VehicleModel)
		}
		if len(stat.CarClass) > maxClassName {
			maxClassName = len(stat.CarClass)
		}
		if len(stat.VehicleNumber) > maxVehicleNumber {
			maxVehicleNumber = len(stat.VehicleNumber)
		}
	}

	if maxDriverName > 40 {
		maxDriverName = 40
	}
	if maxVehicleModel > 40 {
		maxVehicleModel = 40
	}
	if maxClassName > 15 {
		maxClassName = 15
	}
	if maxVehicleNumber > 4 {
		maxVehicleNumber = 4
	}

	headerFormat := fmt.Sprintf("[yellow][::b]%%-%ds %%-%ds %%-%ds %%-%ds %%6s %%8s %%8s %%8s %%8s %%7s %%8s %%8s %%8s %%8s %%6s[::-][-]\n",
		maxDriverName, maxClassName, maxVehicleNumber, maxVehicleModel)
	dataFormat := fmt.Sprintf("%%-%ds %%-%ds %%%ds %%-%ds %%6.1f %%8s %%8s %%8s %%8s %%7.1f %%8s %%8s %%8s %%8s %%6.1f\n",
		maxDriverName, maxClassName, maxVehicleNumber, maxVehicleModel)

	var statsText strings.Builder

	statsText.WriteString(fmt.Sprintf(headerFormat,
		"Driver", "Class", "No.", "Vehicle", "MaxSpd", "BestLap", "BestS1", "BestS2", "BestS3", "MaxSpdC", "BestLapC", "BestS1C", "BestS2C", "BestS3C", "MaxSpdBC"))
	totalWidth := maxDriverName + 1 + maxClassName + 1 + maxVehicleNumber + 1 + maxVehicleModel + 1 + 6 + 1 + 8 + 1 + 8 + 1 + 8 + 1 + 8 + 1 + 7 + 1 + 8 + 1 + 8 + 1 + 8 + 1 + 8 + 1 + 6
	statsText.WriteString(strings.Repeat("-", totalWidth) + "\n")

	for _, stat := range statsList {
		line := fmt.Sprintf(dataFormat,
			truncate(stat.DriverName, maxDriverName),
			truncate(stat.CarClass, maxClassName),
			truncate(stat.VehicleNumber, maxVehicleNumber),
			truncate(stat.VehicleModel, maxVehicleModel),
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

func (d *Display) toggleDriversFullscreen() {
	d.fullscreenStats = false
	if !d.fullscreenDrivers {
		d.prevFocus = d.app.GetFocus()
		fullscreenFrame := tview.NewFrame(d.driversBox).
			SetBorders(0, 0, 0, 0, 0, 0).
			AddText(fmt.Sprintf("[::b]LMU Racing Telemetry[::-] %s", Version), true, tview.AlignCenter, tcell.ColorLightBlue).
			AddText(fmt.Sprintf("Copyright (C) %s Marek Słomnicki <marek@slomnicki.net>", Year), false, tview.AlignCenter, tcell.ColorBlue).
			AddText("Press Ctrl+C or Q to quit | [::b]F - fullscreen drivers[::-] | S - fullscreen stats", false, tview.AlignCenter, tcell.ColorGray)
		d.app.SetRoot(fullscreenFrame, true)
		d.app.SetFocus(d.driversBox)
		d.fullscreenDrivers = true
	} else {
		d.app.SetRoot(d.frame, true)
		if d.prevFocus != nil {
			d.app.SetFocus(d.prevFocus)
		}
		d.fullscreenDrivers = false
	}
}

func (d *Display) toggleStatsFullscreen() {
	d.fullscreenDrivers = false
	if !d.fullscreenStats {
		d.prevFocus = d.app.GetFocus()
		fullscreenFrame := tview.NewFrame(d.statsBox).
			SetBorders(0, 0, 0, 0, 0, 0).
			AddText(fmt.Sprintf("[::b]LMU Racing Telemetry[::-] %s", Version), true, tview.AlignCenter, tcell.ColorLightBlue).
			AddText(fmt.Sprintf("Copyright (C) %s Marek Słomnicki <marek@slomnicki.net>", Year), false, tview.AlignCenter, tcell.ColorBlue).
			AddText("Press Ctrl+C or Q to quit | F - fullscreen drivers | [::b]S - fullscreen stats[::-]", false, tview.AlignCenter, tcell.ColorGray)
		d.app.SetRoot(fullscreenFrame, true)
		d.app.SetFocus(d.statsBox)
		d.fullscreenStats = true
	} else {
		d.app.SetRoot(d.frame, true)
		if d.prevFocus != nil {
			d.app.SetFocus(d.prevFocus)
		}
		d.fullscreenStats = false
	}
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
