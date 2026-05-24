package usecase

import "context"

type MusicTrackOutput struct {
	Provider    string
	ID          string
	Title       string
	ArtistName  string
	PreviewURL  string
	ExternalURL string
	ArtworkURL  string
}

type ListMusicCollectionTracksInput struct {
	Provider     string
	AccessToken  string
	CollectionID string
}

type ListMusicCollectionTracksUsecase struct {
	registry MusicSourceRegistry
}

func NewListMusicCollectionTracksUsecase(registry MusicSourceRegistry) *ListMusicCollectionTracksUsecase {
	return &ListMusicCollectionTracksUsecase{registry: registry}
}

func (u *ListMusicCollectionTracksUsecase) Execute(
	ctx context.Context,
	input ListMusicCollectionTracksInput,
) ([]MusicTrackOutput, error) {
	if input.Provider == "" || input.AccessToken == "" || input.CollectionID == "" {
		return nil, ErrInvalidInput
	}

	client, err := u.registry.Get(input.Provider)
	if err != nil {
		return nil, err
	}

	tracks, err := client.ListCollectionTracks(ctx, input.AccessToken, input.CollectionID)
	if err != nil {
		return nil, err
	}

	outputs := make([]MusicTrackOutput, 0, len(tracks))
	for _, track := range tracks {
		outputs = append(outputs, toMusicTrackOutput(track))
	}

	return outputs, nil
}
