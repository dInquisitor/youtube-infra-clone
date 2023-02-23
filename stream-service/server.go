package main

import (
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
)

const VIDEO_DIRECTORY = "../data/video-store"

// const STREAM_CHUNK_SIZE int64 = 524_288 // 512KB

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

		// --> turns out c.File handles range requests automatically

		// requestedRange := c.Request().Header[http.CanonicalHeaderKey("range")][0]
		// inputVideoFidelity := "1080"

		// videoStart, videoEnd, videoSize, videoPath, err := calculateChunkDetails(inputVideoID, inputVideoFidelity, requestedRange)

		// if err != nil {
		// 	return c.NoContent(http.StatusNotFound)
		// }

		// contentLength := videoEnd - videoStart + 1

		// c.Response().Header().Set("Content-Range", fmt.Sprintf("bytes %d-%d/%d", videoStart, videoEnd, videoSize))
		// c.Response().Header().Set("Accept-Ranges", "bytes")
		// c.Response().Header().Set("Content-Length", fmt.Sprintf("%d", contentLength))
		// c.Response().Header().Set("Cache-Control", "no-cache")
		// // c.Response().Header().Set("Content-Type", "video/mp4")

		// videoFile, err := os.Open(videoPath)

		// if err != nil {
		// 	log.Println("Could not read video file while streaming", err)
		// 	return c.NoContent(http.StatusNotFound)
		// }

		// videoBuffer := make([]byte, contentLength)

		// _, err = videoFile.ReadAt(videoBuffer, videoStart)

		// if err != nil {
		// 	log.Println("Could not read video file into response buffer", err)
		// 	return c.NoContent(http.StatusNotFound)
		// }

		// return c.Blob(http.StatusPartialContent, "video/mp4", videoBuffer)
	})

	e.Logger.Fatal(e.Start(":7000"))
}

// func calculateChunkDetails(inputVideoID string, inputVideoFidelity string, requestedRange string) (int64, int64, int64, string, error) {
// 	removeNonDigitsRegex := regexp.MustCompile(`\D`)
// 	videoStartRangeFromInputRange := removeNonDigitsRegex.ReplaceAllString(requestedRange, "")

// 	videoStart, err := strconv.ParseInt(videoStartRangeFromInputRange, 10, 64)

// 	if err != nil {
// 		log.Println("Could not parse requested video range while streaming", err)
// 		return 0, 0, 0, "", err
// 	}

// 	videoPath, err := filepath.Abs(fmt.Sprintf("%s/%s/%s", VIDEO_DIRECTORY, inputVideoID, inputVideoFidelity))

// 	if err != nil {
// 		log.Println("Could not expand video path while streaming video", err)
// 		return 0, 0, 0, "", err
// 	}

// 	video, err := os.Stat(videoPath)

// 	if err != nil {
// 		log.Println("Could not get video details while streaming", err)
// 		return 0, 0, 0, "", err
// 	}

// 	videoSize := video.Size()
// 	videoEnd := minInt(videoStart+STREAM_CHUNK_SIZE-1, videoSize-1)

// 	return videoStart, videoEnd, videoSize, videoPath, nil
// }

// func minInt(a, b int64) int64 {
// 	if a < b {
// 		return a
// 	}

// 	return b
// }
