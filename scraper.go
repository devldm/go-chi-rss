package main

import (
	"context"
	"database/sql"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/devldm/go-server-rss/internal/database"
	"github.com/google/uuid"
)

func startScraping(
	db *database.Queries,
	concurrency int,
	timeBetweenRequest time.Duration,
) {
	log.Printf("Scraping on %v goroutines every %s duration", concurrency, timeBetweenRequest)

	ticker := time.NewTicker(timeBetweenRequest)

	for ; ; <-ticker.C {
		feeds, err := db.GetNextFeedsToFetch(
			context.Background(), // Global background / used when no scoped context is available
			int32(concurrency),
		)
		if err != nil {
			log.Println("error fetching feeds", err)
			continue
		}

		wg := &sync.WaitGroup{}
		for _, feed := range feeds { //iterate over feeds on main go routine
			wg.Add(1) // add 1 to wg for every feed. 30 feeds = 30 WG
			go scrapeFeed(db, wg, feed)
		}
		wg.Wait() // wait for 30 calls to wg.Done() (in above example)
	}
}

func scrapeFeed(db *database.Queries, wg *sync.WaitGroup, feed database.Feed) {
	defer wg.Done() // deferring wg.Done() will be called after the function returns

	_, err := db.MarkFeedAsFetched(context.Background(), feed.ID)
	if err != nil {
		log.Println("Error marking feed as fetched:", err)
		return
	}

	rssFeed, err := urlToFeed(feed.Url)
	if err != nil {
		log.Println("Error fetching feed:", err)
		return
	}

	for _, item := range rssFeed.Channel.Item {
		description := sql.NullString{}
		if item.Description != "" {
			description.String = item.Description
			description.Valid = true
		}

		pubAt, err := time.Parse(time.RFC1123Z, item.PubDate) // TODO: Update this time. to allow for various formats
		if err != nil {
			log.Printf("couldn't parse data %v with err %v", item.PubDate, err)
		}

		_, err = db.CreatePost(context.Background(),
			database.CreatePostParams{
				ID:          uuid.New(),
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
				Description: description,
				Title:       item.Title,
				PublishedAt: pubAt,
				Url:         item.Link,
				FeedID:      feed.ID,
			})

		if err != nil {
			if strings.Contains(err.Error(), "duplicate key") {
				continue
			}
			log.Println("failed to create post:", err)
		}
	}

	log.Printf("Feed %s collected, %v posts found", feed.Name, len(rssFeed.Channel.Item))

}
