package handler

import (
	"context"
	"cqrs-sample/command"
	"cqrs-sample/handler/presenter"
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

	ArtistReader struct {
		q GetArtistQuery
	}

	ArtistWriter struct {
		cmd SubscribeArtistCommand
	}

	AlbumReader struct {
		q GetAlbumQuery
	}

	AlbumWriter struct {
		cmd PublishAlbumCommand
	}

	SongReader struct {
		q GetSongQuery
	}

	SongWriter struct {
		cmd PublishSongCommand
	}
)

func NewArtistReader(q GetArtistQuery) *ArtistReader {
	return &ArtistReader{
		q: q,
	}
}

func NewArtistWriter(cmd SubscribeArtistCommand) *ArtistWriter {
	return &ArtistWriter{
		cmd: cmd,
	}
}

func NewAlbumReader(q GetAlbumQuery) *AlbumReader {
	return &AlbumReader{
		q: q,
	}
}

func NewAlbumWriter(cmd PublishAlbumCommand) *AlbumWriter {
	return &AlbumWriter{
		cmd: cmd,
	}
}

func NewSongReader(q GetSongQuery) *SongReader {
	return &SongReader{
		q: q,
	}
}

func NewSongWriter(cmd PublishSongCommand) *SongWriter {
	return &SongWriter{
		cmd: cmd,
	}
}

func (ar ArtistReader) GetByID(w http.ResponseWriter, r *http.Request) {
	artistID := chi.URLParam(r, "artistID")
	artist, err := ar.q.Execute(r.Context(), artistID)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(artist); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (aw ArtistWriter) Create(w http.ResponseWriter, r *http.Request) {
	var request presenter.SubscribeArtistRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	artist, err := aw.cmd.Execute(r.Context(), request.ToCommand())
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	response := presenter.NewSubscribeArtistResponseFromDomain(artist)
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (ar AlbumReader) GetByID(w http.ResponseWriter, r *http.Request) {
	albumID := chi.URLParam(r, "albumID")
	album, err := ar.q.Execute(r.Context(), albumID)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(album); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (aw AlbumWriter) Create(w http.ResponseWriter, r *http.Request) {
	var request presenter.PublishAlbumRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	album, err := aw.cmd.Execute(r.Context(), request.ToCommand())
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	response := presenter.NewPublishAlbumResponseFromDomain(album)
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (sr SongReader) GetByID(w http.ResponseWriter, r *http.Request) {
	albumID := chi.URLParam(r, "songID")
	s, err := sr.q.Execute(r.Context(), albumID)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(s); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (sw SongWriter) Create(w http.ResponseWriter, r *http.Request) {
	var request presenter.PublishSongRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	s, err := sw.cmd.Execute(r.Context(), request.ToCommand())
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	response := presenter.NewPublishSongResponseFromDomain(s)
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
