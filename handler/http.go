package handler

import (
	"context"
	"cqrs-sample/command"
	"cqrs-sample/pkg/song"
	"cqrs-sample/query"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"net/http"
)

type (
	GetArtistQuery interface {
		Execute(ctx context.Context, id string) (query.ArtistResponse, error)
	}

	SubscribeArtistCommand interface {
		Execute(ctx context.Context, artist command.SubscribeArtistCommand) (song.Artist, error)
	}

	GetAlbumQuery interface {
		Execute(ctx context.Context, id string) (query.AlbumResponse, error)
	}

	PublishAlbumCommand interface {
		Execute(ctx context.Context, cmd command.PublishAlbumCommand) (song.Album, error)
	}

	GetSongQuery interface {
		Execute(ctx context.Context, id string) (query.SongResponse, error)
	}

	PublishSongCommand interface {
		Execute(ctx context.Context, cmd command.PublishSongCommand) (song.Song, error)
	}

	Artist struct {
		q GetArtistQuery
	}

	Album struct {
		q GetAlbumQuery
	}

	Song struct {
		q GetSongQuery
	}
)

func NewGetArtist(q GetArtistQuery) *GetArtist {
	return &GetArtist{
		q: q,
	}
}

func NewGetAlbum(q GetAlbumQuery) *GetAlbum {
	return &GetAlbum{
		q: q,
	}
}

func NewGetSong(q GetSongQuery) *GetSong {
	return &GetSong{
		q: q,
	}
}

func (ga Artist) GetByID(w http.ResponseWriter, r *http.Request) {
	artistID := chi.URLParam(r, "artistID")
	artist, err := ga.q.Execute(r.Context(), artistID)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(artist); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (ga Album) GetByID(w http.ResponseWriter, r *http.Request) {
	albumID := chi.URLParam(r, "albumID")
	album, err := ga.q.Execute(r.Context(), albumID)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(album); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (ga Song) GetByID(w http.ResponseWriter, r *http.Request) {
	albumID := chi.URLParam(r, "songID")
	s, err := ga.q.Execute(r.Context(), albumID)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(s); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
