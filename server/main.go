package main

import (
    "fmt"
    "log"
    "net/http"
    "os"
    "time"

    "github.com/joho/godotenv"
)

func init() {
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
    if err := http.ListenAndServe(fmt.Sprintf(":%s", serverPort), logRequest(http.DefaultServeMux)); err != nil {
        log.Fatalf("Failed to start HTTP server: %v", err)
    }
}

func configureHTTPRoutes() {
    http.HandleFunc("/upload", handleSlmackFileUpload)

    http.HandleFunc("/health", healthCheckHandler)
}

func handleSlackFileUpload(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "Slack File Upload Endpoint")
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    fmt.Fprintln(w, "OK")
}

func logRequest(handler http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        handler.ServeHTTP(w, r)

        log.Printf(
            "%s %s %s %s",
            r.RemoteAddr,
            r.Method,
            r.URL.EscapedPath(),
            time.Since(start),
        )
    })
}