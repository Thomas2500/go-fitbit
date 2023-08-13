package fitbit

import (
	"encoding/json"
	"errors"
	"net/url"
	"strconv"
)

// ECG data
type ECGLogList struct {
	EcgReadings []struct {
		StartTime               string `json:"startTime"`
		AverageHeartRate        int    `json:"averageHeartRate"`
		ResultClassification    string `json:"resultClassification"`
		WaveformSamples         []int  `json:"waveformSamples"`
		SamplingFrequencyHz     int    `json:"samplingFrequencyHz"`
		ScalingFactor           int    `json:"scalingFactor"`
		NumberOfWaveformSamples int    `json:"numberOfWaveformSamples"`
		LeadNumber              int    `json:"leadNumber"`
		FeatureVersion          string `json:"featureVersion"`
		DeviceName              string `json:"deviceName"`
		FirmwareVersion         string `json:"firmwareVersion"`
	} `json:"ecgReadings"`
	Pagination struct {
		AfterDate string `json:"afterDate"`
		Limit     int    `json:"limit"`
		Next      string `json:"next"`
		Offset    int    `json:"offset"`
		Previous  string `json:"previous"`
		Sort      string `json:"sort"`
	} `json:"pagination"`
}

// ECGLog returns the ECG log list
func (m *Session) ECGLog(params LogListParameters) (ECGLogList, error) {
	parameterList := url.Values{}

	//nolint:gocritic
	if params.BeforeDate != "" {
		parameterList.Add("beforeDate", params.BeforeDate)
		parameterList.Add("sort", "desc")
	} else if params.AfterDate != "" {
		parameterList.Add("afterDate", params.BeforeDate)
		parameterList.Add("sort", "asc")
	} else {
		return ECGLogList{}, errors.New("beforeDate or afterDate must be given")
	}

	if params.Limit > 0 {
		if params.Limit > 10 {
			return ECGLogList{}, errors.New("limit must be 10 or less")
		}
		parameterList.Add("limit", strconv.Itoa(params.Limit))
	}

	parameterList.Add("offset", strconv.Itoa(params.Offset))

	contents, err := m.makeRequest("https://api.fitbit.com/1/user/-/ecg/list.json?" + parameterList.Encode())
	if err != nil {
		return ECGLogList{}, err
	}

	ecgResponse := ECGLogList{}
	if err := json.Unmarshal(contents, &ecgResponse); err != nil {
		return ECGLogList{}, err
	}

	return ecgResponse, nil
}
