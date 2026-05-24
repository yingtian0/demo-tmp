package domain

import "context"

type UserRepository interface {
	FindByID(ctx context.Context, id UserID) (*User, error)
	Save(ctx context.Context, user *User) error
}

type FavoriteTrackRepository interface {
	FindByID(ctx context.Context, id FavoriteTrackID) (*FavoriteTrack, error)
	FindLatestByUserID(ctx context.Context, userID UserID) (*FavoriteTrack, error)
	Save(ctx context.Context, favoriteTrack *FavoriteTrack) error
}

type PartianalityCardRepository interface {
	FindByID(ctx context.Context, id PartianalityCardID) (*PartianalityCard, error)
	FindLatestByUserID(ctx context.Context, userID UserID) (*PartianalityCard, error)
	FindByUserID(ctx context.Context, userID UserID, limit int) ([]*PartianalityCard, error)
	Save(ctx context.Context, card *PartianalityCard) error
}

type EncounterRepository interface {
	FindByID(ctx context.Context, id EncounterID) (*Encounter, error)
	FindByUserID(ctx context.Context, userID UserID, limit int) ([]*Encounter, error)
	Save(ctx context.Context, encounter *Encounter) error

	ExistsInTimeWindow(
		ctx context.Context,
		userAID UserID,
		userBID UserID,
		source EncounterSource,
	) (bool, error)
}

type ListeningGuideRepository interface {
	FindByID(ctx context.Context, id ListeningGuideID) (*ListeningGuide, error)

	FindByEncounterAndViewer(
		ctx context.Context,
		encounterID EncounterID,
		viewerUserID UserID,
	) (*ListeningGuide, error)

	Save(ctx context.Context, guide *ListeningGuide) error
}

type ReactionRepository interface {
	FindByUserAndCard(
		ctx context.Context,
		userID UserID,
		targetPartianalityCardID PartianalityCardID,
	) ([]*Reaction, error)

	Save(ctx context.Context, reaction *Reaction) error
}

type SavedPartianalityCardRepository interface {
	FindByUserID(ctx context.Context, userID UserID, limit int) ([]*SavedPartianalityCard, error)

	Exists(
		ctx context.Context,
		userID UserID,
		targetPartianalityCardID PartianalityCardID,
	) (bool, error)

	Save(ctx context.Context, saved *SavedPartianalityCard) error
}
