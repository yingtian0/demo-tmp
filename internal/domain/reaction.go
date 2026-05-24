package domain

import "time"

type Reaction struct {
	id                       ReactionID
	userID                   UserID
	targetPartianalityCardID PartianalityCardID
	reactionType             ReactionType
	createdAt                time.Time
}

func NewReaction(
	id ReactionID,
	userID UserID,
	targetPartianalityCardID PartianalityCardID,
	reactionType ReactionType,
	now time.Time,
) (*Reaction, error) {
	if id.IsZero() ||
		userID.IsZero() ||
		targetPartianalityCardID.IsZero() ||
		!reactionType.IsValid() ||
		now.IsZero() {
		return nil, ErrInvalidReaction
	}

	return &Reaction{
		id:                       id,
		userID:                   userID,
		targetPartianalityCardID: targetPartianalityCardID,
		reactionType:             reactionType,
		createdAt:                now,
	}, nil
}

func (r *Reaction) ID() ReactionID {
	return r.id
}

func (r *Reaction) UserID() UserID {
	return r.userID
}

func (r *Reaction) TargetPartianalityCardID() PartianalityCardID {
	return r.targetPartianalityCardID
}

func (r *Reaction) Type() ReactionType {
	return r.reactionType
}

func (r *Reaction) CreatedAt() time.Time {
	return r.createdAt
}
