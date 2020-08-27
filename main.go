package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/akaseris/event_api/event"
	"github.com/akaseris/event_api/session"
	"github.com/gorilla/mux"
)

// ActiveSession Global var for active session
var ActiveSession string

func handleSessionStart(jsonData []map[string]interface{}, w http.ResponseWriter, paramID []string, paramIDIndex int) bool {
	// Check if the first JSON is of type SESSION_START and it matches the id on the parameters
	if jsonData[0]["type"] == "SESSION_START" && jsonData[0]["session_id"] == paramID[0] {
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
				timestamp, ok := jsonData[0]["timestamp"].(float64)
				if !ok {
					w.WriteHeader(http.StatusBadRequest)
					w.Write([]byte(`{"error": "Timestamp is not an int"}`))
					return false
				}
				sesOk := session.Add(str)
				ActiveSession = str
				evOk := event.AddSessionStart(ActiveSession, int(timestamp))
				if !sesOk || !evOk {
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte(`{"error": "Error processing data"}`))
					return false
				}
				return true
			}
		}
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "Start session is not a string"}`))
		return false
		// Else catch he case where it does not match the one on the parameters
	} else if jsonData[0]["type"] == "SESSION_START" && jsonData[0]["session_id"] != paramID[0] {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "session_id in body does not match session_id in parameters"}`))
		return false
		// Else catch the case where the first JSON was not of type SESSION_START and the id on the parameters is not in the active sessions
	} else if paramIDIndex < 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "session_id is not registered and there is no SESSION_START EVENT"}`))
		return false
	}
	// Return true if no condition is met to continue with the request
	ActiveSession = paramID[0]
	return true
}

func handleSessionEnd(jsonData []map[string]interface{}, w http.ResponseWriter, paramID []string) bool {
	// Check if the first JSON is of type SESSION_END and it matches the id on the parameters
	if jsonData[len(jsonData)-1]["type"] == "SESSION_END" && jsonData[len(jsonData)-1]["session_id"] == paramID[0] {
		if str, ok := jsonData[len(jsonData)-1]["session_id"].(string); ok {
			found := session.Find(str)
			if found == -1 {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(`{"error": "End session does not exist"}`))
				return false
			} else if found < -1 {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`{"error": "Error searching for sessions"}`))
				return false
			} else {
				timestamp, ok := jsonData[len(jsonData)-1]["timestamp"].(float64)
				if !ok {
					w.WriteHeader(http.StatusBadRequest)
					w.Write([]byte(`{"error": "Timestamp is not an int"}`))
					return false
				}
				sesOk := session.Remove(found)
				evOk := event.AddSessionEnd(ActiveSession, int(timestamp))
				if !sesOk || !evOk {
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte(`{"error": "Error processing data"}`))
					return false
				}
				return true
			}
		}
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "Session is not a string"}`))
		return false
		// Else catch he case where it does not match the one on the parameters
	} else if jsonData[len(jsonData)-1]["type"] == "SESSION_END" && jsonData[len(jsonData)-1]["session_id"] != paramID[0] {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "session_id in body does not match session_id in parameters"}`))
		return false
	}
	// Return true if no condition is met to continue with the request
	return true
}

func handleEvents(jsonData []map[string]interface{}, w http.ResponseWriter) bool {
	for i := 0; i < len(jsonData); i++ {
		if jsonData[i]["type"] == "SESSION_START" || jsonData[i]["type"] == "SESSION_END" {
			continue
		}

		timestamp, ok := jsonData[i]["timestamp"].(float64)
		if !ok || jsonData[i]["name"] == nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"error": "Wrong event object structure"}`))
			return false
		}
		strName, ok := jsonData[i]["name"].(string)
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"error": "Event name is not a string"}`))
			return false
		}

		// Adding event to children list sorted in the coresponding sessionId
		if ok := event.AddChildren(ActiveSession, int(timestamp), strName); !ok {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error": "Error processing the data"}`))
			return false
		}
	}
	return true
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

	// Check for id in params
	paramID, ok := r.URL.Query()["session_id"]
	if !ok {
		w.WriteHeader(http.StatusNonAuthoritativeInfo)
		w.Write([]byte(`{"error": "session_id does not exist in parameters"}`))
		return
	}
	paramIDIndex := session.Find(string(paramID[0]))

	// Handle SESSION_START
	if ok := handleSessionStart(jsonData, w, paramID, paramIDIndex); !ok {
		return
	}

	// Handle EVENT
	if ok := handleEvents(jsonData, w); !ok {
		return
	}

	// Handle SESSION_END
	if ok := handleSessionEnd(jsonData, w, paramID); !ok {
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"message": "post called"}`))
}

func main() {
	r := mux.NewRouter()

	api := r.PathPrefix("/api/v1").Subrouter()
	api.HandleFunc("", get).Methods(http.MethodGet)
	api.HandleFunc("", post).Methods(http.MethodPost)

	log.Fatal(http.ListenAndServe(":8080", r))
}
