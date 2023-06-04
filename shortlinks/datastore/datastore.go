package datastore

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
	"github.com/CSXL/lab.csxlabs.org/shortlinks/config"
	"github.com/google/uuid"
)

var ProjectID string
var CollectionName string

var firestoreClient *firestore.Client

func init() {
	config := config.LoadConfig()
	ProjectID = config.Firestore.ProjectID
	CollectionName = config.Firestore.CollectionName
	client, err := firestore.NewClient(context.Background(), ProjectID)
	if err != nil {
		panic(fmt.Errorf("failed to instantiate Firestore client: %w", err))
	}
	firestoreClient = client
}

// AddURL creates a short link (generate if the short link is "") for the given destination URL and stores it in Firestore.
func AddURL(shortURL string, destinationURL string) (string, error) {
	if shortURL == "" {
		shortURL = uuid.NewString()
	}
	_, err := firestoreClient.Collection(CollectionName).Doc(shortURL).Set(context.Background(), map[string]string{
		"destination_url": destinationURL,
	})

	if err != nil {
		return "", fmt.Errorf("failed to store short link in Firestore: %w", err)
	}

	return shortURL, nil
}

// GetURL fetches the destination URL associated with the given short link from Firestore.
func GetURL(shortURL string) (string, error) {
	doc, err := firestoreClient.Collection(CollectionName).Doc(shortURL).Get(context.Background())
	if err != nil {
		return "", fmt.Errorf("failed to fetch URL from Firestore: %w", err)
	}

	destinationURL, err := doc.DataAt("destination_url")
	if err != nil {
		return "", fmt.Errorf("failed to retrieve destination URL from document: %w", err)
	}

	return destinationURL.(string), nil
}