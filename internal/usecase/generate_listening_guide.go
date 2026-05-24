package usecase

import (
	"context"

	"demo-tmp/internal/domain"
)

type GenerateListeningGuideInput struct {
	ViewerUserID domain.UserID
	EncounterID  domain.EncounterID
}

type GenerateListeningGuideUsecase struct {
	encounterRepository        domain.EncounterRepository
	partianalityCardRepository domain.PartianalityCardRepository
	listeningGuideRepository   domain.ListeningGuideRepository
	listeningGuideGenerator    domain.ListeningGuideGenerator
	idGenerator                IDGenerator
	clock                      Clock
}

func NewGenerateListeningGuideUsecase(
	encounterRepository domain.EncounterRepository,
	partianalityCardRepository domain.PartianalityCardRepository,
	listeningGuideRepository domain.ListeningGuideRepository,
	listeningGuideGenerator domain.ListeningGuideGenerator,
	idGenerator IDGenerator,
	clock Clock,
) *GenerateListeningGuideUsecase {
	return &GenerateListeningGuideUsecase{
		encounterRepository:        encounterRepository,
		partianalityCardRepository: partianalityCardRepository,
		listeningGuideRepository:   listeningGuideRepository,
		listeningGuideGenerator:    listeningGuideGenerator,
		idGenerator:                idGenerator,
		clock:                      clock,
	}
}

func (u *GenerateListeningGuideUsecase) Execute(
	ctx context.Context,
	input GenerateListeningGuideInput,
) (*ListeningGuideOutput, error) {
	if input.ViewerUserID.IsZero() || input.EncounterID.IsZero() {
		return nil, ErrInvalidInput
	}

	encounter, err := u.encounterRepository.FindByID(ctx, input.EncounterID)
	if err != nil {
		return nil, err
	}
	if encounter == nil {
		return nil, ErrNotFound
	}
	if !encounter.HasParticipant(input.ViewerUserID) {
		return nil, domain.ErrNotEncounterParticipant
	}

	viewerCardID, err := encounter.ViewerPartianalityCardID(input.ViewerUserID)
	if err != nil {
		return nil, err
	}
	targetCardID, err := encounter.TargetPartianalityCardID(input.ViewerUserID)
	if err != nil {
		return nil, err
	}

	viewerCard, err := u.partianalityCardRepository.FindByID(ctx, viewerCardID)
	if err != nil {
		return nil, err
	}
	if viewerCard == nil {
		return nil, ErrNotFound
	}

	targetCard, err := u.partianalityCardRepository.FindByID(ctx, targetCardID)
	if err != nil {
		return nil, err
	}
	if targetCard == nil {
		return nil, ErrNotFound
	}

	draft, err := u.listeningGuideGenerator.GenerateListeningGuide(ctx, *viewerCard, *targetCard)
	if err != nil {
		return nil, err
	}
	if draft == nil {
		return nil, ErrInvalidInput
	}

	guide, err := domain.NewListeningGuide(
		domain.ListeningGuideID(u.idGenerator.NewID()),
		encounter.ID(),
		input.ViewerUserID,
		viewerCard.ID(),
		targetCard.ID(),
		draft.Summary,
		draft.ConnectionPoint,
		draft.ListeningTips,
		draft.FirstFocusPoint,
		u.clock.Now(),
	)
	if err != nil {
		return nil, ErrInvalidInput
	}

	if err := u.listeningGuideRepository.Save(ctx, guide); err != nil {
		return nil, err
	}

	output := toListeningGuideOutput(guide)
	return &output, nil
}
