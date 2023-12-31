package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Zaul594/Chirspy-project/internal/auth"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}
	type response struct {
		Users
		Token         string `json:"token"`
		Refresh_Token string `json:"refresh_token"`
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

	accessToken, err := auth.MakeJWT(
		user.ID,
		cfg.jwtSecret,
		time.Hour,
		auth.TokenTypeAccess,
	)
	if err != nil {
		respWithErr(w, http.StatusInternalServerError, "Couldn't create access JWT")
		return
	}

	refreshToken, err := auth.MakeJWT(
		user.ID,
		cfg.jwtSecret,
		time.Hour*24*30*6,
		auth.TokenTypeRefresh,
	)
	if err != nil {
		respWithErr(w, http.StatusInternalServerError, "Couldn't create refresh JWT")
		return
	}

	respJson(w, http.StatusOK, response{
		Users: Users{
			ID:            user.ID,
			Email:         user.Email,
			Is_Chirpy_Red: user.Is_Chirpy_Red,
		},
		Token:         accessToken,
		Refresh_Token: refreshToken,
	})
}
