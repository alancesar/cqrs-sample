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
		Plays       int         `bson:"plays"`
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

func (s Song) ToDomain() song.Song {
	return song.Song{
		ID:          s.ID,
		TrackNumber: s.TrackNumber,
		Title:       s.Title,
		Plays:       s.Plays,
		Album:       s.Album.ToDomain(),
		Artist:      s.Artist.ToDomain(),
	}
}

func (s SongInAlbum) ToDomain() song.Song {
	return song.Song{
		ID:          s.ID,
		TrackNumber: s.TrackNumber,
		Title:       s.Title,
	}
}

func (a Album) ToDomain() song.Album {
	songs := make([]song.Song, len(a.Songs), len(a.Songs))
	for i, s := range a.Songs {
		songs[i] = s.ToDomain()
	}

	return song.Album{
		ID:          a.ID,
		Title:       a.Title,
		Artist:      a.Artist.ToDomain(),
		ReleaseYear: a.ReleaseYear,
		Songs:       songs,
	}
}

func (a AlbumInSong) ToDomain() song.Album {
	return song.Album{
		ID:          a.ID,
		Title:       a.Title,
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
		Plays:       s.Plays,
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
