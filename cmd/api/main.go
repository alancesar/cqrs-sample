package main

import (
	"context"
	"cqrs-sample/handler"
	"cqrs-sample/internal/database"
	"cqrs-sample/query"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)

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

	getAlbumQuery := query.NewGetAlbum(mongoDB)
	getArtistQuery := query.NewGetArtist(mongoDB)
	getSongQuery := query.NewGetSong(mongoDB)

	albumHandler := handler.NewGetAlbum(getAlbumQuery)
	artistHandler := handler.NewGetArtist(getArtistQuery)
	songHandler := handler.NewGetSong(getSongQuery)

	r.Get("/artist/{artistID}", artistHandler.Handle)
	r.Get("/album/{albumID}", albumHandler.Handle)
	r.Get("/song/{songID}", songHandler.Handle)

	if err := http.ListenAndServe(":3000", r); err != nil {
		log.Fatalln(err)
	}
}
