package fitbit

import (
	"encoding/json"
	"strconv"
)

// BadgesList contains a list of badges
type BadgesList struct {
	Badges []Badge `json:"badges"`
}

// Badge contains information about a badge
type Badge struct {
	BadgeGradientEndColor   string        `json:"badgeGradientEndColor"`
	BadgeGradientStartColor string        `json:"badgeGradientStartColor"`
	BadgeType               string        `json:"badgeType"`
	Category                string        `json:"category"`
	Cheers                  []interface{} `json:"cheers"` // FIXME: unknown data
	DateTime                string        `json:"dateTime"`
	Description             string        `json:"description"`
	EarnedMessage           string        `json:"earnedMessage,omitempty"`
	EncodedID               string        `json:"encodedId"`
	Image100Px              string        `json:"image100px"`
	Image125Px              string        `json:"image125px"`
	Image300Px              string        `json:"image300px"`
	Image50Px               string        `json:"image50px"`
	Image75Px               string        `json:"image75px"`
	MarketingDescription    string        `json:"marketingDescription"`
	MobileDescription       string        `json:"mobileDescription"`
	Name                    string        `json:"name"`
	ShareImage640Px         string        `json:"shareImage640px"`
	ShareText               string        `json:"shareText"`
	ShortDescription        string        `json:"shortDescription"`
	ShortName               string        `json:"shortName"`
	TimesAchieved           int           `json:"timesAchieved"`
	Value                   int           `json:"value,omitempty"`
	Unit                    string        `json:"unit,omitempty"`
}

// Badges returns a list of user badges
func (m *Session) Badges(userID uint64) (BadgesList, error) {
	// Default "-" is current logged in user
	requestID := "-"
	if userID > 0 {
		requestID = strconv.FormatUint(userID, 10)
	}
	contents, err := m.makeRequest("https://api.fitbit.com/1/user/" + requestID + "/badges.json")
	if err != nil {
		return BadgesList{}, err
	}

	badgeList := BadgesList{}
	if err := json.Unmarshal(contents, &badgeList); err != nil {
		return BadgesList{}, err
	}

	return badgeList, nil
}
