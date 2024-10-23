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
	mux.HandleFunc("/tube/{username}", LoadTube)
	mux.HandleFunc("/tube/inbox", HandleInbox)
	mux.HandleFunc("/tube/{username}/inbox", HandleUserInbox)
	mux.HandleFunc("/tube/{username}/outbox", HandleOutbox)
	mux.HandleFunc("/tubes/{username}/collections/{collection}", HandleCollections)
	mux.HandleFunc("/tubes/{username}/img/{imgGUID}", HandleImg)
	mux.HandleFunc("/portfolio/img", HandlePortfolioImg)
	mux.HandleFunc("/portfolio/create", HandleCreateTile)
	mux.HandleFunc("/portfolio/load", HandleLoadTiles)
	mux.HandleFunc("/portfolio/edit", HandleEditTile)
	mux.HandleFunc("/portfolio/delete", HandleDelete)
	mux.HandleFunc("/login", Login)

	mux.HandleFunc("/.well-known/webfinger", WebfingerHandler) ///.well-known/webfinger?resource=user@example.com&rel=http%3A%2F%2Fopenid.net%2Fspecs%2Fconnect%2F1.0%2Fissuer

	log.Print("Listening on 8080")

	http.ListenAndServe(":8080", mux)
}
