package fitbit

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"golang.org/x/oauth2"
)

// API REFERENCE AT: https://dev.fitbit.com/build/reference/web-api/heart-rate/
// Application Info: https://dev.fitbit.com/apps/details/22D5RX

const fitbitAuthURL = "https://www.fitbit.com/oauth2/authorize"
const fitbitTokenURL = "https://api.fitbit.com/oauth2/token"

// Different types of useable scopes
const (
	ScopeActivity  = "activity"
	ScopeSettings  = "settings"
	ScopeLocation  = "location"
	ScopeSocial    = "social"
	ScopeHeartrate = "heartrate"
	ScopeProfile   = "profile"
	ScopeSleep     = "sleep"
	ScopeNutrition = "nutrition"
	ScopeWeight    = "weight"
)

// Session is the main object with user data
type Session struct {
	OAuthConfg *oauth2.Config
	Token      *oauth2.Token
	Ratelimit  Ratelimit
	httpClient *http.Client
	mutex      sync.RWMutex
}

// Ratelimit includes the rate limit information provided on every request
type Ratelimit struct {
	RateLimitAvailable int
	RateLimitUsed      int
	RateLimitReset     time.Time
}

// New creates a new fitbit oauth session
func New(clientID string, clientSecret string, redirectURL string, scopes []string) Session {
	return Session{
		OAuthConfg: &oauth2.Config{
			RedirectURL:  redirectURL,
			ClientID:     clientID,
			ClientSecret: clientSecret,
			Scopes:       scopes,
			Endpoint: oauth2.Endpoint{
				AuthURL:  fitbitAuthURL,
				TokenURL: fitbitTokenURL,
			},
		},
	}
}

// LoginURL returns an OAuth login url to obtain an access token
func (m *Session) LoginURL() string {
	return m.OAuthConfg.AuthCodeURL(m.OAuthConfg.ClientID)
}

// SetToken sets a oauth2 token and tries to renew it if expired
func (m *Session) SetToken(token *oauth2.Token) (*oauth2.Token, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	ctx := context.Background()
	conf := &oauth2.Config{}
	if token.Expiry.Before(time.Now()) {
		src := conf.TokenSource(ctx, token)
		newToken, err := src.Token()
		if err != nil {
			return nil, err
		}
		if newToken.AccessToken != token.AccessToken {
			token = newToken
		}
	}
	m.httpClient = m.OAuthConfg.Client(ctx, token)
	m.Token = token
	return token, nil
}

// GetToken returns an oauth2.Token object
func (m *Session) GetToken() *oauth2.Token {
	return m.Token
}

// GetRatelimit returns the current rate limit information
// Only available after a request to the API endpoint
func (m *Session) GetRatelimit() Ratelimit {
	return m.Ratelimit
}

// HandleCallback performs an oauth2 exchange to receive the
// needed token. Has to be called when code was received
func (m *Session) HandleCallback(code string) (*oauth2.Token, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	token, err := m.OAuthConfg.Exchange(context.TODO(), code)
	m.Token = token
	if err != nil {
		return nil, err
	}
	m.httpClient = m.OAuthConfg.Client(context.Background(), token)

	return token, nil
}

// makeRequest creates a new request to a given url using given
// OAuth token of an user
func (m *Session) makeRequest(url string) ([]byte, error) {
	// Build request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// Set custom header
	req.Header.Set("User-Agent", "go-fitbit")
	req.Header.Set("Accept-Language", "de_DE")
	req.Header.Set("Accept-Locale", "de_DE")

	// Fire request
	response, err := m.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	// Get rate limit data of request
	rateLimitData := response.Header.Get("fitbit-rate-limit-remaining")
	if rateLimitData != "" {
		m.Ratelimit.RateLimitUsed, _ = strconv.Atoi(rateLimitData)
	}
	rateLimitData = response.Header.Get("fitbit-rate-limit-limit")
	if rateLimitData != "" {
		m.Ratelimit.RateLimitAvailable, _ = strconv.Atoi(rateLimitData)
	}
	rateLimitData = response.Header.Get("fitbit-rate-limit-reset")
	if rateLimitData != "" {
		remSec, _ := strconv.Atoi(rateLimitData)
		m.Ratelimit.RateLimitReset = time.Now().Add(time.Second * time.Duration(remSec))
	}

	// Read all data from request
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return contents, nil
}

// makePOSTRequest creates a new request to a given url using given
// OAuth token of an user
func (m *Session) makePOSTRequest(targetURL string, param map[string]string) ([]byte, error) {
	// Build post params
	form := url.Values{}
	for name, value := range param {
		form.Add(name, value)
	}

	// Build request
	req, err := http.NewRequest("POST", targetURL, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}

	// Set custom header
	req.Header.Set("User-Agent", "go-fitbit")
	req.Header.Set("Accept-Language", "de_DE")
	req.Header.Set("Accept-Locale", "de_DE")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Fire request
	response, err := m.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	// Get rate limit data of request
	rateLimitData := response.Header.Get("fitbit-rate-limit-remaining")
	if rateLimitData != "" {
		m.Ratelimit.RateLimitUsed, _ = strconv.Atoi(rateLimitData)
	}
	rateLimitData = response.Header.Get("fitbit-rate-limit-limit")
	if rateLimitData != "" {
		m.Ratelimit.RateLimitAvailable, _ = strconv.Atoi(rateLimitData)
	}
	rateLimitData = response.Header.Get("fitbit-rate-limit-reset")
	if rateLimitData != "" {
		remSec, _ := strconv.Atoi(rateLimitData)
		m.Ratelimit.RateLimitReset = time.Now().Add(time.Second * time.Duration(remSec))
	}

	// Read all data from request
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return contents, nil
}

// makeDELETERequest creates a new request to a given url using given
// OAuth token of an user
func (m *Session) makeDELETERequest(url string) ([]byte, error) {
	// Build request
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return nil, err
	}

	// Set custom header
	req.Header.Set("User-Agent", "go-fitbit")
	req.Header.Set("Accept-Language", "de_DE")
	req.Header.Set("Accept-Locale", "de_DE")

	// Fire request
	response, err := m.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	// Get rate limit data of request
	rateLimitData := response.Header.Get("fitbit-rate-limit-remaining")
	if rateLimitData != "" {
		m.Ratelimit.RateLimitUsed, _ = strconv.Atoi(rateLimitData)
	}
	rateLimitData = response.Header.Get("fitbit-rate-limit-limit")
	if rateLimitData != "" {
		m.Ratelimit.RateLimitAvailable, _ = strconv.Atoi(rateLimitData)
	}
	rateLimitData = response.Header.Get("fitbit-rate-limit-reset")
	if rateLimitData != "" {
		remSec, _ := strconv.Atoi(rateLimitData)
		m.Ratelimit.RateLimitReset = time.Now().Add(time.Second * time.Duration(remSec))
	}

	// Read all data from request
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return contents, nil
}
