package handlers

import (
	"encoding/json"
	"html/template"
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
		// Get the token from cookie
		cookie, err := r.Cookie("token")
		if err != nil {
			if err == http.ErrNoCookie {
				http.Error(w, "No token provided", http.StatusUnauthorized)
				return
			}
			http.Error(w, "Failed to get token from cookie", http.StatusBadRequest)
			return
		}
		token := cookie.Value
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
			w.Header().Set("Set-Cookie", "token="+token+"; Path=/; HttpOnly")
			http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
		}
	}
	http.Error(w, "Access denied", http.StatusForbidden)
}

type DashboardPage struct {
	Title string
	ShortLinks map[string]string
}

// Dashboard handles the dashboard request.
func (p DashboardPage) Dashboard(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.New("dashboard.html").ParseFiles("templates/dashboard.html"))
	datastoreShortLinks, err := datastore.ListURLs()
	if err != nil {
		http.Error(w, "Failed to get short links", http.StatusInternalServerError)
		return
	}
	p.ShortLinks = datastoreShortLinks
	tmpl.Execute(w, p)
}

// CreateShortLink creates a short link for a given destination URL and sends the short link in response.
func CreateShortLink(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	shortlink := r.Form.Get("shortlink")
	destinationURL := r.Form.Get("destination_url")
	returnedShortURL, err := datastore.AddURL(shortlink, destinationURL)
	if err != nil {
		http.Error(w, "Failed to create short link", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"short_url": returnedShortURL,
	})
}

// EditShortLink changes the destination URL of a given short link or ShortURL.
func EditShortLink(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	shortURL := r.Form.Get("shortlink")
	newShortURL := r.Form.Get("newShortlink")
	destinationURL, err := datastore.GetURL(shortURL)
	if err != nil {
		http.Error(w, "URL not found", http.StatusNotFound)
		return
	}
	newDestinationURL := r.Form.Get("newDestinationURL")
	err = datastore.EditURL(shortURL, destinationURL, newShortURL, newDestinationURL)
	if err != nil {
		http.Error(w, "Failed to edit short link", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"short_url": shortURL,
	})
}

// RemoveShortLink removes a given short link or ShortURL.
func RemoveShortLink(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	shortURL := r.Form.Get("shortlink")
	err = datastore.RemoveURL(shortURL)
	if err != nil {
		http.Error(w, "Failed to remove short link", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusFound)
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