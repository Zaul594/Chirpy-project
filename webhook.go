package main

//useing webhook to polka

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Zaul594/Chirspy-project/internal/auth"
	"github.com/Zaul594/Chirspy-project/internal/database"
)

func (cfg *apiConfig) upgradeHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Event string `json:"event"`
		Data  struct {
			UserID int `json:"user_id"`
		}
	}

	apiKey, err := auth.GatAPIKey(r.Header)
	if err != nil {
		respWithErr(w, http.StatusUnauthorized, "Couldn't find api key")
		return
	}
	if apiKey != cfg.api_key {
		respWithErr(w, http.StatusUnauthorized, "API key is invalid")
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respWithErr(w, http.StatusInternalServerError, "Couldn't decode	parameters")
		return
	}

	if params.Event != "user.upgraded" {
		respJson(w, http.StatusOK, struct{}{})
		return
	}

	_, err = cfg.DB.UpgradeChirpyRed(params.Data.UserID)
	if err != nil {
		if errors.Is(err, database.ErrNotExist) {
			respWithErr(w, http.StatusNotFound, "Couldn't find user")
			return
		}
		respWithErr(w, http.StatusInternalServerError, "Couldn't update user")
		return
	}

	respJson(w, http.StatusOK, struct{}{})
}
