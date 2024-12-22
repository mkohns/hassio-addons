package main

import (
	"encoding/json"
	"log"
	"os"
	"sync"
)

var slideSessions map[string]*SlideSessionInfo
var slideSessionMutex sync.RWMutex

func initSessions() {
	sessionsFilename := configfolder + "sessions.json"
	// Open the file for reading
	if fileExists(sessionsFilename) {
		file, err := os.Open(configfolder + "sessions.json")
		if err != nil {
			log.Panicln("Error opening file:", err)
		}
		defer file.Close()

		// Create a new JSON decoder
		decoder := json.NewDecoder(file)

		// Decode the session from the file
		if err := decoder.Decode(&slideSessions); err != nil {
			log.Panicln("Error decoding sessions:", err.Error())
		}
	} else {
		log.Printf("No sessions.json found, creating empty session")
		slideSessionMutex.Lock()
		slideSessions = make(map[string]*SlideSessionInfo)
		// Create default session
		slideSessions[""] = &SlideSessionInfo{
			LastConfig:        nil,
			LastSlideIndex:    -1,
			PrioNewSlides:     true,
			NewSlidesPriority: nil,
		}
		saveSessions(slideSessions)
		slideSessionMutex.Unlock()
		log.Println("Created empty sessions.json")
		return
	}
}

func saveSessions(sessions map[string]*SlideSessionInfo) {
	// Open the file for writing
	file, err := os.Create(configfolder + "sessions.json")
	if err != nil {
		log.Panicln("Error creating file:", err.Error())
	}
	defer file.Close()

	// Create a new JSON encoder
	encoder := json.NewEncoder(file)

	// Encode the slides to the file
	if err := encoder.Encode(sessions); err != nil {
		log.Panicln("Error encoding slides:", err.Error())
	}
}

func getSession(sessionID string) *SlideSessionInfo {
	slideSessionMutex.RLock()
	defer slideSessionMutex.RUnlock()

	if slideSessions[sessionID] != nil {
		return slideSessions[sessionID]
	}
	session := &SlideSessionInfo{
		LastConfig:        nil,
		LastSlideIndex:    -1,
		NewSlidesPriority: nil,
		PrioNewSlides:     false,
	}
	slideSessions[sessionID] = session
	saveSessions(slideSessions)
	return session

}

func handleNewPrioritySession(newIndex int) {
	slideSessionMutex.Lock()
	// Iterate over all slideSessions
	for i, session := range slideSessions {
		// If the session is not active, set the new slide as the active slide
		if session.PrioNewSlides && session.NewSlidesPriority != nil {
			log.Printf("Adding new prio slide <%d> to session <%s>", newIndex, i)
			*slideSessions[i].NewSlidesPriority = append(*slideSessions[i].NewSlidesPriority, newIndex)
		}
	}
	slideSessionMutex.Unlock()
}

func updateSession(sessionID string, session *SlideSessionInfo) {
	if session == nil {
		return
	}

	slideSessionMutex.Lock()
	defer slideSessionMutex.Unlock()

	slideSessions[sessionID] = session
	saveSessions(slideSessions)
}
