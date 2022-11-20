package fitbit

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"golang.org/x/oauth2"
)

// Define used Fitbit endpoints
const (
	fitbitAuthURL  = "https://www.fitbit.com/oauth2/authorize" //nolint:gosec
	fitbitTokenURL = "https://api.fitbit.com/oauth2/token"     //nolint:gosec
)

// Scope describes an oauth2 scope for Fitbit
type Scope = string

// Different types of useable scopes
const (
	ScopeActivity      Scope = "activity"
	ScopeCardioFitness Scope = "cardio_fitness"
	ScopeBreathingRate Scope = "respiratory_rate"
	ScopeHeartrate     Scope = "heartrate"
	ScopeLocation      Scope = "location"
	ScopeNutrition     Scope = "nutrition"
	ScopeProfile       Scope = "profile"
	ScopeSettings      Scope = "settings"
	ScopeSleep         Scope = "sleep"
	ScopeSocial        Scope = "social"
	ScopeSpO2          Scope = "oxygen_saturation"
	ScopeTemperature   Scope = "temperature"
	ScopeWeight        Scope = "weight"
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

	// locale is the locale used for this session
	locale string

	mutex sync.RWMutex
}

// Config describes the configuration of a fitbit API configuration
type Config struct {
	ClientID     string  // ClientID is the client id (OAuth 2.0 Client ID) of the application (required)
	ClientSecret string  // ClientSecret is the client secret (Client Secret) of the application (required)
	RedirectURL  string  // RedirectURL is the redirect url of the application (required)
	Scopes       []Scope // Scopes is a list of scopes to request
	Locale       string  // en_AU, fr_FR, de_DE, ja_JP, en_NZ, es_ES, en_GB, en_US (default: de_DE)
}

// Ratelimit includes the rate limit information provided on every request
type Ratelimit struct {
	RateLimitAvailable int       // RateLimitAvailable is the number of requests available for the current rate limit window
	RateLimitUsed      int       // RateLimitUsed is the number of requests used for the current rate limit window
	RateLimitReset     time.Time // RateLimitReset is the time when the rate limit window resets
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

	// determine locale, if not used set to de_DE (this was the previous default)
	// list of locales: https://dev.fitbit.com/build/reference/web-api/developer-guide/application-design/#Localization
	locale := config.Locale
	switch locale {
	case "en_AU", "fr_FR", "de_DE", "ja_JP", "en_NZ", "es_ES", "en_GB", "en_US":
	default:
		locale = "de_DE"
	}

	// return session
	return &Session{
		config:      config,
		oAuthConfig: oAuthConfig,
		locale:      locale,
	}
}

// LoginURL returns an OAuth login url to obtain an access token
func (m *Session) LoginURL(csrf string) string {
	return m.oAuthConfig.AuthCodeURL(csrf, oauth2.AccessTypeOffline)
}

// cacherTransport is a transport which intercepts RoundTrip to check if the token changed on HTTP requests
type cacherTransport struct {
	Base    *oauth2.Transport
	Session *Session
}

// RoundTrip overrides the http.Client RoundTrip method to check if the token changed
func (c *cacherTransport) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	// get the current used token and determine if it is already expired
	if _, err := c.Base.Source.Token(); err != nil {
		return nil, errExpiredToken
	}
	// authorize based on oauth2.Transport.RoundTrip
	resp, err = c.Base.RoundTrip(req)
	if err != nil {
		return nil, err
	}
	// get the current used token to compare it with the previous one
	newTok, err := c.Base.Source.Token()
	if err != nil {
		// a error appeared which means the token is invalid
		// return the response, the user has to decide how to proceed here
		// probably a new authentication is required which is a manual task for the user
		return resp, nil
	}

	// check if the token changed from old to new one
	// if it changed update token copy and trigger TokenChange user function to allow persisting the token
	if c.Session.token == nil ||
		c.Session.token.AccessToken != newTok.AccessToken ||
		c.Session.token.RefreshToken != newTok.RefreshToken {
		// Save new token if it differs from the old one
		c.Session.mutex.Lock()
		c.Session.token = newTok
		c.Session.mutex.Unlock()

		// Call token change hook if defined by user
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

// GetRatelimit returns the current ratelimit information obtained by the last API request
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
	req.Header.Set("Accept-Language", m.locale)
	req.Header.Set("Accept-Locale", m.locale)

	// Fire request
	response, err := m.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	// Parse rate limit headers
	m.parseRatelimit(&response.Header)

	// Read all data from request
	contents, err := io.ReadAll(response.Body)
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
	req.Header.Set("Accept-Language", m.locale)
	req.Header.Set("Accept-Locale", m.locale)
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
	contents, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return contents, nil
}

// makeDELETERequest creates a new request to a given url using given
// OAuth token of an user
//
//nolint:unparam
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
	req.Header.Set("Accept-Language", m.locale)
	req.Header.Set("Accept-Locale", m.locale)

	// Fire request
	response, err := m.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	// Parse rate limit headers
	m.parseRatelimit(&response.Header)

	// Read all data from request
	contents, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return contents, nil
}

// parseRatelimit parses the rate limit headers of fitbit API
func (m *Session) parseRatelimit(header *http.Header) {
	// Get rate limit data of request
	// fist header returns the remaining API requests until reset time is reached
	rateLimitData := header.Get("fitbit-rate-limit-remaining")
	if rateLimitData != "" {
		m.ratelimit.RateLimitUsed, _ = strconv.Atoi(rateLimitData)
	}
	// second header returns the number of API requests allowed in genral
	rateLimitData = header.Get("fitbit-rate-limit-limit")
	if rateLimitData != "" {
		m.ratelimit.RateLimitAvailable, _ = strconv.Atoi(rateLimitData)
	}

	// rate limit reset returns when the rate limit will be reset
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
