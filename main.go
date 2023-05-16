package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"cloud.google.com/go/storage"
	"golang.org/x/oauth2"
	"google.golang.org/api/iamcredentials/v1"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

// ImpersonatedTokenSource is a TokenSource implementation that generates an impersonated access token.
type ImpersonatedTokenSource struct {
	ctx              context.Context
	targetServiceAcc string
}

// Token returns a new token by generating an impersonated access token.
func (ts *ImpersonatedTokenSource) Token() (*oauth2.Token, error) {
	req := &iamcredentials.GenerateAccessTokenRequest{
		Lifetime: "3600s",
		Scope:    []string{"https://www.googleapis.com/auth/cloud-platform"},
	}

	iamClient, err := iamcredentials.NewService(ts.ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create IAM client: %v", err)
	}

	resp, err := iamClient.Projects.ServiceAccounts.GenerateAccessToken(fmt.Sprintf("projects/-/serviceAccounts/%s", ts.targetServiceAcc), req).Do()
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %v", err)
	}

	token := &oauth2.Token{
		AccessToken: resp.AccessToken,
		TokenType:   "Bearer",
	}

	return token, nil
}

func main() {

	ctx := context.Background()

	targetServiceAcc := os.Getenv("TARGET_SERVICE_ACCOUNT")
	if targetServiceAcc == "" {
		log.Fatal("TARGET_SERVICE_ACCOUNT environment variable is not set")
	}

	ts := &ImpersonatedTokenSource{
		ctx:              ctx,
		targetServiceAcc: targetServiceAcc,
	}

	// Use the ImpersonatedTokenSource to create a Storage client
	storageClient, err := storage.NewClient(ctx, option.WithTokenSource(ts))
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

}
