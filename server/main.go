package main

import (
    "bytes"
    "fmt"
    "io"
    "io/ioutil"
    "log"
    "mime/multipart"
    "net/http"
    "os"
    "strings"
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
    http.HandleFunc("/upload", handleSlackFileUpload)

    http.HandleFunc("/health", healthCheckHandler)
}

func handleSlackFileUpload(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Unsupported method. Please use POST.", http.StatusMethodNotAllowed)
        return
    }
    
    err := r.ParseMultipartForm(10 << 20) // Limit upload size
    if err != nil {
        http.Error(w, "Could not parse multipart form: "+err.Error(), http.StatusBadRequest)
        return
    }
    
    file, fileHeader, err := r.FormFile("file")
    if err != nil {
        http.Error(w, "Could not get uploaded file: "+err.Error(), http.StatusBadRequest)
        return
    }
    defer file.Close()

    slackChannel := r.FormValue("channel")
    if slackChannel == "" {
        slackChannel = "default-channel" // Set to your default Slack channel
    }

    token := os.Getenv("SLACK_TOKEN")
    buffer := &bytes.Buffer{}
    writer := multipart.NewWriter(buffer)

    part, err := writer.CreateFormFile("file", fileHeader.Filename)
    if err != nil {
        http.Error(w, "Could not create form file: "+err.Error(), http.StatusInternalServerError)
        return
    }

    _, err = io.Copy(part, file)
    if err != nil {
        http.Error(w, "Could not copy file: "+err.Error(), http.StatusInternalServerError)
        return
    }

    writer.WriteField("channels", slackChannel)
    writer.WriteField("token", token)
    writer.Close()

    request, err := http.NewRequest("POST", "https://slack.com/api/files.upload", buffer)
    if err != nil {
        http.Error(w, "Could not create request: "+err.Error(), http.StatusInternalServerError)
        return
    }
    
    request.Header.Set("Authorization", "Bearer "+token)
    request.Header.Set("Content-Type", writer.FormDataContentType())

    client := &http.Client{}
    response, err := client.Do(request)
    if err != nil {
        http.Error(w, "Error uploading file: "+err.Error(), http.StatusInternalServerError)
        return
    }
    defer response.Body.Close()
    
    body, err := ioutil.ReadAll(response.Body)
    if err != nil {
        http.Error(w, "Could not read response body: "+err.Error(), http.StatusInternalServerError)
        return
    }

    if !strings.Contains(string(body), "ok") {
        http.Error(w, "Error from Slack API: "+string(body), http.StatusInternalServerError)
        return
    }
    
    fmt.Fprintln(w, "File uploaded successfully")
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