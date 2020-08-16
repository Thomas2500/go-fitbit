package fitbit

import (
	"encoding/json"
	"errors"
	"net/url"
	"strconv"
	"time"
)

// SleepDay contains data of a sleep day
type SleepDay struct {
	Sleep []struct {
		DateOfSleep string `json:"dateOfSleep"`
		Duration    int    `json:"duration"`
		Efficiency  int    `json:"efficiency"`
		EndTime     string `json:"endTime"`
		InfoCode    int    `json:"infoCode"`
		IsMainSleep bool   `json:"isMainSleep"`
		Levels      struct {
			Data []struct {
				DateTime string `json:"dateTime"`
				Level    string `json:"level"`
				Seconds  int    `json:"seconds"`
			} `json:"data"`
			ShortData []struct {
				DateTime string `json:"dateTime"`
				Level    string `json:"level"`
				Seconds  int    `json:"seconds"`
			} `json:"shortData"`
			Summary struct {
				Deep struct {
					Count               int `json:"count"`
					Minutes             int `json:"minutes"`
					ThirtyDayAvgMinutes int `json:"thirtyDayAvgMinutes"`
				} `json:"deep"`
				Light struct {
					Count               int `json:"count"`
					Minutes             int `json:"minutes"`
					ThirtyDayAvgMinutes int `json:"thirtyDayAvgMinutes"`
				} `json:"light"`
				Rem struct {
					Count               int `json:"count"`
					Minutes             int `json:"minutes"`
					ThirtyDayAvgMinutes int `json:"thirtyDayAvgMinutes"`
				} `json:"rem"`
				Wake struct {
					Count               int `json:"count"`
					Minutes             int `json:"minutes"`
					ThirtyDayAvgMinutes int `json:"thirtyDayAvgMinutes"`
				} `json:"wake"`
				Asleep struct {
					Count               int `json:"count"`
					Minutes             int `json:"minutes"`
					ThirtyDayAvgMinutes int `json:"thirtyDayAvgMinutes,omitempty"`
				} `json:"asleep,omitempty"`
				Awake struct {
					Count               int `json:"count"`
					Minutes             int `json:"minutes"`
					ThirtyDayAvgMinutes int `json:"thirtyDayAvgMinutes,omitempty"`
				} `json:"awake,omitempty"`
				Restless struct {
					Count               int `json:"count"`
					Minutes             int `json:"minutes"`
					ThirtyDayAvgMinutes int `json:"thirtyDayAvgMinutes,omitempty"`
				} `json:"restless,omitempty"`
			} `json:"summary"`
		} `json:"levels,omitempty"`
		LogID               int64  `json:"logId"`
		MinutesAfterWakeup  int    `json:"minutesAfterWakeup"`
		MinutesAsleep       int    `json:"minutesAsleep"`
		MinutesAwake        int    `json:"minutesAwake"`
		MinutesToFallAsleep int    `json:"minutesToFallAsleep"`
		StartTime           string `json:"startTime"`
		TimeInBed           int    `json:"timeInBed"`
		Type                string `json:"type"`
	} `json:"sleep"`
	Summary struct {
		Stages struct {
			Deep  int `json:"deep"`
			Light int `json:"light"`
			Rem   int `json:"rem"`
			Wake  int `json:"wake"`
		} `json:"stages"`
		TotalMinutesAsleep int `json:"totalMinutesAsleep"`
		TotalSleepRecords  int `json:"totalSleepRecords"`
		TotalTimeInBed     int `json:"totalTimeInBed"`
	} `json:"summary,omitempty"`
	Meta struct {
		RetryDuration int    `json:"retryDuration"`
		State         string `json:"state"`
	} `json:"meta,omitempty"`
}

// SleepByDay returns the sleep data for a given date
// date must be in the format yyyy-MM-dd
func (m *Session) SleepByDay(day string) (SleepDay, error) {
	contents, err := m.makeRequest("https://api.fitbit.com/1.2/user/-/sleep/date/" + day + ".json")
	if err != nil {
		return SleepDay{}, err
	}

	sleep := SleepDay{}
	if err := json.Unmarshal(contents, &sleep); err != nil {
		return SleepDay{}, err
	}

	return sleep, nil
}

// SleepByDayRange returns the sleep data for a given date range
// date must be in the format yyyy-MM-dd
func (m *Session) SleepByDayRange(startDay string, endDay string) (SleepDay, error) {
	contents, err := m.makeRequest("https://api.fitbit.com/1.2/user/-/sleep/date/" + startDay + "/" + endDay + ".json")
	if err != nil {
		return SleepDay{}, err
	}

	sleep := SleepDay{}
	if err := json.Unmarshal(contents, &sleep); err != nil {
		return SleepDay{}, err
	}

	return sleep, nil
}

func (m *Session) SleepLogList(params LogListParameters) (SleepLogList, error) {
	parameterList := url.Values{}
	if params.BeforeDate != "" {
		parameterList.Add("beforeDate", params.BeforeDate)
		parameterList.Add("sort", "desc")
	} else if params.AfterDate != "" {
		parameterList.Add("afterDate", params.BeforeDate)
		parameterList.Add("sort", "asc")
	} else {
		return SleepLogList{}, errors.New("beforeDate or afterDate must be given")
	}

	if params.Limit > 0 {
		if params.Limit > 20 {
			return SleepLogList{}, errors.New("limit must be 20 or less")
		}
		parameterList.Add("limit", strconv.Itoa(params.Limit))
	}

	parameterList.Add("offset", strconv.Itoa(params.Offset))

	contents, err := m.makeRequest("https://api.fitbit.com/1/user/-/sleep/list.json?" + parameterList.Encode())
	if err != nil {
		return SleepLogList{}, err
	}

	activityResponse := SleepLogList{}
	if err := json.Unmarshal(contents, &activityResponse); err != nil {
		return SleepLogList{}, err
	}

	return activityResponse, nil
}

type SleepLogList struct {
	Pagination struct {
		BeforeDate string `json:"beforeDate"`
		Limit      int    `json:"limit"`
		Next       string `json:"next"`
		Offset     int    `json:"offset"`
		Previous   string `json:"previous"`
		Sort       string `json:"sort"`
	} `json:"pagination"`
	Sleep []struct {
		AwakeCount      int    `json:"awakeCount,omitempty"`
		AwakeDuration   int    `json:"awakeDuration,omitempty"`
		AwakeningsCount int    `json:"awakeningsCount,omitempty"`
		DateOfSleep     string `json:"dateOfSleep"`
		Duration        int    `json:"duration"`
		Efficiency      int    `json:"efficiency"`
		EndTime         string `json:"endTime"`
		InfoCode        int    `json:"infoCode"`
		IsMainSleep     bool   `json:"isMainSleep"`
		Levels          struct {
			Data []struct {
				DateTime string `json:"dateTime"`
				Level    string `json:"level"`
				Seconds  int    `json:"seconds"`
			} `json:"data,omitempty"`
			ShortData []struct {
				DateTime string `json:"dateTime"`
				Level    string `json:"level"`
				Seconds  int    `json:"seconds"`
			} `json:"shortData,omitempty"`
			Summary struct {
				Deep struct {
					Count               int `json:"count"`
					Minutes             int `json:"minutes"`
					ThirtyDayAvgMinutes int `json:"thirtyDayAvgMinutes"`
				} `json:"deep,omitempty"`
				Light struct {
					Count               int `json:"count"`
					Minutes             int `json:"minutes"`
					ThirtyDayAvgMinutes int `json:"thirtyDayAvgMinutes"`
				} `json:"light,omitempty"`
				Rem struct {
					Count               int `json:"count"`
					Minutes             int `json:"minutes"`
					ThirtyDayAvgMinutes int `json:"thirtyDayAvgMinutes"`
				} `json:"rem,omitempty"`
				Wake struct {
					Count               int `json:"count"`
					Minutes             int `json:"minutes"`
					ThirtyDayAvgMinutes int `json:"thirtyDayAvgMinutes"`
				} `json:"wake,omitempty"`
			} `json:"summary,omitempty"`
		} `json:"levels,omitempty"`
		LogID               int64  `json:"logId"`
		MinutesAfterWakeup  int    `json:"minutesAfterWakeup"`
		MinutesAsleep       int    `json:"minutesAsleep"`
		MinutesAwake        int    `json:"minutesAwake"`
		MinutesToFallAsleep int    `json:"minutesToFallAsleep"`
		RestlessCount       int    `json:"restlessCount,omitempty"`
		RestlessDuration    int    `json:"restlessDuration,omitempty"`
		StartTime           string `json:"startTime"`
		TimeInBed           int    `json:"timeInBed"`
		Type                string `json:"type"`
		MinuteData          []struct {
			DateTime string `json:"dateTime"`
			Value    string `json:"value"`
		} `json:"minuteData,omitempty"`
	} `json:"sleep"`
}

// AddSleep adds a new sleep record
// date in form of yyyy-MM-dd
// startTime in form of HH:mm
// duration in milliseconds
func (m *Session) AddSleep(date string, startTime string, duration int64) (SleepDay, error) {
	parameterList := url.Values{}
	if date != "" {
		parameterList.Add("date", date)
	} else {
		return SleepDay{}, errors.New("date must be given")
	}
	if startTime != "" {
		parameterList.Add("startTime", startTime)
	} else {
		return SleepDay{}, errors.New("startTime must be given")
	}

	if duration > 0 {
		parameterList.Add("duration", strconv.FormatInt(duration, 10))
	}

	contents, err := m.makeRequest("https://api.fitbit.com/1.2/user/-/sleep.json?" + parameterList.Encode())
	if err != nil {
		return SleepDay{}, err
	}

	activityResponse := SleepDay{}
	if err := json.Unmarshal(contents, &activityResponse); err != nil {
		return SleepDay{}, err
	}

	return activityResponse, nil
}

// RemoveSleep removes a sleep entry
func (m *Session) RemoveSleep(sleepID uint64) error {
	_, err := m.makeDELETERequest("https://api.fitbit.com/1/user/-/sleep/" + strconv.FormatUint(sleepID, 10) + ".json")
	if err != nil {
		return err
	}

	return nil
}

// SleepGoal requests the sleep goal of the user
func (m *Session) SleepGoal() (SleepGoal, error) {
	contents, err := m.makeRequest("https://api.fitbit.com/1.2/user/-/sleep/goal.json")
	if err != nil {
		return SleepGoal{}, err
	}

	sleepGoalResponse := SleepGoal{}
	if err := json.Unmarshal(contents, &sleepGoalResponse); err != nil {
		return SleepGoal{}, err
	}

	return sleepGoalResponse, nil
}

type SleepGoal struct {
	Consistency struct {
		AwakeRestlessPercentage float64 `json:"awakeRestlessPercentage"`
		FlowID                  int     `json:"flowId"`
		RecommendedSleepGoal    int     `json:"recommendedSleepGoal"`
		TypicalDuration         int     `json:"typicalDuration"`
		TypicalWakeupTime       string  `json:"typicalWakeupTime"`
	} `json:"consistency"`
	Goal struct {
		Bedtime     string    `json:"bedtime"`
		MinDuration int       `json:"minDuration"`
		UpdatedOn   time.Time `json:"updatedOn"`
		WakeupTime  string    `json:"wakeupTime"`
	} `json:"goal"`
}

// SetSleepGoal requests the sleep goal of the user
func (m *Session) SetSleepGoal(minDuration int) (SleepGoal, error) {
	contents, err := m.makePOSTRequest("https://api.fitbit.com/1.2/user/-/sleep/goal.json", map[string]string{
		"duration": strconv.Itoa(minDuration),
	})
	if err != nil {
		return SleepGoal{}, err
	}

	sleepGoalResponse := SleepGoal{}
	if err := json.Unmarshal(contents, &sleepGoalResponse); err != nil {
		return SleepGoal{}, err
	}

	return sleepGoalResponse, nil
}
