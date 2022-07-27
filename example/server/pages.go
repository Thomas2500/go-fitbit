package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Thomas2500/go-fitbit"
	"github.com/google/uuid"
)

func handleFitbitLogin(w http.ResponseWriter, r *http.Request) {
	csrf := uuid.New().String()
	redirectURL := fca.LoginURL(csrf)
	http.Redirect(w, r, redirectURL, http.StatusTemporaryRedirect)
}
func httpFitbitGetProfile(w http.ResponseWriter, r *http.Request) {
	d, err := fca.Profile(0)
	if err != nil {
		fmt.Println(err)
		return
	}
	jb, _ := json.Marshal(d)
	fmt.Fprint(w, string(jb))
}
func httpFitbitGetDevices(w http.ResponseWriter, r *http.Request) {
	d, err := fca.Devices(0)
	if err != nil {
		fmt.Println(err)
		return
	}
	jb, _ := json.Marshal(d)
	fmt.Fprint(w, string(jb))
}
func httpFitbitGetFoodLog(w http.ResponseWriter, r *http.Request) {
	d, err := fca.FoodLogByDay(time.Now().Format("2006-01-02"))
	if err != nil {
		fmt.Println(err)
		return
	}
	jb, _ := json.Marshal(d)
	fmt.Fprint(w, string(jb))
}
func httpFitbitGetFoodGoal(w http.ResponseWriter, r *http.Request) {
	d, err := fca.FoodGoal()
	if err != nil {
		fmt.Println(err)
		return
	}
	jb, _ := json.Marshal(d)
	fmt.Fprint(w, string(jb))
}
func httpFitbitGetWaterLog(w http.ResponseWriter, r *http.Request) {
	d, err := fca.WaterLogByDay(time.Now().Format("2006-01-02"))
	if err != nil {
		fmt.Println(err)
		return
	}
	jb, _ := json.Marshal(d)
	fmt.Fprint(w, string(jb))
}
func httpFitbitGetWaterGoal(w http.ResponseWriter, r *http.Request) {
	d, err := fca.WaterGoal()
	if err != nil {
		fmt.Println(err)
		return
	}
	jb, _ := json.Marshal(d)
	fmt.Fprint(w, string(jb))
}
func httpFitbitGetHeatDay(w http.ResponseWriter, r *http.Request) {
	//d, err := fca.HeartLogByDay(time.Now().Add(time.Hour * 24 * -1).Format("2006-01-02"))
	d, err := fca.HeartLogByDateRange(time.Now().Add(time.Hour*24*-4).Format("2006-01-02"), time.Now().Add(time.Hour*24*-1).Format("2006-01-02"))
	if err != nil {
		fmt.Println(err)
		return
	}
	jb, _ := json.Marshal(d)
	fmt.Fprint(w, string(jb))
}
func httpFitbitGetHeatIntraday(w http.ResponseWriter, r *http.Request) {
	d, err := fca.HeartIntraday(time.Now().Format("2006-01-02"), "1sec", "", "")
	if err != nil {
		fmt.Println(err)
		return
	}
	jb, _ := json.Marshal(d)
	fmt.Fprint(w, string(jb))
}
func httpFitbitGetBodyWeight(w http.ResponseWriter, r *http.Request) {
	d, err := fca.BodyWeightLogByDay(time.Now().Format("2006-01-02"))
	if err != nil {
		fmt.Println(err)
		return
	}
	jb, _ := json.Marshal(d)
	fmt.Fprint(w, string(jb))
}
func httpFitbitGetSleepLog(w http.ResponseWriter, r *http.Request) {
	d, err := fca.SleepLogList(fitbit.LogListParameters{
		BeforeDate: time.Now().Format("2006-01-02"),
		Limit:      10,
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	jb, _ := json.Marshal(d)
	fmt.Fprint(w, string(jb))
}
func httpFitbitGetActivitiesDaySummary(w http.ResponseWriter, r *http.Request) {
	d, err := fca.ActivitiesDaySummary(time.Now().Format("2006-01-02"))
	if err != nil {
		fmt.Println(err)
		return
	}
	jb, _ := json.Marshal(d)
	fmt.Fprint(w, string(jb))
}
