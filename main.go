package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/deanishe/awgo"
	"google.golang.org/api/googleapi/transport"
	"google.golang.org/api/youtube/v3"
)

var (
	query      = flag.String("query", "Google", "Search term")
	maxResults = flag.Int64("max-results", 10, "Max YouTube results")
)

// Workflow is the main API
var wf *aw.Workflow

func init() {
	// Create a new Workflow using default settings.
	// Critical settings are provided by Alfred via environment variables,
	// so this *will* die in flames if not run in an Alfred-like environment.
	wf = aw.New()
}

func run() {
	flag.Parse()

	developerKey, ok := os.LookupEnv("YOUTUBE_DEVELOPER_KEY")

	if !ok {
		log.Fatalf("YOUTUBE_DEVELOPER_KEY value is required")
	}

	client := &http.Client{
		Transport: &transport.APIKey{Key: developerKey},
	}

	service, err := youtube.New(client)
	if err != nil {
		log.Fatalf("Error creating new YouTube client: %v", err)
	}

	// Make the API call to YouTube.
	call := service.Search.List("id,snippet").
		Q(*query).
		MaxResults(*maxResults)
	response, _ := call.Do()

	if err != nil {
		log.Fatalf("Error : %v", err)
	}

	// Iterate through each item and add it to the correct list.
	for _, item := range response.Items {
		switch item.Id.Kind {
		case "youtube#video":
			quicklookURL := fmt.Sprintf("https://www.youtube.com/watch?v=%s", item.Id.VideoId)
			wf.NewItem(item.Snippet.Title).
				Subtitle(item.Snippet.Description).
				Quicklook(quicklookURL).
				Arg(item.Id.VideoId).
				Autocomplete(item.Snippet.Title).
				Valid(true)
		}
	}

	wf.SendFeedback()

}

func main() {
	wf.Run(run)
}
