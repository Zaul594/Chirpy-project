package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"
)

func getCleanBody(body string, badWords map[string]struct{}) string {
	words := strings.Split(body, " ")
	for i, word := range words {
		lcWord := strings.ToLower(word)
		if _, ok := badWords[lcWord]; ok {
			words[i] = "****"
		}
	}
	cleaned := strings.Join(words, " ")
	return cleaned
}

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

func postValidtor(body string) (string, error) {

	const maxChirpLength = 140
	if len(body) > maxChirpLength {
		return "", errors.New("Chirp is too long")
	}

	badWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}
	cleaned := getCleanBody(body, badWords)
	return cleaned, nil
}
