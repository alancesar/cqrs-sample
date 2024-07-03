package command

import (
	"context"
	"cqrs-sample/pkg/event"
	"cqrs-sample/pkg/message"
	"cqrs-sample/pkg/song"
	"github.com/google/uuid"
)

type (
	ArtistDatabase interface {
		CreateArtist(ctx context.Context, artist *song.Artist) error
		GetArtistByID(ctx context.Context, id string) (song.Artist, error)
	}

	AlbumDatabase interface {
		ArtistDatabase
		CreateAlbum(ctx context.Context, album *song.Album) error
		GetAlbumByID(ctx context.Context, id string) (song.Album, error)
	}

	SongDatabase interface {
		AlbumDatabase
		ArtistDatabase
		CreateSong(ctx context.Context, s *song.Song) error
	}

	Publisher interface {
		Publish(ctx context.Context, ev event.Message, key event.Event) error
	}

	SubscribeArtistCommand struct {
		Name   string
		Gender song.Gender
	}

	PublishAlbumCommand struct {
		Title       string
		ArtistID    string
		ReleaseYear int
	}

	PublishSongCommand struct {
		TrackNumber int
		Title       string
		AlbumID     string
	}

	SubscribeArtist struct {
		db  ArtistDatabase
		pub Publisher
	}

	PublishAlbum struct {
		db  AlbumDatabase
		pub Publisher
	}

	PublishSong struct {
		db  SongDatabase
		pub Publisher
	}
)

func NewSubscribeArtist(db ArtistDatabase, pub Publisher) *SubscribeArtist {
	return &SubscribeArtist{
		db:  db,
		pub: pub,
	}
}

func NewPublishAlbum(db AlbumDatabase, pub Publisher) *PublishAlbum {
	return &PublishAlbum{
		db:  db,
		pub: pub,
	}
}

func NewPublishSong(db SongDatabase, pub Publisher) *PublishSong {
	return &PublishSong{
		db:  db,
		pub: pub,
	}
}

func (ca SubscribeArtist) Execute(ctx context.Context, cmd SubscribeArtistCommand) (song.Artist, error) {
	artist := &song.Artist{
		ID:     uuid.NewString(),
		Name:   cmd.Name,
		Gender: cmd.Gender,
	}
	if err := ca.db.CreateArtist(ctx, artist); err != nil {
		return song.Artist{}, err
	}

	m := event.NewMessage(message.NewArtistFromDomain(*artist))
	if err := ca.pub.Publish(ctx, m, event.ArtistSubscribedEvent); err != nil {
		return song.Artist{}, err
	}

	return *artist, nil
}

func (ca PublishAlbum) Execute(ctx context.Context, cmd PublishAlbumCommand) (song.Album, error) {
	artist, err := ca.db.GetArtistByID(ctx, cmd.ArtistID)
	if err != nil {
		return song.Album{}, err
	}

	album := &song.Album{
		ID:          uuid.NewString(),
		Title:       cmd.Title,
		Artist:      artist,
		ReleaseYear: cmd.ReleaseYear,
	}
	if err := ca.db.CreateAlbum(ctx, album); err != nil {
		return song.Album{}, err
	}

	m := event.NewMessage(message.NewAlbumFromDomain(*album))
	if err := ca.pub.Publish(ctx, m, event.AlbumPublishedEvent); err != nil {
		return song.Album{}, err
	}

	return *album, nil
}

func (cs PublishSong) Execute(ctx context.Context, cmd PublishSongCommand) (song.Song, error) {
	album, err := cs.db.GetAlbumByID(ctx, cmd.AlbumID)
	if err != nil {
		return song.Song{}, err
	}

	artist, err := cs.db.GetArtistByID(ctx, album.Artist.ID)
	if err != nil {
		return song.Song{}, err
	}

	album.Artist = artist
	s := &song.Song{
		ID:          uuid.NewString(),
		TrackNumber: cmd.TrackNumber,
		Title:       cmd.Title,
		Album:       album,
		Artist:      artist,
	}

	if err := cs.db.CreateSong(ctx, s); err != nil {
		return song.Song{}, err
	}

	m := event.NewMessage(message.NewSongFromDomain(*s))
	if err := cs.pub.Publish(ctx, m, event.SongPublishedEvent); err != nil {
		return song.Song{}, err
	}

	return *s, nil
}
