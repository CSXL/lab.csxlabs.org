package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/CSXL/lab.csxlabs.org/shortlinks/auth"
	"github.com/CSXL/lab.csxlabs.org/shortlinks/config"
	"github.com/CSXL/lab.csxlabs.org/shortlinks/handlers"
)

func main() {
	config := config.LoadConfig()
	authorizer := auth.NewAuthorizer(config.Authentication.SigningKey)
	VerifyJWT := handlers.VerifyJWT{
		Authorizer: authorizer,
	}
	LoginPage := handlers.LoginPage{
		Title: config.Website.Name,
		AllowedUsers: config.Authentication.AllowedUsers,
		Authorizer: authorizer,
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

	fmt.Println("Short links server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}