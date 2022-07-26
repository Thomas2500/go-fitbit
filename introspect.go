package fitbit

import (
	"encoding/json"
)

// IntrospectResponse contains the response of the introspect request
type IntrospectResponse struct {
	Active    bool   `json:"active"`
	Scope     string `json:"scope,omitempty"`
	ClientID  string `json:"client_id,omitempty"`
	UserID    string `json:"user_id,omitempty"`
	TokenType string `json:"token_type,omitempty"`
	Exp       int64  `json:"exp,omitempty"`
	Iat       int64  `json:"iat,omitempty"`
}

// Introspect checks if the currently used oauth token is still valid
func (m *Session) Introspect() (IntrospectResponse, error) {
	// Build request
	postRequestBody := map[string]string{
		"token": m.token.AccessToken,
	}

	contents, err := m.makePOSTRequest("https://api.fitbit.com/1.1/oauth2/introspect", postRequestBody)
	if err != nil {
		return IntrospectResponse{}, err
	}

	intro := IntrospectResponse{}
	if err := json.Unmarshal(contents, &intro); err != nil {
		return IntrospectResponse{}, err
	}

	// As of rfc7662 exp and iat should be seconds since epoch
	// Fitbit doesn't respect this and uses milliseconds
	// correct it to the correct unit
	intro.Exp = intro.Exp / 1000
	intro.Iat = intro.Iat / 1000

	return intro, nil
}
