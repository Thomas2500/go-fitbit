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
	"golang.org/x/oauth2"
)

// Temporary constant
const clientID = ""
const clientSecret = ""
const subscriberCode = ""

var fca fitbit.Session

func main() {
	fca = fitbit.New(clientID, clientSecret, "https://fitbit.bella.pm/callback", []string{
		fitbit.ScopeActivity,
		fitbit.ScopeSettings,
		fitbit.ScopeLocation,
		fitbit.ScopeSocial,
		fitbit.ScopeHeartrate,
		fitbit.ScopeProfile,
		fitbit.ScopeSleep,
		fitbit.ScopeNutrition,
		fitbit.ScopeWeight,
	})

	// Print OAuth2 access url
	fmt.Println(fca.LoginURL())

	// We already have a token which can be loaded
	jsonFile, err := ioutil.ReadFile("token.json") // just pass the file name
	if err != nil {
		fmt.Print(err)
		return
	}
	_ = jsonFile

	token := oauth2.Token{}
	err = json.Unmarshal(jsonFile, &token)
	if err != nil {
		fmt.Println("Error parsing token", err)
		return
	}
	fca.SetToken(&token)

	fmt.Println(fca.GetToken())

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
		writeToken()
	}()

	select {}
}

func writeToken() {
	file, _ := json.MarshalIndent(fca.GetToken(), "", " ")
	_ = ioutil.WriteFile("token.json", file, 0644)
	log.Println("Write current token to token.json")
}

func init() {
	go func() {
		time.Sleep(time.Second * 5)
		log.Println(fca.AddSubscription("", 1))
		time.Sleep(time.Second * 10)
		log.Println(fca.GetSubscriptions(""))
		time.Sleep(time.Second * 30)
		log.Println(fca.GetSubscriptions(""))
	}()

	// TODO: Periodic write/write on token change to save session state
	go func() {
		for {
			time.Sleep(time.Minute)
			writeToken()
		}
	}()
}

func handleFitbitCallback(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("code") == "" {
		if _, err := w.Write([]byte("No code given!")); err != nil {
			fmt.Println("can't write to client on fitbit callback", err)
		}
		return
	}

	// Get token
	if _, err := fca.HandleCallback(r.FormValue("code")); err != nil {
		log.Println("error handling callback", err)
	}

	token := fca.GetToken()
	js, err := json.Marshal(token)
	fmt.Println(string(js), err)

	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}
