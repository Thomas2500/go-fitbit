package fitbit

import (
	"encoding/json"
)

type ActivitiesFrequent struct {
	ActivityID  int     `json:"activityId"`
	Calories    int     `json:"calories"`
	Description string  `json:"description"`
	Distance    float64 `json:"distance"`
	Duration    int     `json:"duration"`
	Name        string  `json:"name"`
}

// ActivityFrequent returns a list of frequent user activities
func (m *Session) ActivityFrequent() ([]ActivitiesFrequent, error) {
	contents, err := m.makeRequest("https://api.fitbit.com/1/user/-/activities/frequent.json")
	if err != nil {
		return []ActivitiesFrequent{}, err
	}

	activities := []ActivitiesFrequent{}
	if err := json.Unmarshal(contents, &activities); err != nil {
		return []ActivitiesFrequent{}, err
	}

	return activities, nil
}

// ActivityRecent returns a list of fRecent user activities
func (m *Session) ActivityRecent() ([]ActivitiesFrequent, error) {
	contents, err := m.makeRequest("https://api.fitbit.com/1/user/-/activities/recent.json")
	if err != nil {
		return []ActivitiesFrequent{}, err
	}

	activities := []ActivitiesFrequent{}
	if err := json.Unmarshal(contents, &activities); err != nil {
		return []ActivitiesFrequent{}, err
	}

	return activities, nil
}
