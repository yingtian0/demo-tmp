package domain

import (
	"strings"
	"time"
)

type PartianalityCard struct {
	id              PartianalityCardID
	userID          UserID
	favoriteTrackID FavoriteTrackID

	// すれ違い時点の表示を固定するため、Trackをスナップショットとして持つ
	track Track

	title            string
	favoritePoint    string
	listeningPoint   string
	recommendedScene string
	tags             []string

	createdAt time.Time
	updatedAt time.Time
}

func NewPartianalityCard(
	id PartianalityCardID,
	userID UserID,
	favoriteTrackID FavoriteTrackID,
	track Track,
	title string,
	favoritePoint string,
	listeningPoint string,
	recommendedScene string,
	tags []string,
	now time.Time,
) (*PartianalityCard, error) {
	if id.IsZero() ||
		userID.IsZero() ||
		favoriteTrackID.IsZero() ||
		isBlank(track.Title()) ||
		isBlank(track.ArtistName()) ||
		isBlank(title) ||
		isBlank(favoritePoint) ||
		isBlank(listeningPoint) ||
		now.IsZero() {
		return nil, ErrInvalidPartianalityCard
	}

	return &PartianalityCard{
		id:               id,
		userID:           userID,
		favoriteTrackID:  favoriteTrackID,
		track:            track,
		title:            strings.TrimSpace(title),
		favoritePoint:    strings.TrimSpace(favoritePoint),
		listeningPoint:   strings.TrimSpace(listeningPoint),
		recommendedScene: strings.TrimSpace(recommendedScene),
		tags:             NormalizeTags(tags),
		createdAt:        now,
		updatedAt:        now,
	}, nil
}

func (c *PartianalityCard) ID() PartianalityCardID {
	return c.id
}

func (c *PartianalityCard) UserID() UserID {
	return c.userID
}

func (c *PartianalityCard) FavoriteTrackID() FavoriteTrackID {
	return c.favoriteTrackID
}

func (c *PartianalityCard) Track() Track {
	return c.track
}

func (c *PartianalityCard) Title() string {
	return c.title
}

func (c *PartianalityCard) FavoritePoint() string {
	return c.favoritePoint
}

func (c *PartianalityCard) ListeningPoint() string {
	return c.listeningPoint
}

func (c *PartianalityCard) RecommendedScene() string {
	return c.recommendedScene
}

func (c *PartianalityCard) Tags() []string {
	return append([]string{}, c.tags...)
}

func (c *PartianalityCard) CreatedAt() time.Time {
	return c.createdAt
}

func (c *PartianalityCard) UpdatedAt() time.Time {
	return c.updatedAt
}

func (c *PartianalityCard) UpdateContent(
	title string,
	favoritePoint string,
	listeningPoint string,
	recommendedScene string,
	tags []string,
	now time.Time,
) error {
	if isBlank(title) ||
		isBlank(favoritePoint) ||
		isBlank(listeningPoint) ||
		now.IsZero() {
		return ErrInvalidPartianalityCard
	}

	c.title = strings.TrimSpace(title)
	c.favoritePoint = strings.TrimSpace(favoritePoint)
	c.listeningPoint = strings.TrimSpace(listeningPoint)
	c.recommendedScene = strings.TrimSpace(recommendedScene)
	c.tags = NormalizeTags(tags)
	c.updatedAt = now

	return nil
}

func (c *PartianalityCard) IsOwnedBy(userID UserID) bool {
	return c.userID == userID
}
