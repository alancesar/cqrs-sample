package handler

import (
	"context"
	"cqrs-sample/pkg/command"
	"cqrs-sample/pkg/handler/presenter"
	"cqrs-sample/pkg/query"
	"cqrs-sample/pkg/song"
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

	GetAlbumsByArtistQuery interface {
		Execute(ctx context.Context, artistID string) ([]query.AlbumResponse, error)
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

	PlaySongCommand interface {
		Execute(ctx context.Context, songID string) error
	}

	ArtistReader struct {
		artistQuery GetArtistQuery
		albumsQuery GetAlbumsByArtistQuery
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
		publishCmd PublishSongCommand
		playCmd    PlaySongCommand
	}
)

func NewArtistReader(artistQuery GetArtistQuery, albumsQuery GetAlbumsByArtistQuery) *ArtistReader {
	return &ArtistReader{
		artistQuery: artistQuery,
		albumsQuery: albumsQuery,
	}
}

func NewArtistWriter(cmd SubscribeArtistCommand) *ArtistWriter {
	return &ArtistWriter{
		cmd: cmd,
	}
}

func NewAlbumReader(albumQuery GetAlbumQuery) *AlbumReader {
	return &AlbumReader{
		q: albumQuery,
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

func NewSongWriter(publishCmd PublishSongCommand, playCmd PlaySongCommand) *SongWriter {
	return &SongWriter{
		publishCmd: publishCmd,
		playCmd:    playCmd,
	}
}

func (ar ArtistReader) Get(w http.ResponseWriter, r *http.Request) {
	artistID := chi.URLParam(r, "artistID")
	artist, err := ar.artistQuery.Execute(r.Context(), artistID)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	writeJsonResponse(w, artist, http.StatusOK)
}

func (ar ArtistReader) GetAlbums(w http.ResponseWriter, r *http.Request) {
	artistID := chi.URLParam(r, "artistID")
	albums, err := ar.albumsQuery.Execute(r.Context(), artistID)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	writeJsonResponse(w, albums, http.StatusOK)

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
	writeJsonResponse(w, response, http.StatusCreated)
}

func (ar AlbumReader) Get(w http.ResponseWriter, r *http.Request) {
	albumID := chi.URLParam(r, "albumID")
	album, err := ar.q.Execute(r.Context(), albumID)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	writeJsonResponse(w, album, http.StatusOK)
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
	writeJsonResponse(w, response, http.StatusCreated)
}

func (sr SongReader) Get(w http.ResponseWriter, r *http.Request) {
	songID := chi.URLParam(r, "songID")
	s, err := sr.q.Execute(r.Context(), songID)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	writeJsonResponse(w, s, http.StatusOK)
}

func (sw SongWriter) Create(w http.ResponseWriter, r *http.Request) {
	var request presenter.PublishSongRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	s, err := sw.publishCmd.Execute(r.Context(), request.ToCommand())
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	response := presenter.NewPublishSongResponseFromDomain(s)
	writeJsonResponse(w, response, http.StatusCreated)
}

func (sw SongWriter) Play(w http.ResponseWriter, r *http.Request) {
	var request presenter.PlaySongRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	if err := sw.playCmd.Execute(r.Context(), request.SongID); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func writeJsonResponse(w http.ResponseWriter, output any, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(output); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
