package model

type SongDetail struct {
	ID          int64   `db:"id"`
	SongID      int64   `db:"song_id"`
	ReleaseDate *string `db:"release_date"`
	Text        *string `db:"text"`
	Link        *string `db:"link"`
}
