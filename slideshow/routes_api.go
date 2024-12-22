package main

import (
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/labstack/echo"
)

func imageHandler(c echo.Context) error {
	filename := c.Param("filename")
	return c.File(outputfolder + filename)
}

func thumbnailHandler(c echo.Context) error {
	filename := c.Param("filename")
	return c.File(thumbnailfolder + filename)
}

var currentPortraitMode = false

type PortraitPatchBody struct {
	PortraitMode bool `json:"PortraitMode"`
}

func portraitModeHandler(c echo.Context) error {
	// Get the body
	body := new(PortraitPatchBody)
	err := c.Bind(body)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "couldNotRetrieveBody")
	}

	currentPortraitMode = body.PortraitMode

	currentMode := "Portrait"
	if !currentPortraitMode {
		currentMode = "Landscape"
	}
	sendMessage("The Bilderrahmen was rearranged to mode: " + currentMode)

	return c.NoContent(http.StatusOK)
}

func slidesPatchHandler(c echo.Context) error {
	filename := c.Param("filename")

	// Get the body
	body := new(SlidePatchBody)
	err := c.Bind(body)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "couldNotRetrieveBody")
	}

	// check if favorite or enabled is set
	if body.Favorite == nil && body.Enabled == nil {
		return echo.NewHTTPError(http.StatusBadRequest, "emptyBody")
	}

	slideMutex.Lock()
	defer slideMutex.Unlock()

	for i, slide := range slides {
		if slide.Filename == filename {

			if body.Favorite != nil {
				slides[i].Favorite = *body.Favorite
			}

			if body.Enabled != nil {
				slides[i].Enabled = *body.Enabled
			}

			saveSlides(slides)
			return c.NoContent(http.StatusOK)
		}
	}
	return c.String(http.StatusNotFound, "Slide not found")
}

func infoHandler(c echo.Context) error {
	remoteIP := c.RealIP()
	// get the size of the image directory
	slidesSize := getDirSize(outputfolder)
	// get the size of the thumbnail directory
	thumbnailSize := getDirSize(thumbnailfolder)
	// get favorite and enabled count
	favoriteCount := 0
	enabledCount := 0
	for _, slide := range slides {
		if slide.Favorite {
			favoriteCount++
		}
		if slide.Enabled {
			enabledCount++
		}
	}

	info := SlideInfo{
		SlidesCount:   len(slides),
		RemoteIP:      remoteIP,
		SlidesSize:    slidesSize,
		ThumbnailSize: thumbnailSize,
		Version:       addon_version,
		GitCommit:     addon_githash,
		FavoriteCount: favoriteCount,
		ActiveCount:   enabledCount,
	}
	return c.JSON(http.StatusOK, info)
}

func getDirSize(path string) int {
	var size int64
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return nil
	})
	if err != nil {
		return 0
	}
	return int(size)
}

func slidesDeleteHandler(c echo.Context) error {
	filename := c.Param("filename")
	slideMutex.RLock()

	for i, slide := range slides {
		if slide.Filename == filename {
			// Remove the slide from the list
			slides = append(slides[:i], slides[i+1:]...)
			saveSlides(slides)
			slideMutex.RUnlock()
			// Delete the image and thumbnail files
			os.Remove(outputfolder + slide.Filename)
			os.Remove(thumbnailfolder + slide.Filename)

			return c.NoContent(http.StatusOK)
		}
	}
	slideMutex.RUnlock()
	return c.String(http.StatusNotFound, "Slide not found")
}

func slidesHandler(c echo.Context) error {
	slideMutex.RLock()
	defer slideMutex.RUnlock()
	return c.JSON(http.StatusOK, slides)
}

type SlideShowConfig struct {
	SessionID                 string
	ShowOnlyFavorites         bool
	ShowOnlyActive            bool
	ShowOnlyInTimeFrame       bool
	ShowNewImagesWithPriority bool
	ModeRandom                bool
	ModeChronological         bool
	ModeReverseChronological  bool
	StartDate                 *time.Time
	EndDate                   *time.Time
}

const (
	ModeRandom               string = "random"
	ModeChronological        string = "chronological"
	ModeReverseChronological string = "reverseChronological"
)

type SlideSessionInfo struct {
	lastConfig        *SlideShowConfig
	randomGenerator   *rand.Rand
	lastSlideIndex    int
	newSlidesPriority *[]int
	prioNewSlides     bool
}

var slideSessions = make(map[string]*SlideSessionInfo)

func init() {
	log.Println("Initializing default slideSession")
	slideSessions[""] = &SlideSessionInfo{
		lastConfig:        nil,
		randomGenerator:   rand.New(rand.NewSource(time.Now().UnixNano())),
		lastSlideIndex:    -1,
		prioNewSlides:     false,
		newSlidesPriority: nil,
	}
}

func (config SlideShowConfig) getMode() string {
	if config.ModeRandom {
		return ModeRandom
	} else if config.ModeChronological {
		return ModeChronological
	} else if config.ModeReverseChronological {
		return ModeReverseChronological
	}
	log.Println("No mode detected, defaulting to random")
	return ModeRandom
}

func getSlideShowConfig(c echo.Context) SlideShowConfig {
	showOnlyFavorites := c.QueryParam("showOnlyFavorites")
	showOnlyActive := c.QueryParam("showOnlyActive")
	showOnlyInTimeFrame := c.QueryParam("showOnlyInTimeFrame")
	showNewImagesWithPriority := c.QueryParam("ShowNewImagesWithPriority")
	modeRandom := c.QueryParam("modeRandom")
	modeChronological := c.QueryParam("modeChronological")
	modeReverseChronological := c.QueryParam("modeReverseChronological")
	startDate := c.QueryParam("startDate")
	endDate := c.QueryParam("endDate")

	config := SlideShowConfig{
		ShowOnlyFavorites:         showOnlyFavorites == "true",
		ShowOnlyActive:            showOnlyActive == "true",
		ShowOnlyInTimeFrame:       showOnlyInTimeFrame == "true",
		ShowNewImagesWithPriority: showNewImagesWithPriority == "true",
		ModeRandom:                modeRandom == "true",
		ModeChronological:         modeChronological == "true",
		ModeReverseChronological:  modeReverseChronological == "true",
	}

	if startDate != "" {
		startDateParsed, err := time.Parse("02/01/2006", startDate)
		if err == nil {
			config.StartDate = &startDateParsed
		}
	}

	if endDate != "" {
		endDateParsed, err := time.Parse("02/01/2006", endDate)
		if err == nil {
			config.EndDate = &endDateParsed
		}
	}

	return config
}

func nextSlideHandler(c echo.Context) error {
	if len(slides) == 0 {
		return c.String(http.StatusNotFound, "No pictures available 😢")
	}

	// get slide config
	config := getSlideShowConfig(c)

	// get current session
	session := slideSessions[config.SessionID]
	if session == nil {
		session = &SlideSessionInfo{
			lastConfig:        nil,
			randomGenerator:   rand.New(rand.NewSource(time.Now().UnixNano())),
			lastSlideIndex:    -1,
			newSlidesPriority: nil,
			prioNewSlides:     false,
		}
		slideSessions[config.SessionID] = session
	}

	if (config.ShowNewImagesWithPriority && session.lastConfig == nil) ||
		(config.ShowNewImagesWithPriority && session.lastConfig != nil && !session.lastConfig.ShowNewImagesWithPriority) {
		log.Printf("activating new slides priority for session <%s> \n", config.SessionID)
		session.prioNewSlides = true
		newArray := make([]int, 0)
		session.newSlidesPriority = &newArray
	} else if (!config.ShowNewImagesWithPriority && session.lastConfig == nil) ||
		(!config.ShowNewImagesWithPriority && session.lastConfig != nil && session.lastConfig.ShowNewImagesWithPriority) {
		log.Printf("deactivating new slides priority for session <%s> \n", config.SessionID)
		session.prioNewSlides = false
		session.newSlidesPriority = nil
	}
	slideSessions[config.SessionID] = session

	if session.prioNewSlides && session.newSlidesPriority != nil && len(*session.newSlidesPriority) > 0 {
		// get first entry
		nextIndex := (*session.newSlidesPriority)[0]
		// remove first entry
		*session.newSlidesPriority = (*session.newSlidesPriority)[1:]
		// persist
		session.lastConfig = &config
		slideSessions[config.SessionID] = session
		// return slide
		log.Printf("Returning slide with priority: %d, remaining prios <%d> \n", nextIndex, len(*session.newSlidesPriority))
		return c.JSON(http.StatusOK, slides[nextIndex])
	} else if session.prioNewSlides {
		log.Printf("No slides with priority available for session: <%s> \n", config.SessionID)
	}

	slideMutex.RLock()
	defer slideMutex.RUnlock()
	// filter slides
	filteredSlides := getFilteredSlides(&config)
	if len(filteredSlides) == 0 {
		return c.String(http.StatusNotFound, "No slides available after filtering 🤔")
	}

	// check if Slide mode changed
	if session.lastConfig == nil || session.lastConfig.getMode() != config.getMode() {
		// start over
		session.lastSlideIndex = -1
		session.lastConfig = &config
		slideSessions[config.SessionID] = session
	}

	session.lastConfig = &config
	slideSessions[config.SessionID] = session

	if config.getMode() == ModeRandom {
		// Fixes endlessloop if only one slide is available
		if len(filteredSlides) == 1 {
			return c.JSON(http.StatusOK, slides[filteredSlides[0]])
		}
		var newSlideIndex int
		for {
			newSlideIndex = session.randomGenerator.Intn(len(filteredSlides))
			// never return the same slide twice in a row
			if newSlideIndex != session.lastSlideIndex {
				break
			}
		}
		session.lastSlideIndex = newSlideIndex
		slideSessions[config.SessionID] = session
		return c.JSON(http.StatusOK, slides[filteredSlides[newSlideIndex]])
	} else if config.getMode() == ModeChronological {
		if session.lastSlideIndex+1 >= len(filteredSlides) {
			session.lastSlideIndex = -1
		}
		session.lastSlideIndex++
		slideSessions[config.SessionID] = session
		return c.JSON(http.StatusOK, slides[filteredSlides[session.lastSlideIndex]])
	} else if config.getMode() == ModeReverseChronological {
		if session.lastSlideIndex-1 < 0 {
			session.lastSlideIndex = len(filteredSlides)
		}
		session.lastSlideIndex--
		slideSessions[config.SessionID] = session
		return c.JSON(http.StatusOK, slides[filteredSlides[session.lastSlideIndex]])
	}

	// this should never happen
	log.Println("No mode detected, this should never happen")
	return c.NoContent(http.StatusNotFound)
}

func getFilteredSlides(config *SlideShowConfig) []int {
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
