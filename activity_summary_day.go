package fitbit

import (
	"encoding/json"
	"time"
)

// ActivitiesSummaryDay contains a summary of activities of a requested day
type ActivitiesSummaryDay struct {
	Activities []struct {
		ActivityID           int       `json:"activityId"`
		ActivityParentID     int       `json:"activityParentId"`
		ActivityParentName   string    `json:"activityParentName"`
		Calories             int       `json:"calories"`
		Description          string    `json:"description"`
		DetailsLink          string    `json:"detailsLink,omitempty"`
		Distance             float64   `json:"distance"`
		Duration             int       `json:"duration"`
		HasActiveZoneMinutes bool      `json:"hasActiveZoneMinutes"`
		HasStartTime         bool      `json:"hasStartTime"`
		IsFavorite           bool      `json:"isFavorite"`
		LastModified         time.Time `json:"lastModified"`
		LogID                int64     `json:"logId"`
		Name                 string    `json:"name"`
		StartDate            string    `json:"startDate"`
		StartTime            string    `json:"startTime"`
		Steps                int       `json:"steps"`
	} `json:"activities"`
	Goals struct {
		ActiveMinutes int     `json:"activeMinutes"`
		CaloriesOut   int     `json:"caloriesOut"`
		Distance      float64 `json:"distance"`
		Floors        int     `json:"floors"`
		Steps         int     `json:"steps"`
	} `json:"goals"`
	Summary struct {
		ActiveScore            int `json:"activeScore"`
		ActivityCalories       int `json:"activityCalories"`
		CalorieEstimationMu    int `json:"calorieEstimationMu"`
		CaloriesBMR            int `json:"caloriesBMR"`
		CaloriesOut            int `json:"caloriesOut"`
		CaloriesOutUnestimated int `json:"caloriesOutUnestimated"`
		Distances              []struct {
			Activity string  `json:"activity"`
			Distance float64 `json:"distance"`
		} `json:"distances"`
		Elevation           float64 `json:"elevation"`
		FairlyActiveMinutes int     `json:"fairlyActiveMinutes"`
		Floors              int     `json:"floors"`
		HeartRateZones      []struct {
			CaloriesOut float64 `json:"caloriesOut"`
			Max         int     `json:"max"`
			Min         int     `json:"min"`
			Minutes     int     `json:"minutes"`
			Name        string  `json:"name"`
		} `json:"heartRateZones"`
		LightlyActiveMinutes int  `json:"lightlyActiveMinutes"`
		MarginalCalories     int  `json:"marginalCalories"`
		RestingHeartRate     int  `json:"restingHeartRate"`
		SedentaryMinutes     int  `json:"sedentaryMinutes"`
		Steps                int  `json:"steps"`
		UseEstimation        bool `json:"useEstimation"`
		VeryActiveMinutes    int  `json:"veryActiveMinutes"`
	} `json:"summary"`
}

// ActivitiesDaySummary returns the summary of activities and made exercises
// date must be in the format yyyy-MM-dd
func (m *Session) ActivitiesDaySummary(day string) (ActivitiesSummaryDay, error) {
	contents, err := m.makeRequest("https://api.fitbit.com/1/user/-/activities/date/" + day + ".json")
	if err != nil {
		return ActivitiesSummaryDay{}, err
	}

	summary := ActivitiesSummaryDay{}
	if err := json.Unmarshal(contents, &summary); err != nil {
		return ActivitiesSummaryDay{}, err
	}

	return summary, nil
}
