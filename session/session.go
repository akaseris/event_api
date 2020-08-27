package session

// importing the required packages
import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

var fileName string = "session\\active_sessions.txt"

// Find id
func Find(id string) int {
	// Reading File
	wDir, err := os.Getwd()
	data, err := ioutil.ReadFile(wDir + "\\" + fileName)
	if err != nil {
		log.Printf("failed reading data from file: %s", err)
		return -2
	}

	// Converting binary data to array of strings
	var arrData []string
	err = json.Unmarshal([]byte(data), &arrData)
	if err != nil {
		log.Printf("failed processing data from file: %s", err)
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

// Add id
func Add(id string) bool {
	// Reading File
	wDir, err := os.Getwd()
	data, err := ioutil.ReadFile(wDir + "\\" + fileName)
	if err != nil {
		log.Printf("failed reading data from file: %s", err)
		return false
	}

	// Converting binary data to array of strings
	var arrData []string
	err = json.Unmarshal([]byte(data), &arrData)
	if err != nil {
		log.Printf("failed processing data from file: %s", err)
		return false
	}

	// Adding id to array
	arrData = append(arrData, id)

	// Creating file to overwrite the old one
	newFile, err := os.Create(fileName)
	if err != nil {
		log.Printf("failed creating file: %s", err)
		return false
	}

	// Converting to prefered form and writing data to file
	len, err := newFile.WriteString("[\"" + strings.Join(arrData, "\",\"") + "\"]")
	if err != nil {
		log.Printf("failed writing to file: %s with length %d", err, len)
		return false
	}

	return true
}

// Remove id
func Remove(index int) bool {
	// Reading File
	wDir, err := os.Getwd()
	data, err := ioutil.ReadFile(wDir + "\\" + fileName)
	if err != nil {
		log.Printf("failed reading data from file: %s", err)
		return false
	}

	// Converting binary data to array of strings
	var arrData []string
	err = json.Unmarshal([]byte(data), &arrData)
	if err != nil {
		log.Printf("failed processing data from file: %s", err)
		return false
	}

	// Removing id from array
	arrData = append(arrData[:index], arrData[index+1:]...)

	// Creating file to overwrite the old one
	newFile, err := os.Create(fileName)
	if err != nil {
		log.Printf("failed creating file: %s", err)
		return false
	}

	// Converting to prefered form and writing data to file
	var str string
	if len(arrData) == 0 {
		str = "[]"
	} else {
		str = "[\"" + strings.Join(arrData, "\",\"") + "\"]"
	}
	len, err := newFile.WriteString(str)
	if err != nil {
		log.Printf("failed writing to file: %s with length %d", err, len)
		return false
	}

	return true
}
