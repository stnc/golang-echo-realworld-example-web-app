package main

/*
license : www.selmantunc.com.tr

*/

import (
	"net/http"

	"github.com/flosch/pongo2"

	"github.com/labstack/echo"

	"github.com/labstack/echo/middleware"
	"github.com/siredwin/pongorenderer/renderer" // this package is not publicly available
)

var (
	data         = pongo2.Context{}
	MainRenderer = renderer.Renderer{Debug: true} // use any renderer
)

func index() func(echo.Context) error {
	return func(c echo.Context) error {
		return c.Render(http.StatusOK, "resources/index.html", data)
	}
}

func main() {

	// Echo instance
	e := echo.New()
	e.Renderer = MainRenderer //pongo init

	e.Static("/assets", "resources/assets") //https://echo.labstack.com/guide/static-files

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// initialMigration(dbConn)

	// Route => handler
	e.GET("/", index())

	// Start server
	e.Logger.Fatal(e.Start(":8000"))
}
