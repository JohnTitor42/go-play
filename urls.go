package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"google.golang.org/api/googleapi/transport"
	"google.golang.org/api/youtube/v3"
)

var songUrls []string
var maxResults = flag.Int64("max-results", 25, "Max YouTube results")

const developerKey = "AIzaSyDs-JNazrMlfMle0u4LSOXidEbFJZ45u7s"

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
	channels := make(map[string]string)
	playlists := make(map[string]string)

	// Iterate through each item and add it to the correct list.
	for _, item := range response.Items {
		switch item.Id.Kind {
		case "youtube#video":
			videos[item.Id.VideoId] = item.Snippet.Title
		case "youtube#channel":
			channels[item.Id.ChannelId] = item.Snippet.Title
		case "youtube#playlist":
			playlists[item.Id.PlaylistId] = item.Snippet.Title
		}
	}

	returnIDs("Search result:\n", videos)
	//printIDs("Channels", channels)Max
	//printIDs("Playlists", playlists)
}

// Print the ID and title of each result in a list as well as a name that
// identifies the list. For example, print the word section name "Videos"
// above a list of video search results, followed by the video ID and title
// of each matching video.
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

func main() {
	//fmt.Scan(&querty)
	//fmt.Println(querty)
	retURL()
	var i int
	fmt.Scan(&i)
	fmt.Println("\n", songUrls[i-1])

}
