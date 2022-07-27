package fitbit

import (
	"encoding/json"
	"fmt"
)

// ! ATTENTION !
// As of writing this file, there was no official Swagger documentation available.
// Most parts of this file are based on the not very accurate documentation which may provide different data.
// https://dev.fitbit.com/build/reference/web-api/temperature/

type TemperatureCore struct {
	TempCore []struct {
		DateTime string  `json:"dateTime"`
		Value    float64 `json:"value"`
	} `json:"tempCore"`
}

// TemperatureCoreByDay returns the core temperature data for a given date
// date must be in the format yyyy-MM-dd or today
func (m *Session) TemperatureCoreByDay(day string) (TemperatureCore, error) {
	contents, err := m.makeRequest(fmt.Sprintf("https://api.fitbit.com/1/user/-/temp/core/date/%s.json", day))
	if err != nil {
		return TemperatureCore{}, err
	}

	temperature := TemperatureCore{}
	if err := json.Unmarshal(contents, &temperature); err != nil {
		return TemperatureCore{}, err
	}

	return temperature, nil
}

// TemperatureCoreByDateRange returns the core temperature data for a given date range
// date must be in the format yyyy-MM-dd or today
func (m *Session) TemperatureCoreByDateRange(startDay string, endDay string) (TemperatureCore, error) {
	contents, err := m.makeRequest(fmt.Sprintf("https://api.fitbit.com/1/user/-/temp/core/date/%s/%s.json", startDay, endDay))
	if err != nil {
		return TemperatureCore{}, err
	}

	temperature := TemperatureCore{}
	if err := json.Unmarshal(contents, &temperature); err != nil {
		return TemperatureCore{}, err
	}

	return temperature, nil
}

type TemperatureSkin struct {
	TempSkin []struct {
		DateTime string `json:"dateTime"`
		Value    struct {
			NightlyRelative float64 `json:"nightlyRelative"`
		} `json:"value"`
		LogType string `json:"logType"`
	} `json:"tempSkin"`
}

// TemperatureSkinByDay returns the skin temperature data for a given date
// date must be in the format yyyy-MM-dd or today
func (m *Session) TemperatureSkinByDay(day string) (TemperatureSkin, error) {
	contents, err := m.makeRequest(fmt.Sprintf("https://api.fitbit.com/1/user/-/temp/skin/date/%s.json", day))
	if err != nil {
		return TemperatureSkin{}, err
	}

	temperature := TemperatureSkin{}
	if err := json.Unmarshal(contents, &temperature); err != nil {
		return TemperatureSkin{}, err
	}

	return temperature, nil
}

// TemperatureSkinByDateRange returns the skin temperature data for a given date range
// date must be in the format yyyy-MM-dd or today
func (m *Session) TemperatureSkinByDateRange(startDay string, endDay string) (TemperatureSkin, error) {
	contents, err := m.makeRequest(fmt.Sprintf("https://api.fitbit.com/1/user/-/temp/skin/date/%s/%s.json", startDay, endDay))
	if err != nil {
		return TemperatureSkin{}, err
	}

	temperature := TemperatureSkin{}
	if err := json.Unmarshal(contents, &temperature); err != nil {
		return TemperatureSkin{}, err
	}

	return temperature, nil
}
