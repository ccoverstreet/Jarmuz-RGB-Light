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
	"strings"
	"sync"

	"io/ioutil"

	"github.com/ccoverstreet/Jarmuz-RGB-Light/jablkodev"
	"github.com/gorilla/websocket"
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
	sync.RWMutex
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
	// Get passed jmodKey. Used for authenticating jmods with Jablko
	globalJMODKey = os.Getenv("JABLKO_MOD_KEY")
	globalJablkoCorePort = os.Getenv("JABLKO_CORE_PORT")
	globalJMODPort = os.Getenv("JABLKO_MOD_PORT")

	// Get Passed config daata
	initConfig()
	log.Println(globalConfig)

	// Handles called by Jablko
	http.HandleFunc("/webComponent", WebComponentHandler)
	http.HandleFunc("/instanceData", InstanceDataHandler)
	http.HandleFunc("/jmod/socket", SocketHandler)

	log.Println(http.ListenAndServe(":"+globalJMODPort, nil))
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

	err = jablkodev.JablkoSaveConfig(globalJablkoCorePort, globalJMODPort, globalJMODKey, configBytes)
	if err != nil {
		log.Printf("ERROR: Unable to save config - %v", err)
	}
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

// WebSocketHandler
var upgrader = websocket.Upgrader{CheckOrigin: func(*http.Request) bool { return true }}

func SocketHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Websocket handler called")

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("ERROR: Unable to upgrade WebSocket - %v", err)
	}
	defer conn.Close()

	log.Println("Websocket connection established")
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("ERROR: Error reading WebSocket message - %v", err)
			return
		}

		sliceEnd := strings.Index(string(message), ",")
		log.Println(string(message[:sliceEnd]))

	}
}
