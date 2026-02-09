// Package auth provides Google service account authentication.
package auth

import (
	"context"
	"net/http"

	"github.com/grokify/goauth"
	"github.com/grokify/goauth/google"
	slides "google.golang.org/api/slides/v1"
)

// Scopes returns the OAuth2 scopes required for read-only access to Google Slides.
func Scopes() []string {
	return []string{
		slides.PresentationsReadonlyScope,
		slides.DriveReadonlyScope,
	}
}

// NewClient creates an authenticated HTTP client using a Google service account credentials file.
// This uses the standard Google Cloud service account JSON format.
func NewClient(ctx context.Context, credentialsFile string) (*http.Client, error) {
	return google.NewClientSvcAccountFromFile(ctx, credentialsFile, Scopes()...)
}

// NewClientFromCredentialsSet creates an authenticated HTTP client using a goauth CredentialsSet file.
// The CredentialsSet should contain a credential entry with the specified account key.
// The credential entry should be of type "gcpsa" with appropriate scopes configured.
func NewClientFromCredentialsSet(ctx context.Context, credentialsFile, accountKey string) (*http.Client, error) {
	return goauth.NewClient(ctx, credentialsFile, accountKey)
}
