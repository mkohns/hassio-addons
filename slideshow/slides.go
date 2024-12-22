package main

import (
	"encoding/json"
	"log"
	"os"
	"sync"
	"time"
)

// slides holder
var slides []Slide
var slideMutex sync.RWMutex

func getFilteredSlides(config *SlideShowConfig) []int {
	slideMutex.RLock()
	defer slideMutex.RUnlock()
	var filteredSlides []int = make([]int, 0)
	// There are three filters that can be applied
	// all filters are combined with AND
	// 1. Show only favorites
	// 2. Show only active
	// 3. Show only in time frame
	// 3.1 The slides createdAt must be after the start date
	// 3.2 If the end date is not set, the end date is now
	// 3.3 The slides createdAt must be before the end date
	// The return value should be a list of indexes of the slides that match the filters
	// If no slides match the filters, an empty list should be returned
	// If no filters are set, all slides should be returned

	// If no filters are set, return all slides
	if !config.ShowOnlyFavorites && !config.ShowOnlyActive && !config.ShowOnlyInTimeFrame {
		for i := range slides {
			filteredSlides = append(filteredSlides, i)
		}
		return filteredSlides
	}

	for i, slide := range slides {
		// 1. Show only favorites
		if config.ShowOnlyFavorites && !slide.Favorite {
			continue
		}
		// 2. Show only active
		if config.ShowOnlyActive && !slide.Enabled {
			continue
		}
		// 3. Show only in time frame
		if config.ShowOnlyInTimeFrame {
			// 3.1 The slides createdAt must be after the start date
			if config.StartDate != nil && slide.CreatedAt.Before(*config.StartDate) {
				continue
			}
			// 3.2 If the end date is not set, the end date is now
			endDate := config.EndDate
			if endDate == nil {
				now := time.Now()
				endDate = &now
			}
			// 3.3 The slides createdAt must be before the end date
			if slide.CreatedAt.After(*endDate) {
				continue
			}
		}
		// Add the index of the slide that matches all filters
		filteredSlides = append(filteredSlides, i)
	}

	return filteredSlides
}

func getSlideCount() int {
	slideMutex.RLock()
	defer slideMutex.RUnlock()
	return len(slides)
}

func updateSlideEnabled(attachmentID string, enabled bool) {
	slideMutex.Lock()
	defer slideMutex.Unlock()

	for index, slide := range slides {
		if slide.AttachmentID == attachmentID {
			slides[index].Enabled = enabled
			break
		}
	}
	saveSlides(slides)
}

func updateSlideFavorite(attachmentID string, favorite bool) {
	slideMutex.Lock()
	defer slideMutex.Unlock()

	for index, slide := range slides {
		if slide.AttachmentID == attachmentID {
			slides[index].Favorite = favorite
			break
		}
	}
	saveSlides(slides)
}

func getSlideByIndex(index int) *Slide {
	slideMutex.RLock()
	defer slideMutex.RUnlock()
	if index < 0 || index >= len(slides) {
		return nil
	}
	return &slides[index]
}

func getSlideByAttachmentID(attachmentID string) *Slide {
	slideMutex.RLock()
	defer slideMutex.RUnlock()
	for _, slide := range slides {
		if slide.AttachmentID == attachmentID {
			return &slide
		}
	}
	return nil
}

func getFavoriteSlideCount() int {
	slideMutex.RLock()
	defer slideMutex.RUnlock()
	count := 0
	for _, slide := range slides {
		if slide.Favorite {
			count++
		}
	}
	return count
}

func getActiveSlideCount() int {
	slideMutex.RLock()
	defer slideMutex.RUnlock()
	count := 0
	for _, slide := range slides {
		if slide.Enabled {
			count++
		}
	}
	return count
}

func removeSlideByAttachementID(attachmentID string) {
	// remove the slide from the list
	slideMutex.RLock()
	newSlides := make([]Slide, 0)
	for _, slide := range slides {
		if slide.AttachmentID == attachmentID {
			log.Printf("Removing slide <%s> from list\n", slide.AttachmentID)
			// Delete the image and thumbnail files
			os.Remove(outputfolder + slide.AttachmentID)
			os.Remove(thumbnailfolder + slide.AttachmentID)
		} else {
			newSlides = append(newSlides, slide)
		}
	}
	slideMutex.RUnlock()
	slideMutex.RLock()
	slides = newSlides
	saveSlides(slides)
	slideMutex.RUnlock()
}

func removeSlideByTimestamp(timestamp int64) {
	// remove the slide from the list
	slideMutex.RLock()
	newSlides := make([]Slide, 0)
	for i, slide := range slides {
		if slide.MsgTimestamp == timestamp {
			log.Printf("Removing slide <%d> from list\n", i)
			// Delete the image and thumbnail files
			os.Remove(outputfolder + slide.AttachmentID)
			os.Remove(thumbnailfolder + slide.AttachmentID)
		} else {
			newSlides = append(newSlides, slide)
		}
	}
	slideMutex.RUnlock()
	slideMutex.RLock()
	slides = newSlides
	saveSlides(slides)
	slideMutex.RUnlock()
}

func addSlide(s Slide) int {
	slideMutex.Lock()
	slides = append(slides, s)
	saveSlides(slides)
	newIndex := len(slides) - 1
	slideMutex.Unlock()
	return newIndex
}

func readSlides() []Slide {
	// Open the file for reading
	file, err := os.Open(configfolder + "slides.json")
	if err != nil {
		log.Println("Error opening file:", err)
		slideMutex.Lock()
		slides = make([]Slide, 0)
		saveSlides(slides)
		slideMutex.Unlock()
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

// Caution, when calling this, the mutex must be locked
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

func initSlides() {
	slides = readSlides()
}
