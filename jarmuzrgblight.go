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
	"fmt"
	"log"
	"os"
	"net/http"
	"encoding/json"

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

var globalConfig jmodConfig

func main() {
	http.HandleFunc("/webComponent", WebComponentHandler)
	http.HandleFunc("/instanceData", InstanceDataHandler)

	// Get Passed config daata
	initConfig()
	log.Println(globalConfig)

	// Get port to start HTTP server
	port := os.Getenv("JABLKO_MOD_PORT")
	log.Printf("Jablko Mod Port: %s", port)

	http.ListenAndServe(":" + port, nil)
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

		err := json.Unmarshal([]byte(defaultConfig), &globalConfig)
		if err != nil {
			log.Printf("FATAL ERROR: Default config is invalid")
			panic(err)
		}
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

func WebComponentHandler(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadFile("./webcomponent.js")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	fmt.Println("Not Implemented")
	fmt.Fprintf(w, "%s", b);
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
	fmt.Fprintf(w, `%s`, b);
}
