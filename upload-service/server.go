package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"upload-service/middlewares"

	"github.com/segmentio/kafka-go"

	"github.com/go-redis/redis/v9"
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	"gopkg.in/boj/redistore.v1"
)

type VideoUploadResponse struct {
	Status string `json:"status"`
}

const UPLOAD_KEY_PREFIX = "upload-lease"
const VIDEO_DIRECTORY = "../data/video-store"

const VIDEO_PROCESSING_TOPIC = "process-video"
const VIDEO_PROCESSING_PARTITION = 0

var ctx = context.Background()

func main() {
	// Init session store
	sessionHandler, err := redistore.NewRediStore(10, "tcp", "redis:6379", "", []byte("secret-key"))

	if err != nil {
		log.Fatal(err)
	}

	const SECONDS_IN_FIVE_DAYS = 5 * 24 * 3600
	sessionHandler.SetMaxAge(SECONDS_IN_FIVE_DAYS)

	defer sessionHandler.Close()

	// Init redis connection for other uses
	redisHandle := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	// Init kafka
	kafkaWriter := &kafka.Writer{
		Addr:     kafka.TCP("kafka:9092"),
		Topic:    VIDEO_PROCESSING_TOPIC,
		Balancer: &kafka.LeastBytes{},
	}

	defer kafkaWriter.Close()

	// Init http server
	e := echo.New()
	e.Pre(echoMiddleware.RemoveTrailingSlash())

	e.POST("/api/upload", func(c echo.Context) error {
		const (
			STATUS_SUCCESS       = "SUCCESS"
			STATUS_UNKNOWN_ERROR = "UNKNOWN_ERROR"
		)

		unknownErrorResponse := c.JSON(http.StatusOK, &VideoUploadResponse{
			Status: STATUS_UNKNOWN_ERROR,
		})

		inputVideoID := c.FormValue("videoID")

		userID := c.Get("authContext").(*middlewares.AuthContext).UserID

		if err := confirmUploadAccess(inputVideoID, userID, redisHandle); err != nil {
			log.Println("User does not have upload access", err)
			return unknownErrorResponse
		}

		// Source
		file, err := c.FormFile("videoFile")

		if err != nil {
			log.Println("File not present in http request", err)
			return unknownErrorResponse
		}
		src, err := file.Open()
		if err != nil {
			log.Println("Could not open uploaded file", err)
			return unknownErrorResponse
		}
		defer src.Close()

		if err = os.Mkdir(fmt.Sprintf("%s/%s/", VIDEO_DIRECTORY, inputVideoID), os.ModePerm); err != nil {
			log.Println("Could not make directory for video", err)
			return unknownErrorResponse
		}

		destFilePath, err := filepath.Abs(fmt.Sprintf("%s/%s/original", VIDEO_DIRECTORY, inputVideoID))
		if err != nil {
			log.Println("Could not expand destination file path", err)
			return unknownErrorResponse
		}

		// Destination
		dst, err := os.Create(destFilePath)
		if err != nil {
			log.Println("Could not open destination for writing", err)
			return unknownErrorResponse
		}
		defer dst.Close()

		log.Println("Uploading...")

		// Copy
		if _, err = io.Copy(dst, src); err != nil {
			log.Println("Could not copy uploaded file to destination", err)
			return unknownErrorResponse
		}

		err = notifyVideoProcessingWorkers(inputVideoID, kafkaWriter)

		if err != nil {
			log.Println(err)
			return unknownErrorResponse
		}

		log.Println("Written to kafka")

		return c.JSON(http.StatusOK, &VideoUploadResponse{
			Status: STATUS_SUCCESS,
		})
	}, func(next echo.HandlerFunc) echo.HandlerFunc {
		return middlewares.BasicAuth(next, sessionHandler)
	})

	e.Logger.Fatal(e.Start(":5000"))
}

func confirmUploadAccess(videoID string, userID string, redisHandle *redis.Client) (err error) {
	_, err = redisHandle.Get(ctx, fmt.Sprintf("%s:%s:%s", UPLOAD_KEY_PREFIX, userID, videoID)).Result()
	if err != nil {
		return err
	}

	return
}

func notifyVideoProcessingWorkers(inputVideoID string, kafkaWriter *kafka.Writer) error {
	err := kafkaWriter.WriteMessages(
		context.Background(),
		kafka.Message{Value: []byte(inputVideoID)},
	)

	if err != nil {
		return err
	}

	return nil
}
