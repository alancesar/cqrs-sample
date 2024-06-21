package database

import (
	"context"
	"cqrs-sample/internal/database/document"
	"cqrs-sample/pkg/song"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	artistCollectionName = "artists"
	albumsCollectionName = "albums"
	songCollectionName   = "songs"
)

type (
	Mongo struct {
		db *mongo.Database
	}
)

func NewMongo(db *mongo.Database) (*Mongo, error) {
	return &Mongo{
		db: db,
	}, nil
}

func (m Mongo) CreateArtist(ctx context.Context, artist song.Artist) error {
	doc := document.NewArtistFromDomain(artist)
	_, err := m.db.Collection(artistCollectionName).InsertOne(ctx, doc)
	return err
}

func (m Mongo) CreateAlbum(ctx context.Context, album song.Album) error {
	doc := document.NewAlbumFromDomain(album)
	_, err := m.db.Collection(albumsCollectionName).InsertOne(ctx, doc)
	return err
}

func (m Mongo) CreateSong(ctx context.Context, song song.Song) error {
	doc := document.NewSongFromDomain(song)
	_, err := m.db.Collection(songCollectionName).InsertOne(ctx, doc)
	return err
}

func (m Mongo) AddSongToAlbum(ctx context.Context, song song.Song) error {
	result := m.db.Collection(albumsCollectionName).FindOne(ctx, bson.M{"_id": song.Album.ID})
	if err := result.Err(); err != nil {
		return err
	}

	doc := document.NewSongInAlbumFromDomain(song)
	_, err := m.db.Collection(albumsCollectionName).
		UpdateOne(ctx, bson.M{"_id": song.Album.ID}, bson.M{"$push": bson.M{"songs": doc}})
	return err
}

func (m Mongo) GetSongByID(ctx context.Context, id string) (song.Song, error) {
	result := m.db.Collection(songCollectionName).FindOne(ctx, bson.M{"_id": id})
	if err := result.Err(); err != nil {
		return song.Song{}, err
	}

	var doc document.Song
	if err := result.Decode(&doc); err != nil {
		return song.Song{}, err
	}

	return doc.ToDomain(), nil
}

func (m Mongo) GetArtistByID(ctx context.Context, id string) (song.Artist, error) {
	result := m.db.Collection(artistCollectionName).FindOne(ctx, bson.M{"_id": id})
	if err := result.Err(); err != nil {
		return song.Artist{}, err
	}

	var doc document.Artist
	if err := result.Decode(&doc); err != nil {
		return song.Artist{}, err
	}

	return doc.ToDomain(), nil
}

func (m Mongo) GetAlbumByID(ctx context.Context, id string) (song.Album, error) {
	result := m.db.Collection(albumsCollectionName).FindOne(ctx, bson.M{"_id": id})
	if err := result.Err(); err != nil {
		return song.Album{}, err
	}

	var doc document.Album
	if err := result.Decode(&doc); err != nil {
		return song.Album{}, err
	}

	return doc.ToDomain(), nil
}
