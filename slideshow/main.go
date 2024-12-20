package main

import (
	"encoding/json"
	"log"
	"os"
)

// slides holder
var slides []Slide

// some global variables
var username string
var password string
var accountNo string
var outputfolder string
var thumbnailfolder string
var configfolder string
var port string
var groupId string
var frontenddist string
var signalws string
var signalapi string

func main() {
	log.Println("Starting Slideshow Server")

	// check if running in hassio
	if FileExists("/data/options.json") {
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
		loadEnvVariable("SIGNAL_SIGNALWS", &signalws)
		loadEnvVariable("SIGNAL_SIGNALAPI", &signalapi)

		loadEnvVariable("SLIDESHOW_PORT", &port)
		loadEnvVariable("SLIDESHOW_CONFIGDIR", &configfolder)
		loadEnvVariable("SLIDESHOW_FRONTEND_DIST", &frontenddist)
	}

	// init signal rest client
	InitClient(signalapi)

	// Read the slides from the file
	slides = readSlides()

	// Open a websocket to signal api
	go connectToWebSocket(signalws+accountNo, username, password)

	// Start the Echo server
	startEchoServer()
}

// FileExists checks if a file exists and is not a directory.
func FileExists(filePath string) bool {
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

func readSlides() []Slide {
	// Open the file for reading
	file, err := os.Open(configfolder + "slides.json")
	if err != nil {
		log.Println("Error opening file:", err)
		slides = make([]Slide, 0)
		saveSlides(slides)
		log.Println("Created empty slides.json")
		return slides
	}
	defer file.Close()

	// Create a new JSON decoder
	decoder := json.NewDecoder(file)

	// Decode the slides from the file
	var slides []Slide
	if err := decoder.Decode(&slides); err != nil {
		log.Panicln("Error decoding slides:", err.Error())
	}

	return slides
}

func saveSlides(slides []Slide) {
	// Open the file for writing
	file, err := os.Create(configfolder + "slides.json")
	if err != nil {
		log.Panicln("Error creating file:", err.Error())
	}
	defer file.Close()

	// Create a new JSON encoder
	encoder := json.NewEncoder(file)

	// Encode the slides to the file
	if err := encoder.Encode(slides); err != nil {
		log.Panicln("Error encoding slides:", err.Error())
	}
}
