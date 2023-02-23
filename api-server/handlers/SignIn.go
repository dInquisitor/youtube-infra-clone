package handlers

import (
	"api-server/session"
	"database/sql"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/boj/redistore.v1"
)

type UserData struct {
	Email string `json:"email"`
}

type SignInResponse struct {
	Status string   `json:"status"`
	Data   UserData `json:"userData,omitempty"`
}

type User struct {
	id       string
	password string
}

func SignIn(c echo.Context, dbPool *sql.DB, sessionHandler *redistore.RediStore) error {
	const (
		STATUS_LOGIN_INFO_INCORRECT = "LOGIN_INFO_INCORRECT"
		STATUS_UNKNOWN_ERROR        = "UNKNOWN_ERROR"
		STATUS_SUCCESS              = "SUCCESS"
	)

	inputEmail := c.FormValue("email")
	inputPassword := c.FormValue("password")

	// If user exists, log them in, else sign them up and log them in (will open up to enumeration attacks but is sufficient for now)
	// May separate log in and sign up later
	userQuery, err := dbPool.Query("select id, password from users where email = $1", inputEmail)

	if err != nil {
		log.Println("Could not query database while signing user in", "->", c.Path())
		return c.JSON(http.StatusOK, &SignInResponse{
			Status: STATUS_UNKNOWN_ERROR,
		})
	}

	userExists := userQuery.Next()

	if userExists {
		var user User
		if err := userQuery.Scan(&user.id, &user.password); err != nil {
			log.Println(err, "->", c.Path())
			return c.JSON(http.StatusOK, &SignInResponse{
				Status: STATUS_UNKNOWN_ERROR,
			})
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.password), []byte(inputPassword)); err == nil {
			if err := session.SetUserSession(sessionHandler, c, "userID", user.id); err != nil {
				log.Println(err)
				return c.JSON(http.StatusOK, &SignInResponse{
					Status: STATUS_UNKNOWN_ERROR,
				})
			}

			return c.JSON(http.StatusOK, &SignInResponse{
				Status: STATUS_SUCCESS,
				Data: UserData{
					Email: inputEmail,
				},
			})
		}

		return c.JSON(http.StatusOK, &SignInResponse{
			Status: STATUS_LOGIN_INFO_INCORRECT,
		})
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(inputPassword), bcrypt.MinCost)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusOK, &SignInResponse{
			Status: STATUS_UNKNOWN_ERROR,
		})
	}

	insertUserResult := dbPool.QueryRow("insert into users (email, password) values($1, $2) returning id", inputEmail, passwordHash)

	var newUserID string
	err = insertUserResult.Scan(&newUserID)

	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusOK, &SignInResponse{
			Status: STATUS_UNKNOWN_ERROR,
		})
	}

	if err := session.SetUserSession(sessionHandler, c, "userID", newUserID); err != nil {
		log.Println(err)
		return c.JSON(http.StatusOK, &SignInResponse{
			Status: STATUS_UNKNOWN_ERROR,
		})
	}

	return c.JSON(http.StatusOK, &SignInResponse{
		Status: STATUS_SUCCESS,
		Data: UserData{
			Email: inputEmail,
		},
	})
}
