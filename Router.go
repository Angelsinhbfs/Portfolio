package main

import (
	"log"
	"net/http"
)

func Route() {
	mux := http.NewServeMux()

	//rh := http.RedirectHandler("http://dontremember.me", 307)
	//
	//mux.Handle("/foo", rh)
	mux.HandleFunc("GET /portfolio/img/", GetImg)
	mux.HandleFunc("POST /portfolio/img/", CheckToken(PostImg))
	mux.HandleFunc("/portfolio/create/", CheckToken(HandleCreateTile))
	mux.HandleFunc("/portfolio/load/", HandleLoadTiles)
	mux.HandleFunc("/portfolio/edit/", CheckToken(HandleEditTile))
	mux.HandleFunc("/portfolio/delete/", CheckToken(HandleDelete))
	mux.HandleFunc("/login/", BasicAuth(HandleLogin))
	mux.HandleFunc("/", HandleRoot)

	log.Print("Listening on 8080")

	http.ListenAndServe(":8080", mux)
}
