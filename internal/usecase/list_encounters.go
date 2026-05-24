package usecase

import (
	"context"

	"demo-tmp/internal/domain"
)

type ListEncountersInput struct {
	ViewerUserID domain.UserID
	Limit        int
}

type ListEncountersUsecase struct {
	encounterRepository        domain.EncounterRepository
	partianalityCardRepository domain.PartianalityCardRepository
}

func NewListEncountersUsecase(
	encounterRepository domain.EncounterRepository,
	partianalityCardRepository domain.PartianalityCardRepository,
) *ListEncountersUsecase {
	return &ListEncountersUsecase{
		encounterRepository:        encounterRepository,
		partianalityCardRepository: partianalityCardRepository,
	}
}

func (u *ListEncountersUsecase) Execute(
	ctx context.Context,
	input ListEncountersInput,
) ([]EncounterListItemOutput, error) {
	if input.ViewerUserID.IsZero() {
		return nil, ErrInvalidInput
	}

	encounters, err := u.encounterRepository.FindByUserID(ctx, input.ViewerUserID, input.Limit)
	if err != nil {
		return nil, err
	}

	outputs := make([]EncounterListItemOutput, 0, len(encounters))
	for _, encounter := range encounters {
		if encounter == nil {
			continue
		}

		targetCardID, err := encounter.TargetPartianalityCardID(input.ViewerUserID)
		if err != nil {
			return nil, err
		}
		targetUserID, err := encounter.OtherUserID(input.ViewerUserID)
		if err != nil {
			return nil, err
		}

		targetCard, err := u.partianalityCardRepository.FindByID(ctx, targetCardID)
		if err != nil {
			return nil, err
		}
		if targetCard == nil {
			return nil, ErrNotFound
		}

		outputs = append(outputs, EncounterListItemOutput{
			EncounterID:            string(encounter.ID()),
			OccurredAt:             encounter.OccurredAt(),
			Source:                 string(encounter.Source()),
			TargetUserID:           string(targetUserID),
			TargetPartianalityCard: toPartianalityCardOutput(targetCard),
		})
	}

	return outputs, nil
}
