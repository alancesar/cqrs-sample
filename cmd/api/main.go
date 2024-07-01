package main

import (
	"context"
	"cqrs-sample/handler"
	"cqrs-sample/internal/database"
	"cqrs-sample/query"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

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

	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Get("/artist/{artistID}", artistHandler.GetByID)
	r.Get("/album/{albumID}", albumHandler.GetByID)
	r.Get("/song/{songID}", songHandler.GetByID)

	server := http.Server{
		Addr:    ":3000",
		Handler: r,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalln(err)
		}
	}()

	log.Println("all systems go!")

	<-ctx.Done()
	stop()

	log.Println("shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Println(err)
	}

	log.Println("good bye")
}
