package main

import (
	"log"
	"net/http"

	"github.com/coreos/go-systemd/v22/login1"
)

func shutdownHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 1. Initialize the managed systemd logind connection
	conn, err := login1.New()
	if err != nil {
		http.Error(w, "Failed to connect to systemd logind: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer conn.Close() // Automatically clean up the connection when done

	log.Println("Hou Yi is drawing the bow... Sending systemd PowerOff request.")

	// 2. Trigger the power off sequence cleanly.
	// The parameter 'false' tells systemd not to prompt the user interactively for policy-kit authorization.
	conn.PowerOff(false)

	w.Write([]byte("System is shutting down gracefully..."))
}

func main() {
	http.HandleFunc("/shutdown", shutdownHandler)
	
	log.Println("Hou Yi API listening on :8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}