package domain

import "time"

type SavedPartianalityCard struct {
	id                       SavedPartianalityCardID
	userID                   UserID
	targetPartianalityCardID PartianalityCardID
	encounterID              EncounterID
	savedAt                  time.Time
}

func NewSavedPartianalityCard(
	id SavedPartianalityCardID,
	userID UserID,
	targetPartianalityCardID PartianalityCardID,
	encounterID EncounterID,
	now time.Time,
) (*SavedPartianalityCard, error) {
	if id.IsZero() ||
		userID.IsZero() ||
		targetPartianalityCardID.IsZero() ||
		encounterID.IsZero() ||
		now.IsZero() {
		return nil, ErrInvalidSavedPartianalityCard
	}

	return &SavedPartianalityCard{
		id:                       id,
		userID:                   userID,
		targetPartianalityCardID: targetPartianalityCardID,
		encounterID:              encounterID,
		savedAt:                  now,
	}, nil
}

func (s *SavedPartianalityCard) ID() SavedPartianalityCardID {
	return s.id
}

func (s *SavedPartianalityCard) UserID() UserID {
	return s.userID
}

func (s *SavedPartianalityCard) TargetPartianalityCardID() PartianalityCardID {
	return s.targetPartianalityCardID
}

func (s *SavedPartianalityCard) EncounterID() EncounterID {
	return s.encounterID
}

func (s *SavedPartianalityCard) SavedAt() time.Time {
	return s.savedAt
}
