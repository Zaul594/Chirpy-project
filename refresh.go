package main

import (
	"net/http"

	"github.com/Zaul594/Chirspy-project/internal/auth"
)

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Token string `json:"token"`
	}

	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respWithErr(w, http.StatusBadRequest, "Couldn't find JWT")
		return
	}

	isRevoked, err := cfg.DB.IsTokenRevoked(refreshToken)
	if err != nil {
		respWithErr(w, http.StatusInternalServerError, "Couldn't check session")
	}
	if isRevoked {
		respWithErr(w, http.StatusUnauthorized, "refresh token is revoked")
		return
	}

	accessToken, err := auth.RefreshToken(refreshToken, cfg.jwtSecret)
	if err != nil {
		respWithErr(w, http.StatusUnauthorized, "Couldn't validate JWT")
		return
	}

	respJson(w, http.StatusOK, response{
		Token: accessToken,
	})
}

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respWithErr(w, http.StatusBadRequest, "Couldn't find JWT")
		return
	}

	err = cfg.DB.RevokeToken(refreshToken)
	if err != nil {
		respWithErr(w, http.StatusInternalServerError, "Couldn't find JWT")
		return
	}

	respJson(w, http.StatusOK, struct{}{})
}
