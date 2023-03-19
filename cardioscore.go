package fitbit

import (
	"encoding/json"
	"fmt"
)

// CardioFitnessScoreLog contains the cardio fitness score (VO2Max) for a given date
type CardioFitnessScoreLog struct {
	CardioScore []struct {
		DateTime string `json:"dateTime"`
		Value    struct {
			Vo2Max string `json:"vo2Max"` // VO2 Max in mL/kg/min
		} `json:"value"`
	} `json:"cardioScore"`
}

// CardioFitnessScoreByDay returns the cardio fitness score (VO2Max) for a given date
// date must be in the format yyyy-MM-dd, scope ScopeCardioFitness must be granted
func (m *Session) CardioFitnessScoreByDay(date string) (CardioFitnessScoreLog, error) {
	contents, err := m.makeRequest(fmt.Sprintf("https://api.fitbit.com/1/user/-/cardioscore/date/%s.json", date))
	if err != nil {
		return CardioFitnessScoreLog{}, err
	}

	summary := CardioFitnessScoreLog{}
	if err := json.Unmarshal(contents, &summary); err != nil {
		return CardioFitnessScoreLog{}, err
	}

	return summary, nil
}

// CardioFitnessScoreByDateRange returns the cardio fitness score (VO2Max) for the given date range
// date must be in the format yyyy-MM-dd, scope ScopeCardioFitness must be granted
func (m *Session) CardioFitnessScoreByDateRange(startDate, endDate string) (CardioFitnessScoreLog, error) {
	contents, err := m.makeRequest(fmt.Sprintf("https://api.fitbit.com/1/user/-/cardioscore/date/%s/%s.json", startDate, endDate))
	if err != nil {
		return CardioFitnessScoreLog{}, err
	}

	summary := CardioFitnessScoreLog{}
	if err := json.Unmarshal(contents, &summary); err != nil {
		return CardioFitnessScoreLog{}, err
	}

	return summary, nil
}
