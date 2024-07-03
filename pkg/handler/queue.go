package handler

import (
	"context"
	"cqrs-sample/pkg/event"
	"cqrs-sample/pkg/message"
	"cqrs-sample/pkg/song"
	"encoding/json"
	"fmt"
)

type (
	ArtistDatabase interface {
		CreateArtist(ctx context.Context, artist song.Artist) error
	}

	AlbumDatabase interface {
		CreateAlbum(ctx context.Context, album song.Album) error
	}

	SongDatabase interface {
		CreateSong(ctx context.Context, song song.Song) error
		AddSongToAlbum(ctx context.Context, song song.Song) error
	}

	ArtistSubscribed struct {
		db ArtistDatabase
	}

	AlbumPublished struct {
		db AlbumDatabase
	}

	SongPublished struct {
		db SongDatabase
	}
)

func NewArtistSubscribed(db ArtistDatabase) *ArtistSubscribed {
	return &ArtistSubscribed{
		db: db,
	}
}

func NewAlbumPublished(db AlbumDatabase) *AlbumPublished {
	return &AlbumPublished{
		db: db,
	}
}

func NewSongPublished(db SongDatabase) *SongPublished {
	return &SongPublished{
		db: db,
	}
}

func (ah ArtistSubscribed) Handle(ctx context.Context, body []byte, _ map[string]interface{}) error {
	artist, err := unmarshal[message.Artist](body)
	if err != nil {
		return err
	}

	return ah.db.CreateArtist(ctx, artist.ToDomain())
}

func (ap AlbumPublished) Handle(ctx context.Context, body []byte, _ map[string]interface{}) error {
	album, err := unmarshal[message.Album](body)
	if err != nil {
		return err
	}

	return ap.db.CreateAlbum(ctx, album.ToDomain())
}

func (sp SongPublished) Handle(ctx context.Context, body []byte, _ map[string]interface{}) error {
	s, err := unmarshal[message.Song](body)
	if err != nil {
		return err
	}

	if err := sp.db.CreateSong(ctx, s.ToDomain()); err != nil {
		return err
	}

	return sp.db.AddSongToAlbum(ctx, s.ToDomain())
}

func unmarshal[T any](body []byte) (T, error) {
	var output T
	if err := json.Unmarshal(body, &output); err != nil {
		return output, fmt.Errorf("%w: %s", event.InvalidPayloadErr, err)
	}

	return output, nil
}
