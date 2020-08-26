package session

// importing the required packages
import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

var fileName string = "active_sessions.txt"

func find(file string, id string) int {
	// Reading File
	wDir, err := os.Getwd()
	data, err := ioutil.ReadFile(wDir + "\\" + file)
	if err != nil {
		log.Panicf("failed reading data from file: %s", err)
		return -2
	}

	// Converting binary data to array of strings
	var arrData []string
	err = json.Unmarshal([]byte(data), &arrData)
	if err != nil {
		log.Panicf("failed processing data from file: %s", err)
		return -2
	}

	// Searching for id
	for i := 0; i < len(arrData); i++ {
		if arrData[i] == id {
			return i
		}
	}
	return -1
}

func add(file string, id string) bool {
	// Reading File
	wDir, err := os.Getwd()
	data, err := ioutil.ReadFile(wDir + "\\" + file)
	if err != nil {
		log.Panicf("failed reading data from file: %s", err)
		return false
	}

	// Converting binary data to array of strings
	var arrData []string
	err = json.Unmarshal([]byte(data), &arrData)
	if err != nil {
		log.Panicf("failed processing data from file: %s", err)
		return false
	}

	// Adding id to array
	arrData = append(arrData, id)

	// Creating file to overwrite the old one
	newFile, err := os.Create(file)
	if err != nil {
		log.Fatalf("failed creating file: %s", err)
		return false
	}

	// Converting to prefered form and writing data to file
	len, err := newFile.WriteString("[\"" + strings.Join(arrData, "\",\"") + "\"]")
	if err != nil {
		log.Fatalf("failed writing to file: %s with length %d", err, len)
		return false
	}

	return true
}

func remove(file string, index int) bool {
	// Reading File
	wDir, err := os.Getwd()
	data, err := ioutil.ReadFile(wDir + "\\" + file)
	if err != nil {
		log.Panicf("failed reading data from file: %s", err)
		return false
	}

	// Converting binary data to array of strings
	var arrData []string
	err = json.Unmarshal([]byte(data), &arrData)
	if err != nil {
		log.Panicf("failed processing data from file: %s", err)
		return false
	}

	// Removing id from array
	arrData = append(arrData[:index], arrData[index+1:]...)

	// Creating file to overwrite the old one
	newFile, err := os.Create(file)
	if err != nil {
		log.Fatalf("failed creating file: %s", err)
		return false
	}

	// Converting to prefered form and writing data to file
	len, err := newFile.WriteString("[\"" + strings.Join(arrData, "\",\"") + "\"]")
	if err != nil {
		log.Fatalf("failed writing to file: %s with length %d", err, len)
		return false
	}

	return true
}
