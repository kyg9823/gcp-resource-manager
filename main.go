package main

import (
	"log"
	"net/http"
	"os"

	"github.com/kyg9823/gcp-resource-manager/api"
)

func main() {
	http.HandleFunc("/healthcheck", api.Healthcheck)
	http.HandleFunc("/gce", api.GceManager)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
