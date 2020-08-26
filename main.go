package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func get(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "get called"}`))
}

func post(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"message": "post called"}`))
}

func params(w http.ResponseWriter, r *http.Request) {
	//pathParams := mux.Vars(r)
	w.Header().Set("Content-Type", "application/json")

	w.Write([]byte(fmt.Sprintf("Not handling yet")))
}

func main() {
	r := mux.NewRouter()

	api := r.PathPrefix("/api/v1").Subrouter()
	api.HandleFunc("", params).Methods(http.MethodGet)
	api.HandleFunc("", params).Methods(http.MethodPost)

	log.Fatal(http.ListenAndServe(":8080", r))
}
