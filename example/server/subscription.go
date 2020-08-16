package main

import (
	"io/ioutil"
	"log"
	"net/http"
)

// https://dev.fitbit.com/build/reference/web-api/subscriptions/index.html
func handleFitbitSubscriber(w http.ResponseWriter, r *http.Request) {
	if keys, ok := r.URL.Query()["verify"]; ok && len(keys[0]) > 0 {
		// drop verify messages not containing out subscriber token
		if keys[0] != subscriberCode {
			w.WriteHeader(http.StatusNotFound)
			return
		}
	}
	w.WriteHeader(http.StatusNoContent)

	log.Println("subscriber incoming request")
	log.Println(r.URL.Query())
	log.Println(r.Header)
	by, err := ioutil.ReadAll(r.Body)
	log.Println(err, string(by))
}
