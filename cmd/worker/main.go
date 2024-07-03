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
	artistHandler := handler.NewArtistSubscribed(mongoDB)
	albumHandler := handler.NewAlbumPublished(mongoDB)
	songHandler := handler.NewSongPublished(mongoDB)

	go func() {
		artistSubscribedQueue := os.Getenv("ARTIST_SUBSCRIBED_QUEUE")
		if err := subscriber.Subscribe(ctx, artistSubscribedQueue, artistHandler); err != nil {
			log.Fatalln(err)
		}
	}()

	go func() {
		albumPublishedQueue := os.Getenv("ALBUM_PUBLISHED_QUEUE")
		if err := subscriber.Subscribe(ctx, albumPublishedQueue, albumHandler); err != nil {
			log.Fatalln(err)
		}
	}()

	go func() {
		songPublishedQueue := os.Getenv("SONG_PUBLISHED_QUEUE")
		if err := subscriber.Subscribe(ctx, songPublishedQueue, songHandler); err != nil {
			log.Fatalln(err)
		}
	}()

	done := make(chan struct{})
	fmt.Println("listening...")
	<-done
}
