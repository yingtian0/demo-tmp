package usecase

import (
	"context"

	"backend/internal/domain"
)

type CreateEncounterInput struct {
	UserAID domain.UserID
	UserBID domain.UserID
	Source  domain.EncounterSource
}

type CreateEncounterUsecase struct {
	encounterRepository        domain.EncounterRepository
	partianalityCardRepository domain.PartianalityCardRepository
	idGenerator                IDGenerator
	clock                      Clock
}

func NewCreateEncounterUsecase(
	encounterRepository domain.EncounterRepository,
	partianalityCardRepository domain.PartianalityCardRepository,
	idGenerator IDGenerator,
	clock Clock,
) *CreateEncounterUsecase {
	return &CreateEncounterUsecase{
		encounterRepository:        encounterRepository,
		partianalityCardRepository: partianalityCardRepository,
		idGenerator:                idGenerator,
		clock:                      clock,
	}
}

func (u *CreateEncounterUsecase) Execute(
	ctx context.Context,
	input CreateEncounterInput,
) (*EncounterDetailOutput, error) {
	if input.UserAID.IsZero() || input.UserBID.IsZero() || input.UserAID == input.UserBID || !input.Source.IsValid() {
		return nil, ErrInvalidInput
	}

	exists, err := u.encounterRepository.ExistsInTimeWindow(ctx, input.UserAID, input.UserBID, input.Source)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrEncounterDuplicated
	}

	userACard, err := u.partianalityCardRepository.FindLatestByUserID(ctx, input.UserAID)
	if err != nil {
		return nil, err
	}
	if userACard == nil {
		return nil, ErrNotFound
	}

	userBCard, err := u.partianalityCardRepository.FindLatestByUserID(ctx, input.UserBID)
	if err != nil {
		return nil, err
	}
	if userBCard == nil {
		return nil, ErrNotFound
	}

	encounter, err := domain.NewEncounter(
		domain.EncounterID(u.idGenerator.NewID()),
		input.UserAID,
		input.UserBID,
		userACard.ID(),
		userBCard.ID(),
		u.clock.Now(),
		input.Source,
	)
	if err != nil {
		return nil, ErrInvalidInput
	}

	if err := u.encounterRepository.Save(ctx, encounter); err != nil {
		return nil, err
	}

	output := EncounterDetailOutput{
		EncounterID:            string(encounter.ID()),
		OccurredAt:             encounter.OccurredAt(),
		Source:                 string(encounter.Source()),
		ViewerPartianalityCard: toPartianalityCardOutput(userACard),
		TargetPartianalityCard: toPartianalityCardOutput(userBCard),
		Track:                  toTrackOutput(userBCard.Track()),
	}
	return &output, nil
}
