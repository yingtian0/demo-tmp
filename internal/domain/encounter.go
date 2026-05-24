package domain

import "time"

type Encounter struct {
	id EncounterID

	userAID UserID
	userBID UserID

	userAPartianalityCardID PartianalityCardID
	userBPartianalityCardID PartianalityCardID

	occurredAt time.Time
	source     EncounterSource
}

func NewEncounter(
	id EncounterID,
	userAID UserID,
	userBID UserID,
	userAPartianalityCardID PartianalityCardID,
	userBPartianalityCardID PartianalityCardID,
	occurredAt time.Time,
	source EncounterSource,
) (*Encounter, error) {
	if id.IsZero() ||
		userAID.IsZero() ||
		userBID.IsZero() ||
		userAID == userBID ||
		userAPartianalityCardID.IsZero() ||
		userBPartianalityCardID.IsZero() ||
		userAPartianalityCardID == userBPartianalityCardID ||
		occurredAt.IsZero() ||
		!source.IsValid() {
		return nil, ErrInvalidEncounter
	}

	return &Encounter{
		id:                      id,
		userAID:                 userAID,
		userBID:                 userBID,
		userAPartianalityCardID: userAPartianalityCardID,
		userBPartianalityCardID: userBPartianalityCardID,
		occurredAt:              occurredAt,
		source:                  source,
	}, nil
}

func (e *Encounter) ID() EncounterID {
	return e.id
}

func (e *Encounter) UserAID() UserID {
	return e.userAID
}

func (e *Encounter) UserBID() UserID {
	return e.userBID
}

func (e *Encounter) UserAPartianalityCardID() PartianalityCardID {
	return e.userAPartianalityCardID
}

func (e *Encounter) UserBPartianalityCardID() PartianalityCardID {
	return e.userBPartianalityCardID
}

func (e *Encounter) OccurredAt() time.Time {
	return e.occurredAt
}

func (e *Encounter) Source() EncounterSource {
	return e.source
}

func (e *Encounter) HasParticipant(userID UserID) bool {
	return e.userAID == userID || e.userBID == userID
}

func (e *Encounter) OtherUserID(viewerID UserID) (UserID, error) {
	switch viewerID {
	case e.userAID:
		return e.userBID, nil
	case e.userBID:
		return e.userAID, nil
	default:
		return "", ErrNotEncounterParticipant
	}
}

func (e *Encounter) ViewerPartianalityCardID(viewerID UserID) (PartianalityCardID, error) {
	switch viewerID {
	case e.userAID:
		return e.userAPartianalityCardID, nil
	case e.userBID:
		return e.userBPartianalityCardID, nil
	default:
		return "", ErrNotEncounterParticipant
	}
}

func (e *Encounter) TargetPartianalityCardID(viewerID UserID) (PartianalityCardID, error) {
	switch viewerID {
	case e.userAID:
		return e.userBPartianalityCardID, nil
	case e.userBID:
		return e.userAPartianalityCardID, nil
	default:
		return "", ErrNotEncounterParticipant
	}
}
