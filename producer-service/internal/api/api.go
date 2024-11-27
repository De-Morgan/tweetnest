package api

// import (
// 	"context"
// 	"fmt"
// 	"log"
// 	"os"
// 	"os/signal"
// 	"producer-service"
// 	"producer-service/internal/models"
// 	"syscall"

// 	"github.com/dghubble/go-twitter/twitter"
// 	"github.com/dghubble/oauth1"
// )

// type Config struct {
// 	apiKey       string
// 	apiSecret    string
// 	accessToken  string
// 	accessSecret string
// }

// var _ producer.TweetProducer = (*TwitterApi)(nil)

// type TwitterApi struct {
// 	client *twitter.Client
// 	tracks []string
// }

// func NewTwitterApi(tracks ...string) (*TwitterApi, error) {
// 	cfg := createConfig()
// 	if cfg.apiKey == "" || cfg.apiSecret == "" || cfg.accessToken == "" || cfg.accessSecret == "" {
// 		return nil, fmt.Errorf("%q,%q,%q,%q are required", "API_KEY", "API_SECRET", "ACCESS_TOKEN", "ACCESS_SECRET")
// 	}
// 	config := oauth1.NewConfig(cfg.apiKey, cfg.apiSecret)
// 	token := oauth1.NewToken(cfg.accessToken, cfg.accessSecret)
// 	httpClient := config.Client(oauth1.NoContext, token)
// 	client := twitter.NewClient(httpClient)
// 	_, _, err := client.Accounts.VerifyCredentials(&twitter.AccountVerifyParams{})
// 	if err != nil {
// 		return nil, err
// 	}
// 	api := TwitterApi{
// 		client: client,
// 		tracks: tracks,
// 	}
// 	return &api, nil
// }

// func (a *TwitterApi) StreamTweet() <-chan models.Tweet {
// 	result := make(chan models.Tweet, 10)
// 	go func() {
// 		a.stream(result)
// 	}()
// 	return result
// }

// func (a *TwitterApi) stream(sink chan<- models.Tweet) error {
// 	params := &twitter.StreamFilterParams{
// 		Track:         a.tracks,
// 		StallWarnings: twitter.Bool(true),
// 	}
// 	stream, err := a.client.Streams.Filter(params)
// 	if err != nil {
// 		log.Fatalf("Error starting stream: %v", err)
// 		return err
// 	}
// 	defer stream.Stop()
// 	demux := twitter.NewSwitchDemux()
// 	demux.Tweet = func(tweet *twitter.Tweet) {
// 		fmt.Printf("Tweet: %s\n", tweet.Text)
// 		sink <- models.ApiToTweet(*tweet)
// 	}
// 	demux.Warning = func(warning *twitter.StallWarning) {
// 		fmt.Printf("Warning: %s\n", warning.Message)
// 	}
// 	demux.All = func(message interface{}) {
// 		fmt.Printf("Message: %s\n", message)

// 	}
// 	// Handle interrupts for clean shutdown
// 	ctx, cancel := context.WithCancel(context.Background())
// 	go func() {
// 		c := make(chan os.Signal, 1)
// 		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
// 		<-c
// 		cancel()
// 	}()
// 	log.Println("Starting stream...")
// 	go demux.HandleChan(stream.Messages)
// 	<-ctx.Done()
// 	log.Println("Stopping stream...")
// 	close(sink)
// 	return nil
// }

// // ==================================================
// func createConfig() Config {
// 	return Config{
// 		apiKey:       os.Getenv("API_KEY"),
// 		apiSecret:    os.Getenv("API_SECRET"),
// 		accessToken:  os.Getenv("ACCESS_TOKEN"),
// 		accessSecret: os.Getenv("ACCESS_SECRET"),
// 	}
// }
