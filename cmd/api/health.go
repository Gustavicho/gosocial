package main

import "net/http"

func (app *application) handlerHealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("It's OK!"))
	w.WriteHeader(200)
}
