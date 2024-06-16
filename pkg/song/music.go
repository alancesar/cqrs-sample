package song

type (
	Gender string

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
		Songs       []Song `json:"songs"`
	}

	Artist struct {
		ID     string  `json:"id"`
		Name   string  `json:"name"`
		Gender Gender  `json:"gender"`
		Albums []Album `json:"albums"`
	}
)
