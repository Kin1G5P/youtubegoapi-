package main

import (
        "context"
        "fmt"
        "io"
        "log"
        "net/http"
        "os"
        "path/filepath"

        "google.golang.org/api/option"
        "google.golang.org/api/youtube/v3"
)

func main() {
        apiKey := os.Getenv("YOUTUBE_API_KEY")
        if apiKey == "" {
                log.Fatalf("Please set the YOUTUBE_API_KEY environment variable")
        }

        ctx := context.Background()

        youtubeService, err := youtube.NewService(ctx, option.WithAPIKey(apiKey))
        if err != nil {
                log.Fatalf("Error creating YouTube service: %v", err)
        }

        // Example: Search for videos.
        call := youtubeService.Search.List([]string{"snippet"}).
                Q("golang tutorial").
                MaxResults(1) // Download the thumbnail of the first result

        response, err := call.Do()
        if err != nil {
                log.Fatalf("Error making API call: %v", err)
        }

        if len(response.Items) > 0 && response.Items[0].Id.Kind == "youtube#video" {
                videoID := response.Items[0].Id.VideoId
                thumbnailURL := response.Items[0].Snippet.Thumbnails.High.Url // Get the thumbnail URL.
                downloadThumbnail(thumbnailURL, videoID+".jpg")
        } else {
                fmt.Println("No video found.")
        }

}

func downloadThumbnail(url, filename string) {

        resp, err := http.Get(url)
        if err != nil {
                log.Fatalf("Error downloading thumbnail: %v", err)
        }
        defer resp.Body.Close()

        if resp.StatusCode != http.StatusOK {
                log.Fatalf("Unexpected status code: %d", resp.StatusCode)
        }

        // Create the file.
        out, err := os.Create(filename)
        if err != nil {
                log.Fatalf("Error creating file: %v", err)
        }
        defer out.Close()

        // Write the body to file.
        _, err = io.Copy(out, resp.Body)
        if err != nil {
                log.Fatalf("Error writing to file: %v", err)
        }

        fmt.Printf("Thumbnail downloaded to %s\n", filename)
}

func getFileNameFromURL(url string) string {
        return filepath.Base(url)
}
