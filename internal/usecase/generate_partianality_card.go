package usecase

import (
	"context"

	"backend/internal/domain"
)

type GeneratePartianalityCardInput struct {
	UserID          domain.UserID
	FavoriteTrackID domain.FavoriteTrackID
}

type GeneratePartianalityCardUsecase struct {
	favoriteTrackRepository    domain.FavoriteTrackRepository
	partianalityCardRepository domain.PartianalityCardRepository
	partianalityCardGenerator  domain.PartianalityCardGenerator
	idGenerator                IDGenerator
	clock                      Clock
}

func NewGeneratePartianalityCardUsecase(
	favoriteTrackRepository domain.FavoriteTrackRepository,
	partianalityCardRepository domain.PartianalityCardRepository,
	partianalityCardGenerator domain.PartianalityCardGenerator,
	idGenerator IDGenerator,
	clock Clock,
) *GeneratePartianalityCardUsecase {
	return &GeneratePartianalityCardUsecase{
		favoriteTrackRepository:    favoriteTrackRepository,
		partianalityCardRepository: partianalityCardRepository,
		partianalityCardGenerator:  partianalityCardGenerator,
		idGenerator:                idGenerator,
		clock:                      clock,
	}
}

func (u *GeneratePartianalityCardUsecase) Execute(
	ctx context.Context,
	input GeneratePartianalityCardInput,
) (*PartianalityCardOutput, error) {
	favoriteTrack, err := u.favoriteTrackRepository.FindByID(ctx, input.FavoriteTrackID)
	if err != nil {
		return nil, err
	}
	if favoriteTrack == nil {
		return nil, ErrNotFound
	}
	if favoriteTrack.UserID() != input.UserID {
		return nil, ErrPermissionDenied
	}

	draft, err := u.partianalityCardGenerator.GeneratePartianalityCard(ctx, *favoriteTrack)
	if err != nil {
		return nil, err
	}
	if draft == nil {
		return nil, ErrInvalidInput
	}

	card, err := domain.NewPartianalityCard(
		domain.PartianalityCardID(u.idGenerator.NewID()),
		input.UserID,
		favoriteTrack.ID(),
		favoriteTrack.Track(),
		draft.Title,
		draft.FavoritePoint,
		draft.ListeningPoint,
		draft.RecommendedScene,
		draft.Tags,
		u.clock.Now(),
	)
	if err != nil {
		return nil, ErrInvalidInput
	}

	if err := u.partianalityCardRepository.Save(ctx, card); err != nil {
		return nil, err
	}

	output := toPartianalityCardOutput(card)
	return &output, nil
}
