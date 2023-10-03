package main

import (
	"encoding/json"
	"net/http"
)

type Users struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
}

func (cfg *apiConfig) createUserHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respWithErr(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	user, err := cfg.DB.CreateUsers(params.Email)
	if err != nil {
		respWithErr(w, http.StatusInternalServerError, "Couldn't create user")
		return
	}

	respJson(w, http.StatusCreated, Users{
		ID:    user.ID,
		Email: user.Email,
	})

}

//func (cfg *apiConfig) getAllUsersHandler(w http.ResponseWriter, r *http.Request) {

//}

//func (cfg *apiConfig) getOneUser(){}
