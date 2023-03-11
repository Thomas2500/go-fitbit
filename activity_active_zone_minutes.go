package fitbit

import (
	"encoding/json"
	"fmt"
)

// ActiveZoneMinutesDay contains the active zone minutes for a given day
type ActiveZoneMinutesDay struct {
	ActivitiesActiveZoneMinutes []struct {
		DateTime string `json:"dateTime"`
		Value    struct {
			FatBurnActiveZoneMinutes int `json:"fatBurnActiveZoneMinutes"`
			ActiveZoneMinutes        int `json:"activeZoneMinutes"`
		} `json:"value"`
	} `json:"activities-active-zone-minutes"`
}

// ActiveZoneMinutesLogByDay returns the active zone minutes log by a given date
// date must be in the format yyyy-MM-dd
func (m *Session) ActiveZoneMinutesLogByDay(day string) (ActiveZoneMinutesDay, error) {
	// If not day is given assume today
	if day == "" {
		day = "today"
	}

	contents, err := m.makeRequest(fmt.Sprintf("https://api.fitbit.com/1/user/-/activities/heart/date/%s/1d.json", day))
	if err != nil {
		return ActiveZoneMinutesDay{}, err
	}

	azm := ActiveZoneMinutesDay{}
	if err := json.Unmarshal(contents, &azm); err != nil {
		return ActiveZoneMinutesDay{}, err
	}

	return azm, nil
}

// ActiveZoneMinutesLogByDateRange returns the active zone minutes log by a given date
// date must be in the format yyyy-MM-dd
func (m *Session) ActiveZoneMinutesLogByDateRange(startDay string, endDay string) (ActiveZoneMinutesDay, error) {
	contents, err := m.makeRequest(fmt.Sprintf("https://api.fitbit.com/1/user/-/activities/heart/date/%s/%s.json", startDay, endDay))
	if err != nil {
		return ActiveZoneMinutesDay{}, err
	}

	azm := ActiveZoneMinutesDay{}
	if err := json.Unmarshal(contents, &azm); err != nil {
		return ActiveZoneMinutesDay{}, err
	}

	return azm, nil
}

type ActiveZoneMinutesIntraday struct {
	ActivitiesActiveZoneMinutesIntraday []struct {
		DateTime string `json:"dateTime"`
		Minutes  []struct {
			Minute string `json:"minute"`
			Value  struct {
				FatBurnActiveZoneMinutes int `json:"fatBurnActiveZoneMinutes,omitempty"`
				ActiveZoneMinutes        int `json:"activeZoneMinutes"`
			} `json:"value,omitempty"`
		} `json:"minutes"`
	} `json:"activities-active-zone-minutes-intraday"`
}

// ActiveZoneMinutesIntraday returns the active zone minutes log intraday by a given date in the given resolution
// date must be in the format yyyy-MM-dd
// resolution can be 1min, 5min, or 15min 1min is default
func (m *Session) ActiveZoneMinutesIntraday(day string, resolution string) (ActiveZoneMinutesIntraday, error) {
	// If not day is given assume today
	if day == "" {
		day = "today"
	}

	// default to 1sec if resolution dos not match to 1min
	if resolution != "5min" && resolution != "15min" {
		resolution = "1min"
	}

	contents, err := m.makeRequest(fmt.Sprintf("https://api.fitbit.com/1/user/-/activities/active-zone-minutes/date/%s/1d/%s.json", day, resolution))
	if err != nil {
		return ActiveZoneMinutesIntraday{}, err
	}

	azm := ActiveZoneMinutesIntraday{}
	if err := json.Unmarshal(contents, &azm); err != nil {
		return ActiveZoneMinutesIntraday{}, err
	}

	return azm, nil
}

// ActiveZoneMinutesIntradayByDateRange returns the active zone minutes log intraday by a given date in the given resolution
// date must be in the format yyyy-MM-dd
// resolution can be 1min, 5min, or 15min 1min is default
func (m *Session) ActiveZoneMinutesIntradayByDateRange(startDay string, endDay string, resolution string) (ActiveZoneMinutesIntraday, error) {
	// default to 1sec if resolution dos not match to 1min
	if resolution != "5min" && resolution != "15min" {
		resolution = "1min"
	}

	contents, err := m.makeRequest(fmt.Sprintf("https://api.fitbit.com/1/user/-/activities/active-zone-minutes/date/%s/%s/%s.json", startDay, endDay, resolution))
	if err != nil {
		return ActiveZoneMinutesIntraday{}, err
	}

	azm := ActiveZoneMinutesIntraday{}
	if err := json.Unmarshal(contents, &azm); err != nil {
		return ActiveZoneMinutesIntraday{}, err
	}

	return azm, nil
}
