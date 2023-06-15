package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/CSXL/lab.csxlabs.org/shortlinks/auth"
	"github.com/CSXL/lab.csxlabs.org/shortlinks/config"
	"github.com/CSXL/lab.csxlabs.org/shortlinks/handlers"
)

func main() {
	config := config.LoadConfig()
	authorizer := auth.NewAuthorizer(config.Authentication.SigningKey)
	VerifyJWT := handlers.VerifyJWT{
		Authorizer:                  authorizer,
		ReservedManagementEndpoints: config.ReservedManagementEndpoints,
	}
	LoginPage := handlers.LoginPage{
		Title:        config.Website.Name,
		AllowedUsers: config.Authentication.AllowedUsers,
		Authorizer:   authorizer,
		DashboardURL: config.ReservedManagementEndpoints.Dashboard,
	}
	Dashboard := handlers.DashboardPage{
		Title: config.Website.Name,
	}
	http.HandleFunc(config.ReservedManagementEndpoints.Login, LoginPage.Login)
	http.HandleFunc(config.ReservedManagementEndpoints.Dashboard, VerifyJWT.VerifyJWT(Dashboard.Dashboard))
	http.HandleFunc("/create", VerifyJWT.VerifyJWT(handlers.CreateShortLink))
	http.HandleFunc("/delete", VerifyJWT.VerifyJWT(handlers.RemoveShortLink))
	http.HandleFunc("/edit", VerifyJWT.VerifyJWT(handlers.EditShortLink))
	http.HandleFunc("/", handlers.RedirectToDestinationURL)

	// PORT environment variable is provided by Cloud Run.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Println("Short links server started on :" + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
