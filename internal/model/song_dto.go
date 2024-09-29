package model

type SongWithDetailsDTO struct {
	ID          int64   `json:"song_id,omitempty" db:"song_id"`
	Group       string  `json:"group_name,omitempty" db:"group_name"`
	Title       string  `json:"song_name,omitempty" db:"song_name"`
	ReleaseDate *string `json:"release_date,omitempty" db:"release_date"`
	Text        *string `json:"text,omitempty" db:"text"`
	Link        *string `json:"link,omitempty" db:"link"`
}
