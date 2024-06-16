package library

import (
	"cqrs-sample/pkg/song"
	"cqrs-sample/pkg/user"
	"time"
)

type (
	Playlist struct {
		ID       string      `json:"id"`
		Name     string      `json:"name"`
		Owner    user.User   `json:"owner"`
		Songs    []song.Song `json:"songs"`
		CreateAt time.Time   `json:"create_at"`
	}

	Favorites struct {
		ID    string      `json:"id"`
		Owner user.User   `json:"owner"`
		Songs []song.Song `json:"songs"`
	}
)
