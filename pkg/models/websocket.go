package models

type WSMessage struct {
	Type string      `json:"type"`
	Body interface{} `json:"body"`
}

type StandingsData struct {
	AttackMode         AttackMode `json:"attackMode"`
	BestLapSectorTime1 float64    `json:"bestLapSectorTime1"`
	BestLapSectorTime2 float64    `json:"bestLapSectorTime2"`
	BestLapTime        float64    `json:"bestLapTime"`
	BestSectorTime1    float64    `json:"bestSectorTime1"`
	BestSectorTime2    float64    `json:"bestSectorTime2"`
	CarAcceleration    CarVector  `json:"carAcceleration"`
	CarClass           string     `json:"carClass"`
	CarId              string     `json:"carId"`
	CarNumber          string     `json:"carNumber"`
	CarPosition        CarVector  `json:"carPosition"`
	CarVelocity        CarVector  `json:"carVelocity"`
	CountLapFlag       string     `json:"countLapFlag"`
	CurrentSectorTime1 float64    `json:"currentSectorTime1"`
	CurrentSectorTime2 float64    `json:"currentSectorTime2"`
	DriverName         string     `json:"driverName"`
	DrsActive          bool       `json:"drsActive"`
	EstimatedLapTime   float64    `json:"estimatedLapTime"`
	FinishStatus       string     `json:"finishStatus"`
	Flag               string     `json:"flag"`
	Focus              bool       `json:"focus"`
	FuelFraction       float64    `json:"fuelFraction"`
	FullTeamName       string     `json:"fullTeamName"`
	GamePhase          string     `json:"gamePhase"`
	HasFocus           bool       `json:"hasFocus"`
	Headlights         bool       `json:"headlights"`
	InControl          int        `json:"inControl"`
	InGarageStall      bool       `json:"inGarageStall"`
	LapDistance        float64    `json:"lapDistance"`
	LapStartET         float64    `json:"lapStartET"`
	LapsBehindLeader   int        `json:"lapsBehindLeader"`
	LapsBehindNext     int        `json:"lapsBehindNext"`
	LapsCompleted      int        `json:"lapsCompleted"`
	LastLapTime        float64    `json:"lastLapTime"`
	LastSectorTime1    float64    `json:"lastSectorTime1"`
	LastSectorTime2    float64    `json:"lastSectorTime2"`
	PathLateral        float64    `json:"pathLateral"`
	Penalties          int        `json:"penalties"`
	PitGroup           string     `json:"pitGroup"`
	PitLapDistance     float64    `json:"pitLapDistance"`
	PitState           string     `json:"pitState"`
	Pitstops           int        `json:"pitstops"`
	Pitting            bool       `json:"pitting"`
	Player             bool       `json:"player"`
	Position           int        `json:"position"`
	Qualification      int        `json:"qualification"`
	Sector             string     `json:"sector"`
	ServerScored       bool       `json:"serverScored"`
	SlotID             int        `json:"slotID"`
	SteamID            int64      `json:"steamID"`
	TimeBehindLeader   float64    `json:"timeBehindLeader"`
	TimeBehindNext     float64    `json:"timeBehindNext"`
	TimeIntoLap        float64    `json:"timeIntoLap"`
	TrackEdge          float64    `json:"trackEdge"`
	UnderYellow        bool       `json:"underYellow"`
	UpgradePack        string     `json:"upgradePack"`
	VehicleFilename    string     `json:"vehicleFilename"`
	VehicleName        string     `json:"vehicleName"`
	VehicleModel       string
	VehicleNumber      string
}

type AttackMode struct {
	RemainingCount int     `json:"remainingCount"`
	TimeRemaining  float64 `json:"timeRemaining"`
	TotalCount     int     `json:"totalCount"`
}

type CarVector struct {
	Velocity float64 `json:"velocity"`
	X        float64 `json:"x"`
	Y        float64 `json:"y"`
	Z        float64 `json:"z"`
}

type SessionData struct {
	AmbientTemp        float64     `json:"ambientTemp"`
	AveragePathWetness float64     `json:"averagePathWetness"`
	CurrentEventTime   float64     `json:"currentEventTime"`
	DarkCloud          float64     `json:"darkCloud"`
	EndEventTime       float64     `json:"endEventTime"`
	GameMode           string      `json:"gameMode"`
	GamePhase          int         `json:"gamePhase"`
	InRealtime         bool        `json:"inRealtime"`
	LapDistance        float64     `json:"lapDistance"`
	MaxPathWetness     float64     `json:"maxPathWetness"`
	MaxPlayers         int         `json:"maxPlayers"`
	MaximumLaps        int         `json:"maximumLaps"`
	MinPathWetness     float64     `json:"minPathWetness"`
	NumRedLights       int         `json:"numRedLights"`
	NumberOfVehicles   int         `json:"numberOfVehicles"`
	PasswordProtected  bool        `json:"passwordProtected"`
	PlayerFileName     string      `json:"playerFileName"`
	PlayerName         string      `json:"playerName"`
	RaceCompletion     interface{} `json:"raceCompletion"`
	Raining            float64     `json:"raining"`
	SectorFlag         []string    `json:"sectorFlag"`
	ServerName         string      `json:"serverName"`
	ServerPort         int         `json:"serverPort"`
	Session            string      `json:"session"`
	StartEventTime     float64     `json:"startEventTime"`
	StartLightFrame    int         `json:"startLightFrame"`
	TrackName          string      `json:"trackName"`
	TrackTemp          float64     `json:"trackTemp"`
	WindSpeed          CarVector   `json:"windSpeed"`
	YellowFlagState    string      `json:"yellowFlagState"`
}
