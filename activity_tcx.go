package fitbit

import (
	"errors"
	"fmt"
)

// ActivityTCX returns the activity TCX for the given activity
func (m *Session) ActivityTCX(logID int64) ([]byte, error) {
	if logID == 0 {
		return nil, errors.New("logID must be given")
	}

	// Fetch data from Fitbit
	contents, err := m.makeRequest(fmt.Sprintf("https://api.fitbit.com/1/user/-/activities/%d.tcx?includePartialTCX=true", logID))
	if err != nil {
		return nil, err
	}

	// TCX data isn't parsed here, return it as-is
	return contents, nil
}
