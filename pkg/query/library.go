package query

import (
	"context"
	"cqrs-sample/pkg/song"
)

type (
	AlbumDatabase interface {
		GetAlbumByID(ctx context.Context, id string) (song.Album, error)
	}

	ArtistDatabase interface {
		GetArtistByID(ctx context.Context, id string) (song.Artist, error)
		GetAlbumsByArtistID(ctx context.Context, artistID string) ([]song.Album, error)
	}

	SongDatabase interface {
		GetSongByID(ctx context.Context, id string) (song.Song, error)
	}

	GetAlbum struct {
		db AlbumDatabase
	}

	GetAlbumsByArtist struct {
		db ArtistDatabase
	}

	GetArtist struct {
		db ArtistDatabase
	}

	GetSong struct {
		db SongDatabase
	}
)

func NewGetAlbum(db AlbumDatabase) *GetAlbum {
	return &GetAlbum{
		db: db,
	}
}

func NewGetAlbumsByArtist(db ArtistDatabase) *GetAlbumsByArtist {
	return &GetAlbumsByArtist{
		db: db,
	}
}

func NewGetArtist(db ArtistDatabase) *GetArtist {
	return &GetArtist{
		db: db,
	}
}

func NewGetSong(db SongDatabase) *GetSong {
	return &GetSong{
		db: db,
	}
}

func (ga GetAlbum) Execute(ctx context.Context, id string) (AlbumResponse, error) {
	album, err := ga.db.GetAlbumByID(ctx, id)
	if err != nil {
		return AlbumResponse{}, err
	}

	return NewAlbumResponseFromDomain(album), nil
}

func (ga GetAlbumsByArtist) Execute(ctx context.Context, artistID string) ([]AlbumResponse, error) {
	albums, err := ga.db.GetAlbumsByArtistID(ctx, artistID)
	if err != nil {
		return nil, err
	}

	output := make([]AlbumResponse, len(albums), len(albums))
	for i, album := range albums {
		output[i] = NewAlbumResponseFromDomain(album)
	}

	return output, nil
}

func (ga GetArtist) Execute(ctx context.Context, id string) (ArtistResponse, error) {
	artist, err := ga.db.GetArtistByID(ctx, id)
	if err != nil {
		return ArtistResponse{}, err
	}

	return NewArtistResponseFromDomain(artist), nil
}

func (gs GetSong) Execute(ctx context.Context, id string) (SongResponse, error) {
	s, err := gs.db.GetSongByID(ctx, id)
	if err != nil {
		return SongResponse{}, err
	}

	return NewSongResponseFromDomain(s), err
}
