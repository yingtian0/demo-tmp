package domain

import "strings"

type UserID string
type FavoriteTrackID string
type PartianalityCardID string
type EncounterID string
type ListeningGuideID string
type ReactionID string
type SavedPartianalityCardID string

func isBlank(s string) bool {
	return strings.TrimSpace(s) == ""
}

func (id UserID) IsZero() bool {
	return isBlank(string(id))
}

func (id FavoriteTrackID) IsZero() bool {
	return isBlank(string(id))
}

func (id PartianalityCardID) IsZero() bool {
	return isBlank(string(id))
}

func (id EncounterID) IsZero() bool {
	return isBlank(string(id))
}

func (id ListeningGuideID) IsZero() bool {
	return isBlank(string(id))
}

func (id ReactionID) IsZero() bool {
	return isBlank(string(id))
}

func (id SavedPartianalityCardID) IsZero() bool {
	return isBlank(string(id))
}
