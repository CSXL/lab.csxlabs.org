package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/CSXL/lab.csxlabs.org/shortlinks/handlers"
)

func main() {
	http.HandleFunc("/create", handlers.CreateShortLink)
	http.HandleFunc("/", handlers.RedirectToDestinationURL)

	fmt.Println("Short links server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}