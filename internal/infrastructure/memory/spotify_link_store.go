package memory

import (
	"context"
	"sync"
	"time"

	"demo-tmp/internal/usecase"
)

type SpotifyAccountLink struct {
	AppUserID          string
	SpotifyUserID      string
	SpotifyDisplayName string
	AccessToken        string
	RefreshToken       string
	Scope              string
	ExpiresAt          time.Time
}

type SpotifyOAuthProvider interface {
	RefreshToken(ctx context.Context, refreshToken string) (*SpotifyOAuthToken, error)
}

type SpotifyOAuthToken struct {
	AccessToken  string
	RefreshToken string
	Scope        string
	ExpiresAt    time.Time
}

type SpotifyLinkStore struct {
	mu    sync.RWMutex
	links map[string]*SpotifyAccountLink
}

func NewSpotifyLinkStore() *SpotifyLinkStore {
	return &SpotifyLinkStore{
		links: make(map[string]*SpotifyAccountLink),
	}
}

func (s *SpotifyLinkStore) Save(link *SpotifyAccountLink) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.links[link.AppUserID] = link
}

func (s *SpotifyLinkStore) Get(appUserID string) (*SpotifyAccountLink, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	link, ok := s.links[appUserID]
	if !ok {
		return nil, usecase.ErrNotFound
	}

	cloned := *link
	return &cloned, nil
}

func (s *SpotifyLinkStore) GetValidAccessToken(
	ctx context.Context,
	appUserID string,
	provider SpotifyOAuthProvider,
) (string, error) {
	s.mu.RLock()
	link, ok := s.links[appUserID]
	if !ok {
		s.mu.RUnlock()
		return "", usecase.ErrNotFound
	}

	if time.Until(link.ExpiresAt) > 30*time.Second {
		token := link.AccessToken
		s.mu.RUnlock()
		return token, nil
	}
	refreshToken := link.RefreshToken
	s.mu.RUnlock()

	refreshed, err := provider.RefreshToken(ctx, refreshToken)
	if err != nil {
		return "", err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	current, ok := s.links[appUserID]
	if !ok {
		return "", usecase.ErrNotFound
	}
	current.AccessToken = refreshed.AccessToken
	if refreshed.RefreshToken != "" {
		current.RefreshToken = refreshed.RefreshToken
	}
	if refreshed.Scope != "" {
		current.Scope = refreshed.Scope
	}
	current.ExpiresAt = refreshed.ExpiresAt
	return current.AccessToken, nil
}
