package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Thomas2500/go-fitbit"
	"github.com/google/uuid"
	"golang.org/x/oauth2"
)

// Temporary constant - replace with your own clientID, secret and subscriber code
const clientID = "22D5RX"
const clientSecret = "982ad62d23b286f308c5b5e7e4eebffd"
const subscriberCode = "e12b47c93c2c7f2eb562fb05005884e57a118d16100fa00969be49382393d61d"

// fitbit Session established for this demo
var fca *fitbit.Session

func main() {
	// Create a new fitbit session
	fca = fitbit.New(fitbit.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  fmt.Sprintf("https://%s/callback", "fitbit.bella.pm"),
		Scopes: []string{
			fitbit.ScopeActivity,
			fitbit.ScopeBreathingRate,
			fitbit.ScopeHeartrate,
			fitbit.ScopeLocation,
			fitbit.ScopeNutrition,
			fitbit.ScopeProfile,
			fitbit.ScopeSettings,
			fitbit.ScopeSleep,
			fitbit.ScopeSocial,
			fitbit.ScopeSpO2,
			fitbit.ScopeTemperature,
			fitbit.ScopeWeight,
		},
	})

	// Define fitbit hook function to save token changes to file
	fca.TokenChange = fitbitSaveToken

	// Print OAuth2 access url to grant permissions to use API requests on behalf of the user
	csrf := uuid.New().String()
	fmt.Println(fca.LoginURL(csrf))

	// We already have a token which can be loaded
	jsonFile, err := ioutil.ReadFile("token.json") // just pass the file name
	if err != nil {
		fmt.Print(err)
		return
	}

	// load file contents into token struct to initialize the session using SetTocken
	token := oauth2.Token{}
	err = json.Unmarshal(jsonFile, &token)
	if err != nil {
		fmt.Println("Error parsing token", err)
		return
	}
	fca.SetToken(&token)

	// Execute some functions async to test some API functionality - may fail if not authorized using token
	go func() {
		time.Sleep(time.Second * 5)
		log.Println(fca.GetSubscriptions(""))
		time.Sleep(time.Second * 10)
		log.Println(fca.AddSubscription("", 1))
		time.Sleep(time.Second * 10)
		log.Println(fca.GetSubscriptions(""))
		time.Sleep(time.Second * 15)
		log.Println("Core")
		log.Println(fca.TemperatureCoreByDay("today"))
		log.Println("Skin")
		log.Println(fca.TemperatureSkinByDay("today"))
	}()

	// Create a webserver to allow simple naviation and exploration of the fitbit API
	http.HandleFunc("/", handleMain)
	http.HandleFunc("/style.css", handleStyleFile)

	// Login and callback page
	http.HandleFunc("/callback", handleFitbitCallback)
	http.HandleFunc("/subscriber", handleFitbitSubscriber)

	// API pages with data
	http.HandleFunc("/login", handleFitbitLogin)
	http.HandleFunc("/profile", httpFitbitGetProfile)
	http.HandleFunc("/devices", httpFitbitGetDevices)
	http.HandleFunc("/food/log", httpFitbitGetFoodLog)
	http.HandleFunc("/food/goal", httpFitbitGetFoodGoal)
	http.HandleFunc("/water/log", httpFitbitGetWaterLog)
	http.HandleFunc("/water/goal", httpFitbitGetWaterGoal)
	http.HandleFunc("/heart/day", httpFitbitGetHeatDay)
	http.HandleFunc("/heart/intraday", httpFitbitGetHeatIntraday)
	http.HandleFunc("/body/weight", httpFitbitGetBodyWeight)
	http.HandleFunc("/sleep/log", httpFitbitGetSleepLog)
	http.HandleFunc("/activities/summary", httpFitbitGetActivitiesDaySummary)

	// Start listener
	//http.ListenAndServe("127.0.0.1:48558", nil)
	if err := http.ListenAndServe("0.0.0.0:48558", nil); err != nil {
		log.Fatal("error starting http server", err)
	}

	// Listen for stop signals to force a token write in exit
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	go func() {
		s := <-sigc
		_ = s
		err := fca.SaveToken()
		if err != nil {
			log.Println("error saving token", err)
		}
	}()

	// endless wait to provide webserver
	select {}
}

// fitbitSaveToken as hook funktion automatically called when the Token is changed
func fitbitSaveToken(token *oauth2.Token) {
	jsonFile, err := os.Create("token.json")
	if err != nil {
		log.Println("error creating token file", err)
	}
	defer jsonFile.Close()
	err = json.NewEncoder(jsonFile).Encode(token)
	if err != nil {
		log.Println("error encoding token", err)
	}
	log.Println("FITBIT: token saved to file")
}

// handleFitbitCallback handles the oAuth2 callback visited by the user after granting permissions
func handleFitbitCallback(w http.ResponseWriter, r *http.Request) {
	// check if the request is a callback from fitbit and a code is given
	// Advice: CSRF (state) token can be checked here if the request is legitimate
	if r.FormValue("code") == "" {
		if _, err := w.Write([]byte("No code given!")); err != nil {
			fmt.Println("can't write to client on fitbit callback", err)
		}
		return
	}

	// Exchange the code for an access token
	token, err := fca.Exchange(r.FormValue("code"))
	if err != nil {
		log.Println("FITBIT: error persisting initial token", err)
		return
	}

	// Prettify token for logging
	js, err := json.Marshal(token)
	log.Printf("FITBIT: token: %s - E: %s", string(js), err.Error())

	// save token to file for recurring use
	fitbitSaveToken(token)
	log.Println("FITBIT: token saved")

	// redirect to main page
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}
