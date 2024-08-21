package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"sync"
)

const (
	SlackUploadAPI = "https://slack.com/api/files.upload"
)

var (
	slackToken     = os.Getenv("SLACK_BOT_TOKEN")
	slackChannelID = os.Getenv("SLACK_CHANNEL_ID")
)

type FileInfo struct {
	FileName        string
	FileType        string
	Title           string
	InitialComment  string
}

type SlackResponse struct {
	OK    bool   `json:"ok"`
	Error string `json:"error,omitempty"`
}

func UploadFileToSlack(fileInfo FileInfo, fileBuffer *bytes.Buffer) (bool, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", fileInfo.FileName)
	if err != nil {
		return false, fmt.Errorf("error creating form file: %w", err)
	}
	if _, err := io.Copy(part, fileBuffer); err != nil {
		return false, fmt.Errorf("error copying file buffer: %w", err)
	}

	fields := map[string]string{
		"channels":        slackChannelID,
		"filename":        fileInfo.FileName,
		"filetype":        fileInfo.FileType,
		"initial_comment": fileInfo.InitialComment,
		"title":           fileInfo.Title,
	}
	for k, v := range fields {
		if err := writer.WriteField(k, v); err != nil {
			return false, fmt.Errorf("error writing field %s: %w", k, err)
		}
	}

	if err := writer.Close(); err != nil {
		return false, fmt.Errorf("error closing writer: %w", err)
	}

	req, err := http.NewRequest("POST", SlackUploadAPI, body)
	if err != nil {
		return false, fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+slackToken)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return false, fmt.Errorf("error performing request: %w", err)
	}
	defer resp.Body.Close()

	var slackResp SlackResponse
	if err := json.NewDecoder(resp.Body).Decode(&slackResp); err != nil {
		return false, fmt.Errorf("error decoding response: %w", err)
	}

	if !slackResp.OK {
		return false, fmt.Errorf("failed to upload file: %s", slackResp.Error)
	}

	return true, nil
}

func main() {
	fileInfos := []FileInfo{
		{
			FileName:       "example.txt",
			FileType:       "text",
			Title:          "Example File",
			InitialComment: "Here's the file you requested!",
		},
		// Add more FileInfo structs to upload multiple files concurrently.
	}

	var wg sync.WaitGroup
	for _, fileInfo := range fileInfos {
		wg.Add(1)
		go func(fileInfo FileInfo) {
			defer wg.Done()
			fileBuffer := bytes.NewBufferString("This is an example file content") // This should ideally be unique per file
			success, err := UploadFileToSlack(fileInfo, fileBuffer)
			if !success || err != nil {
				fmt.Printf("File upload failed: %v\n", err)
			} else {
				fmt.Println("File uploaded successfully.")
			}
		}(fileInfo)
	}
	wg.Wait()
}