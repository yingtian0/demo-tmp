package main

import (
	"log"
	"net/http"

	"backend/internal/httpapi"
	"backend/internal/infrastructure/memory"
	"backend/internal/infrastructure/spotify"
	"backend/internal/usecase"
)

func main() {
	cfg, err := loadConfig()
	if err != nil {
		log.Fatal(err)
	}

	spotifyHTTPClient := &http.Client{}

	spotifyClient, err := spotify.NewClient(spotifyHTTPClient)
	if err != nil {
		log.Fatal(err)
	}

	oauthClient := spotify.NewOAuthClient(spotify.OAuthConfig{
		ClientID:     cfg.Spotify.ClientID,
		ClientSecret: cfg.Spotify.ClientSecret,
		RedirectURI:  cfg.Spotify.RedirectURI,
		Scopes:       cfg.Spotify.Scopes,
	}, spotifyHTTPClient)

	registry := usecase.NewStaticMusicSourceRegistry(map[string]usecase.MusicSourceClient{
		spotify.ProviderName: spotifyClient,
	})

	favoriteTrackRepository := memory.NewFavoriteTrackRepository()
	idGenerator := memory.NewIncrementalIDGenerator("fav_")
	clock := usecase.RealClock{}
	spotifyLinks := memory.NewSpotifyLinkStore()

	createFavoriteTrackFromMusicSource := usecase.NewCreateFavoriteTrackFromMusicSourceUsecase(
		registry,
		favoriteTrackRepository,
		idGenerator,
		clock,
	)

	server := httpapi.NewServer(
		oauthClient,
		spotifyClient,
		spotifyLinks,
		createFavoriteTrackFromMusicSource,
	)

	log.Printf("listening on %s", cfg.Server.Addr)
	log.Fatal(http.ListenAndServe(cfg.Server.Addr, server.Handler()))
}
