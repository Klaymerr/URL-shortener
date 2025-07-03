package main

import (
	"fmt"
	"io"
	"math/rand/v2"
	"mime"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

var shortToOriginal = make(map[string]string)

func handler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		fmt.Println(r.Header.Get("Content-Type"))
		mediaType, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
		if err != nil || mediaType != "text/plain" {
			http.Error(w, "Content-Type not supported", http.StatusUnsupportedMediaType)
			return
		}

		text, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read body", http.StatusInternalServerError)
			return
		}

		_, err = url.ParseRequestURI(string(text))
		if err != nil {
			http.Error(w, "Invalid URL", http.StatusInternalServerError)
			return
		}
		fmt.Println(string(text))
		newURL := strconv.FormatInt(rand.Int64(), 36)
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("http://localhost:8080/" + newURL))
		shortToOriginal[newURL] = string(text)
		return
	}

	if r.Method == http.MethodGet {
		id := strings.TrimPrefix(r.URL.Path, "/")
		orig, ok := shortToOriginal[id]
		if !ok {
			http.Error(w, "Not Found", http.StatusNotFound)
			return
		}
		http.Redirect(w, r, orig, http.StatusMovedPermanently)
		return
	}

	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc(`/`, handler)

	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}
}
