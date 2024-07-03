package main

import (
	"context"
	"cqrs-sample/handler"
	"cqrs-sample/internal/database"
	"cqrs-sample/internal/queue"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

func main() {
	ctx := context.Background()
	mongoURI := "mongodb://root:Pa55w0rd@localhost:27017/"
	amqpDial := "amqp://rabbitmq:Pa55w0rd@localhost:5672/"

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

	mongoDB, err := database.NewMongo(mongoClient.Database("library"))
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

		artistSubscriber := queue.NewRabbitMQSubscriber(channel, "artist.subscribed")
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

		albumSubscriber := queue.NewRabbitMQSubscriber(channel, "album.published")
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

		songSubscriber := queue.NewRabbitMQSubscriber(channel, "song.published")
		if err := songSubscriber.Subscribe(ctx, songHandler); err != nil {
			log.Fatalln(err)
		}

		_ = channel.Close()
	}()

	done := make(chan struct{})
	fmt.Println("listening...")
	<-done
}
