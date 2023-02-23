package handlers

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"api-server/middlewares"

	"github.com/go-redis/redis/v9"
	"github.com/labstack/echo/v4"
)

type VideoData struct {
	VideoID string `json:"id"`
}

type BeginUploadResponse struct {
	Status string    `json:"status"`
	Data   VideoData `json:"videoData,omitempty"`
}

var ctx = context.Background()

const UPLOAD_KEY_PREFIX = "upload-lease"
const UPLOAD_ALLOW_VALUE = 1

const UPLOAD_LEASE_DURATION_MINUTES = 5

// This may be used to give access to an aws server upload
func BeginUpload(c echo.Context, dbPool *sql.DB, redisHandle *redis.Client) error {
	const (
		STATUS_SUCCESS       = "SUCCESS"
		STATUS_UNKNOWN_ERROR = "UNKNOWN_ERROR"
	)

	userID := c.Get("authContext").(*middlewares.AuthContext).UserID

	// create video
	videoID, err := createVideo(userID, dbPool)

	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusOK, &BeginUploadResponse{
			Status: STATUS_UNKNOWN_ERROR,
		})
	}

	// give access to upload server
	if err := redisHandle.Set(ctx, fmt.Sprintf("%s:%s:%s", UPLOAD_KEY_PREFIX, userID, videoID), UPLOAD_ALLOW_VALUE, time.Duration(time.Duration.Minutes(UPLOAD_LEASE_DURATION_MINUTES))).Err(); err != nil {
		log.Printf("Could not acquire upload lease %v -> /api/begin-upload \n", err)
		return c.JSON(http.StatusOK, &BeginUploadResponse{
			Status: STATUS_UNKNOWN_ERROR,
		})
	}

	return c.JSON(http.StatusOK, &BeginUploadResponse{
		Status: STATUS_SUCCESS,
		Data: VideoData{
			VideoID: videoID,
		},
	})
}

func createVideo(userID string, dbPool *sql.DB) (videoID string, err error) {
	// insertVideoResult := dbPool.QueryRow("insert into videos(author_id) values($1) returning id", userID)

	// err = insertVideoResult.Scan(&videoID)

	// if err != nil {
	// 	return "", err
	// }

	const DUMMY_VIDEO_ID = "9fdb5aa2-5540-469e-8eb0-c0932980784f"

	return DUMMY_VIDEO_ID, nil
}
