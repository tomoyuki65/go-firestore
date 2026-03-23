package database

import (
	"context"
	"os"

	"cloud.google.com/go/firestore"
)

func NewFirestoreClient(ctx context.Context) (*firestore.Client, error) {
	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT_ID")

	return firestore.NewClient(ctx, projectID)
}
