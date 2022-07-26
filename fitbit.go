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

const (
	fitbitAuthURL  = "https://www.fitbit.com/oauth2/authorize"
	fitbitTokenURL = "https://api.fitbit.com/oauth2/token"
)

// Scope describes an oauth2 scope for Fitbit
type Scope = string

// Different types of useable scopes
const (
	ScopeActivity  Scope = "activity"
	ScopeSettings  Scope = "settings"
	ScopeLocation  Scope = "location"
	ScopeSocial    Scope = "social"
	ScopeHeartrate Scope = "heartrate"
	ScopeProfile   Scope = "profile"
	ScopeSleep     Scope = "sleep"
	ScopeNutrition Scope = "nutrition"
	ScopeWeight    Scope = "weight"
)

// ApplicationType represents the type of registered application.
type ApplicationType int64

const (
	PersonalApplication ApplicationType = iota // PersonalApplication represents personal type application
	ClientApplication                          // ClientApplication represents client type application
	ServerApplication                          // ServerApplication represents server type application
)

// Session is the main object with user data
type Session struct {
	// HookTokenChange is a function that is called when the refresh_token changes
	TokenChange func(token *oauth2.Token)

	ratelimit Ratelimit

	// config is the configuration for this session
	config Config

	// oauth config
	oAuthConfig *oauth2.Config

	// current token
	token *oauth2.Token

	// httpClient is the authenticated http client used for this oAuth session
	httpClient *http.Client

	mutex sync.RWMutex
}

// Config describes the configuration of a fitbit API configuration
type Config struct {
	ClientID        string
	ClientSecret    string
	RedirectURL     string
	Scopes          []Scope
	ApplicationType ApplicationType
	TokenPath       string
}

// Ratelimit includes the rate limit information provided on every request
type Ratelimit struct {
	RateLimitAvailable int
	RateLimitUsed      int
	RateLimitReset     time.Time
}

// New creates a new fitbit oauth session
func New(config Config) *Session {
	// Create new oauth configuation
	oAuthConfig := &oauth2.Config{
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecret,
		Scopes:       config.Scopes,
		Endpoint: oauth2.Endpoint{
			AuthURL:  fitbitAuthURL,
			TokenURL: fitbitTokenURL,
		},
	}

	return &Session{
		config:      config,
		oAuthConfig: oAuthConfig,
	}
}

// LoginURL returns an OAuth login url to obtain an access token
func (m *Session) LoginURL(csrf string) string {
	return m.oAuthConfig.AuthCodeURL(csrf, oauth2.AccessTypeOffline)
}

type cacherTransport struct {
	Base    *oauth2.Transport
	Session *Session
}

func (c *cacherTransport) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	if _, err := c.Base.Source.Token(); err != nil {
		return nil, errExpiredToken
	}
	resp, err = c.Base.RoundTrip(req)
	if err != nil {
		return nil, err
	}
	newTok, err := c.Base.Source.Token()
	if err != nil {
		// While we’re unable to obtain a new token, the request was still
		// successful, so let’s gracefully handle this error by not caching a
		// new token. In either case, the user will need to re-authenticate.
		return resp, nil
	}
	if c.Session.token == nil ||
		c.Session.token.AccessToken != newTok.AccessToken ||
		c.Session.token.RefreshToken != newTok.RefreshToken {
		// Save new token if it differs from the old one
		c.Session.mutex.Lock()
		c.Session.token = newTok
		c.Session.mutex.Unlock()

		// Call token change hook
		if c.Session.TokenChange != nil {
			go c.Session.TokenChange(newTok)
		}
	}
	return resp, nil
}

// Like oauth2.Config.Client(), but using cacherTransport to persist tokens.
func (m *Session) newHTTPClient() *http.Client {
	return &http.Client{
		Transport: &cacherTransport{
			Session: m,
			Base: &oauth2.Transport{
				Source: m.oAuthConfig.TokenSource(context.Background(), m.token),
			},
		},
	}
}

// SetToken sets the token to use for the session
func (m *Session) SetToken(token *oauth2.Token) {
	m.mutex.Lock()
	m.token = token
	m.httpClient = m.newHTTPClient()
	m.mutex.Unlock()
}

// SaveToken triggers the TokenChange function to manually save the token
func (m *Session) SaveToken() error {
	if m.TokenChange == nil {
		return errTokenChangeNotDefined
	}
	m.TokenChange(m.token)
	return nil
}

func (m *Session) GetRatelimit() Ratelimit {
	return m.ratelimit
}

// makeRequest creates a new request to a given url using given
// OAuth token of an user
func (m *Session) makeRequest(url string) ([]byte, error) {
	// if httpClient is nil build a new one
	if m.httpClient == nil {
		m.httpClient = m.newHTTPClient()
	}

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

	// Parse rate limit headers
	m.parseRatelimit(&response.Header)

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
	// if httpClient is nil build a new one
	if m.httpClient == nil {
		m.httpClient = m.newHTTPClient()
	}

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

	// Parse rate limit headers
	m.parseRatelimit(&response.Header)

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
	// if httpClient is nil build a new one
	if m.httpClient == nil {
		m.httpClient = m.newHTTPClient()
	}

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

	// Parse rate limit headers
	m.parseRatelimit(&response.Header)

	// Read all data from request
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return contents, nil
}

// parseRatelimit parses the rate limit headers of fitbit API
func (m *Session) parseRatelimit(header *http.Header) {
	// Get rate limit data of request
	rateLimitData := header.Get("fitbit-rate-limit-remaining")
	if rateLimitData != "" {
		m.ratelimit.RateLimitUsed, _ = strconv.Atoi(rateLimitData)
	}
	rateLimitData = header.Get("fitbit-rate-limit-limit")
	if rateLimitData != "" {
		m.ratelimit.RateLimitAvailable, _ = strconv.Atoi(rateLimitData)
	}
	rateLimitData = header.Get("fitbit-rate-limit-reset")
	if rateLimitData != "" {
		remSec, _ := strconv.Atoi(rateLimitData)
		m.ratelimit.RateLimitReset = time.Now().Add(time.Second * time.Duration(remSec))
	}
}

// Exchange uses an authorization code to retrieve an access token and refresh token.
// sets them in the current session using SetToken and rebuilds the httpClient
func (m *Session) Exchange(code string) (*oauth2.Token, error) {
	token, err := m.oAuthConfig.Exchange(context.Background(), code)
	if err != nil {
		return nil, err
	}
	m.SetToken(token)
	return token, nil
}
