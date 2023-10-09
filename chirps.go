package main

import (
	"encoding/json"
	"net/http"
	"sort"
	"strconv"

	"github.com/Zaul594/Chirspy-project/internal/auth"
	"github.com/go-chi/chi/v5"
)

type Chirp struct {
	ID        int    `json:"id"`
	Author_ID int    `json:"author_id"`
	Body      string `json:"body"`
}

func (cfg *apiConfig) chirpCreateHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
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

	userID, err := strconv.Atoi(subject)
	if err != nil {
		respWithErr(w, http.StatusBadRequest, "Couldn't parse user ID")
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respWithErr(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	cleaned, err := postValidtor(params.Body)
	if err != nil {
		respWithErr(w, http.StatusBadRequest, err.Error())
		return
	}

	chirp, err := cfg.DB.CreateChirps(cleaned, userID)
	if err != nil {
		respWithErr(w, http.StatusInternalServerError, "Couldn't create chirp")
		return
	}

	respJson(w, http.StatusCreated, Chirp{
		ID:        chirp.ID,
		Author_ID: chirp.Author_ID,
		Body:      chirp.Body,
	})
}

func (cfg *apiConfig) chirpGetHandaler(w http.ResponseWriter, r *http.Request) {
	dbChirps, err := cfg.DB.GetChirps()
	if err != nil {
		respWithErr(w, http.StatusInternalServerError, "Could Couldn't retrieve chirps")
		return
	}

	chirps := []Chirp{}
	for _, dbChirp := range dbChirps {
		chirps = append(chirps, Chirp{
			ID:        dbChirp.ID,
			Author_ID: dbChirp.Author_ID,
			Body:      dbChirp.Body,
		})
	}

	sort.Slice(chirps, func(i, j int) bool {
		return chirps[i].ID < chirps[j].ID
	})

	respJson(w, 200, chirps)
}

func (cfg *apiConfig) GetOneChirpHandler(w http.ResponseWriter, r *http.Request) {
	chirpIDString := chi.URLParam(r, "chirpID")
	chirpID, err := strconv.Atoi(chirpIDString)
	if err != nil {
		respWithErr(w, http.StatusBadRequest, "Invalid chirp ID")
		return
	}

	dbChirp, err := cfg.DB.GetChirp(chirpID)
	if err != nil {
		respWithErr(w, http.StatusNotFound, "Couldn't get chirp")
		return
	}

	respJson(w, 200, Chirp{
		ID:        dbChirp.ID,
		Author_ID: dbChirp.Author_ID,
		Body:      dbChirp.Body,
	})
}

func (cfg *apiConfig) deleteChirpyHandler(w http.ResponseWriter, r *http.Request) {
	chirpIDString := chi.URLParam(r, "chirpID")
	chirpID, err := strconv.Atoi(chirpIDString)
	if err != nil {
		respWithErr(w, http.StatusBadRequest, "Invalid chirp ID")
		return
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

	userID, err := strconv.Atoi(subject)
	if err != nil {
		respWithErr(w, http.StatusBadRequest, "Couldn't parse user ID")
		return
	}

	dbChirp, err := cfg.DB.GetChirp(chirpID)
	if err != nil {
		respWithErr(w, http.StatusBadRequest, "Could not get chirpy")
		return
	}
	if dbChirp.Author_ID != userID {
		respWithErr(w, 403, "You can't delete this chirp")
	}

	err = cfg.DB.DeleteChirpy(chirpID)
	if err != nil {
		respWithErr(w, http.StatusInternalServerError, "Couldn't delete chirp")
		return
	}

	respJson(w, http.StatusOK, struct{}{})
}
