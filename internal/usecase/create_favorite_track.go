package usecase

import (
	"context"

	"backend/internal/domain"
)

type CreateFavoriteTrackInput struct {
	UserID         domain.UserID
	Title          string
	ArtistName     string
	PreviewURL     string
	ExternalURL    string
	ArtworkURL     string
	Reason         string
	ListeningScene string
}

type CreateFavoriteTrackUsecase struct {
	favoriteTrackRepository domain.FavoriteTrackRepository
	idGenerator             IDGenerator
	clock                   Clock
}

func NewCreateFavoriteTrackUsecase(
	favoriteTrackRepository domain.FavoriteTrackRepository,
	idGenerator IDGenerator,
	clock Clock,
) *CreateFavoriteTrackUsecase {
	return &CreateFavoriteTrackUsecase{
		favoriteTrackRepository: favoriteTrackRepository,
		idGenerator:             idGenerator,
		clock:                   clock,
	}
}

func (u *CreateFavoriteTrackUsecase) Execute(
	ctx context.Context,
	input CreateFavoriteTrackInput,
) (*FavoriteTrackOutput, error) {
	track, err := domain.NewTrack(
		input.Title,
		input.ArtistName,
		input.PreviewURL,
		input.ExternalURL,
		input.ArtworkURL,
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
