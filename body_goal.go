package fitbit

import (
	"encoding/json"
	"fmt"
)

// BodyWeightGoal contains the currently set body goal by the user
type BodyWeightGoal struct {
	Goal struct {
		GoalType        string  `json:"goalType"`
		StartDate       string  `json:"startDate"`
		StartWeight     float64 `json:"startWeight"`
		Weight          int     `json:"weight"`
		WeightThreshold float64 `json:"weightThreshold"`
	} `json:"goal"`
}

// BodyFatGoal contains the currently set fat goal of the user
type BodyFatGoal struct {
	Goal struct {
		Fat int `json:"fat"`
	} `json:"goal"`
}

// BodyWeightGoal requests the weight goal of the user
func (m *Session) BodyWeightGoal() (BodyWeightGoal, error) {
	contents, err := m.makeRequest("https://api.fitbit.com/1/user/-/body/log/weight/goal.json")
	if err != nil {
		return BodyWeightGoal{}, err
	}

	weightGoal := BodyWeightGoal{}
	if err := json.Unmarshal(contents, &weightGoal); err != nil {
		return BodyWeightGoal{}, err
	}

	return weightGoal, nil
}

// SetBodyWeightGoal sets the users body fat goal
func (m *Session) SetBodyWeightGoal(startDate string, startWeight float64, weightGoal float64) (BodyWeightGoal, error) {
	contents, err := m.makePOSTRequest("https://api.fitbit.com/1/user/-/body/log/weight/goal.json", map[string]string{
		"startDate":   startDate,
		"startWeight": fmt.Sprintf("%f", startWeight),
		"weight":      fmt.Sprintf("%f", weightGoal),
	})
	if err != nil {
		return BodyWeightGoal{}, err
	}

	weightGoalResponse := BodyWeightGoal{}
	if err := json.Unmarshal(contents, &weightGoalResponse); err != nil {
		return BodyWeightGoal{}, err
	}

	return weightGoalResponse, nil
}

// BodyFatGoal requests the fat goal of the user
func (m *Session) BodyFatGoal() (BodyFatGoal, error) {
	contents, err := m.makeRequest("https://api.fitbit.com/1/user/-/body/log/fat/goal.json")
	if err != nil {
		return BodyFatGoal{}, err
	}

	fatGoal := BodyFatGoal{}
	if err := json.Unmarshal(contents, &fatGoal); err != nil {
		return BodyFatGoal{}, err
	}

	return fatGoal, nil
}

// SetBodyFatGoal sets the users body fat goal
func (m *Session) SetBodyFatGoal(targetFat float64) (FoodGoal, error) {
	contents, err := m.makePOSTRequest("https://api.fitbit.com/1/user/-/body/log/fat/goal.json", map[string]string{
		"fat": fmt.Sprintf("%f", targetFat),
	})
	if err != nil {
		return FoodGoal{}, err
	}

	foodgoal := FoodGoal{}
	if err := json.Unmarshal(contents, &foodgoal); err != nil {
		return FoodGoal{}, err
	}

	return foodgoal, nil
}
