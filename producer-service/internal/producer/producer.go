package producer

import (
	"encoding/json"
	"log"
	"producer-service/internal/models"
	"time"

	"github.com/IBM/sarama"
)

type TweetProducer interface {
	StreamTweet() <-chan models.Tweet
}

const (
	KafkaServerAddress = "kafka:9092"
	KafkaTopic         = "tweets"
)

type Producer struct {
	config *sarama.Config
	tweets TweetProducer
}

func New(tweets TweetProducer) *Producer {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForLocal
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	config.Producer.Compression = sarama.CompressionSnappy
	return &Producer{
		config: config,
		tweets: tweets,
	}
}

func (p *Producer) SendToTopic() error {

	producer, err := sarama.NewAsyncProducer([]string{KafkaServerAddress}, p.config)
	if err != nil {
		log.Printf("Failed to start Kafka producer: %v\n", err)
		return err
	}
	defer func() {
		if err := producer.Close(); err != nil {
			log.Printf("Failed to close Kafka producer: %v\n", err)
		}
	}()

	// Goroutine to handle successful messages
	go func() {
		for msg := range producer.Successes() {
			log.Printf("Message sent to topic %s, partition %d, offset %d\n",
				msg.Topic, msg.Partition, msg.Offset)
		}
	}()

	// Goroutine to handle errors
	go func() {
		for err := range producer.Errors() {
			log.Printf("Failed to send message: %v\n", err)
		}
	}()

	for tweet := range p.tweets.StreamTweet() {
		log.Println("Proccessing tweet: ", tweet.ID)
		twtByte, err := json.Marshal(tweet)
		if err != nil {
			log.Printf("Failed to marshal tweet[%q]: %v\n", tweet.ID, err)
			return err
		}
		mes := sarama.ProducerMessage{
			Topic: KafkaTopic,
			Key:   sarama.StringEncoder(tweet.ID),
			Value: sarama.ByteEncoder(twtByte),
		}
		producer.Input() <- &mes
		<-time.After(200 * time.Millisecond) // Delay for 200 miliseconds
	}
	// Give the producer some time to process messages
	log.Println("All messages sent, waiting for acknowledgment...")
	time.Sleep(5 * time.Second) // Adjust this for longer processing
	log.Println("Shutting down producer...")

	return nil
}
