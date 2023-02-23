package handlers

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"gopkg.in/boj/redistore.v1"
)

type SignOutResponse struct {
	Status string `json:"status"`
}

func SignOut(c echo.Context, dbPool *sql.DB, sessionHandler *redistore.RediStore) error {
	const SESSION_KEY = "up_sessid"

	const (
		STATUS_SUCCESS       = "SUCCESS"
		STATUS_UNKNOWN_ERROR = "UNKNOWN_ERROR"
	)

	currentSession, err := sessionHandler.Get(c.Request(), SESSION_KEY)
	if err != nil {
		return err
	}

	currentSession.Options.MaxAge = -1
	if err = currentSession.Save(c.Request(), c.Response().Writer); err != nil {
		log.Printf("Error saving session while logging user out: %v\n", err)
		return c.JSON(http.StatusOK, &SignOutResponse{
			Status: STATUS_UNKNOWN_ERROR,
		})
	}

	return c.JSON(http.StatusOK, &SignOutResponse{
		Status: STATUS_SUCCESS,
	})
}
