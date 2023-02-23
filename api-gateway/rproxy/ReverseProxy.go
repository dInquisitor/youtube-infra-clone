package rproxy

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// Adapted from https://itnext.io/why-should-you-write-your-own-api-gateway-from-scratch-378074bfc49e

func ReverseProxy(c echo.Context) error {
	path := c.Request().URL.Path
	target, err := Target(path)

	log.Println("Final path", target)

	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusNotFound, nil)
	}

	if targetUrl, err := url.Parse(target); err != nil {
		return c.JSON(http.StatusInternalServerError, nil)
	} else {
		Proxy(targetUrl).ServeHTTP(c.Response().Writer, c.Request())
	}
	return nil
}

func Target(path string) (string, error) {
	parts := strings.Split(strings.TrimPrefix(path, "/"), "/")

	if len(parts) < 2 {
		return "", fmt.Errorf("path must have at least two parts. found %d", len(parts))
	}

	if parts[0] != "api" {
		return "", errors.New("path has to begin with 'api'")
	}

	var service string
	var port int

	if parts[1] == "upload" {
		service = "upload-service"
		port = 5000
	}

	if parts[1] == "stream" {
		service = "stream-service"
		port = 7000
	}

	if service == "" {
		service = "api-server"
		port = 4000
	}

	return fmt.Sprintf(
		"http://%s:%d/%s",
		service, port, strings.Join(parts, "/")), nil
}

func Proxy(address *url.URL) *httputil.ReverseProxy {
	p := httputil.NewSingleHostReverseProxy(address)
	p.Director = func(request *http.Request) {
		request.Host = address.Host
		request.URL.Scheme = address.Scheme
		request.URL.Host = address.Host
		request.URL.Path = address.Path
	}
	p.ModifyResponse = func(response *http.Response) error {
		if response.StatusCode == http.StatusInternalServerError {
			u, s := readBody(response)
			log.Printf("%s ,req %s ,with error %d, body:%s", u.String(), address, response.StatusCode, s)
			response.Body = io.NopCloser(bytes.NewReader([]byte(fmt.Sprintf("error %s", u.String()))))
		} else if response.StatusCode > 300 {
			_, s := readBody(response)
			log.Printf("req %s ,with error %d, body:%s", address, response.StatusCode, s)
			response.Body = io.NopCloser(bytes.NewReader([]byte(s)))
		}
		return nil
	}
	return p
}

func readBody(response *http.Response) (uuid.UUID, string) {
	defer response.Body.Close()
	all, _ := io.ReadAll(response.Body)
	u := uuid.New()
	var s string
	if len(all) > 0 {
		s = string(all)
	}
	return u, s
}
