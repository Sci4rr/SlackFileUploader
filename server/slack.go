package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
)

const SlackAPIURL = "https://slack.com/api/files.upload"

func uploadFileToSlack(filePath, channels string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", filePath)
	if err != nil {
		return fmt.Errorf("error creating form file: %w", err)
	}

	fileContents, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("error reading file contents: %w", err)
	}

	_, err = part.Write(fileContents)
	if err != nil {
		return fmt.Errorf("error writing file contents to form: %w", err)
	}

	if err := writer.WriteField("channels", channels); err != nil {
		return fmt.Errorf("error writing field 'channels': %w", err)
	}
	if err := writer.WriteField("token", os.Getenv("SLACK_BOT_TOKEN")); err != nil {
		return fmt.Errorf("error writing field 'token': %w", err)
	}

	if err := writer.Close(); err != nil {
		return fmt.Errorf("error closing writer: %w", err)
	}

	request, err := http.NewRequest("POST", SlackAPIURL, body)
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}
	request.Header.Add("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return fmt.Errorf("error performing request: %w", err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error from Slack API: %s", responseBody)
	}

	fmt.Println("Response from Slack: ", string(responseBody))

	return nil
}

func handleResponse(responseBody []byte) error {
	fmt.Println("Handle the response: ", string(responseBody))
	return nil
}

func main() {
	filePath := "path/to/your/file.jpg"
	channels := "channelID"
	err := uploadFileToSlack(filePath, channels)
	if err != nil {
		log.Fatalf("Failed to upload file: %v", err)
	}
}