package model

import "time"

type Song struct {
	ID    int64
	Name  string
	Group string
}

type SongDetail struct {
	ID          int64
	SongID      int64
	ReleaseDate *time.Time
	Text        *string
	Link        *string
}
