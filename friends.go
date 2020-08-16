package fitbit

import (
	"encoding/json"
	"strconv"
)

// FriendsList contins a list of friends
type FriendsList struct {
	Data []Friend `json:"data"`
}

// Friend describes a friend
type Friend struct {
	Attributes struct {
		Friend bool   `json:"friend"`
		Child  bool   `json:"child"`
		Name   string `json:"name"`
		Avatar string `json:"avatar"`
	} `json:"attributes"`
	ID   string `json:"id"`
	Type string `json:"type"`
}

// GetFriends returns the current friends of the current user
func (m *Session) GetFriends() (FriendsList, error) {
	contents, err := m.makeRequest("https://api.fitbit.com/1.1/user/-/friends.json")
	if err != nil {
		return FriendsList{}, err
	}

	friends := FriendsList{}
	if err := json.Unmarshal(contents, &friends); err != nil {
		return FriendsList{}, err
	}

	return friends, nil
}

// FriendsLeaderboard contains the leaderboard of a user
type FriendsLeaderboard struct {
	Data []struct {
		Attributes struct {
			StepSummary int `json:"step-summary"`
			StepRank    int `json:"step-rank"`
		} `json:"attributes"`
		Relationships struct {
			User struct {
				Data struct {
					ID   string `json:"id"`
					Type string `json:"type"`
				} `json:"data"`
			} `json:"user"`
		} `json:"relationships"`
		ID   string `json:"id"`
		Type string `json:"type"`
	} `json:"data"`
	Included []struct {
		Attributes struct {
			Friend bool   `json:"friend,omitempty"`
			Child  bool   `json:"child"`
			Name   string `json:"name"`
			Avatar string `json:"avatar"`
		} `json:"attributes,omitempty"`
		ID   string `json:"id"`
		Type string `json:"type"`
	} `json:"included"`
}

//GetFriendsLeaderboard returns the leaderbord including the user
func (m *Session) GetFriendsLeaderboard() (FriendsLeaderboard, error) {
	contents, err := m.makeRequest("https://api.fitbit.com/1.1/user/-/leaderboard/friends.json")
	if err != nil {
		return FriendsLeaderboard{}, err
	}

	friends := FriendsLeaderboard{}
	if err := json.Unmarshal(contents, &friends); err != nil {
		return FriendsLeaderboard{}, err
	}

	return friends, nil
}

// FriendInviteByEmail invites another user to be friends by email
// FIXME: untested function, unknown response from server
func (m *Session) FriendInviteByEmail(value string) ([]byte, error) {
	data := map[string]string{
		"invitedUserEmail": value,
	}

	contents, err := m.makePOSTRequest("https://api.fitbit.com/1.1/user/-/friends/invitations", data)
	if err != nil {
		return nil, err
	}

	return contents, nil
}

// FriendInviteByUserID invites another user to be friends by user id
// FIXME: untested function, unknown response from server
func (m *Session) FriendInviteByUserID(value string) ([]byte, error) {
	data := map[string]string{
		"invitedUserId": value,
	}

	contents, err := m.makePOSTRequest("https://api.fitbit.com/1.1/user/-/friends/invitations", data)
	if err != nil {
		return nil, err
	}

	return contents, nil
}

// FriendsInvitations contains a list of open invitations to the user
type FriendsInvitations struct {
	Data []struct {
		Type       string `json:"type"`
		ID         string `json:"id"`
		Attributes struct {
			DateTime string `json:"dateTime"`
			Email    string `json:"email,omitempty"`
			Source   string `json:"source"`
		} `json:"attributes,omitempty"`
		Relationships struct {
			User struct {
				Data struct {
					Type string `json:"type"`
					ID   string `json:"id"`
				} `json:"data"`
			} `json:"user"`
		} `json:"relationships,omitempty"`
	} `json:"data,omitempty"`
	Included []struct {
		Type       string `json:"type"`
		ID         string `json:"id"`
		Attributes struct {
			Friend bool   `json:"friend"`
			Child  bool   `json:"child"`
			Avatar string `json:"avatar"`
			Name   string `json:"name"`
		} `json:"attributes"`
	} `json:"included,omitempty"`
}

// GetFriendInvitations returns a list of open friend invitations
func (m *Session) GetFriendInvitations() (FriendsInvitations, error) {
	contents, err := m.makeRequest("https://api.fitbit.com/1.1/user/-/friends/invitations.json")
	if err != nil {
		return FriendsInvitations{}, err
	}

	friends := FriendsInvitations{}
	if err := json.Unmarshal(contents, &friends); err != nil {
		return FriendsInvitations{}, err
	}

	return friends, nil
}

// FriendsRespondToInvition reacts to a inviation
// FIXME: untested function, unknown response from server
func (m *Session) FriendsRespondToInvition(userID string, accept bool) ([]byte, error) {
	data := map[string]string{
		"accept": strconv.FormatBool(accept),
	}

	contents, err := m.makePOSTRequest("https://api.fitbit.com/1.1/user/-/friends/invitations/"+userID, data)
	if err != nil {
		return nil, err
	}

	return contents, nil
}
