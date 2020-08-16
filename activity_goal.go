package fitbit

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// ActivitiesGoal contains the activities goal of an user
type ActivitiesGoal struct {
	Goals struct {
		ActiveMinutes int     `json:"activeMinutes,omitempty"`
		CaloriesOut   int     `json:"caloriesOut,omitempty"`
		Distance      float64 `json:"distance"`
		Floors        int     `json:"floors"`
		Steps         int     `json:"steps"`
	} `json:"goals"`
}

// ActivitiesGoal returns user set activities goal
func (m *Session) ActivitiesGoal(period string) (ActivitiesGoal, error) {
	if period != "weekly" {
		period = "daily"
	}

	contents, err := m.makeRequest("https://api.fitbit.com/1/user/-/activities/goals/" + period + ".json")
	if err != nil {
		return ActivitiesGoal{}, err
	}

	summary := ActivitiesGoal{}
	if err := json.Unmarshal(contents, &summary); err != nil {
		return ActivitiesGoal{}, err
	}

	return summary, nil
}

// SetActivitiesGoal sets a new activities goal on daily or weekly basis
func (m *Session) SetActivitiesGoal(period string, goals ActivitiesGoal) (ActivitiesGoal, error) {
	if period != "weekly" {
		period = "daily"
	}

	var goalsData map[string]string
	if goals.Goals.CaloriesOut != 0 {
		goalsData["caloriesOut"] = strconv.Itoa(goals.Goals.CaloriesOut)
	}
	if goals.Goals.ActiveMinutes != 0 {
		goalsData["activeMinutes"] = strconv.Itoa(goals.Goals.ActiveMinutes)
	}
	if goals.Goals.Floors != 0 {
		goalsData["floors"] = strconv.Itoa(goals.Goals.Floors)
	}
	if goals.Goals.Steps != 0 {
		goalsData["steps"] = strconv.Itoa(goals.Goals.Steps)
	}
	if goals.Goals.Distance != 0 {
		goalsData["distance"] = fmt.Sprintf("%f", goals.Goals.Distance)
	}

	contents, err := m.makePOSTRequest("https://api.fitbit.com/1/user/-/activities/goals/"+period+".json", goalsData)
	if err != nil {
		return ActivitiesGoal{}, err
	}

	summary := ActivitiesGoal{}
	if err := json.Unmarshal(contents, &summary); err != nil {
		return ActivitiesGoal{}, err
	}

	return summary, nil
}
