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

// some global variables
var randomGenerator = rand.New(rand.NewSource(time.Now().UnixNano()))
var currentPortraitMode = false

func imageHandler(c echo.Context) error {
	attachmentID := c.Param("attachmentID")
	slide := getSlideByAttachmentID(attachmentID)
	if slide == nil {
		return c.String(http.StatusNotFound, "Slide not found")
	}
	return c.File(outputfolder + attachmentID)
}

func thumbnailHandler(c echo.Context) error {
	attachmentID := c.Param("attachmentID")
	slide := getSlideByAttachmentID(attachmentID)
	if slide == nil {
		return c.String(http.StatusNotFound, "Slide not found")
	}
	return c.File(thumbnailfolder + attachmentID)
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
	attachmentID := c.Param("attachmentID")
	slide := getSlideByAttachmentID(attachmentID)
	if slide == nil {
		return c.String(http.StatusNotFound, "Slide not found")
	}

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

	if body.Favorite != nil {
		updateSlideFavorite(slide.AttachmentID, *body.Favorite)
	}

	if body.Enabled != nil {
		updateSlideEnabled(slide.AttachmentID, *body.Enabled)
	}

	return c.JSON(http.StatusOK, getSlideByAttachmentID(attachmentID))
}

func infoHandler(c echo.Context) error {
	remoteIP := c.RealIP()
	// get the size of the image directory
	slidesSize := getDirSize(outputfolder)
	// get the size of the thumbnail directory
	thumbnailSize := getDirSize(thumbnailfolder)
	// get favorite and enabled count
	favoriteCount := getFavoriteSlideCount()
	enabledCount := getActiveSlideCount()

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
	attachmentID := c.Param("attachmentID")
	slide := getSlideByAttachmentID(attachmentID)
	if slide == nil {
		return c.String(http.StatusNotFound, "Slide not found")
	}

	// remove the slide from the list
	removeSlideByAttachementID(attachmentID)

	return c.NoContent(http.StatusOK)
}

func slidesHandler(c echo.Context) error {
	slideMutex.RLock()
	defer slideMutex.RUnlock()
	return c.JSON(http.StatusOK, slides)
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
	showNewImagesWithPriority := c.QueryParam("showNewImagesWithPriority")
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
	if getSlideCount() == 0 {
		return c.String(http.StatusNotFound, "No pictures available 😢")
	}

	// get slide config
	config := getSlideShowConfig(c)

	// get current session
	session := getSession(config.SessionID)

	if (config.ShowNewImagesWithPriority && session.LastConfig == nil) ||
		(config.ShowNewImagesWithPriority && session.LastConfig != nil && !session.LastConfig.ShowNewImagesWithPriority) {
		log.Printf("activating new slides priority for session <%s> \n", config.SessionID)
		session.PrioNewSlides = true
		newArray := make([]int, 0)
		session.NewSlidesPriority = &newArray
		updateSession(config.SessionID, session)
	} else if (!config.ShowNewImagesWithPriority && session.LastConfig == nil) ||
		(!config.ShowNewImagesWithPriority && session.LastConfig != nil && session.LastConfig.ShowNewImagesWithPriority) {
		log.Printf("deactivating new slides priority for session <%s> \n", config.SessionID)
		session.PrioNewSlides = false
		session.NewSlidesPriority = nil
		updateSession(config.SessionID, session)
	}

	if session.PrioNewSlides && session.NewSlidesPriority != nil && len(*session.NewSlidesPriority) > 0 {
		// get first entry
		nextIndex := (*session.NewSlidesPriority)[0]
		// remove first entry
		*session.NewSlidesPriority = (*session.NewSlidesPriority)[1:]
		// persist
		session.LastConfig = &config
		updateSession(config.SessionID, session)

		// return slide
		log.Printf("Returning slide with priority: %d, remaining prios <%d> \n", nextIndex, len(*session.NewSlidesPriority))
		return c.JSON(http.StatusOK, slides[nextIndex])
	} else if session.PrioNewSlides {
		log.Printf("No slides with priority available for session: <%s> \n", config.SessionID)
	}

	// filter slides
	filteredSlides := getFilteredSlides(&config)
	if len(filteredSlides) == 0 {
		return c.String(http.StatusNotFound, "No slides available after filtering 🤔")
	}

	// check if Slide mode changed
	if session.LastConfig == nil || session.LastConfig.getMode() != config.getMode() {
		log.Println("Slide Mode changed, resetting LastSlideIndex")
		// start over
		session.LastSlideIndex = -1
	}

	// update LastConfig
	session.LastConfig = &config
	updateSession(config.SessionID, session)

	if config.getMode() == ModeRandom {
		// Fixes endlessloop if only one slide is available
		if len(filteredSlides) == 1 {
			return c.JSON(http.StatusOK, slides[filteredSlides[0]])
		}
		var newSlideIndex int
		for {
			newSlideIndex = randomGenerator.Intn(len(filteredSlides))
			// never return the same slide twice in a row
			if newSlideIndex != session.LastSlideIndex {
				break
			}
		}
		session.LastSlideIndex = newSlideIndex
	} else if config.getMode() == ModeChronological {
		if session.LastSlideIndex+1 >= len(filteredSlides) {
			session.LastSlideIndex = -1
		}
		session.LastSlideIndex++
	} else if config.getMode() == ModeReverseChronological {
		if session.LastSlideIndex-1 < 0 {
			session.LastSlideIndex = len(filteredSlides)
		}
		session.LastSlideIndex--
	} else {
		log.Panicf("No mode detected, this should never happen")
	}

	if session.LastSlideIndex == -1 {
		log.Panicf("LastSlideIndex is -1, this should never happen")
	}

	updateSession(config.SessionID, session)
	return c.JSON(http.StatusOK, getSlideByIndex(filteredSlides[session.LastSlideIndex]))
}
