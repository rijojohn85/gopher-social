package main

import (
	"log"
	"net/http"
)

func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("hello")
	w.Write([]byte("OK"))
}
