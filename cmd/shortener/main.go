package main

import (
	"flag"
	"github.com/gin-gonic/gin"
	"io"
	"math/rand/v2"
	"net/http"
	"net/url"
	"strconv"
)

var ServerAddressShort *string
var ServerAddressLong *string

var shortToOriginal = make(map[string]string)

func getting(c *gin.Context) {
	id := c.Param("id")
	orig, ok := shortToOriginal[id]
	if !ok {
		c.String(http.StatusNotFound, "Not Found")
		return
	}
	c.Redirect(http.StatusTemporaryRedirect, orig)
}

func posting(c *gin.Context) {
	mediaType := c.ContentType()
	if mediaType != "text/plain" {
		c.String(http.StatusUnsupportedMediaType, "Content-Type not supported")
		return
	}

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to read body")
		return
	}

	rawURL := string(body)
	_, err = url.ParseRequestURI(rawURL)
	if err != nil {
		c.String(http.StatusInternalServerError, "Invalid URL")
		return
	}

	newID := strconv.FormatInt(rand.Int64(), 36)
	shortToOriginal[newID] = rawURL

	c.String(http.StatusCreated, *ServerAddressShort+newID)
}

func main() {
	ServerAddressLong = flag.String("a", "localhost:8080", "HTTP server address")
	ServerAddressShort = flag.String("b", "http://localhost:8080/", "Base URL for short links")
	flag.Parse()

	router := gin.Default()

	router.HandleMethodNotAllowed = true
	router.NoMethod(func(c *gin.Context) {
		c.String(http.StatusMethodNotAllowed, "Method not allowed")
	})

	router.GET("/:id", getting)
	router.POST("/", posting)

	router.Run(*ServerAddressLong)
}
