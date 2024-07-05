package main

import (
	"context"
	"cqrs-sample/internal/database"
	"cqrs-sample/internal/server"
	"cqrs-sample/pkg/handler"
	"cqrs-sample/pkg/query"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/joho/godotenv/autoload"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
)

func main() {
	ctx := context.Background()
	mongoURI := os.Getenv("MONGO_URI")

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
	getAlbumsByArtistQuery := query.NewGetAlbumsByArtist(mongoDB)
	getSongQuery := query.NewGetSong(mongoDB)

	artistHandler := handler.NewArtistReader(getArtistQuery, getAlbumsByArtistQuery)
	albumHandler := handler.NewAlbumReader(getAlbumQuery)
	songHandler := handler.NewSongReader(getSongQuery)

	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Get("/artist/{artistID}", artistHandler.Get)
	r.Get("/artist/{artistID}/albums", artistHandler.GetAlbums)
	r.Get("/album/{albumID}", albumHandler.Get)
	r.Get("/song/{songID}", songHandler.Get)

	s := server.New(r)
	if err := s.StartWithGracefulShutdown(ctx, ":3031"); err != nil {
		log.Fatalln(err)
	}
}
