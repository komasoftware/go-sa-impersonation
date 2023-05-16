package main

import (
	"context"
	"fmt"
	"log"

	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

func listObjectsWithServiceAccountKey(serviceAccountKeyPath, bucketName string) {
	ctx := context.Background()

	// Authenticate with the service account key
	client, err := storage.NewClient(ctx, option.WithCredentialsFile(serviceAccountKeyPath))
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// List objects in the bucket
	it := client.Bucket(bucketName).Objects(ctx, nil)
	for {
		objAttrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("Failed to list objects: %v", err)
		}
		fmt.Println(objAttrs.Name)
	}
}

func main() {
	serviceAccountKeyPath := "path/to/service-account-key.json"
	bucketName := "your-bucket-name"
	listObjectsWithServiceAccountKey(serviceAccountKeyPath, bucketName)
}
