package usecase

import (
	"context"

	"demo-tmp/internal/domain"
)

type SavePartianalityCardInput struct {
	UserID                   domain.UserID
	TargetPartianalityCardID domain.PartianalityCardID
	EncounterID              domain.EncounterID
}

type SavePartianalityCardUsecase struct {
	partianalityCardRepository      domain.PartianalityCardRepository
	savedPartianalityCardRepository domain.SavedPartianalityCardRepository
	idGenerator                     IDGenerator
	clock                           Clock
}

func NewSavePartianalityCardUsecase(
	partianalityCardRepository domain.PartianalityCardRepository,
	savedPartianalityCardRepository domain.SavedPartianalityCardRepository,
	idGenerator IDGenerator,
	clock Clock,
) *SavePartianalityCardUsecase {
	return &SavePartianalityCardUsecase{
		partianalityCardRepository:      partianalityCardRepository,
		savedPartianalityCardRepository: savedPartianalityCardRepository,
		idGenerator:                     idGenerator,
		clock:                           clock,
	}
}

func (u *SavePartianalityCardUsecase) Execute(
	ctx context.Context,
	input SavePartianalityCardInput,
) (*SavedPartianalityCardOutput, error) {
	if input.UserID.IsZero() || input.TargetPartianalityCardID.IsZero() || input.EncounterID.IsZero() {
		return nil, ErrInvalidInput
	}

	card, err := u.partianalityCardRepository.FindByID(ctx, input.TargetPartianalityCardID)
	if err != nil {
		return nil, err
	}
	if card == nil {
		return nil, ErrNotFound
	}

	exists, err := u.savedPartianalityCardRepository.Exists(ctx, input.UserID, input.TargetPartianalityCardID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrAlreadyExists
	}

	saved, err := domain.NewSavedPartianalityCard(
		domain.SavedPartianalityCardID(u.idGenerator.NewID()),
		input.UserID,
		input.TargetPartianalityCardID,
		input.EncounterID,
		u.clock.Now(),
	)
	if err != nil {
		return nil, ErrInvalidInput
	}

	if err := u.savedPartianalityCardRepository.Save(ctx, saved); err != nil {
		return nil, err
	}

	output := toSavedPartianalityCardOutput(saved)
	return &output, nil
}
