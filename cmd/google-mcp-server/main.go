// google-mcp-server is an MCP server for reading Google Slides presentations.
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/grokify/google-mcp-server/internal/auth"
	"github.com/grokify/google-mcp-server/internal/slides"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const (
	serverName    = "google-mcp-server"
	serverVersion = "v0.1.0"
)

func main() {
	var (
		credentialsFile       string
		goauthCredentialsFile string
		goauthCredentialsKey  string
	)

	flag.StringVar(&credentialsFile, "credentials", os.Getenv("GOOGLE_CREDENTIALS_FILE"),
		"Path to Google service account credentials JSON file")
	flag.StringVar(&goauthCredentialsFile, "goauth-credentials-file", os.Getenv("GOAUTH_CREDENTIALS_FILE"),
		"Path to goauth CredentialsSet JSON file")
	flag.StringVar(&goauthCredentialsKey, "goauth-credentials-account", os.Getenv("GOAUTH_CREDENTIALS_ACCOUNT"),
		"Account key within goauth CredentialsSet file")
	flag.Parse()

	// Validate flags
	hasGoogleCreds := credentialsFile != ""
	hasGoauthCreds := goauthCredentialsFile != "" && goauthCredentialsKey != ""

	if !hasGoogleCreds && !hasGoauthCreds {
		fmt.Fprintln(os.Stderr, "Error: credentials required")
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintln(os.Stderr, "Option 1: Google service account credentials")
		fmt.Fprintln(os.Stderr, "  google-mcp-server -credentials /path/to/service-account.json")
		fmt.Fprintln(os.Stderr, "  or set GOOGLE_CREDENTIALS_FILE environment variable")
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintln(os.Stderr, "Option 2: goauth CredentialsSet")
		fmt.Fprintln(os.Stderr, "  google-mcp-server -goauth-credentials-file /path/to/credentials.json -goauth-credentials-account myaccount")
		fmt.Fprintln(os.Stderr, "  or set GOAUTH_CREDENTIALS_FILE and GOAUTH_CREDENTIALS_ACCOUNT environment variables")
		os.Exit(1)
	}

	if hasGoogleCreds && hasGoauthCreds {
		fmt.Fprintln(os.Stderr, "Error: cannot use both -credentials and -goauth-credentials-file")
		os.Exit(1)
	}

	ctx := context.Background()

	// Create authenticated HTTP client
	var httpClient *http.Client
	var err error

	if hasGoauthCreds {
		httpClient, err = auth.NewClientFromCredentialsSet(ctx, goauthCredentialsFile, goauthCredentialsKey)
	} else {
		httpClient, err = auth.NewClient(ctx, credentialsFile)
	}
	if err != nil {
		log.Fatalf("Failed to create authenticated client: %v", err)
	}

	// Create Slides service
	slidesService, err := slides.NewService(ctx, httpClient)
	if err != nil {
		log.Fatalf("Failed to create Slides service: %v", err)
	}

	// Create MCP server
	server := mcp.NewServer(&mcp.Implementation{
		Name:    serverName,
		Version: serverVersion,
	}, nil)

	// Register tools
	tools := slides.NewTools(slidesService)
	slides.RegisterTools(server, tools)

	// Run server with stdio transport
	if err := server.Run(ctx, &mcp.StdioTransport{}); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
