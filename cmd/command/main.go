package main

import (
	"context"
	"cqrs-sample/command"
	"cqrs-sample/handler"
	"cqrs-sample/internal/database"
	"cqrs-sample/internal/queue"
	"cqrs-sample/internal/server"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
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

	subscribeArtistCommand := command.NewSubscribeArtist(postgresDB, rabbitMQPublisher)
	publishAlbumCommand := command.NewPublishAlbum(postgresDB, rabbitMQPublisher)
	publishSongCommand := command.NewPublishSong(postgresDB, rabbitMQPublisher)

	artistHandler := handler.NewArtistWriter(subscribeArtistCommand)
	albumHandler := handler.NewAlbumWriter(publishAlbumCommand)
	songHandler := handler.NewSongWriter(publishSongCommand)

	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Post("/artists", artistHandler.Create)
	r.Post("/albums", albumHandler.Create)
	r.Post("/songs", songHandler.Create)

	s := server.New(r)
	if err := s.StartWithGracefulShutdown(ctx, ":3030"); err != nil {
		log.Fatalln(err)
	}
}
