package dto

type GetSongTextResponse struct {
	SongID   int64    `json:"song_id"`
	Couplets []string `json:"couplets"`
}
