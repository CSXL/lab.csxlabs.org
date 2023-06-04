package datastore

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
	"github.com/CSXL/lab.csxlabs.org/shortlinks/config"
	"github.com/google/uuid"
	"google.golang.org/api/iterator"
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

// ListURLs returns a map of short link paths to destination URLs from Firestore.
func ListURLs() (map[string]string, error) {
	iter := firestoreClient.Collection(CollectionName).Documents(context.Background())
	urls := make(map[string]string)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to iterate over Firestore documents: %w", err)
		}
		urls[doc.Ref.ID] = doc.Data()["destination_url"].(string)
	}

	return urls, nil
}

// RemoveURL removes the short link associated with the given short link path from Firestore.
func RemoveURL(shortURL string) error {
	_, err := firestoreClient.Collection(CollectionName).Doc(shortURL).Delete(context.Background())
	if err != nil {
		return fmt.Errorf("failed to remove short link from Firestore: %w", err)
	}

	return nil
}

// EditURL updates the destination URL associated with the given short link in Firestore.
func EditURL(shortURL string, destinationURL string, newShortURL string, newDestinationURL string) error {
	if newShortURL == "" && newDestinationURL == "" {
		return fmt.Errorf("no new values provided")
	}
	if newDestinationURL != "" {
		_, err := firestoreClient.Collection(CollectionName).Doc(shortURL).Update(context.Background(), []firestore.Update{
			{
				Path:  "destination_url",
				Value: newDestinationURL,
			},
		})
		if err != nil {
			return fmt.Errorf("failed to update destination URL in Firestore: %w", err)
		}
		destinationURL = newDestinationURL
	}
	if newShortURL != "" {
		err := RemoveURL(shortURL)
		if err != nil {
			return fmt.Errorf("failed to remove old short link from Firestore: %w", err)
		}
		_, err = AddURL(newShortURL, destinationURL)
		if err != nil {
			return fmt.Errorf("failed to add new short link to Firestore: %w", err)
		}
	}
	return nil
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