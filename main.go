package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
)

const MaxContentLength = 1024 * 1024 // 1MB

func main() {
	http.HandleFunc("/", okHandler)
	http.HandleFunc("/expect", expectHandler)
	if err := http.ListenAndServe(":"+os.Getenv("PORT"), nil); err != nil {
		fmt.Printf("error: server exited: %s\n", err.Error())
	}
}

func okHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	_, err := fmt.Fprint(w, "ok")
	if err != nil {
		fmt.Printf("error: handle request: %s\n", err.Error())
	}
}

func expectHandler(w http.ResponseWriter, r *http.Request) {
	contentLength := 0
	if r.Header.Get("Content-Length") != "" {
		var err error
		contentLength, err = strconv.Atoi(r.Header.Get("Content-Length"))
		if err != nil {
			http.Error(w, "Invalid Content-Length", http.StatusBadRequest)
			return
		}
	}

	expectedHeader := r.Header.Get("Expect")
	if expectedHeader != "100-continue" {
		http.Error(w, "Missing Expected Header: "+r.Header.Get("Expect"), http.StatusExpectationFailed)
		return
	}

	// An example to work with 100-continue, restrict the content-contentLength to some value
	if contentLength > MaxContentLength {
		http.Error(w, "Content size is too large", http.StatusRequestEntityTooLarge)
		return
	}

	// Send 100 Continue response
	w.WriteHeader(http.StatusContinue)

	buf := make([]byte, contentLength)
	if _, err := io.ReadFull(r.Body, buf); err != nil {
		http.Error(w, "Error reading body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "%d bytes received and read\n", len(buf))
	fmt.Fprintln(w, "Created")
}
