package dto

type UpdateSongRequest struct {
	Group       string `json:"group,omitempty"`
	Song        string `json:"song,omitempty"`
	Text        string `json:"text,omitempty"`
	Link        string `json:"link,omitempty"`
	ReleaseDate string `json:"release_date,omitempty"`
}

type AddSongRequest struct {
	Group string `json:"group"`
	Song  string `json:"song"`
}

type GetSongDetailsRequest struct {
	Group string `json:"group"`
	Song  string `json:"song"`
}
