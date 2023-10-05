package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Zaul594/Chirspy-project/internal/auth"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password         string `json:"password"`
		Email            string `json:"email"`
		ExpiresInSeconds int    `json:"expires_in_seconds"`
	}
	type response struct {
		Users
		Token string `json:"token"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respWithErr(w, http.StatusInternalServerError, "Could't decode parameters")
		return
	}

	user, err := cfg.DB.GetUserByEmail(params.Email)
	if err != nil {
		respWithErr(w, http.StatusInternalServerError, "Couldn't get users")
		return
	}

	err = auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil {
		respWithErr(w, http.StatusUnauthorized, "Invalid password")
		return
	}

	defaultExpiration := 60 * 60 * 24
	if params.ExpiresInSeconds == 0 {
		params.ExpiresInSeconds = defaultExpiration
	} else if params.ExpiresInSeconds > defaultExpiration {
		params.ExpiresInSeconds = defaultExpiration
	}

	token, err := auth.MakeJWT(user.ID, cfg.jwtSecret, time.Duration(params.ExpiresInSeconds)*time.Second)
	if err != nil {
		respWithErr(w, http.StatusInternalServerError, "Couldn't create JWT")
		return
	}

	respJson(w, http.StatusOK, response{
		Users: Users{
			ID:    user.ID,
			Email: user.Email,
		},
		Token: token,
	})
}
