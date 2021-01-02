package fitbit

import (
	"encoding/json"
	"errors"
	"net/url"
	"strconv"
	"time"
)

func (m *Session) ActivityLog(params LogListParameters) (ActivitiesLogList, error) {

	parameterList := url.Values{}
	if params.BeforeDate != "" {
		parameterList.Add("beforeDate", params.BeforeDate)
		parameterList.Add("sort", "desc")
	} else if params.AfterDate != "" {
		parameterList.Add("afterDate", params.BeforeDate)
		parameterList.Add("sort", "asc")
	} else {
		return ActivitiesLogList{}, errors.New("beforeDate or afterDate must be given")
	}

	if params.Limit > 0 {
		if params.Limit > 20 {
			return ActivitiesLogList{}, errors.New("limit must be 20 or less")
		}
		parameterList.Add("limit", strconv.Itoa(params.Limit))
	}

	parameterList.Add("offset", strconv.Itoa(params.Offset))

	contents, err := m.makeRequest("https://api.fitbit.com/1/user/-/activities/list.json?" + parameterList.Encode())
	if err != nil {
		return ActivitiesLogList{}, err
	}

	activityResponse := ActivitiesLogList{}
	if err := json.Unmarshal(contents, &activityResponse); err != nil {
		return ActivitiesLogList{}, err
	}

	return activityResponse, nil
}

type LogListParameters struct {
	BeforeDate string
	AfterDate  string
	Limit      int
	Offset     int
}

type ActivitiesLogList struct {
	Activities []struct {
		ActiveDuration int `json:"activeDuration"`
		ActivityLevel  []struct {
			Minutes int    `json:"minutes"`
			Name    string `json:"name"`
		} `json:"activityLevel"`
		ActivityName          string    `json:"activityName"`
		ActivityTypeID        int       `json:"activityTypeId"`
		Calories              int       `json:"calories"`
		CaloriesLink          string    `json:"caloriesLink"`
		Distance              float64   `json:"distance"`
		DistanceUnit          string    `json:"distanceUnit"`
		Duration              int       `json:"duration"`
		ElevationGain         float64   `json:"elevationGain"`
		HasActiveZoneMinutes  bool      `json:"hasActiveZoneMinutes"`
		LastModified          time.Time `json:"lastModified"`
		LogID                 int64     `json:"logId"`
		LogType               string    `json:"logType"`
		ManualValuesSpecified struct {
			Calories bool `json:"calories"`
			Distance bool `json:"distance"`
			Steps    bool `json:"steps"`
		} `json:"manualValuesSpecified"`
		OriginalDuration  int       `json:"originalDuration"`
		OriginalStartTime time.Time `json:"originalStartTime"`
		PoolLength        int       `json:"poolLength,omitempty"`
		PoolLengthUnit    string    `json:"poolLengthUnit,omitempty"`
		Source            struct {
			ID              string   `json:"id"`
			Name            string   `json:"name"`
			TrackerFeatures []string `json:"trackerFeatures"`
			Type            string   `json:"type"`
			URL             string   `json:"url"`
		} `json:"source"`
		Speed            float64   `json:"speed"`
		StartTime        time.Time `json:"startTime"`
		AverageHeartRate int       `json:"averageHeartRate,omitempty"`
		DetailsLink      string    `json:"detailsLink,omitempty"`
		HeartRateLink    string    `json:"heartRateLink,omitempty"`
		HeartRateZones   []struct {
			Max     int    `json:"max"`
			Min     int    `json:"min"`
			Minutes int    `json:"minutes"`
			Name    string `json:"name"`
		} `json:"heartRateZones,omitempty"`
		Pace    float64 `json:"pace,omitempty"`
		Steps   int     `json:"steps,omitempty"`
		TcxLink string  `json:"tcxLink,omitempty"`
	} `json:"activities"`
	Pagination struct {
		BeforeDate string `json:"beforeDate"`
		Limit      int    `json:"limit"`
		Next       string `json:"next"`
		Offset     int    `json:"offset"`
		Previous   string `json:"previous"`
		Sort       string `json:"sort"`
	} `json:"pagination"`
}

// LogActivity logs a new activity
// date must be in the format yyyy-MM-dd
// TODO: TESTME
func (m *Session) LogActivity(activity NewActivity) (NewActivityResponse, error) {
	postData := make(map[string]string)
	if activity.ActivityID != 0 {
		postData["activityId"] = strconv.Itoa(activity.ActivityID)
	} else if activity.ActivityName != "" {
		postData["activityName"] = activity.ActivityName
		if activity.ManualCalories != 0 {
			postData["manualCalories"] = strconv.Itoa(activity.ManualCalories)
		} else {
			return NewActivityResponse{}, errors.New("manualCalories must be given if activityName is used")
		}
	} else {
		return NewActivityResponse{}, errors.New("activityId or activityName must be given")
	}

	if activity.StartTime != "" {
		postData["startTime"] = activity.StartTime
	} else {
		return NewActivityResponse{}, errors.New("startTime must be given")
	}

	if activity.DurationMillis > 0 {
		postData["durationMillis"] = strconv.FormatInt(activity.DurationMillis, 10)
	} else {
		return NewActivityResponse{}, errors.New("durationMillis must be given")
	}

	if activity.Date != "" {
		postData["date"] = activity.Date
	} else {
		return NewActivityResponse{}, errors.New("date must be given")
	}

	if activity.Distance > 0 {
		postData["distance"] = strconv.FormatFloat(activity.Distance, 'f', 3, 64)
	}
	if activity.DistanceUnit != "" {
		postData["distanceUnit"] = activity.DistanceUnit
	}

	contents, err := m.makePOSTRequest("https://api.fitbit.com/1/user/-/activities.json", postData)
	if err != nil {
		return NewActivityResponse{}, err
	}

	activityResponse := NewActivityResponse{}
	if err := json.Unmarshal(contents, &activityResponse); err != nil {
		return NewActivityResponse{}, err
	}

	return activityResponse, nil
}

type NewActivity struct {
	ActivityID     int     `json:"activityId,omitempty"`
	ActivityName   string  `json:"activityName,omitempty"`
	ManualCalories int     `json:"manualCalories,omitempty"`
	StartTime      string  `json:"startTime"`              // HH:mm
	DurationMillis int64   `json:"durationMillis"`         // milliseconds
	Date           string  `json:"date"`                   // yyyy-MM-dd
	Distance       float64 `json:"distance,omitempty"`     // x.xx
	DistanceUnit   string  `json:"distanceUnit,omitempty"` // Steps units are available only for "Walking" (activityId=90013) and "Running" (activityId=90009) directory activities and their intensity levels.
}

type NewActivityResponse struct {
	ActivityLog struct {
		ActivityID       int     `json:"activityId"`
		ActivityParentID int     `json:"activityParentId"`
		Calories         int     `json:"calories"`
		Description      string  `json:"description"`
		Distance         float64 `json:"distance"`
		Duration         int     `json:"duration"`
		IsFavorite       bool    `json:"isFavorite"`
		LogID            int     `json:"logId"`
		Name             string  `json:"name"`
		StartTime        string  `json:"startTime"`
		Steps            int     `json:"steps"`
	} `json:"activityLog"`
}

// RemoveActivity deletes an existing activity by activity log id
// TODO: TESTME
func (m *Session) RemoveActivity(id int) error {
	_, err := m.makeDELETERequest("https://api.fitbit.com/1/user/-/activities/" + strconv.Itoa(id) + ".json")
	if err != nil {
		return err
	}

	return nil
}
