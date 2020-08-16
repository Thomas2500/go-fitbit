package fitbit

import (
	"encoding/json"
	"strconv"
)

// https://dev.fitbit.com/build/reference/web-api/devices/

// Device contains information about a fitbit device
type Device struct {
	Battery       string        `json:"battery"`
	BatteryLevel  int           `json:"batteryLevel,omitempty"`
	DeviceVersion string        `json:"deviceVersion"`
	Features      []interface{} `json:"features"`
	ID            string        `json:"id"`
	LastSyncTime  string        `json:"lastSyncTime"`
	Mac           string        `json:"mac,omitempty"`
	Type          string        `json:"type"`
}

// Devices returns
func (m *Session) Devices(userID uint64) ([]Device, error) {
	// Default "-" is current logged in user
	requestID := "-"
	if userID > 0 {
		requestID = strconv.FormatUint(userID, 10)
	}
	contents, err := m.makeRequest("https://api.fitbit.com/1/user/" + requestID + "/devices.json")
	if err != nil {
		return []Device{}, err
	}

	device := []Device{}
	if err := json.Unmarshal(contents, &device); err != nil {
		return []Device{}, err
	}

	return device, nil
}
