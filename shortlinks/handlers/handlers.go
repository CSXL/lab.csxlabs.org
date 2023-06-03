package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/CSXL/lab.csxlabs.org/shortlinks/datastore"
)

// CreateShortLink creates a short link for a given destination URL and sends the short link in response.
func CreateShortLink(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	var payload map[string]string
	err = json.Unmarshal(body, &payload)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	destinationURL, found := payload["destination_url"]
	if !found {
		http.Error(w, "URL not found in request payload", http.StatusBadRequest)
		return
	}

	shortURL, found := payload["short_url"]
	if !found {
		shortURL = ""
	}

	returnedShortURL, err := datastore.AddURL(shortURL, destinationURL)
	if err != nil {
		http.Error(w, "Failed to create short link", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"short_url": returnedShortURL,
	})
}

// RedirectToDestinationURL redirects the user to the original destination URL associated with the given short link path.
func RedirectToDestinationURL(w http.ResponseWriter, r *http.Request) {
	shortURL := serializePath(r.URL.Path)
	destinationURL, err := datastore.GetURL(shortURL)
	if err != nil {
		http.Error(w, "URL not found", http.StatusNotFound)
		return
	}

	http.Redirect(w, r, destinationURL, http.StatusFound)
}

// serializePath returns the serialized relative path of the given request.
func serializePath(path string) string {
	if path[0] == '/' {
		path = path[1:]
	}
	if path[len(path)-1] == '/' {
		path = path[:len(path)-1]
	}
	return path
}