package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func respWithErr(w http.ResponseWriter, code int, msg string) {
	if code > 499 {
		log.Printf("Responding with 5XX error: %s", msg)
	}
	type errResp struct {
		Error string `json:"error"`
	}
	respJson(w, code, errResp{
		Error: msg,
	})

}

func respJson(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	dat, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marchalling Json: %s", err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(code)
	w.Write(dat)
}
