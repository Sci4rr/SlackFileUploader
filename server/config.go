package main

import (
    "sync"
    "time"
)

type YourFileType struct {
}

func UploadFilesToSlack(files []YourFileType) {
    var wg sync.WaitGroup
    throttle := time.NewTicker(time.Second / SlackRateLimit)
    defer throttle.Stop()

    for _, file := range files {
        wg.Add(1)
        go func(file YourFileType) {
            defer wg.Done()
            <-throttle.C
            uploadFile(file)
        }(file)
    }
    wg.Wait()
}

func uploadFile(file YourFileType) {
}

func main() {
    files := []YourFileType{}
    UploadFilesToSlack(files)
}