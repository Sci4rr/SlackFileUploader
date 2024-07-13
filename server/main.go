package main

import (
    "fmt"
    "log"
    "net/http"
    "os"

    "github.com/joho/godotenv"
)

func init() {
    // Load environment variables from .env file if present
    if err := godotenv.Load(); err != nil {
        log.Println("No .env file found, continuing with environment variables")
    }
}

func main() {
    serverPort := os.Getenv("PORT")
    if serverPort == "" {
        log.Fatal("PORT must be set")
    }

    slackAPIKey := os.Getenv("SLACK_TOKEN")
    if slackAPIKey == "" {
        log.Fatal("SLACK_TOKEN must be set")
    }

    configureHTTPRoutes()

    log.Printf("Starting HTTP server on port %s\n", serverPort)
    if err := http.ListenAndServe(fmt.Sprintf(":%s", serverPort), nil); err != nil {
        log.Fatalf("Failed to start HTTP server: %v", err)
    }
}

// configureHTTPRoutes setups the URL routes for the HTTP server
func configureHTTPRoutes() {
    http.HandleFunc("/upload", handleSlackFileUpload)

    http.HandleFunc("/health", healthCheckHandler)
}

// handleSlackFileUpload processes file upload requests to Slack
func handleSlackFileUpload(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "Slack File Upload Endpoint")
}

// healthCheckHandler responds to health check requests
func healthCheckHandler(w http.ResponseWriter, r *helpRequest) {
    w.WriteHeader(http.StatusOK)
    fmt.Fprintln(w, "OK")
}