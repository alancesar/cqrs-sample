package command

import (
	"context"
	"cqrs-sample/internal/database"
	"cqrs-sample/pkg/event"
	"cqrs-sample/pkg/song"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"path/filepath"
	"reflect"
	"testing"
)

func setupDatabase(t *testing.T) *database.Gorm {
	testDBPath := filepath.Join(t.TempDir(), "database.sqlite")
	db, err := gorm.Open(sqlite.Open(testDBPath), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	sqliteDB := database.NewGorm(db)
	return sqliteDB
}

func Test_Create_Artist_Album_And_Song(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	// Arrange
	db := setupDatabase(t)
	publisher := &fakePublisher{}
	ctx := context.Background()
	artistSubscriber := NewSubscribeArtist(db, publisher)
	albumPublisher := NewPublishAlbum(db, publisher)
	songPublisher := NewPublishSong(db, publisher)

	// Act
	artist, err := artistSubscriber.Execute(ctx, SubscribeArtistCommand{
		Name:   "Some Artist",
		Gender: "Some Gender",
	})
	if err != nil {
		t.Fatal(err)
	}

	album, err := albumPublisher.Execute(ctx, PublishAlbumCommand{
		Title:       "Some Album",
		ArtistID:    artist.ID,
		ReleaseYear: 2024,
	})
	if err != nil {
		t.Fatal(err)
	}

	s, err := songPublisher.Execute(ctx, PublishSongCommand{
		TrackNumber: 1,
		Title:       "Some Song",
		AlbumID:     album.ID,
	})
	if err != nil {
		log.Fatal(err)
	}

	// Assert
	wantSong := song.Song{
		ID:          s.ID,
		TrackNumber: 1,
		Title:       "Some Song",
		Album: song.Album{
			ID:    album.ID,
			Title: "Some Album",
			Artist: song.Artist{
				ID:     artist.ID,
				Name:   "Some Artist",
				Gender: "Some Gender",
				Albums: make([]song.Album, 0),
			},
			ReleaseYear: 2024,
		},
		Artist: song.Artist{
			ID:     artist.ID,
			Name:   "Some Artist",
			Gender: "Some Gender",
			Albums: make([]song.Album, 0),
		},
	}
	if !reflect.DeepEqual(s, wantSong) {
		t.Errorf("\n\tgot = %+v\n\twant= %+v", s, wantSong)
	}
}

type (
	fakePublisher struct{}
)

func (f fakePublisher) Publish(_ context.Context, _ event.Message, _ event.Event) error {
	return nil
}
