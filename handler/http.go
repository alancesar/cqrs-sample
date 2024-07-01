package handler

import (
	"context"
	"cqrs-sample/query"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"net/http"
)

type (
	GetArtistQuery interface {
		Execute(ctx context.Context, id string) (query.ArtistResponse, error)
	}

	GetAlbumQuery interface {
		Execute(ctx context.Context, id string) (query.AlbumResponse, error)
	}

	GetSongQuery interface {
		Execute(ctx context.Context, id string) (query.SongResponse, error)
	}

	GetArtist struct {
		q GetArtistQuery
	}

	GetAlbum struct {
		q GetAlbumQuery
	}

	GetSong struct {
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

func (ga GetArtist) Handle(w http.ResponseWriter, r *http.Request) {
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

func (ga GetAlbum) Handle(w http.ResponseWriter, r *http.Request) {
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

func (ga GetSong) Handle(w http.ResponseWriter, r *http.Request) {
	albumID := chi.URLParam(r, "songID")
	song, err := ga.q.Execute(r.Context(), albumID)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(song); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
