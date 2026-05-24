package domain

import (
	"strings"
	"time"
)

type FavoriteTrack struct {
	id             FavoriteTrackID
	userID         UserID
	track          Track
	reason         string
	listeningScene string
	createdAt      time.Time
}

func NewFavoriteTrack(
	id FavoriteTrackID,
	userID UserID,
	track Track,
	reason string,
	listeningScene string,
	now time.Time,
) (*FavoriteTrack, error) {
	if id.IsZero() ||
		userID.IsZero() ||
		isBlank(track.Title()) ||
		isBlank(track.ArtistName()) ||
		isBlank(reason) ||
		now.IsZero() {
		return nil, ErrInvalidFavoriteTrack
	}

	return &FavoriteTrack{
		id:             id,
		userID:         userID,
		track:          track,
		reason:         strings.TrimSpace(reason),
		listeningScene: strings.TrimSpace(listeningScene),
		createdAt:      now,
	}, nil
}

func (f *FavoriteTrack) ID() FavoriteTrackID {
	return f.id
}

func (f *FavoriteTrack) UserID() UserID {
	return f.userID
}

func (f *FavoriteTrack) Track() Track {
	return f.track
}

func (f *FavoriteTrack) Reason() string {
	return f.reason
}

func (f *FavoriteTrack) ListeningScene() string {
	return f.listeningScene
}

func (f *FavoriteTrack) CreatedAt() time.Time {
	return f.createdAt
}
