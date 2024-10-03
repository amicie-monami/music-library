package dto

type Song struct {
	ID    int64  `json:"song_id,omitempty"`
	Group string `json:"group,omitempty"`
	Name  string `json:"title,omitempty"`
}

type SongDetails struct {
	Text        *string `json:"text,omitempty"`
	Link        *string `json:"link,omitempty"`
	ReleaseDate string  `json:"release_date,omitempty"`
}

type SongWithDetails struct {
	ID          int64   `json:"song_id,omitempty" db:"song_id"`
	Group       string  `json:"group,omitempty" db:"group_name"`
	Title       string  `json:"song,omitempty" db:"song_name"`
	ReleaseDate *string `json:"release_date,omitempty" db:"release_date"`
	Text        *string `json:"text,omitempty" db:"text"`
	Link        *string `json:"link,omitempty" db:"link"`
}
