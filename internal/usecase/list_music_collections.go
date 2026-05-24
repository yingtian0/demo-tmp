package usecase

import "context"

type MusicCollectionOutput struct {
	Provider    string
	ID          string
	Name        string
	Description string
	ArtworkURL  string
	ExternalURL string
}

type ListMusicCollectionsInput struct {
	Provider    string
	AccessToken string
}

type ListMusicCollectionsUsecase struct {
	registry MusicSourceRegistry
}

func NewListMusicCollectionsUsecase(registry MusicSourceRegistry) *ListMusicCollectionsUsecase {
	return &ListMusicCollectionsUsecase{registry: registry}
}

func (u *ListMusicCollectionsUsecase) Execute(
	ctx context.Context,
	input ListMusicCollectionsInput,
) ([]MusicCollectionOutput, error) {
	if input.Provider == "" || input.AccessToken == "" {
		return nil, ErrInvalidInput
	}

	client, err := u.registry.Get(input.Provider)
	if err != nil {
		return nil, err
	}

	collections, err := client.ListCollections(ctx, input.AccessToken)
	if err != nil {
		return nil, err
	}

	outputs := make([]MusicCollectionOutput, 0, len(collections))
	for _, collection := range collections {
		outputs = append(outputs, MusicCollectionOutput{
			Provider:    collection.Provider,
			ID:          collection.ID,
			Name:        collection.Name,
			Description: collection.Description,
			ArtworkURL:  collection.ArtworkURL,
			ExternalURL: collection.ExternalURL,
		})
	}

	return outputs, nil
}
