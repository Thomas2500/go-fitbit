package fitbit

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// BodyWeight contains one or multiple records
type BodyWeight struct {
	Weight []struct {
		Bmi    float64 `json:"bmi"`
		Date   string  `json:"date"`
		Fat    float64 `json:"fat"`
		LogID  int64   `json:"logId"`
		Source string  `json:"source"`
		Time   string  `json:"time"`
		Weight float64 `json:"weight"`
	} `json:"weight"`
}

// BodyWeightLogByDay returns the weight log by a given date
// date must be in the format yyyy-MM-dd
func (m *Session) BodyWeightLogByDay(day string) (BodyWeight, error) {
	contents, err := m.makeRequest("https://api.fitbit.com/1/user/-/body/log/weight/date/" + day + ".json")
	if err != nil {
		return BodyWeight{}, err
	}

	weight := BodyWeight{}
	if err := json.Unmarshal(contents, &weight); err != nil {
		return BodyWeight{}, err
	}

	return weight, nil
}

// BodyWeightLogByDateRange returns the weight log of a given time range by date
// date must be in the format yyyy-MM-dd
func (m *Session) BodyWeightLogByDateRange(startDay string, endDay string) (BodyWeight, error) {
	contents, err := m.makeRequest("https://api.fitbit.com/1/user/-/body/log/weight/date/" + startDay + "/" + endDay + ".json")
	if err != nil {
		return BodyWeight{}, err
	}

	weight := BodyWeight{}
	if err := json.Unmarshal(contents, &weight); err != nil {
		return BodyWeight{}, err
	}

	return weight, nil
}

// AddBodyWeight adds a new body weight record
// date must be in the format yyyy-MM-dd
func (m *Session) AddBodyWeight(day string, weight float64) (BodyWeight, error) {
	contents, err := m.makePOSTRequest("https://api.fitbit.com/1/user/-/body/log/weight.json", map[string]string{
		"date":   day,
		"weight": fmt.Sprintf("%f", weight),
	})
	if err != nil {
		return BodyWeight{}, err
	}

	weightResponse := BodyWeight{}
	if err := json.Unmarshal(contents, &weightResponse); err != nil {
		return BodyWeight{}, err
	}

	return weightResponse, nil
}

// RemoveBodyWeight removes a existing record by it's log ID
func (m *Session) RemoveBodyWeight(logID int64) error {
	_, err := m.makeDELETERequest("https://api.fitbit.com/1/user/-/body/log/weight/" + strconv.FormatInt(logID, 10) + ".json")
	if err != nil {
		return err
	}

	return nil
}
