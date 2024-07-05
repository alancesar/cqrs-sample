package message

import (
	"cqrs-sample/pkg/song"
)

type (
	Song struct {
		ID          string `json:"id"`
		TrackNumber int    `json:"track_number"`
		Title       string `json:"title"`
		Album       Album  `json:"album"`
		Artist      Artist `json:"artist"`
	}

	Album struct {
		ID          string `json:"id"`
		Title       string `json:"title"`
		Artist      Artist `json:"artist"`
		ReleaseYear int    `json:"release_year"`
	}

	Artist struct {
		ID     string `json:"id"`
		Name   string `json:"name"`
		Gender string `json:"gender"`
	}

	PlaySong struct {
		SongID string `json:"song_id"`
	}
)

func (s Song) ToDomain() song.Song {
	return song.Song{
		ID:          s.ID,
		TrackNumber: s.TrackNumber,
		Title:       s.Title,
		Album:       s.Album.ToDomain(),
		Artist:      s.Artist.ToDomain(),
	}
}

func (a Album) ToDomain() song.Album {
	return song.Album{
		ID:          a.ID,
		Title:       a.Title,
		Artist:      a.Artist.ToDomain(),
		ReleaseYear: a.ReleaseYear,
	}
}

func (a Artist) ToDomain() song.Artist {
	return song.Artist{
		ID:     a.ID,
		Name:   a.Name,
		Gender: song.Gender(a.Gender),
	}
}

func NewSongFromDomain(s song.Song) Song {
	return Song{
		ID:          s.ID,
		TrackNumber: s.TrackNumber,
		Title:       s.Title,
		Album:       NewAlbumFromDomain(s.Album),
		Artist:      NewArtistFromDomain(s.Artist),
	}
}

func NewAlbumFromDomain(album song.Album) Album {
	return Album{
		ID:          album.ID,
		Title:       album.Title,
		Artist:      NewArtistFromDomain(album.Artist),
		ReleaseYear: album.ReleaseYear,
	}
}

func NewArtistFromDomain(artist song.Artist) Artist {
	return Artist{
		ID:     artist.ID,
		Name:   artist.Name,
		Gender: string(artist.Gender),
	}
}
