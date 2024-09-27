package main

import (
    "sync"
    "time"
)

func UploadFilesToSlack(files []YourFileType) {
    var wg sync.WaitGroup
    throttle := time.Tick(time.Second / SlackRateLimit)

    for _, file := range files {
        wg.Add(1)
        go func(file YourFileType) {
            defer wg.Done()
            <-throttle
            uploadFile(file)
        }(file)
    }
    wg.Wait()
}

func uploadFile(file YourFileType) {
}

type YourFileType struct {
}

func main() {
    files := []YourFileType{}
    UploadFilesToSlack(files)
}