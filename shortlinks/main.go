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
	http.HandleFunc(config.ReservedManagementEndpoints.Login, LoginPage.Login)
	http.HandleFunc("/create", VerifyJWT.VerifyJWT(handlers.CreateShortLink))
	http.HandleFunc("/", handlers.RedirectToDestinationURL)

	fmt.Println("Short links server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}