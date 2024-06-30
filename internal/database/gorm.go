package database

import (
	"context"
	"cqrs-sample/internal/database/model"
	"cqrs-sample/pkg/song"
	"gorm.io/gorm"
)

type (
	Gorm struct {
		db *gorm.DB
	}
)

func NewGorm(db *gorm.DB) *Gorm {
	_ = db.AutoMigrate(
		&model.Artist{},
		&model.Album{},
		&model.Song{},
	)

	return &Gorm{
		db: db,
	}
}

func (g Gorm) CreateArtist(ctx context.Context, artist *song.Artist) error {
	m := model.NewArtistFromDomain(*artist)
	tx := g.db.WithContext(ctx).Create(&m)
	if err := tx.Error; err != nil {
		return err
	}

	artist.ID = m.ID
	return nil
}

func (g Gorm) GetArtistByID(ctx context.Context, id string) (song.Artist, error) {
	m := model.Artist{ID: id}
	tx := g.db.WithContext(ctx).First(&m)
	if err := tx.Error; err != nil {
		return song.Artist{}, err
	}

	return m.ToDomain(), nil
}

func (g Gorm) CreateAlbum(ctx context.Context, album *song.Album) error {
	m := model.NewAlbumFromDomain(*album)
	tx := g.db.WithContext(ctx).Create(&m)
	if err := tx.Error; err != nil {
		return err
	}

	album.ID = m.ID
	return nil
}

func (g Gorm) GetAlbumByID(ctx context.Context, id string) (song.Album, error) {
	m := model.Album{ID: id}
	tx := g.db.WithContext(ctx).First(&m)
	if err := tx.Error; err != nil {
		return song.Album{}, err
	}

	return m.ToDomain(), nil
}

func (g Gorm) CreateSong(ctx context.Context, s *song.Song) error {
	m := model.NewSongFromDomain(*s)
	tx := g.db.WithContext(ctx).Create(&m)
	if err := tx.Error; err != nil {
		return err
	}

	s.ID = m.ID
	return nil
}
