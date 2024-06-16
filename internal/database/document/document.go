package document

import (
	"cqrs-sample/pkg/song"
)

type (
	Song struct {
		ID          string      `bson:"_id"`
		TrackNumber int         `bson:"track_number"`
		Title       string      `bson:"title"`
		Album       AlbumInSong `bson:"album"`
		Artist      Artist      `bson:"artist"`
	}

	SongInAlbum struct {
		ID          string `bson:"_id"`
		TrackNumber int    `bson:"track_number"`
		Title       string `bson:"title"`
	}

	Album struct {
		ID          string        `bson:"_id"`
		Title       string        `bson:"title"`
		Artist      Artist        `bson:"artist"`
		ReleaseYear int           `bson:"release_year"`
		Songs       []SongInAlbum `bson:"songs"`
	}

	AlbumInSong struct {
		ID          string `bson:"_id"`
		Title       string `bson:"title"`
		ReleaseYear int    `bson:"release_year"`
	}

	Artist struct {
		ID     string `bson:"_id"`
		Name   string `bson:"name"`
		Gender string `bson:"gender"`
	}
)

func NewArtistFromDomain(a song.Artist) Artist {
	return Artist{
		ID:     a.ID,
		Name:   a.Name,
		Gender: string(a.Gender),
	}
}

func NewAlbumFromDomain(a song.Album) Album {
	songs := make([]SongInAlbum, len(a.Songs), len(a.Songs))
	for i, s := range a.Songs {
		songs[i] = NewSongInAlbumFromDomain(s)
	}

	return Album{
		ID:          a.ID,
		Title:       a.Title,
		ReleaseYear: a.ReleaseYear,
		Artist:      NewArtistFromDomain(a.Artist),
		Songs:       songs,
	}
}

func NewSongFromDomain(s song.Song) Song {
	return Song{
		ID:          s.ID,
		TrackNumber: s.TrackNumber,
		Title:       s.Title,
		Album:       NewAlbumInSongFromDomain(s.Album),
		Artist:      NewArtistFromDomain(s.Artist),
	}
}

func NewSongInAlbumFromDomain(s song.Song) SongInAlbum {
	return SongInAlbum{
		ID:          s.ID,
		TrackNumber: s.TrackNumber,
		Title:       s.Title,
	}
}

func NewAlbumInSongFromDomain(a song.Album) AlbumInSong {
	songs := make([]Song, len(a.Songs), len(a.Songs))
	for i, s := range a.Songs {
		songs[i] = NewSongFromDomain(s)
	}

	return AlbumInSong{
		ID:          a.ID,
		Title:       a.Title,
		ReleaseYear: a.ReleaseYear,
	}
}
