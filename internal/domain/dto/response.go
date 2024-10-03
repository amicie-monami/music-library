package dto

type AddSongResponse struct {
	Song *Song `json:"song"`
}

type GetSongTextResponse struct {
	SongID   int64    `json:"song_id"`
	Couplets []string `json:"couplets"`
}

type GetSongDetailsResponse struct {
	Song *SongWithDetails `json:"song"`
}

type GetSongsResponse struct {
	Songs []*SongWithDetails `json:"songs"`
}
