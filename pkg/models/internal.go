package models

import "time"

type DriverStats struct {
	DriverName            string
	VehicleName           string
	VehicleModel          string
	VehicleNumber         string
	CarClass              string
	SteamID               int64
	MaxSpeed              float64
	BestLapTime           float64
	BestSector1           float64
	BestSector2           float64
	BestSector3           float64
	MaxSpeedOnBestLap     float64
	BestLapTimeCalculated float64
	BestSector1Calculated float64
	BestSector2Calculated float64
	BestSector3Calculated float64
	MaxSpeedOnBestLapCalc float64
	Position              int
	LapsCompleted         int
	LastUpdate            time.Time
}
