package rss

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"

	"github.com/WaronLimsakul/gator/internal/database"
)

func FetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)

	if err != nil {
		return &RSSFeed{}, err
	}

	req.Header.Set("Accept", "application/xml")
	req.Header.Set("User-Agent", "gator")

	client := http.Client{}

	res, err := client.Do(req)
	if err != nil {
		return &RSSFeed{}, err
	}

	defer res.Body.Close()
	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return &RSSFeed{}, err
	}

	var fetchedFeed RSSFeed
	if err = xml.Unmarshal(bodyBytes, &fetchedFeed); err != nil {
		return &RSSFeed{}, err
	}

	fetchedFeed.Channel.Title = html.UnescapeString(fetchedFeed.Channel.Title)
	fetchedFeed.Channel.Description = html.UnescapeString(fetchedFeed.Channel.Description)

	return &fetchedFeed, nil

}

// 1. get the next feed to fetch from db
// 2. mark it as fetched
// 3. fetches the feed using the url
// 4. print the items in the fetched feed.
func ScrapeFeeds(db *database.Queries) error {
	nextFeed, err := db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return err
	}

	err = db.MarkFeedFetched(context.Background(), nextFeed.ID)
	if err != nil {
		return err
	}

	targetUrl := nextFeed.Url
	fetchedFeed, err := FetchFeed(context.Background(), targetUrl)
	if err != nil {
		return err
	}

	for _, feed := range (*fetchedFeed).Channel.Item {
		fmt.Println("----------------------------------------")
		fmt.Printf("Title: %s\n", feed.Title)
		fmt.Printf("Description: %s\n", feed.Description)
		fmt.Printf("Link: %s\n", feed.Link)
		fmt.Printf("Pub Date: %s\n", feed.PubDate)
	}
	return nil

}
