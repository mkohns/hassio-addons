package main

import (
	"log"
)

func main() {
	// log greetings
	log.Println("Starting Slideshow Server")

	// load the configuration
	initConfig()

	// init signal rest client
	initClient(signalapi)

	// init sessions
	initSessions()

	// Read the slides from the file
	initSlides()

	// Open a websocket to signal api
	go connectToWebSocket(signalws+accountNo, username, password)

	// Start the Echo server
	startEchoServer()
}
