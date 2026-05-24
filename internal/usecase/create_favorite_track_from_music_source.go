package usecase

import (
	"context"

	"backend/internal/domain"
)

type CreateFavoriteTrackFromMusicSourceInput struct {
	UserID         domain.UserID
	Provider       string
	AccessToken    string
	TrackID        string
	Reason         string
	ListeningScene string
}

type CreateFavoriteTrackFromMusicSourceUsecase struct {
	registry                MusicSourceRegistry
	favoriteTrackRepository domain.FavoriteTrackRepository
	idGenerator             IDGenerator
	clock                   Clock
}

func NewCreateFavoriteTrackFromMusicSourceUsecase(
	registry MusicSourceRegistry,
	favoriteTrackRepository domain.FavoriteTrackRepository,
	idGenerator IDGenerator,
	clock Clock,
) *CreateFavoriteTrackFromMusicSourceUsecase {
	return &CreateFavoriteTrackFromMusicSourceUsecase{
		registry:                registry,
		favoriteTrackRepository: favoriteTrackRepository,
		idGenerator:             idGenerator,
		clock:                   clock,
	}
}

func (u *CreateFavoriteTrackFromMusicSourceUsecase) Execute(
	ctx context.Context,
	input CreateFavoriteTrackFromMusicSourceInput,
) (*FavoriteTrackOutput, error) {
	if input.UserID.IsZero() || input.Provider == "" || input.AccessToken == "" || input.TrackID == "" {
		return nil, ErrInvalidInput
	}

	client, err := u.registry.Get(input.Provider)
	if err != nil {
		return nil, err
	}

	sourceTrack, err := client.GetTrack(ctx, input.AccessToken, input.TrackID)
	if err != nil {
		return nil, err
	}
	if sourceTrack == nil {
		return nil, ErrNotFound
	}

	track, err := domain.NewTrack(
		sourceTrack.Title,
		sourceTrack.ArtistName,
		sourceTrack.PreviewURL,
		sourceTrack.ExternalURL,
		sourceTrack.ArtworkURL,
	)
	if err != nil {
		return nil, ErrInvalidInput
	}

	favoriteTrack, err := domain.NewFavoriteTrack(
		domain.FavoriteTrackID(u.idGenerator.NewID()),
		input.UserID,
		track,
		input.Reason,
		input.ListeningScene,
		u.clock.Now(),
	)
	if err != nil {
		return nil, ErrInvalidInput
	}

	if err := u.favoriteTrackRepository.Save(ctx, favoriteTrack); err != nil {
		return nil, err
	}

	output := toFavoriteTrackOutput(favoriteTrack)
	return &output, nil
}
