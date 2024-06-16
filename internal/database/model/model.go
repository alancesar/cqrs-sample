package model

import "cqrs-sample/pkg/song"

type (
	Song struct {
		ID          string `gorm:"primarykey"`
		TrackNumber int
		Title       string
		Album       Album
		AlbumID     string
		Artist      Artist
		ArtistID    string
	}

	Album struct {
		ID          string `gorm:"primarykey"`
		Title       string
		Artist      Artist
		ArtistID    string
		ReleaseYear int
		Songs       []Song
	}

	Artist struct {
		ID     string `gorm:"primarykey"`
		Name   string
		Gender string
		Albums []Album
	}
)

func (a Artist) ToDomain() song.Artist {
	albums := make([]song.Album, len(a.Albums), len(a.Albums))
	for i, album := range a.Albums {
		albums[i] = album.ToDomain()
	}

	return song.Artist{
		ID:     a.ID,
		Name:   a.Name,
		Gender: song.Gender(a.Gender),
		Albums: albums,
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

func (s Song) ToDomain() song.Song {
	return song.Song{
		ID:          s.ID,
		TrackNumber: s.TrackNumber,
		Title:       s.Title,
		Album:       s.Album.ToDomain(),
		Artist:      s.Artist.ToDomain(),
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
	songs := make([]Song, len(a.Songs), len(a.Songs))
	for i, s := range a.Songs {
		songs[i] = NewSongFromDomain(s)
	}

	return Album{
		ID:          a.ID,
		Title:       a.Title,
		Artist:      NewArtistFromDomain(a.Artist),
		ArtistID:    a.Artist.ID,
		ReleaseYear: a.ReleaseYear,
		Songs:       songs,
	}
}

func NewSongFromDomain(s song.Song) Song {
	return Song{
		ID:          s.ID,
		TrackNumber: s.TrackNumber,
		Title:       s.Title,
		Album:       NewAlbumFromDomain(s.Album),
		AlbumID:     s.Album.ID,
		Artist:      NewArtistFromDomain(s.Artist),
		ArtistID:    s.Artist.ID,
	}
}
