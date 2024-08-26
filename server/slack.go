package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
)

const SlackAPIURL = "https://slack.com/api/files.upload"

func uploadFileToSlack(filePath, channels string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", filePath)
	if err != nil {
		return err
	}
	_, err = ioutil.ReadAll(file)
	if err != nil {
		return err
	}
	_, err = part.Write([]byte("file contents here"))
	if err != nil {
		return err
	}
	writer.WriteField("channels", channels)
	writer.WriteField("token", os.Getenv("SLACK_BOT_TOKEN"))

	err = writer.Close()
	if err != nil {
		return err
	}
	request, err := http.NewRequest("POST", SlackAPIURL, body)
	if err != nil {
		return err
	}
	request.Header.Add("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
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
		log.Fatal(err)
	}
}