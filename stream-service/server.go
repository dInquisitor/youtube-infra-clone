package main

import (
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
)

const VIDEO_DIRECTORY = "../data/video-store"

// for now, we assume anonymous access for all videos

const MANIFEST_FILE_NAME = "manifest.mpd"
const MANIFEST_FILE_MIME_TYPE = "application/dash+xml"

const VIDEO_FILE_MIMETYPE = "video/webm"

// ideally should be moved to cdn
func main() {
	// Init http server
	e := echo.New()
	e.Pre(echoMiddleware.RemoveTrailingSlash())

	e.GET("/api/stream/:videoID/:fileToServe", func(c echo.Context) error {
		inputVideoID := c.Param("videoID")

		// TODO? check if video exist, is fully processed and user can see it (cache this fact in redis since multiple calls will be made to this endpoint throughout the course of streaming) -> ensure to remove slashes

		inputFileToServe := c.Param("fileToServe")

		videoPath, err := filepath.Abs(fmt.Sprintf("%s/%s/%s", VIDEO_DIRECTORY, inputVideoID, inputFileToServe))

		if err != nil {
			return c.NoContent(http.StatusNotFound)
		}

		fileMimeType := VIDEO_FILE_MIMETYPE

		if inputFileToServe == MANIFEST_FILE_NAME {
			fileMimeType = MANIFEST_FILE_MIME_TYPE
		}

		c.Response().Header().Set("Content-Type", fileMimeType)
		c.Response().Header().Set("Access-Control-Allow-Origin", "*")

		return c.File(videoPath)
	})

	e.Logger.Fatal(e.Start(":7000"))
}
