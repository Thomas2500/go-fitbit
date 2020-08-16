package fitbit

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// BodyFat contains one or multiple records, similar to BodyFat but without weight
type BodyFat struct {
	Fat []struct {
		Date   string  `json:"date"`
		Fat    float64 `json:"fat"`
		LogID  int64   `json:"logId"`
		Source string  `json:"source"`
		Time   string  `json:"time"`
	} `json:"fat"`
}

// BodyFatLogByDay returns the fat log by a given date
// date must be in the format yyyy-MM-dd
func (m *Session) BodyFatLogByDay(day string) (BodyFat, error) {
	contents, err := m.makeRequest("https://api.fitbit.com/1/user/-/body/log/fat/date/" + day + ".json")
	if err != nil {
		return BodyFat{}, err
	}

	fat := BodyFat{}
	if err := json.Unmarshal(contents, &fat); err != nil {
		return BodyFat{}, err
	}

	return fat, nil
}

// BodyFatLogByDateRange returns the weight log of a given time range by date
// date must be in the format yyyy-MM-dd
func (m *Session) BodyFatLogByDateRange(startDay string, endDay string) (BodyFat, error) {
	contents, err := m.makeRequest("https://api.fitbit.com/1/user/-/body/log/fat/date/" + startDay + "/" + endDay + ".json")
	if err != nil {
		return BodyFat{}, err
	}

	fat := BodyFat{}
	if err := json.Unmarshal(contents, &fat); err != nil {
		return BodyFat{}, err
	}

	return fat, nil
}

// AddBodyFat adds a new body weight record
// date must be in the format yyyy-MM-dd
func (m *Session) AddBodyFat(day string, fat float64) (BodyFat, error) {
	contents, err := m.makePOSTRequest("https://api.fitbit.com/1/user/-/body/log/fat.json", map[string]string{
		"date": day,
		"fat":  fmt.Sprintf("%f", fat),
	})
	if err != nil {
		return BodyFat{}, err
	}

	fatResponse := BodyFat{}
	if err := json.Unmarshal(contents, &fatResponse); err != nil {
		return BodyFat{}, err
	}

	return fatResponse, nil
}

// RemoveBodyFat removes a existing record by it's log ID
func (m *Session) RemoveBodyFat(logID int64) error {
	_, err := m.makeDELETERequest("https://api.fitbit.com/1/user/-/body/log/fat/" + strconv.FormatInt(logID, 10) + ".json")
	if err != nil {
		return err
	}

	return nil
}
