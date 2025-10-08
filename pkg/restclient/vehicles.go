package restclient

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type VehicleInfo struct {
	Id           string `json:"id"`
	FullPathTree string `json:"fullPathTree"`
	Number       string `json:"number"`
}

func GetAllVehicles(host string, port string) ([]VehicleInfo, error) {
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

	var vehicles []VehicleInfo
	if err := json.Unmarshal(body, &vehicles); err != nil {
		return nil, fmt.Errorf("JSON decode error: %w", err)
	}

	return vehicles, nil
}
