package main

import (
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

var addon_version string = "not_injected"
var addon_githash string = "not_injected"

func startEchoServer() {
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Enable CORS for any host
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE},
	}))

	// Routes
	e.GET("/", slideshowHandler)
	e.GET("/manage", slideshowHandler)
	e.GET("/config", slideshowHandler)
	e.GET("/images/:filename", imageHandler)
	e.GET("/thumbnails/:filename", thumbnailHandler)
	e.GET("/assets/:filename", assetHandler)
	e.GET("/nextslide", nextSlideHandler)
	e.GET("/slides", slidesHandler)
	e.DELETE("/slides/:filename", slidesDeleteHandler)
	e.GET("/info", infoHandler)

	// Start server
	e.Logger.Fatal(e.Start(":" + port))
}

func infoHandler(c echo.Context) error {
	remoteIP := c.RealIP()
	// get the size of the image directory
	slidesSize := getDirSize(outputfolder)
	// get the size of the thumbnail directory
	slidesSize += getDirSize(thumbnailfolder)

	info := SlideInfo{
		SlidesCount: len(slides),
		RemoteIP:    remoteIP,
		SlidesSize:  slidesSize,
		Version:     addon_version,
		GitCommit:   addon_githash,
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
	for i, slide := range slides {
		if slide.Filename == filename {
			// Remove the slide from the list
			slides = append(slides[:i], slides[i+1:]...)
			saveSlides(slides)

			// Delete the image and thumbnail files
			os.Remove(outputfolder + slide.Filename)
			os.Remove(thumbnailfolder + slide.Filename)

			return c.NoContent(http.StatusOK)
		}
	}
	return c.String(http.StatusNotFound, "Slide not found")
}

func slidesHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, slides)
}

var lastSlideIndex int = -1

func nextSlideHandler(c echo.Context) error {
	if len(slides) == 0 {
		return c.String(http.StatusNotFound, "No slides available")
	}
	if len(slides) == 1 {
		return c.JSON(http.StatusOK, slides[0])
	}

	// Create a new random generator
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	var newSlideIndex int
	for {
		newSlideIndex = r.Intn(len(slides))
		if newSlideIndex != lastSlideIndex {
			break
		}
	}

	lastSlideIndex = newSlideIndex
	randomSlide := slides[newSlideIndex]

	return c.JSON(http.StatusOK, randomSlide)
}

func assetHandler(c echo.Context) error {
	filename := c.Param("filename")
	return c.File(frontenddist + "assets/" + filename)
}

func imageHandler(c echo.Context) error {
	filename := c.Param("filename")
	return c.File(outputfolder + filename)
}

func thumbnailHandler(c echo.Context) error {
	filename := c.Param("filename")
	return c.File(thumbnailfolder + filename)
}

func slideshowHandler(c echo.Context) error {
	return c.File(frontenddist + "index.html")
}
