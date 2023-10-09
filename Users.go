package main

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Zaul594/Chirspy-project/internal/auth"
	"github.com/Zaul594/Chirspy-project/internal/database"
)

type Users struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"-"`
}

func (cfg *apiConfig) createUserHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}
	type response struct {
		Users
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respWithErr(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respWithErr(w, http.StatusInternalServerError, "Couldn't hash password")
		return
	}

	user, err := cfg.DB.CreateUsers(params.Email, hashedPassword)
	if err != nil {
		if errors.Is(err, database.ErrAlreadyExists) {
			respWithErr(w, http.StatusConflict, "User already exists")
			return
		}
		respWithErr(w, http.StatusInternalServerError, "Couldn't create user")
		return
	}

	respJson(w, http.StatusCreated, response{
		Users: Users{
			ID:    user.ID,
			Email: user.Email,
		},
	})

}
