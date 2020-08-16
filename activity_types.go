package fitbit

import (
	"encoding/json"
)

type ActivitiesTypes struct {
	Categories []struct {
		Activities []struct {
			AccessLevel    string `json:"accessLevel"`
			ActivityLevels []struct {
				ID          int64   `json:"id"`
				MaxSpeedMPH float64 `json:"maxSpeedMPH"`
				Mets        float64 `json:"mets"`
				MinSpeedMPH float64 `json:"minSpeedMPH"`
				Name        string  `json:"name"`
			} `json:"activityLevels,omitempty"`
			HasSpeed bool    `json:"hasSpeed"`
			ID       int     `json:"id"`
			Name     string  `json:"name"`
			Mets     float64 `json:"mets,omitempty"`
		} `json:"activities"`
		ID            int    `json:"id"`
		Name          string `json:"name"`
		SubCategories []struct {
			Activities []struct {
				AccessLevel string  `json:"accessLevel"`
				HasSpeed    bool    `json:"hasSpeed"`
				ID          int64   `json:"id"`
				Mets        float64 `json:"mets"`
				Name        string  `json:"name"`
			} `json:"activities"`
			ID   int64  `json:"id"`
			Name string `json:"name"`
		} `json:"subCategories,omitempty"`
	} `json:"categories"`
}

// ActivityTypes returns a list of activities available, even user created ones
func (m *Session) ActivityTypes() (ActivitiesTypes, error) {
	contents, err := m.makeRequest("https://api.fitbit.com/1/activities.json")
	if err != nil {
		return ActivitiesTypes{}, err
	}

	activities := ActivitiesTypes{}
	if err := json.Unmarshal(contents, &activities); err != nil {
		return ActivitiesTypes{}, err
	}

	return activities, nil
}
