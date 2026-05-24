package usecase

import "context"

type MusicCollection struct {
	Provider    string
	ID          string
	Name        string
	Description string
	ArtworkURL  string
	ExternalURL string
}

type MusicTrack struct {
	Provider    string
	ID          string
	Title       string
	ArtistName  string
	PreviewURL  string
	ExternalURL string
	ArtworkURL  string
}

type MusicSourceClient interface {
	ListCollections(ctx context.Context, accessToken string) ([]MusicCollection, error)
	ListCollectionTracks(ctx context.Context, accessToken string, collectionID string) ([]MusicTrack, error)
	GetTrack(ctx context.Context, accessToken string, trackID string) (*MusicTrack, error)
}

type MusicSourceRegistry interface {
	Get(provider string) (MusicSourceClient, error)
}

type StaticMusicSourceRegistry struct {
	clients map[string]MusicSourceClient
}

func NewStaticMusicSourceRegistry(clients map[string]MusicSourceClient) *StaticMusicSourceRegistry {
	cloned := make(map[string]MusicSourceClient, len(clients))
	for provider, client := range clients {
		cloned[provider] = client
	}

	return &StaticMusicSourceRegistry{clients: cloned}
}

func (r *StaticMusicSourceRegistry) Get(provider string) (MusicSourceClient, error) {
	client, ok := r.clients[provider]
	if !ok || client == nil {
		return nil, ErrUnsupportedProvider
	}

	return client, nil
}
