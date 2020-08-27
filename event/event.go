package main

// importing the required packages
import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"sort"
)

// var fileName string = "event\\event.json"
// Temporary only for standalone use
var fileName string = "event.json"

type child struct {
	Type      string `json:"type"`
	Timestamp int    `json:"timestamp"`
	Name      string `json:"name"`
}

type sessionObject struct {
	Type     string  `json:"type"`
	ID       string  `json:"id"`
	Start    int     `json:"start"`
	End      int     `json:"end"`
	Children []child `json:"children"`
}

func AddSessionStart(id string, start int) bool {
	// Check if there is already a session with refered id
	if existingSession, index := findSession(id); existingSession != nil {
		log.Panicf("Session already exists with index: %d", index)
		return false
	}

	// Prepare session object
	var sessionToInput sessionObject
	sessionToInput.ID, sessionToInput.Type, sessionToInput.Start = id, "SESSION", start

	// Reading File
	wDir, err := os.Getwd()
	data, err := ioutil.ReadFile(wDir + "\\" + fileName)
	if err != nil {
		log.Panicf("failed reading data from file: %s", err)
		return false
	}

	// Converting binary data to array of strings
	var arrData []sessionObject
	err = json.Unmarshal([]byte(data), &arrData)
	if err != nil {
		log.Panicf("failed processing data from file: %s", err)
		return false
	}

	// Adding session object to array
	arrData = append(arrData, sessionToInput)

	// Creating file to overwrite the old one
	newFile, err := os.Create(fileName)
	if err != nil {
		log.Panicf("failed creating file: %s", err)
		return false
	}

	// Converting to compatible format
	jsonData, err := json.Marshal(arrData)
	if err != nil {
		log.Panicf("failed prossesing data")
		return false
	}

	// Converting to prefered form and writing data to file
	len, err := newFile.WriteString(string(jsonData))
	if err != nil {
		log.Fatalf("failed writing to file: %s with length %d", err, len)
		return false
	}
	return true
}

func AddSessionEnd(id string, end int) bool {
	// Check if there is already a session with refered id
	existingSession, index := findSession(id)
	if existingSession == nil {
		log.Panicf("Session does not exist")
		return false
	}

	// Prepare session object
	existingSession.End = end

	// Reading File
	wDir, err := os.Getwd()
	data, err := ioutil.ReadFile(wDir + "\\" + fileName)
	if err != nil {
		log.Panicf("failed reading data from file: %s", err)
		return false
	}

	// Converting binary data to array of strings
	var arrData []sessionObject
	err = json.Unmarshal([]byte(data), &arrData)
	if err != nil {
		log.Panicf("failed processing data from file: %s", err)
		return false
	}

	// Replacing old session object to updated one
	arrData[index] = *existingSession

	// Creating file to overwrite the old one
	newFile, err := os.Create(fileName)
	if err != nil {
		log.Panicf("failed creating file: %s", err)
		return false
	}

	// Converting to compatible format
	jsonData, err := json.Marshal(arrData)
	if err != nil {
		log.Panicf("failed prossesing data")
		return false
	}

	// Converting to prefered form and writing data to file
	len, err := newFile.WriteString(string(jsonData))
	if err != nil {
		log.Fatalf("failed writing to file: %s with length %d", err, len)
		return false
	}
	return true
}

func AddChildren(id string, timestamp int, name string) bool {
	// Check if there is already a session with refered id
	existingSession, index := findSession(id)
	if existingSession == nil {
		log.Panicf("Session does not exist")
		return false
	}

	// Prepare child object
	var childToInput child
	childToInput.Type, childToInput.Name, childToInput.Timestamp = "EVENT", name, timestamp

	// Reading File
	wDir, err := os.Getwd()
	data, err := ioutil.ReadFile(wDir + "\\" + fileName)
	if err != nil {
		log.Panicf("failed reading data from file: %s", err)
		return false
	}

	// Converting binary data to array of strings
	var arrData []sessionObject
	err = json.Unmarshal([]byte(data), &arrData)
	if err != nil {
		log.Panicf("failed processing data from file: %s", err)
		return false
	}

	// Replacing old children object to updated one
	arrData[index].Children = append(arrData[index].Children, childToInput)

	// Sorting children by timestamp
	sort.Slice(arrData[index].Children, func(i, j int) bool {
		return arrData[index].Children[i].Timestamp < arrData[index].Children[j].Timestamp
	})

	// Creating file to overwrite the old one
	newFile, err := os.Create(fileName)
	if err != nil {
		log.Panicf("failed creating file: %s", err)
		return false
	}

	// Converting to compatible format
	jsonData, err := json.Marshal(arrData)
	if err != nil {
		log.Panicf("failed prossesing data")
		return false
	}

	// Converting to prefered form and writing data to file
	len, err := newFile.WriteString(string(jsonData))
	if err != nil {
		log.Fatalf("failed writing to file: %s with length %d", err, len)
		return false
	}
	return true
}

func findSession(id string) (*sessionObject, int) {
	var arrData []sessionObject

	// Reading File
	wDir, err := os.Getwd()
	data, err := ioutil.ReadFile(wDir + "\\" + fileName)
	if err != nil {
		log.Panicf("failed reading data from file: %s", err)
		return nil, -1
	}

	// Converting binary data to array of structs
	err = json.Unmarshal([]byte(data), &arrData)
	if err != nil {
		log.Panicf("failed processing data from file: %s", err)
		return nil, -1
	}

	// Searching for session object with matching id
	for i := 0; i < len(arrData); i++ {
		if arrData[i].ID == id {
			return &arrData[i], i
		}
	}
	return nil, -1
}

func main() {
	AddChildren("tasos", 999999999999, "cart")
	AddChildren("tasos", 888888888888, "cart")
	AddChildren("tasos", 1523, "cart")
}
