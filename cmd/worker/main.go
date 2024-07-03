package main

import (
	"context"
	"cqrs-sample/internal/database"
	"cqrs-sample/internal/queue"
	"cqrs-sample/pkg/handler"
	"fmt"
	_ "github.com/joho/godotenv/autoload"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
)

func main() {
	ctx := context.Background()
	mongoURI := os.Getenv("MONGO_URI")
	amqpDial := os.Getenv("AMQP_DIAL")

	mongoClient, err := mongo.Connect(ctx, options.Client().
		ApplyURI(mongoURI))
	if err != nil {
		log.Fatalln(err)
	}

	defer func() {
		if err := mongoClient.Disconnect(ctx); err != nil {
			log.Fatalln(err)
		}
	}()

	libraryDatabase := os.Getenv("LIBRARY_DATABASE")
	mongoDB, err := database.NewMongo(mongoClient.Database(libraryDatabase))
	if err != nil {
		log.Fatalln(err)
	}

	amqpConnection, err := amqp.Dial(amqpDial)
	if err != nil {
		log.Fatalln(err)
	}
	defer func() {
		_ = amqpConnection.Close()
	}()

	artistHandler := handler.NewArtistSubscribed(mongoDB)
	albumHandler := handler.NewAlbumPublished(mongoDB)
	songHandler := handler.NewSongPublished(mongoDB)

	go func() {
		channel, err := amqpConnection.Channel()
		if err != nil {
			log.Fatalln(err)
		}

		artistSubscribedQueue := os.Getenv("ARTIST_SUBSCRIBED_QUEUE")
		artistSubscriber := queue.NewRabbitMQSubscriber(channel, artistSubscribedQueue)
		if err := artistSubscriber.Subscribe(ctx, artistHandler); err != nil {
			log.Fatalln(err)
		}

		_ = channel.Close()
	}()

	go func() {
		channel, err := amqpConnection.Channel()
		if err != nil {
			log.Fatalln(err)
		}

		albumPublishedQueue := os.Getenv("ALBUM_PUBLISHED_QUEUE")
		albumSubscriber := queue.NewRabbitMQSubscriber(channel, albumPublishedQueue)
		if err := albumSubscriber.Subscribe(ctx, albumHandler); err != nil {
			log.Fatalln(err)
		}

		_ = channel.Close()
	}()

	go func() {
		channel, err := amqpConnection.Channel()
		if err != nil {
			log.Fatalln(err)
		}

		songPublishedQueue := os.Getenv("SONG_PUBLISHED_QUEUE")
		songSubscriber := queue.NewRabbitMQSubscriber(channel, songPublishedQueue)
		if err := songSubscriber.Subscribe(ctx, songHandler); err != nil {
			log.Fatalln(err)
		}

		_ = channel.Close()
	}()

	done := make(chan struct{})
	fmt.Println("listening...")
	<-done
}
