package handlers

import (
	"encoding/json"
	"html/template"
	"io"
	"net/http"

	"github.com/CSXL/lab.csxlabs.org/shortlinks/auth"
	"github.com/CSXL/lab.csxlabs.org/shortlinks/config"
	"github.com/CSXL/lab.csxlabs.org/shortlinks/datastore"
)

type VerifyJWT struct {
	Authorizer *auth.Authorizer
}

func (v VerifyJWT) VerifyJWT(endpointHandler http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Invalid request method", http.StatusBadRequest)
			return
		}
		// Extract the token from the request body
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
		token, ok := payload["token"]
		if !ok {
			http.Error(w, "Missing token", http.StatusBadRequest)
			return
		}
		// Validate the token
		valid, err := v.Authorizer.ValidateToken(token)
		if err != nil {
			http.Error(w, "Failed to validate token", http.StatusBadRequest)
			return
		}
		if !valid {
			http.Error(w, "Invalid token", http.StatusBadRequest)
			return
		}
		endpointHandler(w, r)
	})
}

type LoginPage struct {
	Title string
	AllowedUsers []config.AllowedUser
	Authorizer *auth.Authorizer
}

// Login handles the login request.
func (p LoginPage) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		tmpl := template.Must(template.New("login.html").ParseFiles("templates/login.html"))
		tmpl.Execute(w, p)
		return
	}
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	username := r.Form.Get("username")
	password := r.Form.Get("password")
	// Check if the user is allowed to login
	for _, user := range p.AllowedUsers {
		if user.Username == username && user.Password == password {
			token, err := p.Authorizer.GenerateToken()
			if err != nil {
				http.Error(w, "Failed to generate token", http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("Token", token)
			return
		}
	}
	http.Error(w, "Access denied", http.StatusForbidden)
}

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
	if path == "/" {
		return ""
	}
	if path[0] == '/' {
		path = path[1:]
	}
	if path[len(path)-1] == '/' {
		path = path[:len(path)-1]
	}
	return path
}