package restclient

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/mslomnicki/LMURacingTelemetry/pkg/models"
)

func GetAllVehicles(host string, port string) (map[string]models.VehicleInfo, error) {
	url := fmt.Sprintf("http://%s:%s/rest/sessions/getAllVehicles", host, port)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("GET request error: %w", err)
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			log.Printf("Error closing response body: %v", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("invalid response status: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("response read error: %w", err)
	}

	var vehicles []models.VehicleInfo
	if err := json.Unmarshal(body, &vehicles); err != nil {
		return nil, fmt.Errorf("JSON decode error: %w", err)
	}

	vehicleMap := make(map[string]models.VehicleInfo)
	for _, v := range vehicles {
		vehicleMap[v.Id] = v
	}

	return vehicleMap, nil
}
