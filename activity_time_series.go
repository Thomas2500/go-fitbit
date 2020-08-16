package fitbit

import (
	"encoding/json"
	"errors"
)

// ActivitiesLog contains user activity logs, only one dataset is used and the other ones are empty
type ActivitiesLog struct {
	ActivitiesTrackerSteps                []ActivitiesLogSingleRecord `json:"activities-tracker-steps,omitempty"`
	ActivitiesTrackerCalories             []ActivitiesLogSingleRecord `json:"activities-tracker-calories,omitempty"`
	ActivitiesTrackerDistance             []ActivitiesLogSingleRecord `json:"activities-tracker-distance,omitempty"`
	ActivitiesTrackerFloors               []ActivitiesLogSingleRecord `json:"activities-tracker-floors,omitempty"`
	ActivitiesTrackerElevation            []ActivitiesLogSingleRecord `json:"activities-tracker-elevation,omitempty"`
	ActivitiesTrackerMinutesSedentary     []ActivitiesLogSingleRecord `json:"activities-tracker-minutesSedentary,omitempty"`
	ActivitiesTrackerMinutesLightlyActive []ActivitiesLogSingleRecord `json:"activities-tracker-minutesLightlyActive,omitempty"`
	ActivitiesTrackerMinutesFairlyActive  []ActivitiesLogSingleRecord `json:"activities-tracker-minutesFairlyActive,omitempty"`
	ActivitiesTrackerMinutesVeryActive    []ActivitiesLogSingleRecord `json:"activities-tracker-minutesVeryActive,omitempty"`
	ActivitiesTrackerActivityCalories     []ActivitiesLogSingleRecord `json:"activities-tracker-activityCalories,omitempty"`
}

// ActivitiesLogSingleRecord contains a single record of an activity
type ActivitiesLogSingleRecord struct {
	DateTime string `json:"dateTime"`
	Value    string `json:"value"`
}

type ActivitiesInterdayLog struct {
	ActivitiesCalories          []ActivitiesLogSingleRecord    `json:"activities-calories,omitempty"`
	ActivitiesCaloriesIntraday  ActivitiesIntradaySingleRecord `json:"activities-calories-intraday,omitempty"`
	ActivitiesSteps             []ActivitiesLogSingleRecord    `json:"activities-steps,omitempty"`
	ActivitiesStepsIntraday     ActivitiesIntradaySingleRecord `json:"activities-steps-intraday,omitempty"`
	ActivitiesDistance          []ActivitiesLogSingleRecord    `json:"activities-distance,omitempty"`
	ActivitiesDistanceIntraday  ActivitiesIntradaySingleRecord `json:"activities-distance-intraday,omitempty"`
	ActivitiesFloors            []ActivitiesLogSingleRecord    `json:"activities-floors,omitempty"`
	ActivitiesFloorsIntraday    ActivitiesIntradaySingleRecord `json:"activities-floors-intraday,omitempty"`
	ActivitiesElevation         []ActivitiesLogSingleRecord    `json:"activities-elevation,omitempty"`
	ActivitiesElevationIntraday ActivitiesIntradaySingleRecord `json:"activities-elevation-intraday,omitempty"`
}

type ActivitiesIntradaySingleRecord struct {
	Dataset []struct {
		Time  string  `json:"time"`
		Value float64 `json:"value"`
	} `json:"dataset,omitempty"`
	DatasetInterval int    `json:"datasetInterval,omitempty"`
	DatasetType     string `json:"datasetType,omitempty"`
}

// ActivitiesLogByDay returns the activities recorded for a given day and type
// date must be in the format yyyy-MM-dd and describes the end date
// activity is type of data to be fetched and returned
// fetchRange defines the timespan the data reaches back from day
func (m *Session) ActivitiesLogByDay(day string, activity string, fetchRange string) (ActivitiesLog, error) {
	// Supported activities: https://dev.fitbit.com/build/reference/web-api/activity/#resource-path-options:~:text=1y-,Resource%20Path%20Options
	switch activity {
	case "calories", "steps", "distance", "floors", "elevation", "minutesSedentary", "minutesLightlyActive", "minutesFairlyActive", "minutesVeryActive", "activityCalories":
		// noting to do here, client provided correct activity
	case "caloriesBMR":
		return ActivitiesLog{}, errors.New("currently not supported, error within fitbit api")
	default:
		return ActivitiesLog{}, errors.New("unknown activity given")
	}

	switch fetchRange {
	case "1d", "7d", "30d", "1w", "1m", "3m", "6m", "1y":
		// noting to do here, client provided correct range
	default:
		fetchRange = "1d"
	}

	contents, err := m.makeRequest("https://api.fitbit.com/1/user/-/activities/" + activity + "/date/" + day + "/" + fetchRange + ".json")
	if err != nil {
		return ActivitiesLog{}, err
	}

	summary := ActivitiesLog{}
	if err := json.Unmarshal(contents, &summary); err != nil {
		return ActivitiesLog{}, err
	}

	return summary, nil
}

// ActivitiesLogInterdayByDay returns the interday activities recorded for a given day and type
// date must be in the format yyyy-MM-dd and describes the end date
// activity is type of data to be fetched and returned
func (m *Session) ActivitiesLogInterdayByDay(day string, activity string) (ActivitiesInterdayLog, error) {
	// Supported activities: https://dev.fitbit.com/build/reference/web-api/activity/#resource-path-options:~:text=1y-,Resource%20Path%20Options
	switch activity {
	case "calories", "steps", "distance", "floors", "elevation":
		// noting to do here, client provided correct activity
	default:
		return ActivitiesInterdayLog{}, errors.New("unknown activity given")
	}

	contents, err := m.makeRequest("https://api.fitbit.com/1/user/-/activities/" + activity + "/date/" + day + "/1d/1min.json")
	if err != nil {
		return ActivitiesInterdayLog{}, err
	}

	interday := ActivitiesInterdayLog{}
	if err := json.Unmarshal(contents, &interday); err != nil {
		return ActivitiesInterdayLog{}, err
	}

	return interday, nil
}
