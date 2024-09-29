package model

type Song struct {
	ID    int64
	Name  string
	Group string
}

type SongDetail struct {
	ID          int64
	SongID      int64
	ReleaseDate *string
	Text        *string
	Link        *string
}
