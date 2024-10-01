package dto

type UpdateSongRequest struct {
	Song        *Song        `json:"song"`
	SongDetails *SongDetails `json:"song_details"`
}

type AddSongRequest struct {
	Group string `json:"group"`
	Song  string `json:"song"`
}
