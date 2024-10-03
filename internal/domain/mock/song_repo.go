package mock

import (
	"context"
	"fmt"

	"github.com/amicie-monami/music-library/internal/domain/dto"
	"github.com/amicie-monami/music-library/internal/domain/model"
)

var (
	ValidSongName         = "Song12"
	ValidGroupName        = "Group12"
	ValidSongID           = int64(12)
	SongIDWithoutTextData = int64(89)
)

type SongRepo struct{}

///

func (m *SongRepo) Create(ctx context.Context, song *model.Song) error {
	if song.Name == "" || song.Group == "" {
		return &dto.Error{Code: 500, Message: "internal server error", Details: "database", DebugMsg: ""}
	}
	song.ID = 1
	return nil
}

///

func (m *SongRepo) GetSongWithDetails(ctx context.Context, group string, song string) (*dto.SongWithDetails, error) {

	if group == ValidGroupName && song == ValidSongName {
		return &dto.SongWithDetails{}, nil
	}

	return nil, &dto.Error{Code: 400, Message: "song not found", Details: fmt.Sprintf("name=%s group=%s", song, group)}
}

///

func (m *SongRepo) Delete(ctx context.Context, id int64) error {

	fmt.Println("song id:", id)

	if id != ValidSongID {
		return &dto.Error{Code: 400, Message: "song doesn't exist"}
	}

	return nil
}

///

func (m *SongRepo) GetSongText(ctx context.Context, id int64) (*string, error) {

	if id == SongIDWithoutTextData {
		return nil, &dto.Error{Code: 400, Message: "no info about song text", Details: fmt.Sprintf("id=%d", id)}
	}

	if id != ValidSongID {
		return nil, &dto.Error{Code: 400, Message: "song not found"}
	}

	songText := "boundaries, key..."
	return &songText, nil
}

///

func (m *SongRepo) GetSongs(ctx context.Context, aggregation map[string]any) ([]*dto.SongWithDetails, error) {
	return nil, nil
}

///

func (m *SongRepo) Tx(ctx context.Context, txActions func() error) error {
	return txActions()
}

func (m *SongRepo) UpdateSong(ctx context.Context, song *model.Song) error {
	if song.ID != ValidSongID {
		return &dto.Error{Code: 400}
	}
	return nil
}

func (m *SongRepo) UpdateSongDetails(ctx context.Context, details *model.SongDetail) error {
	if details.SongID != ValidSongID {
		return &dto.Error{Code: 400}
	}
	return nil
}
