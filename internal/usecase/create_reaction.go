package usecase

import (
	"context"

	"backend/internal/domain"
)

type CreateReactionInput struct {
	UserID                   domain.UserID
	TargetPartianalityCardID domain.PartianalityCardID
	ReactionType             domain.ReactionType
}

type CreateReactionUsecase struct {
	partianalityCardRepository domain.PartianalityCardRepository
	reactionRepository         domain.ReactionRepository
	idGenerator                IDGenerator
	clock                      Clock
}

func NewCreateReactionUsecase(
	partianalityCardRepository domain.PartianalityCardRepository,
	reactionRepository domain.ReactionRepository,
	idGenerator IDGenerator,
	clock Clock,
) *CreateReactionUsecase {
	return &CreateReactionUsecase{
		partianalityCardRepository: partianalityCardRepository,
		reactionRepository:         reactionRepository,
		idGenerator:                idGenerator,
		clock:                      clock,
	}
}

func (u *CreateReactionUsecase) Execute(
	ctx context.Context,
	input CreateReactionInput,
) (*ReactionOutput, error) {
	if !input.ReactionType.IsValid() || input.UserID.IsZero() || input.TargetPartianalityCardID.IsZero() {
		return nil, ErrInvalidInput
	}

	card, err := u.partianalityCardRepository.FindByID(ctx, input.TargetPartianalityCardID)
	if err != nil {
		return nil, err
	}
	if card == nil {
		return nil, ErrNotFound
	}

	reaction, err := domain.NewReaction(
		domain.ReactionID(u.idGenerator.NewID()),
		input.UserID,
		input.TargetPartianalityCardID,
		input.ReactionType,
		u.clock.Now(),
	)
	if err != nil {
		return nil, ErrInvalidInput
	}

	if err := u.reactionRepository.Save(ctx, reaction); err != nil {
		return nil, err
	}

	output := toReactionOutput(reaction)
	return &output, nil
}
