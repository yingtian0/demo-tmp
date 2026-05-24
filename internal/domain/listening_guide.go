package domain

import (
	"strings"
	"time"
)

type ListeningGuide struct {
	id          ListeningGuideID
	encounterID EncounterID

	viewerUserID UserID

	sourcePartianalityCardID PartianalityCardID
	targetPartianalityCardID PartianalityCardID

	summary         string
	connectionPoint string
	listeningTips   []string
	firstFocusPoint string

	createdAt time.Time
}

func NewListeningGuide(
	id ListeningGuideID,
	encounterID EncounterID,
	viewerUserID UserID,
	sourcePartianalityCardID PartianalityCardID,
	targetPartianalityCardID PartianalityCardID,
	summary string,
	connectionPoint string,
	listeningTips []string,
	firstFocusPoint string,
	now time.Time,
) (*ListeningGuide, error) {
	if id.IsZero() ||
		encounterID.IsZero() ||
		viewerUserID.IsZero() ||
		sourcePartianalityCardID.IsZero() ||
		targetPartianalityCardID.IsZero() ||
		sourcePartianalityCardID == targetPartianalityCardID ||
		isBlank(summary) ||
		isBlank(connectionPoint) ||
		now.IsZero() {
		return nil, ErrInvalidListeningGuide
	}

	return &ListeningGuide{
		id:                       id,
		encounterID:              encounterID,
		viewerUserID:             viewerUserID,
		sourcePartianalityCardID: sourcePartianalityCardID,
		targetPartianalityCardID: targetPartianalityCardID,
		summary:                  strings.TrimSpace(summary),
		connectionPoint:          strings.TrimSpace(connectionPoint),
		listeningTips:            NormalizeTags(listeningTips),
		firstFocusPoint:          strings.TrimSpace(firstFocusPoint),
		createdAt:                now,
	}, nil
}

func (g *ListeningGuide) ID() ListeningGuideID {
	return g.id
}

func (g *ListeningGuide) EncounterID() EncounterID {
	return g.encounterID
}

func (g *ListeningGuide) ViewerUserID() UserID {
	return g.viewerUserID
}

func (g *ListeningGuide) SourcePartianalityCardID() PartianalityCardID {
	return g.sourcePartianalityCardID
}

func (g *ListeningGuide) TargetPartianalityCardID() PartianalityCardID {
	return g.targetPartianalityCardID
}

func (g *ListeningGuide) Summary() string {
	return g.summary
}

func (g *ListeningGuide) ConnectionPoint() string {
	return g.connectionPoint
}

func (g *ListeningGuide) ListeningTips() []string {
	return append([]string{}, g.listeningTips...)
}

func (g *ListeningGuide) FirstFocusPoint() string {
	return g.firstFocusPoint
}

func (g *ListeningGuide) CreatedAt() time.Time {
	return g.createdAt
}
