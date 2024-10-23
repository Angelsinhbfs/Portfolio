package main

import (
	"encoding/json"
	"net/http"
	"time"
)

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func Login(w http.ResponseWriter, r *http.Request) {
	var creds Credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// This is where you would normally validate the credentials
	if creds.Username == "user" && creds.Password == "password" {
		token, err := GenerateJWT(creds.Username)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		http.SetCookie(w, &http.Cookie{
			Name:    "token",
			Value:   token,
			Expires: time.Now().Add(5 * time.Minute),
		})
	} else {
		w.WriteHeader(http.StatusUnauthorized)
	}
}
