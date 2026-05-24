package domain

import "context"

type PartianalityCardDraft struct {
	Title            string
	FavoritePoint    string
	ListeningPoint   string
	RecommendedScene string
	Tags             []string
}

type PartianalityCardGenerator interface {
	GeneratePartianalityCard(
		ctx context.Context,
		favoriteTrack FavoriteTrack,
	) (*PartianalityCardDraft, error)
}

type ListeningGuideDraft struct {
	Summary         string
	ConnectionPoint string
	ListeningTips   []string
	FirstFocusPoint string
}

type ListeningGuideGenerator interface {
	GenerateListeningGuide(
		ctx context.Context,
		viewerCard PartianalityCard,
		targetCard PartianalityCard,
	) (*ListeningGuideDraft, error)
}
