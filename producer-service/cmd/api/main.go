package main

import (
	"log"
	"producer-service/internal/producer"
	"producer-service/internal/scrapper"
	"time"
)

func main() {
	tweetsProducer := scrapper.New("elonmusk", 20, 12*time.Hour)
	producer := producer.New(tweetsProducer)
	err := producer.SendToTopic()
	if err != nil {
		log.Fatal(err)
	}
}
