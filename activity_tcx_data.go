package fitbit

import "encoding/xml"

func ReadTCX(data []byte) (GarminTrainingCenterDatabasev2, error) {
	activity := GarminTrainingCenterDatabasev2{}
	err := xml.Unmarshal(data, &activity)
	return activity, err
}

// ActivityTCX contains the activity TCX for a given activity
// Fitbit is using http://www.garmin.com/xmlschemas/TrainingCenterDatabase/v2 for the TCX format
// Also available at https://www8.garmin.com/xmlschemas/TrainingCenterDatabasev2.xsd
// TODO: Extend struct below to include all available fields
type GarminTrainingCenterDatabasev2 struct {
	XMLName    xml.Name `xml:"TrainingCenterDatabase"`
	Text       string   `xml:",chardata"`
	Xmlns      string   `xml:"xmlns,attr"`
	Activities struct {
		Text     string `xml:",chardata"`
		Activity struct {
			Text  string `xml:",chardata"`
			Sport string `xml:"Sport,attr"`
			ID    string `xml:"Id"`
			Lap   []struct {
				Text             string `xml:",chardata"`
				StartTime        string `xml:"StartTime,attr"`
				TotalTimeSeconds string `xml:"TotalTimeSeconds"`
				DistanceMeters   string `xml:"DistanceMeters"`
				Calories         string `xml:"Calories"`
				Intensity        string `xml:"Intensity"`
				TriggerMethod    string `xml:"TriggerMethod"`
				Track            struct {
					Text       string `xml:",chardata"`
					Trackpoint []struct {
						Text     string `xml:",chardata"`
						Time     string `xml:"Time"`
						Position struct {
							Text             string `xml:",chardata"`
							LatitudeDegrees  string `xml:"LatitudeDegrees"`
							LongitudeDegrees string `xml:"LongitudeDegrees"`
						} `xml:"Position"`
						AltitudeMeters string `xml:"AltitudeMeters"`
						DistanceMeters string `xml:"DistanceMeters"`
						HeartRateBpm   struct {
							Text  string `xml:",chardata"`
							Value string `xml:"Value"`
						} `xml:"HeartRateBpm"`
					} `xml:"Trackpoint"`
				} `xml:"Track"`
			} `xml:"Lap"`
			Creator struct {
				Text      string `xml:",chardata"`
				Type      string `xml:"type,attr"`
				Xsi       string `xml:"xsi,attr"`
				Name      string `xml:"Name"`
				UnitID    string `xml:"UnitId"`
				ProductID string `xml:"ProductID"`
			} `xml:"Creator"`
		} `xml:"Activity"`
	} `xml:"Activities"`
}
