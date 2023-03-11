package fitbit

import (
	"encoding/json"
	"fmt"
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

// HeartRateVariabilitySummary contains a summary of heartrate variability (HRV) values for a given date range
type HeartRateVariabilitySummary struct {
	Hrv []HRVDay `json:"hrv"`
}

type HRVDay struct {
	Value    HRVValue `json:"value"`
	DateTime string   `json:"dateTime"`
}

type HRVValue struct {
	DailyRmssd float64 `json:"dailyRmssd"`
	DeepRmssd  float64 `json:"deepRmssd"`
}

// HeartRateVariabilityIntraday with slightly different structure to HeartRateVariabilitySummary
type HeartRateVariabilityIntraday struct {
	Hrv []HRVIntraday `json:"hrv"`
}

type HRVIntraday struct {
	Minutes  []HRVMinutes `json:"minutes"`
	DateTime string       `json:"dateTime"`
}

type HRVMinutes struct {
	Minute string           `json:"minutes"`
	Value  HRVIntradayValue `json:"value"`
}

type HRVIntradayValue struct {
	Rmssd    float64 `json:"rmssd"`
	Coverage float64 `json:"coverage"`
	Hf       float64 `json:"hf"`
	Lf       float64 `json:"lf"`
}

// HeartLogByDay returns the heart log by a given date
// date must be in the format yyyy-MM-dd
func (m *Session) HeartLogByDay(day string) (HeartDay, error) {
	// If not day is given assume today
	if day == "" {
		day = "today"
	}

	contents, err := m.makeRequest(fmt.Sprintf("https://api.fitbit.com/1/user/-/activities/heart/date/%s/1d.json", day))
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

	contents, err := m.makeRequest(fmt.Sprintf("https://api.fitbit.com/1/user/-/activities/heart/date/%s/1d/%s/time/%s/%s.json", day, resolution, timeFrom, timeTo))
	if err != nil {
		return HeartIntraday{}, err
	}

	heartintra := HeartIntraday{}
	if err := json.Unmarshal(contents, &heartintra); err != nil {
		return HeartIntraday{}, err
	}

	return heartintra, nil
}

// HeartLogByDateRange returns the heart log of a given time range by date in default resolution
// date must be in the format yyyy-MM-dd
func (m *Session) HeartLogByDateRange(startDay string, endDay string) (HeartDay, error) {
	contents, err := m.makeRequest(fmt.Sprintf("https://api.fitbit.com/1/user/-/activities/heart/date/%s/%s.json", startDay, endDay))
	if err != nil {
		return HeartDay{}, err
	}

	heart := HeartDay{}
	if err := json.Unmarshal(contents, &heart); err != nil {
		return HeartDay{}, err
	}

	return heart, nil
}

// HeartLogByDateRangeIntraday returns the heart log by a given date range in the given resolution
// date must be in the format yyyy-MM-dd
// resolution can be 1min or 1sec, 1sec is default
func (m *Session) HeartLogByDateRangeIntraday(startDay string, endDay string, resolution string) (HeartDay, error) {
	// default to 1sec if resolution dos not match to 1min
	if resolution != "1min" {
		resolution = "1sec"
	}

	contents, err := m.makeRequest(fmt.Sprintf("https://api.fitbit.com/1/user/-/activities/heart/date/%s/%s/%s.json", startDay, endDay, resolution))
	if err != nil {
		return HeartDay{}, err
	}

	heart := HeartDay{}
	if err := json.Unmarshal(contents, &heart); err != nil {
		return HeartDay{}, err
	}

	return heart, nil
}

// HRVSummaryByDateRange the Heart Rate Variability (HRV) data for a date range.
// HRV data applies specifically to a user’s “main sleep,” which is the longest single period of time asleep on a given date.
// date must be in the format yyyy-MM-dd
func (m *Session) HRVSummaryByDateRange(startDay string, endDay string) (HeartRateVariabilitySummary, error) {
	contents, err := m.makeRequest(fmt.Sprintf("https://api.fitbit.com/1/user/-/hrv/date/%s/%s.json", startDay, endDay))
	if err != nil {
		return HeartRateVariabilitySummary{}, err
	}

	hrv := HeartRateVariabilitySummary{}
	if err := json.Unmarshal(contents, &hrv); err != nil {
		return HeartRateVariabilitySummary{}, err
	}

	return hrv, nil
}

// HRVSummaryByDate the Heart Rate Variability (HRV) data for a date.
// HRV data applies specifically to a user’s “main sleep,” which is the longest single period of time asleep on a given date.
// date must be in the format yyyy-MM-dd
func (m *Session) HRVSummaryByDate(day string) (HeartRateVariabilitySummary, error) {
	// If not day is given assume today
	if day == "" {
		day = "today"
	}

	contents, err := m.makeRequest(fmt.Sprintf("https://api.fitbit.com/1/user/-/hrv/date/%s.json", day))
	if err != nil {
		return HeartRateVariabilitySummary{}, err
	}

	hrv := HeartRateVariabilitySummary{}
	if err := json.Unmarshal(contents, &hrv); err != nil {
		return HeartRateVariabilitySummary{}, err
	}

	return hrv, nil
}

// HRVSummaryByDateRange the Heart Rate Variability (HRV) data for a date range.
// HRV data applies specifically to a user’s “main sleep,” which is the longest single period of time asleep on a given date.
// date must be in the format yyyy-MM-dd
func (m *Session) HRVIntradayByDateRange(startDay string, endDay string) (HeartRateVariabilityIntraday, error) {
	contents, err := m.makeRequest(fmt.Sprintf("https://api.fitbit.com/1/user/-/hrv/date/%s/%s/all.json", startDay, endDay))
	if err != nil {
		return HeartRateVariabilityIntraday{}, err
	}

	hrv := HeartRateVariabilityIntraday{}
	if err := json.Unmarshal(contents, &hrv); err != nil {
		return HeartRateVariabilityIntraday{}, err
	}

	return hrv, nil
}

// HRVSummaryByDate the Heart Rate Variability (HRV) data for a date.
// HRV data applies specifically to a user’s “main sleep,” which is the longest single period of time asleep on a given date.
// date must be in the format yyyy-MM-dd
func (m *Session) HRVIntradayByDate(day string) (HeartRateVariabilityIntraday, error) {
	// If not day is given assume today
	if day == "" {
		day = "today"
	}

	contents, err := m.makeRequest(fmt.Sprintf("https://api.fitbit.com/1/user/-/hrv/date/%s/all.json", day))
	if err != nil {
		return HeartRateVariabilityIntraday{}, err
	}

	hrv := HeartRateVariabilityIntraday{}
	if err := json.Unmarshal(contents, &hrv); err != nil {
		return HeartRateVariabilityIntraday{}, err
	}

	return hrv, nil
}
