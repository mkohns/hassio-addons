package main

import "github.com/labstack/echo"

func slideshowHandler(c echo.Context) error {
	return c.File(frontenddist + "index.html")
}

func assetHandler(c echo.Context) error {
	filename := c.Param("filename")
	return c.File(frontenddist + "assets/" + filename)
}
