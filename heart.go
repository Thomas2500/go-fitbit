package fitbit

import (
	"encoding/json"
)

// HeartDay contains a summary of heartrates for a given date range
type HeartDay struct {
	ActivitiesHeart []struct {
		DateTime string `json:"dateTime"`
		Value    struct {
			CustomHeartRateZones []interface{}    `json:"customHeartRateZones,omitempty"`
			HeartRateZones       []HeartRateZones `json:"heartRateZones"`
			RestingHeartRate     int              `json:"restingHeartRate"`
		} `json:"value"`
	} `json:"activities-heart"`
	ActivitiesHeartIntraday ActivitiesHeartIntraday `json:"activities-heart-intraday,omitempty"`
}

// HeartIntraday with slightly different structure to HeartDay
type HeartIntraday struct {
	ActivitiesHeart []struct {
		CustomHeartRateZones []interface{}    `json:"customHeartRateZones"`
		DateTime             string           `json:"dateTime"`
		HeartRateZones       []HeartRateZones `json:"heartRateZones"`
		Value                string           `json:"value"`
	} `json:"activities-heart"`
	ActivitiesHeartIntraday ActivitiesHeartIntraday `json:"activities-heart-intraday,omitempty"`
}

// HeartRateZones contains the heart rate zones of different types like cardio
type HeartRateZones struct {
	CaloriesOut float64 `json:"caloriesOut"`
	Max         int     `json:"max"`
	Min         int     `json:"min"`
	Minutes     int     `json:"minutes"`
	Name        string  `json:"name"`
}

// ActivitiesHeartIntraday intraday data
type ActivitiesHeartIntraday struct {
	Dataset []struct {
		Time  string `json:"time"`
		Value int    `json:"value"`
	} `json:"dataset"`
	DatasetInterval int    `json:"datasetInterval"`
	DatasetType     string `json:"datasetType"`
}

// HeartLogByDay returns the heart log by a given date
// date must be in the format yyyy-MM-dd
func (m *Session) HeartLogByDay(day string) (HeartDay, error) {
	// If not day is given assume today
	if day == "" {
		day = "today"
	}

	contents, err := m.makeRequest("https://api.fitbit.com/1/user/-/activities/heart/date/" + day + "/1d.json")
	if err != nil {
		return HeartDay{}, err
	}

	foods := HeartDay{}
	if err := json.Unmarshal(contents, &foods); err != nil {
		return HeartDay{}, err
	}

	return foods, nil
}

// HeartIntraday returns the heart log by a given date in the given resolution
// date must be in the format yyyy-MM-dd, default is today
// resolution can be 1min or 1sec, 1sec is default
// timeFrom and timeTo are in the format 00:00 for hour:minute, default entire day
func (m *Session) HeartIntraday(day string, resolution string, timeFrom string, timeTo string) (HeartIntraday, error) {
	// If not day is given assume today
	if day == "" {
		day = "today"
	}

	if timeFrom == "" {
		timeFrom = "00:00"
	}
	if timeTo == "" {
		timeTo = "23:59"
	}

	// default to 1sec if resolution dos not match to 1min
	if resolution != "1min" {
		resolution = "1sec"
	}

	contents, err := m.makeRequest("https://api.fitbit.com/1/user/-/activities/heart/date/" + day + "/1d/" + resolution + "/time/" + timeFrom + "/" + timeTo + ".json")
	if err != nil {
		return HeartIntraday{}, err
	}

	heartintra := HeartIntraday{}
	if err := json.Unmarshal(contents, &heartintra); err != nil {
		return HeartIntraday{}, err
	}

	return heartintra, nil
}

// HeartLogByDateRange returns the calories log of a given time range by date
// date must be in the format yyyy-MM-dd
func (m *Session) HeartLogByDateRange(startDay string, endDay string) (HeartDay, error) {
	contents, err := m.makeRequest("https://api.fitbit.com/1/user/-/activities/heart/date/" + startDay + "/" + endDay + ".json")
	if err != nil {
		return HeartDay{}, err
	}

	heart := HeartDay{}
	if err := json.Unmarshal(contents, &heart); err != nil {
		return HeartDay{}, err
	}

	return heart, nil
}
