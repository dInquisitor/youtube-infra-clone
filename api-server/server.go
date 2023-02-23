package main

import (
	"api-server/handlers"
	"api-server/middlewares"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-redis/redis/v9"
	_ "github.com/lib/pq"
	"gopkg.in/boj/redistore.v1"

	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
)

// Considerations:
// Might use echo context to pass around db connection singleton in the future, same goes with session and redis
// User should only be able to upload a video at once? so that people can't spam database with begin-upload requests without actually uploading anything
// Wrapper library to reconnect to kafka for the producer side (maybe the consumer side?)

// TODO:
// Bind expect input on every endpoint: https://echo.labstack.com/guide/binding/

func main() {
	// Init Database
	fmt.Println(os.Getenv("POSTGRES_PASSWORD"))
	dbPool, err := sql.Open("postgres", fmt.Sprintf("host=postgres user=postgres password=%s dbname=video_uploader sslmode=disable", os.Getenv("POSTGRES_PASSWORD")))

	if err != nil {
		log.Fatal(err)
	}

	// Init session store
	// Secret key is for encrypting session cookies
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

	// Init http server
	e := echo.New()
	e.Pre(echoMiddleware.RemoveTrailingSlash())

	e.Any("/", func(c echo.Context) error {
		log.Println("here got here")
		return c.JSON(http.StatusOK, nil)
	})

	e.POST("/api/sign-in", func(c echo.Context) error {
		return handlers.SignIn(c, dbPool, sessionHandler)
	})

	e.GET("/api/sign-out", func(c echo.Context) error {
		return handlers.SignOut(c, dbPool, sessionHandler)
	})

	e.POST("/api/begin-upload", func(c echo.Context) error {
		return handlers.BeginUpload(c, dbPool, redisHandle)
	}, func(next echo.HandlerFunc) echo.HandlerFunc {
		return middlewares.BasicAuth(next, sessionHandler)
	})

	e.Logger.Fatal(e.Start(":4000"))
}
