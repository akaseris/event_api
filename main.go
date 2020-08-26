package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/akaseris/event_api/session"
	"github.com/gorilla/mux"
)

var acceptedTypes = []string{"SESSION_START", "EVENT", "SESSION_END"}

func handleSessionStart(jsonData []map[string]interface{}, w http.ResponseWriter) bool {
	if str, ok := jsonData[0]["session_id"].(string); ok {
		found := session.Find(str)
		if found > -1 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"error": "Session already exists"}`))
			return false
		} else if found < -1 {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error": "Error searching for sessions"}`))
			return false
		} else {
			session.Add(str)
			return true
		}
	}
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(`{"error": "Session is not a string"}`))
	return false
}

func get(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "get called"}`))
}

func post(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Map body data into an array of JSONs
	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "can not decode body"}`))
		return
	}
	var jsonData []map[string]interface{}
	err = json.Unmarshal([]byte(bodyBytes), &jsonData)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "can not decode body"}`))
		return
	}

	// Check for types
	for i := 0; i < len(jsonData); i++ {
		if jsonData[i]["type"] != "SESSION_START" && jsonData[i]["type"] != "SESSION_END" && jsonData[i]["type"] != "EVENT" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"error": "wrong type"}`))
			return
		}
	}

	// Handle SESSION_START
	if jsonData[0]["type"] == "SESSION_START" {
		if ok := handleSessionStart(jsonData, w); !ok {
			return
		}
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"message": "post called"}`))
}

func checkParams(r *http.Request) bool {
	query := r.URL.Query()
	id := query.Get("session_ID")
	if id == "123456789" {
		return true
	}
	return false
}

func main() {
	r := mux.NewRouter()

	api := r.PathPrefix("/api/v1").Subrouter()
	api.HandleFunc("", get).Methods(http.MethodGet)
	api.HandleFunc("", post).Methods(http.MethodPost)

	log.Fatal(http.ListenAndServe(":8080", r))
}
