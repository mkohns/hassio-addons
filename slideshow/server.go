package main

import (
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

	// Routes for the SPA
	e.GET("/", slideshowHandler)
	e.GET("/manage", slideshowHandler)
	e.GET("/config", slideshowHandler)
	e.GET("/assets/:filename", assetHandler)

	// API routes
	e.GET("/images/:filename", imageHandler)
	e.GET("/thumbnails/:filename", thumbnailHandler)
	e.GET("/nextslide", nextSlideHandler)
	e.GET("/slides", slidesHandler)
	e.DELETE("/slides/:filename", slidesDeleteHandler)
	e.PATCH("/slides/:filename", slidesPatchHandler)
	e.PATCH("portraitMode", portraitModeHandler)
	e.GET("/info", infoHandler)

	// Start server
	e.Logger.Fatal(e.Start(":" + port))
}
