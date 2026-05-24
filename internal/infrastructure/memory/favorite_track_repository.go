package memory

import (
	"context"
	"sync"

	"demo-tmp/internal/domain"
	"demo-tmp/internal/usecase"
)

type FavoriteTrackRepository struct {
	mu     sync.RWMutex
	byID   map[domain.FavoriteTrackID]*domain.FavoriteTrack
	latest map[domain.UserID]domain.FavoriteTrackID
}

func NewFavoriteTrackRepository() *FavoriteTrackRepository {
	return &FavoriteTrackRepository{
		byID:   make(map[domain.FavoriteTrackID]*domain.FavoriteTrack),
		latest: make(map[domain.UserID]domain.FavoriteTrackID),
	}
}

func (r *FavoriteTrackRepository) FindByID(_ context.Context, id domain.FavoriteTrackID) (*domain.FavoriteTrack, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	track, ok := r.byID[id]
	if !ok {
		return nil, usecase.ErrNotFound
	}

	return track, nil
}

func (r *FavoriteTrackRepository) FindLatestByUserID(_ context.Context, userID domain.UserID) (*domain.FavoriteTrack, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	id, ok := r.latest[userID]
	if !ok {
		return nil, usecase.ErrNotFound
	}

	track, ok := r.byID[id]
	if !ok {
		return nil, usecase.ErrNotFound
	}

	return track, nil
}

func (r *FavoriteTrackRepository) Save(_ context.Context, favoriteTrack *domain.FavoriteTrack) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.byID[favoriteTrack.ID()] = favoriteTrack
	r.latest[favoriteTrack.UserID()] = favoriteTrack.ID()
	return nil
}
