package main

import (
	"context"
	"fmt"
	"log"

	"cloud.google.com/go/storage"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/iamcredentials/v1"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

func main() {
	ctx := context.Background()

	// Use the Application Default Credentials to get an access token
	creds, err := google.FindDefaultCredentials(ctx, iamcredentials.CloudPlatformScope)
	if err != nil {
		log.Fatalf("Failed to find default credentials: %v", err)
	}

	token, err := creds.TokenSource.Token()
	if err != nil {
		log.Fatalf("Failed to find default credentials: %v", err)
	}
	log.Printf("got token of type: %s", token.TokenType)
	tokenSource := oauth2.StaticTokenSource(token)

	authorizedClient := oauth2.NewClient(ctx, tokenSource)

	// Create the IAM Credentials API client
	iamClient, err := iamcredentials.NewService(ctx, option.WithHTTPClient(authorizedClient))
	if err != nil {
		log.Fatalf("Failed to create IAM client: %v", err)
	}

	// Replace with your target service account and desired lifetime
	targetSA := "target-service-account@koen-gcompany-demo.iam.gserviceaccount.com"
	lifetime := "3600s"

	// Create the request to generate an access token for the target service account
	req := &iamcredentials.GenerateAccessTokenRequest{
		Lifetime: lifetime,
		Scope:    []string{"https://www.googleapis.com/auth/cloud-platform"},
	}

	// Generate an access token for the target service account
	resp, err := iamClient.Projects.ServiceAccounts.GenerateAccessToken(fmt.Sprintf("projects/-/serviceAccounts/%s", targetSA), req).Do()
	if err != nil {
		log.Fatalf("Failed to generate access token: %v", err)
	}

	// Use the access token to create a storage client and list objects in a bucket
	storageClient, err := storage.NewClient(ctx, option.WithCredentials(creds))
	if err != nil {
		log.Fatalf("Failed to create storage client: %v", err)
	}
	defer storageClient.Close()

	bucket := "core-infra-demo-bucket"
	it := storageClient.Bucket(bucket).Objects(ctx, nil)
	for {
		objAttrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("Failed to get object: %v", err)
		}
		fmt.Println(objAttrs.Name)
	}

	fmt.Println(resp.AccessToken)
}
