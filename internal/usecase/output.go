package usecase

import (
	"time"

	"demo-tmp/internal/domain"
)

type TrackOutput struct {
	Title       string
	ArtistName  string
	PreviewURL  string
	ExternalURL string
	ArtworkURL  string
}

func toMusicTrackOutput(track MusicTrack) MusicTrackOutput {
	return MusicTrackOutput{
		Provider:    track.Provider,
		ID:          track.ID,
		Title:       track.Title,
		ArtistName:  track.ArtistName,
		PreviewURL:  track.PreviewURL,
		ExternalURL: track.ExternalURL,
		ArtworkURL:  track.ArtworkURL,
	}
}

type FavoriteTrackOutput struct {
	ID             string
	UserID         string
	Track          TrackOutput
	Reason         string
	ListeningScene string
	CreatedAt      time.Time
}

type PartianalityCardOutput struct {
	ID               string
	UserID           string
	FavoriteTrackID  string
	Track            TrackOutput
	Title            string
	FavoritePoint    string
	ListeningPoint   string
	RecommendedScene string
	Tags             []string
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

type EncounterListItemOutput struct {
	EncounterID            string
	OccurredAt             time.Time
	Source                 string
	TargetUserID           string
	TargetPartianalityCard PartianalityCardOutput
}

type ListeningGuideOutput struct {
	ID                       string
	EncounterID              string
	ViewerUserID             string
	SourcePartianalityCardID string
	TargetPartianalityCardID string
	Summary                  string
	ConnectionPoint          string
	ListeningTips            []string
	FirstFocusPoint          string
	CreatedAt                time.Time
}

type EncounterDetailOutput struct {
	EncounterID            string
	OccurredAt             time.Time
	Source                 string
	ViewerPartianalityCard PartianalityCardOutput
	TargetPartianalityCard PartianalityCardOutput
	ListeningGuide         *ListeningGuideOutput
	Track                  TrackOutput
}

type ReactionOutput struct {
	ID                       string
	UserID                   string
	TargetPartianalityCardID string
	ReactionType             string
	CreatedAt                time.Time
}

type SavedPartianalityCardOutput struct {
	ID                       string
	UserID                   string
	TargetPartianalityCardID string
	EncounterID              string
	SavedAt                  time.Time
}

func toTrackOutput(track domain.Track) TrackOutput {
	return TrackOutput{
		Title:       track.Title(),
		ArtistName:  track.ArtistName(),
		PreviewURL:  track.PreviewURL(),
		ExternalURL: track.ExternalURL(),
		ArtworkURL:  track.ArtworkURL(),
	}
}

func toFavoriteTrackOutput(favoriteTrack *domain.FavoriteTrack) FavoriteTrackOutput {
	return FavoriteTrackOutput{
		ID:             string(favoriteTrack.ID()),
		UserID:         string(favoriteTrack.UserID()),
		Track:          toTrackOutput(favoriteTrack.Track()),
		Reason:         favoriteTrack.Reason(),
		ListeningScene: favoriteTrack.ListeningScene(),
		CreatedAt:      favoriteTrack.CreatedAt(),
	}
}

func toPartianalityCardOutput(card *domain.PartianalityCard) PartianalityCardOutput {
	return PartianalityCardOutput{
		ID:               string(card.ID()),
		UserID:           string(card.UserID()),
		FavoriteTrackID:  string(card.FavoriteTrackID()),
		Track:            toTrackOutput(card.Track()),
		Title:            card.Title(),
		FavoritePoint:    card.FavoritePoint(),
		ListeningPoint:   card.ListeningPoint(),
		RecommendedScene: card.RecommendedScene(),
		Tags:             card.Tags(),
		CreatedAt:        card.CreatedAt(),
		UpdatedAt:        card.UpdatedAt(),
	}
}

func toListeningGuideOutput(guide *domain.ListeningGuide) ListeningGuideOutput {
	return ListeningGuideOutput{
		ID:                       string(guide.ID()),
		EncounterID:              string(guide.EncounterID()),
		ViewerUserID:             string(guide.ViewerUserID()),
		SourcePartianalityCardID: string(guide.SourcePartianalityCardID()),
		TargetPartianalityCardID: string(guide.TargetPartianalityCardID()),
		Summary:                  guide.Summary(),
		ConnectionPoint:          guide.ConnectionPoint(),
		ListeningTips:            guide.ListeningTips(),
		FirstFocusPoint:          guide.FirstFocusPoint(),
		CreatedAt:                guide.CreatedAt(),
	}
}

func toReactionOutput(reaction *domain.Reaction) ReactionOutput {
	return ReactionOutput{
		ID:                       string(reaction.ID()),
		UserID:                   string(reaction.UserID()),
		TargetPartianalityCardID: string(reaction.TargetPartianalityCardID()),
		ReactionType:             string(reaction.Type()),
		CreatedAt:                reaction.CreatedAt(),
	}
}

func toSavedPartianalityCardOutput(saved *domain.SavedPartianalityCard) SavedPartianalityCardOutput {
	return SavedPartianalityCardOutput{
		ID:                       string(saved.ID()),
		UserID:                   string(saved.UserID()),
		TargetPartianalityCardID: string(saved.TargetPartianalityCardID()),
		EncounterID:              string(saved.EncounterID()),
		SavedAt:                  saved.SavedAt(),
	}
}
