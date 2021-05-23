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

// -------------------- END GLOBALS --------------------

func main() {
	http.HandleFunc("/webComponent", WebComponentHandler)
	http.HandleFunc("/instanceData", InstanceDataHandler)

	// Get passed jmodKey. Used for authenticating jmods with Jablko
	globalJMODKey = os.Getenv("JABLKO_MOD_KEY")
	log.Println(globalJMODKey)

	// Get Passed config daata
	initConfig()
	log.Println(globalConfig)

	// Get port to start HTTP server
	port := os.Getenv("JABLKO_MOD_PORT")
	log.Printf("Jablko Mod Port: %s", port)

	http.ListenAndServe(":"+port, nil)
}

func initConfig() {
	confStr := os.Getenv("JABLKO_MOD_CONFIG")
	log.Printf("\"%s\"", confStr)

	// Check if config was provided. Replace confStr with default
	// if not.
	if confStr == "" {
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
}

// This function sends a JSON of the current config to Jablko
// which then triggers a config save on Jablko.
func saveConfig() {

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
	fmt.Println("Not Implemented")
	fmt.Fprintf(w, `%s`, b)
}
