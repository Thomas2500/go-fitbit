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

// Temporary constant
const clientID = "22D5RX"
const clientSecret = "982ad62d23b286f308c5b5e7e4eebffd"
const subscriberCode = "e12b47c93c2c7f2eb562fb05005884e57a118d16100fa00969be49382393d61d"

var fca *fitbit.Session

func main() {
	fca = fitbit.New(fitbit.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  fmt.Sprintf("https://%s/callback", "fitbit.bella.pm"),
		Scopes: []string{
			fitbit.ScopeActivity,
			fitbit.ScopeSettings,
			fitbit.ScopeLocation,
			fitbit.ScopeSocial,
			fitbit.ScopeHeartrate,
			fitbit.ScopeProfile,
			fitbit.ScopeSleep,
			fitbit.ScopeNutrition,
			fitbit.ScopeWeight,
		},
	})

	// Define fitbit hook function to save token changes to file
	fca.TokenChange = fitbitSaveToken

	// Print OAuth2 access url
	csrf := uuid.New().String()
	fmt.Println(fca.LoginURL(csrf))

	// We already have a token which can be loaded
	jsonFile, err := ioutil.ReadFile("token.json") // just pass the file name
	if err != nil {
		fmt.Print(err)
		return
	}

	token := oauth2.Token{}
	err = json.Unmarshal(jsonFile, &token)
	if err != nil {
		fmt.Println("Error parsing token", err)
		return
	}
	fca.SetToken(&token)

	http.HandleFunc("/", handleMain)
	//http.HandleFunc("/style.css", handleStyleFile)

	// Login and callback page
	//http.HandleFunc("/login", handleFitbitLogin)
	http.HandleFunc("/callback", handleFitbitCallback)
	http.HandleFunc("/subscriber", handleFitbitSubscriber)

	// Information about this API backend
	//http.HandleFunc("/quota", handleHTTPQuota)

	// API pages with data
	http.HandleFunc("/profile", httpFitbitGetProfile)

	http.HandleFunc("/badges", httpFitbitGetBadges)
	http.HandleFunc("/devices", httpFitbitGetDevices)
	http.HandleFunc("/food/goal", httpFitbitGetFoodGoal)
	http.HandleFunc("/food/log", httpFitbitGetFoodLog)
	http.HandleFunc("/water/goal", httpFitbitGetWaterGoal)
	http.HandleFunc("/water/log", httpFitbitGetWaterLog)
	http.HandleFunc("/heart/day", httpFitbitGetHeatDay)
	http.HandleFunc("/heart/intraday", httpFitbitGetHeatIntraday)
	http.HandleFunc("/body/weight", httpFitbitGetBodyWeight)
	http.HandleFunc("/body/weightrange", BodyWeightLogByDateRange)
	http.HandleFunc("/body/fat", httpFitbitGetBodyFat)
	http.HandleFunc("/activities/types", httpFitbitGetActivitiesTypes)
	http.HandleFunc("/activities/interday", httpFitbitGetActivitiesInterday)
	http.HandleFunc("/activities/frequent", httpFitbitGetActivitiesFrequent)
	http.HandleFunc("/activities/recent", httpFitbitGetActivitiesRecent)
	http.HandleFunc("/activities/day", httpFitbitGetActivitiesDaySummary)
	http.HandleFunc("/sleep/today", httpFitbitGetSleepToday)
	http.HandleFunc("/sleep/log", httpFitbitGetSleepLog)
	http.HandleFunc("/sleep/goal", httpFitbitGetSleepGoal)

	// Start listener
	//http.ListenAndServe("127.0.0.1:48558", nil)
	if err := http.ListenAndServe("0.0.0.0:48558", nil); err != nil {
		log.Fatal("error starting http server", err)
	}

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

	select {}
}

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

func init() {
	go func() {
		time.Sleep(time.Second * 5)
		log.Println(fca.GetSubscriptions(""))
		time.Sleep(time.Second * 10)
		log.Println(fca.AddSubscription("", 1))
		time.Sleep(time.Second * 30)
		log.Println(fca.GetSubscriptions(""))
		time.Sleep(time.Second * 45)
		log.Println(fca.Introspect())
	}()
}

func handleFitbitCallback(w http.ResponseWriter, r *http.Request) {
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
	fitbitSaveToken(token)
	log.Println("FITBIT: token saved")
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}
