package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"

	"google.golang.org/api/googleapi/transport"
	"google.golang.org/api/youtube/v3"
)

var songUrls []string
var i int
var maxResults = flag.Int64("max-results", 19, "Youtube results")

const developerKey = "AIzaSyDs-JNazrMlfMle0u4LSOXidEbFJZ45u7s"

type streamBuffer struct {
	stream io.ReadCloser
	url    string
}

func retURL() {
	reader := bufio.NewReader(os.Stdin)
	query, _ := reader.ReadString('\n')
	flag.Parse()
	client := &http.Client{
		Transport: &transport.APIKey{Key: developerKey},
	}

	service, err := youtube.New(client)
	if err != nil {
		log.Fatalf("Error creating new YouTube client: %v", err)
	}

	// Make the API call to YouTube.
	call := service.Search.List("id,snippet").
		Q(query).
		MaxResults(*maxResults)
	response, err := call.Do()
	if err != nil {
		log.Fatalf("Error making search API call: %v", err)
	}

	// Group video, channel, and playlist results in separate lists.
	videos := make(map[string]string)

	// Iterate through each item and add it to the correct list.
	for _, item := range response.Items {
		switch item.Id.Kind {
		case "youtube#video":
			videos[item.Id.VideoId] = item.Snippet.Title
		}
	}

	returnIDs("Search result:\n", videos)
}

func returnIDs(sectionName string, matches map[string]string) {
	fmt.Printf("%v\n", sectionName)

	for id, title := range matches {
		songUrls = append(songUrls, getURL(id))
		fmt.Printf(" %v \n", title)
	}
	fmt.Printf("\n")
}

func getURL(id string) string {
	url := "https://www.youtube.com/watch?v=" + id
	return url
}

func newStreamBuffer(stream io.ReadCloser, url string) *streamBuffer {
	return &streamBuffer{stream, url}
}

// Read input URLS from stdin.
func reader(c chan<- string) {
	defer close(c)
	stdin := bufio.NewScanner(os.Stdin)
	for stdin.Scan() {
		video := songUrls[i-1]
		c <- video
	}
	if err := stdin.Err(); err != nil {
		log.Print(err)
	}
}

// Buffer video streams for input URLs.
func bufferer(in <-chan string, out chan<- *streamBuffer) {
	defer close(out)
	for video := range in {
		cmd := exec.Command("youtube-dl", "-q", "-o", "-", video)
		stream, err := cmd.StdoutPipe()
		if err != nil {
			log.Print(err)
			break
		}
		fmt.Printf("Buffering %s\n", video)
		cmd.Start()
		out <- newStreamBuffer(stream, video)
	}
}

// Play buffered streams one by one as they come in.
func player(streams <-chan *streamBuffer) {
	for stream := range streams {
		cmd := exec.Command("mpv", "--no-terminal", "--no-video", "-")
		cmd.Stdin = stream.stream
		fmt.Printf("Playing %s\n", stream.url)
		cmd.Run()
		return
	}
}

func main() {
	retURL()
	fmt.Scan(&i)

	videos := make(chan string)
	streams := make(chan *streamBuffer)

	go reader(videos)
	go bufferer(videos, streams)
	player(streams)
}
