package main

import (
	"api-gateway/rproxy"

	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
)

func main() {
	// Init http server
	e := echo.New()
	e.Pre(echoMiddleware.RemoveTrailingSlash())

	e.Any("/*", rproxy.ReverseProxy)
	// e.Any("/api", func(c echo.Context) error {
	// 	log.Println("here got here")
	// 	return c.JSON(http.StatusOK, nil)
	// })

	e.Logger.Fatal(e.Start(":4100"))
}
