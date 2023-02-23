package session

import (
	"log"

	"github.com/labstack/echo/v4"
	"gopkg.in/boj/redistore.v1"
)

const SESSION_KEY = "up_sessid"

func SetUserSession(sessionHandler *redistore.RediStore, c echo.Context, key string, value string) error {

	currentSession, err := sessionHandler.Get(c.Request(), SESSION_KEY)
	if err != nil {
		return err
	}

	currentSession.Values[key] = value
	if err = currentSession.Save(c.Request(), c.Response().Writer); err != nil {
		return err
	}

	return nil
}

func GetUserIDFromSession(c echo.Context, sessionHandler *redistore.RediStore) (userID string, err error) {

	currentSession, err := sessionHandler.Get(c.Request(), SESSION_KEY)
	if err != nil {
		return "", err
	}

	userIDFromSession := currentSession.Values["userID"]
	userID, ok := userIDFromSession.(string)

	if !ok {
		return "", err
	}

	log.Println(userID)
	return userID, nil
}

func IsUserLoggedIn(c echo.Context, sessionHandler *redistore.RediStore) bool {
	_, err := GetUserIDFromSession(c, sessionHandler)
	return err == nil
}
