package query

import "cqrs-sample/pkg/song"

type (
	AlbumInSongResponse struct {
		ID          string `json:"id"`
		Title       string `json:"title"`
		ReleaseYear int    `json:"release_year"`
	}

	AlbumResponse struct {
		ID          string                `json:"id"`
		Title       string                `json:"title"`
		Artist      ArtistResponse        `json:"artist"`
		ReleaseYear int                   `json:"release_year"`
		Songs       []SongInAlbumResponse `json:"songs"`
	}

	ArtistResponse struct {
		ID     string `json:"id"`
		Name   string `json:"name"`
		Gender string `json:"gender"`
	}

	SongResponse struct {
		ID          string              `json:"id"`
		TrackNumber int                 `json:"track_number"`
		Title       string              `json:"title"`
		Plays       int                 `json:"plays"`
		Album       AlbumInSongResponse `json:"album"`
		Artist      ArtistResponse      `json:"artist"`
	}

	SongInAlbumResponse struct {
		ID          string `json:"id"`
		TrackNumber int    `json:"track_number"`
		Title       string `json:"title"`
	}
)

func NewAlbumInSongResponseFromDomain(album song.Album) AlbumInSongResponse {
	return AlbumInSongResponse{
		ID:          album.ID,
		Title:       album.Title,
		ReleaseYear: album.ReleaseYear,
	}
}

func NewSongInAlbumResponseFromDomain(s song.Song) SongInAlbumResponse {
	return SongInAlbumResponse{
		ID:          s.ID,
		TrackNumber: s.TrackNumber,
		Title:       s.Title,
	}
}

func NewAlbumResponseFromDomain(album song.Album) AlbumResponse {
	songs := make([]SongInAlbumResponse, len(album.Songs), len(album.Songs))
	for i, s := range album.Songs {
		songs[i] = NewSongInAlbumResponseFromDomain(s)
	}

	return AlbumResponse{
		ID:          album.ID,
		Title:       album.Title,
		Artist:      NewArtistResponseFromDomain(album.Artist),
		ReleaseYear: album.ReleaseYear,
		Songs:       songs,
	}
}

func NewArtistResponseFromDomain(artist song.Artist) ArtistResponse {
	return ArtistResponse{
		ID:     artist.ID,
		Name:   artist.Name,
		Gender: string(artist.Gender),
	}
}

func NewSongResponseFromDomain(s song.Song) SongResponse {
	return SongResponse{
		ID:          s.ID,
		TrackNumber: s.TrackNumber,
		Title:       s.Title,
		Plays:       s.Plays,
		Album:       NewAlbumInSongResponseFromDomain(s.Album),
		Artist:      NewArtistResponseFromDomain(s.Artist),
	}
}
