package presenter

import (
	"cqrs-sample/pkg/command"
	"cqrs-sample/pkg/song"
)

type (
	SubscribeArtistRequest struct {
		Name   string `json:"name"`
		Gender string `json:"gender"`
	}

	SubscribeArtistResponse struct {
		ID     string `json:"id"`
		Name   string `json:"name"`
		Gender string `json:"gender"`
	}

	PublishAlbumRequest struct {
		Title       string `json:"title"`
		ArtistID    string `json:"artist_id"`
		ReleaseYear int    `json:"release_year"`
	}

	AlbumResponse struct {
		ID          string `json:"id"`
		Title       string `json:"title"`
		ReleaseYear int    `json:"release_year"`
	}

	PublishAlbumResponse struct {
		AlbumResponse
		Artist SubscribeArtistResponse `json:"artist"`
	}

	PublishSongRequest struct {
		TrackNumber int    `json:"track_number"`
		Title       string `json:"title"`
		AlbumID     string `json:"album_id"`
	}

	PublishSongResponse struct {
		ID          string                  `json:"id"`
		TrackNumber int                     `json:"track_number"`
		Title       string                  `json:"title"`
		Album       AlbumResponse           `json:"album"`
		Artist      SubscribeArtistResponse `json:"artist"`
	}

	PlaySongRequest struct {
		SongID string `json:"song_id"`
	}
)

func (r SubscribeArtistRequest) ToCommand() command.SubscribeArtistCommand {
	return command.SubscribeArtistCommand{
		Name:   r.Name,
		Gender: song.Gender(r.Gender),
	}
}

func (r PublishAlbumRequest) ToCommand() command.PublishAlbumCommand {
	return command.PublishAlbumCommand{
		Title:       r.Title,
		ArtistID:    r.ArtistID,
		ReleaseYear: r.ReleaseYear,
	}
}

func (r PublishSongRequest) ToCommand() command.PublishSongCommand {
	return command.PublishSongCommand{
		TrackNumber: r.TrackNumber,
		Title:       r.Title,
		AlbumID:     r.AlbumID,
	}
}

func NewSubscribeArtistResponseFromDomain(artist song.Artist) SubscribeArtistResponse {
	return SubscribeArtistResponse{
		ID:     artist.ID,
		Name:   artist.Name,
		Gender: string(artist.Gender),
	}
}

func NewPublishAlbumResponseFromDomain(album song.Album) PublishAlbumResponse {
	return PublishAlbumResponse{
		AlbumResponse: AlbumResponse{
			ID:          album.ID,
			Title:       album.Title,
			ReleaseYear: album.ReleaseYear,
		},
		Artist: NewSubscribeArtistResponseFromDomain(album.Artist),
	}
}

func NewPublishSongResponseFromDomain(song song.Song) PublishSongResponse {
	return PublishSongResponse{
		ID:          song.ID,
		TrackNumber: song.TrackNumber,
		Title:       song.Title,
		Album: AlbumResponse{
			ID:          song.Album.ID,
			Title:       song.Album.Title,
			ReleaseYear: song.Album.ReleaseYear,
		},
		Artist: NewSubscribeArtistResponseFromDomain(song.Artist),
	}
}
