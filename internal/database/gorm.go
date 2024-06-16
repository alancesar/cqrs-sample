package database

import (
	"context"
	"cqrs-sample/internal/database/model"
	"cqrs-sample/pkg/song"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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

func (g Gorm) CreateArtist(ctx context.Context, artist song.Artist) (song.Artist, error) {
	m := model.NewArtistFromDomain(artist)
	tx := g.db.WithContext(ctx).Create(&m)
	if err := tx.Error; err != nil {
		return song.Artist{}, err
	}

	return m.ToDomain(), nil
}

func (g Gorm) CreateAlbum(ctx context.Context, album song.Album) (song.Album, error) {
	m := model.NewAlbumFromDomain(album)
	tx := g.db.WithContext(ctx).Preload(clause.Associations).Create(&m).First(&m)
	if err := tx.Error; err != nil {
		return song.Album{}, err
	}

	return m.ToDomain(), nil
}

func (g Gorm) CreateSong(ctx context.Context, s song.Song) (song.Song, error) {
	albumModel := model.Album{
		ID: s.Album.ID,
	}

	tx := g.db.WithContext(ctx).First(&albumModel)
	if err := tx.Error; err != nil {
		return song.Song{}, err
	}

	artistModel := model.Artist{
		ID: albumModel.ArtistID,
	}
	tx = g.db.WithContext(ctx).First(&artistModel)
	if err := tx.Error; err != nil {
		return song.Song{}, err
	}

	albumModel.Artist = artistModel
	songModel := model.NewSongFromDomain(s)
	songModel.Album = albumModel
	songModel.Artist = artistModel
	tx = g.db.WithContext(ctx).Create(&songModel).First(&songModel)
	if err := tx.Error; err != nil {
		return song.Song{}, err
	}

	return songModel.ToDomain(), nil
}
