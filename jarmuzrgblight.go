// Jarmuz RGB Light
// Cale Overstreet
// May 10, 2021

/*
A Jablko Mod that communicates using UDP to RGB Lights. This
module is to serve as a demo for the full development chain
of Jablko.
*/

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"io/ioutil"
)

const defaultConfig = `{
	"instances": {
		"inst1": {
			"lightIPs": [
				"10.0.0.5"
			]
		}	
	}
}
`

type jmodConfig struct {
	Instances map[string]instanceData `json:"instances"`
}

type instanceData struct {
	LightIPs []string `json:"lightIPs"`
}

// -------------------- GLOBALS --------------------
var globalConfig jmodConfig
var globalJMODKey string
var globalJMODPort string
var globalJablkoCorePort string

// -------------------- END GLOBALS --------------------

func main() {
	http.HandleFunc("/webComponent", WebComponentHandler)
	http.HandleFunc("/instanceData", InstanceDataHandler)

	// Get passed jmodKey. Used for authenticating jmods with Jablko
	globalJMODKey = os.Getenv("JABLKO_MOD_KEY")
	globalJablkoCorePort = os.Getenv("JABLKO_CORE_PORT")
	globalJMODPort = os.Getenv("JABLKO_MOD_PORT")

	// Get Passed config daata
	initConfig()
	log.Println(globalConfig)

	// Get port to start HTTP server
	log.Printf("Jablko Mod Port: %s", globalJablkoCorePort)

	http.ListenAndServe(":"+globalJMODPort, nil)
}

func initConfig() {
	confStr := os.Getenv("JABLKO_MOD_CONFIG")
	log.Printf("\"%s\"", confStr)

	// Check if config was provided. Replace confStr with default
	// if not.
	if len(confStr) < 3 {
		log.Println("No config provided. Starting with default config")
		loadDefaultConfig()
		// Should also send a request to Jablko with updated config
		return
	}

	err := json.Unmarshal([]byte(confStr), &globalConfig)
	if err != nil {
		log.Printf("Provided config is invalid. Loading default config: %v", err)
	}
}

func loadDefaultConfig() {
	err := json.Unmarshal([]byte(defaultConfig), &globalConfig)
	if err != nil {
		log.Printf("FATAL ERROR: Default config is invalid")
		panic(err)
	}

	saveConfig()
}

// This function sends a JSON of the current config to Jablko
// which then triggers a config save on Jablko.
func saveConfig() {
	configBytes, err := json.Marshal(globalConfig)
	if err != nil {
		log.Printf("Unable to marshal config: %v", err)
		return
	}

	client := &http.Client{}

	req, err := http.NewRequest("POST", "http://localhost:"+globalJablkoCorePort+"/service/saveConfig", bytes.NewBuffer(configBytes))
	if err != nil {
		log.Println(err)
		return
	}

	req.Header.Add("JMOD-KEY", globalJMODKey)
	req.Header.Add("JMOD-PORT", globalJMODPort)

	log.Println(globalJMODPort)

	client.Do(req)
}

func WebComponentHandler(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadFile("./webcomponent.js")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "%s", b)
}

func InstanceDataHandler(w http.ResponseWriter, r *http.Request) {
	b, err := json.Marshal(globalConfig.Instances)
	if err != nil {
		log.Printf("Unable to generate JSON from globalConfig.Instances: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Unable to generate JSON string for instances: %v", err)
		return
	}
	fmt.Fprintf(w, `%s`, b)
}
