package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Zaul594/Chirspy-project/internal/auth"
)

func (cfg *apiConfig) userUpdateHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"Email"`
	}

	type response struct {
		Users
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respWithErr(w, http.StatusUnauthorized, "Couldn't find JWT")
		return
	}

	subject, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respWithErr(w, http.StatusUnauthorized, "Couldn't validate JWT")
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)

	if err != nil {
		respWithErr(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respWithErr(w, http.StatusInternalServerError, "Couldn't hash password")
		return
	}

	userIDInt, err := strconv.Atoi(subject)
	if err != nil {
		respWithErr(w, http.StatusInternalServerError, "Couldn't parse user ID")
		return
	}

	user, err := cfg.DB.UpdateUser(userIDInt, params.Email, hashedPassword)
	if err != nil {
		respWithErr(w, http.StatusInternalServerError, "Couldn't create user")
	}

	respJson(w, http.StatusOK, response{
		Users: Users{
			ID:    user.ID,
			Email: user.Email,
		},
	})
}
