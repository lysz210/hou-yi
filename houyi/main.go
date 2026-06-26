package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/coreos/go-systemd/v22/login1"
)

type ResponseStatus string

const (
	StatusOk ResponseStatus = "OK"
	StatusError   ResponseStatus = "ERROR"
)

type Response struct {
	Status  ResponseStatus `json:"status"`
	Message string         `json:"message,omitempty"`
}

func writeJSONResponse(w http.ResponseWriter, status ResponseStatus, message string) {
	response := Response{
		Status:  status,
		Message: message,
	}
	var statusCode int
	if status == StatusOk {
		statusCode = http.StatusOK
	} else {
		statusCode = http.StatusInternalServerError
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}

func shutdownHandler(w http.ResponseWriter, r *http.Request) {

	// 1. Initialize the managed systemd logind connection
	conn, err := login1.New()
	if err != nil {
		writeJSONResponse(w, StatusError, "Failed to connect to systemd logind: "+err.Error())
		return
	}
	defer conn.Close() // Automatically clean up the connection when done

	log.Println("Hou Yi is drawing the bow... Sending systemd PowerOff request.")

	// 2. Trigger the power off sequence cleanly.
	// The parameter 'false' tells systemd not to prompt the user interactively for policy-kit authorization.
	conn.PowerOff(false)

	writeJSONResponse(w, StatusOk, "System is shutting down gracefully...")
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /shutdown", shutdownHandler)

	log.Println("Hou Yi API listening on :8080...")

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	log.Fatal(server.ListenAndServe())
}