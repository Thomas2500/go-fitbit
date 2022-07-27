package fitbit

import (
	"encoding/json"
	"fmt"
)

// ! ATTENTION !
// Functions are untested and may not work as intended.

type SpO2 struct {
	DateTime string `json:"dateTime"`
	Value    struct {
		Avg float64 `json:"avg"`
		Min float64 `json:"min"`
		Max float64 `json:"max"`
	} `json:"value"`
}

// SleepByDay returns the sleep data for a given date
// date must be in the format yyyy-MM-dd
func (m *Session) SpO2ByDay(day string) (SpO2, error) {
	contents, err := m.makeRequest(fmt.Sprintf("https://api.fitbit.com/1/user/-/spo2/date/%s.json", day))
	if err != nil {
		return SpO2{}, err
	}

	spo2 := SpO2{}
	if err := json.Unmarshal(contents, &spo2); err != nil {
		return SpO2{}, err
	}

	return spo2, nil
}

type SpO2Intraday struct {
	DateTime string `json:"dateTime"`
	Minutes  []struct {
		Value  float64 `json:"value"`
		Minute string  `json:"minute"`
	} `json:"minutes"`
}

// SpO2ByDayIntraday returns the sleep data for a given date with intraday accuration
// date must be in the format yyyy-MM-dd
func (m *Session) SpO2ByDayIntraday(day string) (SpO2Intraday, error) {
	contents, err := m.makeRequest(fmt.Sprintf("https://api.fitbit.com/1/user/-/spo2/date/%s/all.json", day))
	if err != nil {
		return SpO2Intraday{}, err
	}

	spo2 := SpO2Intraday{}
	if err := json.Unmarshal(contents, &spo2); err != nil {
		return SpO2Intraday{}, err
	}

	return spo2, nil
}

// TODO: SpO2ByDayRange
// TODO: SpO2IntradayByDayRange
