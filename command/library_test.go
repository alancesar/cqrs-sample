package command

import (
	"context"
	"cqrs-sample/internal/database"
	"cqrs-sample/pkg/event"
	"cqrs-sample/pkg/song"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"os"
	"reflect"
	"testing"
)

const (
	fixedArtistID = "b979883e-6195-4440-b7a2-5204e6116c55"
	fixedAlbumID  = "ead31e54-cb29-4c1e-ba18-de75c1d99436"
	fixedSongID   = "fc6fa5a1-5be5-419e-aa5f-93ad34586085"

	testDBPath = "testdata/database.sqlite"
)

func setupEmptyDatabase(_ testing.TB) func(tb testing.TB) {
	_ = os.Remove(testDBPath)

	return func(tb testing.TB) {
		_ = os.Remove(testDBPath)
	}
}

func setupDatabaseWithArtist(t testing.TB) func(tb testing.TB) {
	_ = os.Remove(testDBPath)
	db, err := gorm.Open(sqlite.Open(testDBPath), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}

	sqliteDB := database.NewGorm(db)
	if _, err := sqliteDB.CreateArtist(context.Background(), song.Artist{
		ID:     fixedArtistID,
		Name:   "Some Artist",
		Gender: "Rock",
	}); err != nil {
		t.Fatal(err)
	}

	return func(tb testing.TB) {
		_ = os.Remove(testDBPath)
	}
}

func setupDatabaseWithArtistAndAlbum(t testing.TB) func(tb testing.TB) {
	_ = os.Remove(testDBPath)
	db, err := gorm.Open(sqlite.Open(testDBPath), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}

	sqliteDB := database.NewGorm(db)
	artist, err := sqliteDB.CreateArtist(context.Background(), song.Artist{
		ID:     fixedArtistID,
		Name:   "Some Artist",
		Gender: "Rock",
	})

	if err != nil {
		t.Fatal(err)
	}

	if _, err := sqliteDB.CreateAlbum(context.Background(), song.Album{
		ID:          fixedAlbumID,
		Title:       "Some Album",
		Artist:      artist,
		ReleaseYear: 2024,
	}); err != nil {
		t.Fatal(err)
	}

	return func(tb testing.TB) {
		_ = os.Remove(testDBPath)
	}
}

func TestSubscribeArtist_Execute(t *testing.T) {
	teardownSuite := setupEmptyDatabase(t)
	defer teardownSuite(t)

	db, err := gorm.Open(sqlite.Open(testDBPath), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}

	sqliteDB := database.NewGorm(db)

	type fields struct {
		db  ArtistDatabase
		pub Publisher
	}
	type args struct {
		ctx     context.Context
		request SubscribeArtistCommand
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    song.Artist
		wantErr bool
	}{
		{
			name: "Should create artist properly",
			fields: fields{
				db:  sqliteDB,
				pub: &fakePublisher{},
			},
			args: args{
				ctx: context.Background(),
				request: SubscribeArtistCommand{
					Name:   "Some Artist",
					Gender: "Rock",
				},
			},
			want: song.Artist{
				ID:     fixedArtistID,
				Name:   "Some Artist",
				Gender: "Rock",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ca := NewSubscribeArtist(tt.fields.db, tt.fields.pub, newFakeUUIDGenerator(fixedArtistID))
			got, err := ca.Execute(tt.args.ctx, tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Execute() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPublishAlbum_Execute(t *testing.T) {
	teardownSuite := setupDatabaseWithArtist(t)
	defer teardownSuite(t)

	db, err := gorm.Open(sqlite.Open(testDBPath), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}

	sqliteDB := database.NewGorm(db)

	type fields struct {
		db  AlbumDatabase
		pub Publisher
	}
	type args struct {
		ctx     context.Context
		request PublishAlbumCommand
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    song.Album
		wantErr bool
	}{
		{
			name: "Should create album properly",
			fields: fields{
				db:  sqliteDB,
				pub: &fakePublisher{},
			},
			args: args{
				ctx: context.Background(),
				request: PublishAlbumCommand{
					Title:       "Some Album",
					ArtistID:    fixedArtistID,
					ReleaseYear: 2024,
				},
			},
			want: song.Album{
				ID:    fixedAlbumID,
				Title: "Some Album",
				Artist: song.Artist{
					ID:     fixedArtistID,
					Name:   "Some Artist",
					Gender: "Rock",
				},
				ReleaseYear: 2024,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ca := NewPublishAlbum(tt.fields.db, tt.fields.pub, newFakeUUIDGenerator(fixedAlbumID))
			got, err := ca.Execute(tt.args.ctx, tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Execute() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPublishSong_Execute(t *testing.T) {
	teardownSuite := setupDatabaseWithArtistAndAlbum(t)
	defer teardownSuite(t)

	db, err := gorm.Open(sqlite.Open(testDBPath), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}

	sqliteDB := database.NewGorm(db)

	type fields struct {
		db  SongDatabase
		pub Publisher
	}
	type args struct {
		ctx     context.Context
		request PublishSongCommand
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    song.Song
		wantErr bool
	}{
		{
			name: "Should create song properly",
			fields: fields{
				db:  sqliteDB,
				pub: &fakePublisher{},
			},
			args: args{
				ctx: context.Background(),
				request: PublishSongCommand{
					TrackNumber: 1,
					Title:       "Some Song",
					AlbumID:     fixedAlbumID,
				},
			},
			want: song.Song{
				ID:          fixedSongID,
				TrackNumber: 1,
				Title:       "Some Song",
				Album: song.Album{
					ID:    fixedAlbumID,
					Title: "Some Album",
					Artist: song.Artist{
						ID:     fixedArtistID,
						Name:   "Some Artist",
						Gender: "Rock",
					},
					ReleaseYear: 2024,
				},
				Artist: song.Artist{
					ID:     fixedArtistID,
					Name:   "Some Artist",
					Gender: "Rock",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ca := NewPublishSong(tt.fields.db, tt.fields.pub, newFakeUUIDGenerator(fixedSongID))
			got, err := ca.Execute(tt.args.ctx, tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Execute() got = %v, want %v", got, tt.want)
			}
		})
	}
}

type (
	fakeUUIDGenerator struct {
		uuid string
	}

	fakePublisher struct{}
)

func newFakeUUIDGenerator(uuid string) *fakeUUIDGenerator {
	return &fakeUUIDGenerator{
		uuid: uuid,
	}
}

func (r fakeUUIDGenerator) Generate() string {
	return r.uuid
}

func (f fakePublisher) Publish(_ context.Context, _ event.Message, _ event.Event) error {
	return nil
}
