package main

import (
	"context"
	"cqrs-sample/handler"
	"cqrs-sample/internal/database"
	"cqrs-sample/internal/server"
	"cqrs-sample/query"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

func main() {
	ctx := context.Background()

	mongoURI := "mongodb://root:Pa55w0rd@localhost:27017/"
	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
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

	getArtistQuery := query.NewGetArtist(mongoDB)
	getAlbumQuery := query.NewGetAlbum(mongoDB)
	getSongQuery := query.NewGetSong(mongoDB)

	artistHandler := handler.NewArtistReader(getArtistQuery)
	albumHandler := handler.NewAlbumReader(getAlbumQuery)
	songHandler := handler.NewSongReader(getSongQuery)

	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Get("/artist/{artistID}", artistHandler.GetByID)
	r.Get("/album/{albumID}", albumHandler.GetByID)
	r.Get("/song/{songID}", songHandler.GetByID)

	s := server.New(r)
	if err := s.StartWithGracefulShutdown(ctx, ":3030"); err != nil {
		log.Fatalln(err)
	}
}
