package command

import (
	"context"
	"cqrs-sample/pkg/event"
	"cqrs-sample/pkg/song"
)

type (
	ArtistDatabase interface {
		CreateArtist(ctx context.Context, artist song.Artist) (song.Artist, error)
	}

	AlbumDatabase interface {
		CreateAlbum(ctx context.Context, album song.Album) (song.Album, error)
	}

	SongDatabase interface {
		CreateSong(ctx context.Context, s song.Song) (song.Song, error)
	}

	Publisher interface {
		Publish(ctx context.Context, ev event.Message, key event.Event) error
	}

	UUIDGenerator interface {
		Generate() string
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
		db   ArtistDatabase
		pub  Publisher
		uuid UUIDGenerator
	}

	PublishAlbum struct {
		db   AlbumDatabase
		pub  Publisher
		uuid UUIDGenerator
	}

	PublishSong struct {
		db   SongDatabase
		pub  Publisher
		uuid UUIDGenerator
	}
)

func NewSubscribeArtist(db ArtistDatabase, pub Publisher, uuid UUIDGenerator) *SubscribeArtist {
	return &SubscribeArtist{
		db:   db,
		pub:  pub,
		uuid: uuid,
	}
}

func NewPublishAlbum(db AlbumDatabase, pub Publisher, uuid UUIDGenerator) *PublishAlbum {
	return &PublishAlbum{
		db:   db,
		pub:  pub,
		uuid: uuid,
	}
}

func NewPublishSong(db SongDatabase, pub Publisher, uuid UUIDGenerator) *PublishSong {
	return &PublishSong{
		db:   db,
		pub:  pub,
		uuid: uuid,
	}
}

func (ca SubscribeArtist) Execute(ctx context.Context, cmd SubscribeArtistCommand) (song.Artist, error) {
	artist, err := ca.db.CreateArtist(ctx, song.Artist{
		ID:     ca.uuid.Generate(),
		Name:   cmd.Name,
		Gender: cmd.Gender,
	})
	if err != nil {
		return song.Artist{}, err
	}

	message := event.NewMessage(event.NewArtistFromDomain(artist))
	if err := ca.pub.Publish(ctx, message, event.ArtistSubscribedEvent); err != nil {
		return song.Artist{}, err
	}

	return artist, nil
}

func (ca PublishAlbum) Execute(ctx context.Context, cmd PublishAlbumCommand) (song.Album, error) {
	album, err := ca.db.CreateAlbum(ctx, song.Album{
		ID:    ca.uuid.Generate(),
		Title: cmd.Title,
		Artist: song.Artist{
			ID: cmd.ArtistID,
		},
		ReleaseYear: cmd.ReleaseYear,
	})
	if err != nil {
		return song.Album{}, err
	}

	message := event.NewMessage(event.NewAlbumFromDomain(album))
	if err := ca.pub.Publish(ctx, message, event.AlbumPublishedEvent); err != nil {
		return song.Album{}, err
	}

	return album, nil
}

func (cs PublishSong) Execute(ctx context.Context, cmd PublishSongCommand) (song.Song, error) {
	s, err := cs.db.CreateSong(ctx, song.Song{
		ID:          cs.uuid.Generate(),
		TrackNumber: cmd.TrackNumber,
		Title:       cmd.Title,
		Album: song.Album{
			ID: cmd.AlbumID,
		},
	})
	if err != nil {
		return song.Song{}, err
	}

	message := event.NewMessage(event.NewSongFromDomain(s))
	if err := cs.pub.Publish(ctx, message, event.SongPublishedEvent); err != nil {
		return song.Song{}, err
	}

	return s, nil
}
