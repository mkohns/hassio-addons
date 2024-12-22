package main

import (
	"encoding/json"
	"log"
	"os"
)

// some global variables
var username string
var password string
var accountNo string
var outputfolder string
var thumbnailfolder string
var configfolder string
var port string
var groupId string
var groupIdReal string
var frontenddist string
var signalws string
var signalapi string

func initConfig() {
	// check if running in hassio
	if fileExists("/data/options.json") {
		// open options file
		file, err := os.Open("/data/options.json")
		if err != nil {
			log.Panicln("Error opening file:", err)
		}
		defer file.Close()

		// Create a new JSON decoder
		decoder := json.NewDecoder(file)
		config := HAConfig{}
		if err := decoder.Decode(&config); err != nil {
			log.Panicln("Error decoding options.json:", err.Error())
		}

		username = config.SignalUsername
		password = config.SignalPassword
		accountNo = config.SignalAccountNo
		outputfolder = config.SignalOutputFolder
		thumbnailfolder = config.SignalThumbnailFolder
		groupId = config.SignalGroupID
		groupIdReal = config.SignalGroupIDReal
		signalws = config.SignalSignalWS
		signalapi = config.SignalSignalAPI

		port = config.SlideshowPort
		configfolder = config.SlideshowConfigDir
		frontenddist = config.SlideshowFrontendDist
	} else {
		// loading envrionment variables
		loadEnvVariable("SIGNAL_USERNAME", &username)
		loadEnvVariable("SIGNAL_PASSWORD", &password)
		loadEnvVariable("SIGNAL_ACCOUNTNO", &accountNo)
		loadEnvVariable("SIGNAL_OUTPUTFOLDER", &outputfolder)
		loadEnvVariable("SIGNAL_THUMBNAILFOLDER", &thumbnailfolder)
		loadEnvVariable("SIGNAL_GROUPID", &groupId)
		loadEnvVariable("SIGNAL_GROUPID_REAL", &groupIdReal)
		loadEnvVariable("SIGNAL_SIGNALWS", &signalws)
		loadEnvVariable("SIGNAL_SIGNALAPI", &signalapi)

		loadEnvVariable("SLIDESHOW_PORT", &port)
		loadEnvVariable("SLIDESHOW_CONFIGDIR", &configfolder)
		loadEnvVariable("SLIDESHOW_FRONTEND_DIST", &frontenddist)
	}
}

// fileExists checks if a file exists and is not a directory.
func fileExists(filePath string) bool {
	info, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func loadEnvVariable(envname string, varpointer *string) {
	if os.Getenv(envname) != "" {
		*varpointer = os.Getenv(envname)
	} else {
		log.Panicf("%s not set", envname)
	}
}
