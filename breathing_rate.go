package fitbit

import (
	"encoding/json"
	"fmt"
)

type BreathingRate struct {
	Br []struct {
		Value struct {
			BreathingRate float64 `json:"breathingRate"`
		} `json:"value"`
		DateTime string `json:"dateTime"`
	} `json:"br"`
}

// BreathingRateLogByDay returns the breathing rate log (summary) by a given date
// date must be in the format yyyy-MM-dd
func (m *Session) BreathingRateLogByDay(day string) (BreathingRate, error) {
	contents, err := m.makeRequest(fmt.Sprintf("https://api.fitbit.com/1/user/-/br/date/%s.json", day))
	if err != nil {
		return BreathingRate{}, err
	}

	br := BreathingRate{}
	if err := json.Unmarshal(contents, &br); err != nil {
		return BreathingRate{}, err
	}

	return br, nil
}

// BreathingRateLogByDateRange returns the breathing rate summary log of a given time range by date
// date must be in the format yyyy-MM-dd
func (m *Session) BreathingRateLogByDateRange(startDay string, endDay string) (BreathingRate, error) {
	contents, err := m.makeRequest(fmt.Sprintf("https://api.fitbit.com/1/user/-/br/date/%s/%s.json", startDay, endDay))
	if err != nil {
		return BreathingRate{}, err
	}

	br := BreathingRate{}
	if err := json.Unmarshal(contents, &br); err != nil {
		return BreathingRate{}, err
	}

	return br, nil
}

type BreathingRateIntraday struct {
	Br []struct {
		Value struct {
			DeepSleepSummary struct {
				BreathingRate float64 `json:"breathingRate"`
			} `json:"deepSleepSummary"`
			RemSleepSummary struct {
				BreathingRate float64 `json:"breathingRate"`
			} `json:"remSleepSummary"`
			FullSleepSummary struct {
				BreathingRate float64 `json:"breathingRate"`
			} `json:"fullSleepSummary"`
			LightSleepSummary struct {
				BreathingRate float64 `json:"breathingRate"`
			} `json:"lightSleepSummary"`
		} `json:"value"`
		DateTime string `json:"dateTime"`
	} `json:"br"`
}

// BreathingRateLogByDayIntraday returns the breathing rate log by a given date
// date must be in the format yyyy-MM-dd
func (m *Session) BreathingRateLogByDayIntraday(day string) (BreathingRateIntraday, error) {
	contents, err := m.makeRequest(fmt.Sprintf("https://api.fitbit.com/1/user/-/br/date/%s/all.json", day))
	if err != nil {
		return BreathingRateIntraday{}, err
	}

	br := BreathingRateIntraday{}
	if err := json.Unmarshal(contents, &br); err != nil {
		return BreathingRateIntraday{}, err
	}

	return br, nil
}

// BreathingRateLogByDateRangeIntraday returns the breathing rate log of a given time range by date
// date must be in the format yyyy-MM-dd
func (m *Session) BreathingRateLogByDateRangeIntraday(startDay string, endDay string) (BreathingRateIntraday, error) {
	contents, err := m.makeRequest(fmt.Sprintf("https://api.fitbit.com/1/user/-/br/date/%s/%s/all.json", startDay, endDay))
	if err != nil {
		return BreathingRateIntraday{}, err
	}

	br := BreathingRateIntraday{}
	if err := json.Unmarshal(contents, &br); err != nil {
		return BreathingRateIntraday{}, err
	}

	return br, nil
}
