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

	subscriber := queue.NewRabbitMQSubscriber(amqpConnection)
	artistSubscribedHandler := handler.NewArtistSubscribed(mongoDB)
	albumPublishedHandler := handler.NewAlbumPublished(mongoDB)
	songPublishedHandler := handler.NewSongPublished(mongoDB)
	songPlayedHandler := handler.NewIncrementSongPlays(mongoDB)

	go func() {
		artistSubscribedQueue := os.Getenv("ARTIST_SUBSCRIBED_QUEUE")
		if err := subscriber.Subscribe(ctx, artistSubscribedQueue, artistSubscribedHandler); err != nil {
			log.Fatalln(err)
		}
	}()

	go func() {
		albumPublishedQueue := os.Getenv("ALBUM_PUBLISHED_QUEUE")
		if err := subscriber.Subscribe(ctx, albumPublishedQueue, albumPublishedHandler); err != nil {
			log.Fatalln(err)
		}
	}()

	go func() {
		songPublishedQueue := os.Getenv("SONG_PUBLISHED_QUEUE")
		if err := subscriber.Subscribe(ctx, songPublishedQueue, songPublishedHandler); err != nil {
			log.Fatalln(err)
		}
	}()

	go func() {
		songPlayedQueue := os.Getenv("SONG_PLAYED_QUEUE")
		if err := subscriber.Subscribe(ctx, songPlayedQueue, songPlayedHandler); err != nil {
			log.Fatalln(err)
		}
	}()

	done := make(chan struct{})
	fmt.Println("listening...")
	<-done
}
