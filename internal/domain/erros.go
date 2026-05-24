package domain

import "errors"

var (
	ErrInvalidID                    = errors.New("invalid id")
	ErrInvalidUser                  = errors.New("invalid user")
	ErrInvalidTrack                 = errors.New("invalid track")
	ErrInvalidFavoriteTrack         = errors.New("invalid favorite track")
	ErrInvalidPartianalityCard      = errors.New("invalid partianality card")
	ErrInvalidEncounter             = errors.New("invalid encounter")
	ErrInvalidListeningGuide        = errors.New("invalid listening guide")
	ErrInvalidReaction              = errors.New("invalid reaction")
	ErrInvalidSavedPartianalityCard = errors.New("invalid saved partianality card")

	ErrPermissionDenied        = errors.New("permission denied")
	ErrNotEncounterParticipant = errors.New("user is not a participant of this encounter")
)
