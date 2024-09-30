package dto

type UpdateSong struct {
	Group string `json:"group,omitempty"`
	Title string `json:"title,omitempty"`
}

type SongDetails struct {
	Text        *string `json:"text,omitempty"`
	Link        *string `json:"link,omitempty"`
	ReleaseDate *string `json:"release_date,omitempty"`
}

type UpdateSongRequest struct {
	Song        *UpdateSong  `json:"song"`
	SongDetails *SongDetails `json:"song_details"`
}

type SongWithDetails struct {
	ID          int64   `json:"song_id,omitempty" db:"song_id"`
	Group       string  `json:"group_name,omitempty" db:"group_name"`
	Title       string  `json:"song_name,omitempty" db:"song_name"`
	ReleaseDate *string `json:"release_date,omitempty" db:"release_date"`
	Text        *string `json:"text,omitempty" db:"text"`
	Link        *string `json:"link,omitempty" db:"link"`
}
