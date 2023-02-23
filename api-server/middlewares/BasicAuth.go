package middlewares

import (
	"api-server/session"
	"net/http"

	"github.com/labstack/echo/v4"
	"gopkg.in/boj/redistore.v1"
)

type AuthContext struct {
	UserID string
}

func BasicAuth(next echo.HandlerFunc, sessionHandler *redistore.RediStore) echo.HandlerFunc {
	return func(c echo.Context) error {
		isUserLoggedIn := session.IsUserLoggedIn(c, sessionHandler)
		if !isUserLoggedIn {
			return c.JSON(http.StatusUnauthorized, nil)
		}

		userID, _ := session.GetUserIDFromSession(c, sessionHandler)

		c.Set("authContext", &AuthContext{userID})

		return next(c)
	}
}
