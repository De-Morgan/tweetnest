package scrapper

import (
	"log"
	"producer-service/internal/models"
	"producer-service/internal/producer"
	"slices"
	"time"

	twitterscraper "github.com/imperatrona/twitter-scraper"
)

type Scraper struct {
	username        string
	frequency       time.Duration
	maxPickPerFetch int
	scraper         *twitterscraper.Scraper
}

var _ producer.TweetProducer = (*Scraper)(nil)

func New(username string, maxPickPerFetch int, frequency time.Duration) *Scraper {
	s := twitterscraper.New()
	scraper := Scraper{
		username:        username,
		frequency:       frequency,
		scraper:         s,
		maxPickPerFetch: maxPickPerFetch,
	}
	return &scraper
}

func (s *Scraper) scrap() ([]models.Tweet, error) {
	var cursor string

	tweets, _, err := s.scraper.FetchTweets(s.username, s.maxPickPerFetch, cursor)
	if err != nil {
		return nil, err
	}
	tws := models.ScrapperToTweets(tweets)
	slices.SortFunc(tws, func(a, b models.Tweet) int {
		return b.Time.Compare(a.Time)
	})

	if len(tws) < s.maxPickPerFetch {
		return tws[:], nil
	}
	return tws[:s.maxPickPerFetch], nil
}

func (s *Scraper) StreamTweet() <-chan models.Tweet {
	result := make(chan models.Tweet, s.maxPickPerFetch)
	go func() {
		for {
			twts, err := s.scrap()
			if err != nil {
				log.Println(err)
				continue
			}
			for _, t := range twts {
				result <- t
			}
			<-time.After(s.frequency)
		}
	}()
	return result
}
