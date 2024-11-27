package models

import (
	"time"

	twitterscraper "github.com/imperatrona/twitter-scraper"
)

type Photo struct {
	Url string `json:"url,omitempty"`
}
type Video struct {
	Preview string `json:"preview,omitempty"`
	URL     string `json:"url,omitempty"`
	HLSURL  string `json:"hls_url,omitempty"`
}

type Tweet struct {
	HTML             string    `json:"html,omitempty"`
	ID               string    `json:"id,omitempty"`
	Name             string    `json:"name,omitempty"`
	Photos           []Photo   `json:"photos,omitempty"`
	Text             string    `json:"text,omitempty"`
	Time             time.Time `json:"time,omitempty"`
	Username         string    `json:"username,omitempty"`
	Videos           []Video   `json:"videos,omitempty"`
	SensitiveContent bool      `json:"sensitive,omitempty"`
}

func ScrapperToTweets(tweets []*twitterscraper.Tweet) []Tweet {
	tws := make([]Tweet, len(tweets))
	for i := range tweets {
		tws[i] = scrapperToTweet(tweets[i])
	}
	return tws
}

func scrapperToTweet(tweet *twitterscraper.Tweet) Tweet {
	photos := make([]Photo, len(tweet.Photos))
	for i, photo := range tweet.Photos {
		photos[i] = Photo{
			Url: photo.URL,
		}
	}
	videos := make([]Video, len(tweet.Videos))

	for i, video := range tweet.Videos {
		videos[i] = Video{
			Preview: video.Preview,
			URL:     video.URL,
			HLSURL:  video.HLSURL,
		}
	}
	return Tweet{
		HTML:             tweet.HTML,
		ID:               tweet.ID,
		Name:             tweet.Name,
		Photos:           photos,
		Text:             tweet.Text,
		Time:             tweet.TimeParsed,
		Username:         tweet.Username,
		Videos:           videos,
		SensitiveContent: tweet.SensitiveContent,
	}
}
