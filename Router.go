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
	mux.HandleFunc("/", HandleRoot)

	mux.HandleFunc("/portfolio/img", HandlePortfolioImg)
	mux.HandleFunc("/portfolio/create", HandleCreateTile)
	mux.HandleFunc("/portfolio/load", HandleLoadTiles)
	mux.HandleFunc("/portfolio/edit", HandleEditTile)
	mux.HandleFunc("/portfolio/delete", HandleDelete)
	mux.HandleFunc("/login", Login)

	log.Print("Listening on 8080")

	http.ListenAndServe(":8080", mux)
}
