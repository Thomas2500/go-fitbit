package fitbit

import (
	"encoding/json"
	"errors"
	"strconv"
)

// Subscription contains response and request data of fitbit eventsx
type Subscription struct {
	CollectionType string `json:"collectionType"`
	OwnerID        string `json:"ownerId"`
	OwnerType      string `json:"ownerType"`
	SubscriberID   string `json:"subscriberId"`
	SubscriptionID string `json:"subscriptionId"`
}

// SubscriptionList contains a list of current active subscriptions
type SubscriptionList struct {
	APISubscriptions []struct {
		CollectionType string `json:"collectionType"`
		OwnerID        string `json:"ownerId"`
		OwnerType      string `json:"ownerType"`
		SubscriberID   string `json:"subscriberId"`
		SubscriptionID string `json:"subscriptionId"`
	} `json:"apiSubscriptions"`
}

// AddSubscription adds a new subscription where fitbit sends a request on changes caused by the user
func (m *Session) AddSubscription(collectionPath string, uniqueID int) (bool, Subscription, error) {
	if uniqueID == 0 {
		return false, Subscription{}, errors.New("no unique subscription id given")
	}
	if collectionPath != "" {
		collectionPath += "/"
	}
	contents, err := m.makePOSTRequest("https://api.fitbit.com/1/user/-/"+collectionPath+"apiSubscriptions/"+strconv.Itoa(uniqueID)+".json", nil)
	if err != nil {
		return false, Subscription{}, err
	}

	subscription := Subscription{}
	if err := json.Unmarshal(contents, &subscription); err != nil {
		return false, Subscription{}, err
	}

	return true, subscription, nil
}

// RemoveSubscription removes a previously added subscription
func (m *Session) RemoveSubscription(collectionPath string, uniqueID int) (bool, error) {
	if uniqueID == 0 {
		return false, errors.New("no unique subscription id given")
	}
	if collectionPath != "" {
		collectionPath += "/"
	}
	_, err := m.makeDELETERequest("https://api.fitbit.com/1/user/-/" + collectionPath + "apiSubscriptions/" + strconv.Itoa(uniqueID) + ".json")
	if err != nil {
		return false, err
	}

	return true, nil
}

// GetSubscriptions get's a list of current subscriptions
func (m *Session) GetSubscriptions(collectionPath string) (bool, SubscriptionList, error) {
	if collectionPath != "" {
		collectionPath += "/"
	}
	contents, err := m.makePOSTRequest("https://api.fitbit.com/1/user/-/"+collectionPath+"apiSubscriptions.json", nil)
	if err != nil {
		return false, SubscriptionList{}, err
	}

	subscription := SubscriptionList{}
	if err := json.Unmarshal(contents, &subscription); err != nil {
		return false, SubscriptionList{}, err
	}

	return true, subscription, nil
}
