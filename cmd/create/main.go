package main

import (
	"context"
	"cqrs-sample/command"
	"cqrs-sample/internal/database"
	"cqrs-sample/internal/queue"
	"cqrs-sample/internal/uuid"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

func main() {
	ctx := context.Background()
	postgresDSN := "host=localhost user=postgres password=Pa55w0rd dbname=postgres port=5432 sslmode=disable TimeZone=America/Sao_Paulo"
	amqpDial := "amqp://rabbitmq:Pa55w0rd@localhost:5672/"

	db, err := gorm.Open(postgres.Open(postgresDSN), &gorm.Config{})
	if err != nil {
		log.Fatalln(err)
	}

	postgresDB := database.NewGorm(db)

	amqpConnection, err := amqp.Dial(amqpDial)
	if err != nil {
		log.Fatalln(err)
	}
	defer func() {
		_ = amqpConnection.Close()
	}()

	channel, err := amqpConnection.Channel()
	defer func() {
		_ = channel.Close()
	}()

	rabbitMQPublisher := queue.NewRabbitMQPublisher(channel, "library")
	uuidGenerator := uuid.New()

	subscribeArtist := command.NewSubscribeArtist(postgresDB, rabbitMQPublisher, uuidGenerator)
	publishAlbum := command.NewPublishAlbum(postgresDB, rabbitMQPublisher, uuidGenerator)
	publishSong := command.NewPublishSong(postgresDB, rabbitMQPublisher, uuidGenerator)

	artist, err := subscribeArtist.Execute(ctx, command.SubscribeArtistCommand{
		Name:   "Ramones",
		Gender: "Punk",
	})
	if err != nil {
		log.Fatalln(err)
	}

	album, err := publishAlbum.Execute(ctx, command.PublishAlbumCommand{
		Title:       "Rocket to Russia",
		ArtistID:    artist.ID,
		ReleaseYear: 1977,
	})
	if err != nil {
		log.Fatalln(err)
	}

	songs := []string{
		"Cretin Hop",
		"Rockaway Beach",
		"Here Today, Gone Tomorrow",
		"Locket Love",
		"I Don't Care",
		"Sheena Is a Punk Rocker",
		"We're a Happy Family",
		"Teenage Lobotomy",
		"Do You Wanna Dance?",
		"I Wanna Be Well",
		"I Can't Give You Anything",
		"Ramona",
		"Surfin' Bird",
		"Why Is It Always This Way?",
	}

	for i, title := range songs {
		song, err := publishSong.Execute(ctx, command.PublishSongCommand{
			TrackNumber: i + 1,
			Title:       title,
			AlbumID:     album.ID,
		})
		if err != nil {
			log.Fatalln(err)
		}

		fmt.Println(song.Title, "created successfully")
	}
}
