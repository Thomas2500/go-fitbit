package fitbit

import "encoding/json"

// ActivitiesLifetime contains the account lifetime statistics
type ActivitiesLifetime struct {
	Best struct {
		Total struct {
			Distance struct {
				Date  string  `json:"date"`
				Value float64 `json:"value"`
			} `json:"distance"`
			Floors struct {
				Date  string  `json:"date"`
				Value float64 `json:"value"`
			} `json:"floors"`
			Steps struct {
				Date  string `json:"date"`
				Value int64  `json:"value"`
			} `json:"steps"`
		} `json:"total"`
		Tracker struct {
			Distance struct {
				Date  string  `json:"date"`
				Value float64 `json:"value"`
			} `json:"distance"`
			Floors struct {
				Date  string  `json:"date"`
				Value float64 `json:"value"`
			} `json:"floors"`
			Steps struct {
				Date  string `json:"date"`
				Value int64  `json:"value"`
			} `json:"steps"`
		} `json:"tracker"`
	} `json:"best"`
	Lifetime struct {
		Total struct {
			ActiveScore float64 `json:"activeScore"`
			CaloriesOut float64 `json:"caloriesOut"`
			Distance    float64 `json:"distance"`
			Floors      int64   `json:"floors"`
			Steps       int64   `json:"steps"`
		} `json:"total"`
		Tracker struct {
			ActiveScore float64 `json:"activeScore"`
			CaloriesOut float64 `json:"caloriesOut"`
			Distance    float64 `json:"distance"`
			Floors      int64   `json:"floors"`
			Steps       int64   `json:"steps"`
		} `json:"tracker"`
	} `json:"lifetime"`
}

// ActivitiesLifetime returns the summary of activities and made exercises
// date must be in the format yyyy-MM-dd
func (m *Session) ActivitiesLifetime() (ActivitiesLifetime, error) {
	contents, err := m.makeRequest("https://api.fitbit.com/1/user/-/activities.json")
	if err != nil {
		return ActivitiesLifetime{}, err
	}

	summary := ActivitiesLifetime{}
	if err := json.Unmarshal(contents, &summary); err != nil {
		return ActivitiesLifetime{}, err
	}

	return summary, nil
}
