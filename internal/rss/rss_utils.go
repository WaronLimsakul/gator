package rss

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
	"time"

	"github.com/WaronLimsakul/gator/internal/database"
	"github.com/google/uuid"
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
// 4. save the item in db.
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

	for _, item := range (*fetchedFeed).Channel.Item {
		pubDate, err := time.Parse(time.RFC1123Z, item.PubDate)
		if err != nil {
			return err
		}

		_, err = db.CreatePost(
			context.Background(),
			database.CreatePostParams{
				ID:          uuid.New(),
				Title:       item.Title,
				Url:         item.Link,
				Description: item.Description,
				PublishedAt: pubDate,
				FeedID:      nextFeed.ID,
			})
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("save post")
	}

	return nil
}
