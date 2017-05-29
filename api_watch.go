package main

import (
	"net/http"
)

// apiWatch adds the user to a repository watchlist.
func apiWatch(w http.ResponseWriter, r *http.Request) error {
	var res Result
	return respond(w, &res)
}

// apiUnwatch removes the user from a repo's watchlist.
func apiUnwatch(w http.ResponseWriter, r *http.Request) error {
	var res Result
	return respond(w, &res)
}
