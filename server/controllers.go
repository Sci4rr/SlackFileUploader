package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
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
		return false, err
	}
	_, err = io.Copy(part, fileBuffer)
	if err != nil {
		return false, err
	}

	_ = writer.WriteField("channels", slackChannelID)
	_ = writer.WriteField("filename", fileInfo.FileName)
	_ = writer.WriteField("filetype", fileInfo.FileType)
	_ = writer.WriteField("initial_comment", fileInfo.InitialComment)
	_ = writer.WriteField("title", fileInfo.Title)

	err = writer.Close()
	if err != nil {
		return false, err
	}

	req, err := http.NewRequest("POST", SlackUploadAPI, body)
	req.Header.Set("Authorization", "Bearer "+slackToken)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	if err != nil {
		return false, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	var slackResp SlackResponse
	if err := json.NewDecoder(resp.Body).Decode(&slackResp); err != nil {
		return false, err
	}

	if !slackResp.OK {
		return false, fmt.Errorf("failed to upload file: %s", slackResp.Error)
	}

	return true, nil
}

func main() {
	fileInfo := FileInfo{
		FileName:       "example.txt",
		FileType:       "text",
		Title:          "Example File",
		InitialComment: "Here's the file you requested!",
	}
	fileBuffer := bytes.NewBufferString("This is an example file content")

	success, err := UploadFileToSlack(fileInfo, fileBuffer)
	if !success || err != nil {
		fmt.Printf("File upload failed: %v\n", err)
	} else {
		fmt.Println("File uploaded successfully.")
	}
}