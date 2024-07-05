package song

type (
	Gender string

	Song struct {
		ID          string
		TrackNumber int
		Title       string
		Plays       int
		Album       Album
		Artist      Artist
	}

	Album struct {
		ID          string
		Title       string
		Artist      Artist
		ReleaseYear int
		Songs       []Song
	}

	Artist struct {
		ID     string
		Name   string
		Gender Gender
		Albums []Album
	}
)
