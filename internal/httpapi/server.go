package httpapi

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"backend/internal/domain"
	"backend/internal/infrastructure/memory"
	"backend/internal/infrastructure/spotify"
	"backend/internal/usecase"
)

type Server struct {
	oauthClient                        *spotify.OAuthClient
	spotifyClient                      *spotify.Client
	spotifyLinks                       *memory.SpotifyLinkStore
	createFavoriteTrackFromMusicSource *usecase.CreateFavoriteTrackFromMusicSourceUsecase
}

func NewServer(
	oauthClient *spotify.OAuthClient,
	spotifyClient *spotify.Client,
	spotifyLinks *memory.SpotifyLinkStore,
	createFavoriteTrackFromMusicSource *usecase.CreateFavoriteTrackFromMusicSourceUsecase,
) *Server {
	return &Server{
		oauthClient:                        oauthClient,
		spotifyClient:                      spotifyClient,
		spotifyLinks:                       spotifyLinks,
		createFavoriteTrackFromMusicSource: createFavoriteTrackFromMusicSource,
	}
}

func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", s.handleHealthz)
	mux.HandleFunc("/auth/spotify/login", s.handleSpotifyLogin)
	mux.HandleFunc("/auth/spotify/callback", s.handleSpotifyCallback)
	mux.HandleFunc("/api/spotify/link", s.handleGetSpotifyLink)
	mux.HandleFunc("/api/music/collections", s.handleListCollections)
	mux.HandleFunc("/api/music/collections/", s.handleCollectionRoutes)
	mux.HandleFunc("/api/favorite-tracks/import", s.handleImportFavoriteTrack)
	return mux
}

func (s *Server) handleHealthz(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *Server) handleSpotifyLogin(w http.ResponseWriter, r *http.Request) {
	appUserID := strings.TrimSpace(r.URL.Query().Get("user_id"))
	if appUserID == "" {
		writeError(w, http.StatusBadRequest, "user_id is required")
		return
	}

	url, err := s.oauthClient.AuthorizeURL(appUserID)
	if err != nil {
		writeInternalError(w, err)
		return
	}

	http.Redirect(w, r, url, http.StatusFound)
}

func (s *Server) handleSpotifyCallback(w http.ResponseWriter, r *http.Request) {
	if errParam := r.URL.Query().Get("error"); errParam != "" {
		writeError(w, http.StatusBadRequest, errParam)
		return
	}

	appUserID := strings.TrimSpace(r.URL.Query().Get("state"))
	code := strings.TrimSpace(r.URL.Query().Get("code"))
	if appUserID == "" || code == "" {
		writeError(w, http.StatusBadRequest, "state and code are required")
		return
	}

	token, err := s.oauthClient.ExchangeCode(r.Context(), code)
	if err != nil {
		writeInternalError(w, err)
		return
	}

	profile, err := s.oauthClient.GetCurrentUserProfile(r.Context(), token.AccessToken)
	if err != nil {
		writeInternalError(w, err)
		return
	}

	s.spotifyLinks.Save(&memory.SpotifyAccountLink{
		AppUserID:          appUserID,
		SpotifyUserID:      profile.ID,
		SpotifyDisplayName: profile.DisplayName,
		AccessToken:        token.AccessToken,
		RefreshToken:       token.RefreshToken,
		Scope:              token.Scope,
		ExpiresAt:          token.ExpiresAt,
	})

	writeJSON(w, http.StatusOK, map[string]string{
		"status":               "linked",
		"user_id":              appUserID,
		"spotify_user_id":      profile.ID,
		"spotify_display_name": profile.DisplayName,
	})
}

func (s *Server) handleGetSpotifyLink(w http.ResponseWriter, r *http.Request) {
	appUserID := strings.TrimSpace(r.URL.Query().Get("user_id"))
	if appUserID == "" {
		writeError(w, http.StatusBadRequest, "user_id is required")
		return
	}

	link, err := s.spotifyLinks.Get(appUserID)
	if err != nil {
		writeDomainError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"user_id":              link.AppUserID,
		"spotify_user_id":      link.SpotifyUserID,
		"spotify_display_name": link.SpotifyDisplayName,
		"scope":                link.Scope,
		"expires_at":           link.ExpiresAt,
	})
}

func (s *Server) handleListCollections(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	appUserID, token, err := s.accessTokenForRequest(r.Context(), r)
	if err != nil {
		writeDomainError(w, err)
		return
	}

	registry := usecase.NewStaticMusicSourceRegistry(map[string]usecase.MusicSourceClient{
		spotify.ProviderName: s.spotifyClient,
	})
	uc := usecase.NewListMusicCollectionsUsecase(registry)
	output, err := uc.Execute(r.Context(), usecase.ListMusicCollectionsInput{
		Provider:    spotify.ProviderName,
		AccessToken: token,
	})
	if err != nil {
		writeDomainError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"user_id":     appUserID,
		"provider":    spotify.ProviderName,
		"collections": output,
	})
}

func (s *Server) handleCollectionRoutes(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	trimmed := strings.TrimPrefix(r.URL.Path, "/api/music/collections/")
	parts := strings.Split(strings.Trim(trimmed, "/"), "/")
	if len(parts) != 2 || parts[1] != "tracks" || parts[0] == "" {
		writeError(w, http.StatusNotFound, "not found")
		return
	}

	collectionID := parts[0]
	appUserID, token, err := s.accessTokenForRequest(r.Context(), r)
	if err != nil {
		writeDomainError(w, err)
		return
	}

	registry := usecase.NewStaticMusicSourceRegistry(map[string]usecase.MusicSourceClient{
		spotify.ProviderName: s.spotifyClient,
	})
	uc := usecase.NewListMusicCollectionTracksUsecase(registry)
	output, err := uc.Execute(r.Context(), usecase.ListMusicCollectionTracksInput{
		Provider:     spotify.ProviderName,
		AccessToken:  token,
		CollectionID: collectionID,
	})
	if err != nil {
		writeDomainError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"user_id":       appUserID,
		"provider":      spotify.ProviderName,
		"collection_id": collectionID,
		"tracks":        output,
	})
}

func (s *Server) handleImportFavoriteTrack(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	appUserID, token, err := s.accessTokenForRequest(r.Context(), r)
	if err != nil {
		writeDomainError(w, err)
		return
	}

	var input struct {
		TrackID        string `json:"track_id"`
		Reason         string `json:"reason"`
		ListeningScene string `json:"listening_scene"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json body")
		return
	}

	output, err := s.createFavoriteTrackFromMusicSource.Execute(r.Context(), usecase.CreateFavoriteTrackFromMusicSourceInput{
		UserID:         domain.UserID(appUserID),
		Provider:       spotify.ProviderName,
		AccessToken:    token,
		TrackID:        input.TrackID,
		Reason:         input.Reason,
		ListeningScene: input.ListeningScene,
	})
	if err != nil {
		writeDomainError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, output)
}

func (s *Server) accessTokenForRequest(ctx context.Context, r *http.Request) (string, string, error) {
	appUserID := strings.TrimSpace(r.URL.Query().Get("user_id"))
	if appUserID == "" {
		return "", "", usecase.ErrInvalidInput
	}

	token, err := s.spotifyLinks.GetValidAccessToken(ctx, appUserID, oauthTokenRefresher{s.oauthClient})
	if err != nil {
		return "", "", err
	}

	return appUserID, token, nil
}

type oauthTokenRefresher struct {
	client *spotify.OAuthClient
}

func (o oauthTokenRefresher) RefreshToken(ctx context.Context, refreshToken string) (*memory.SpotifyOAuthToken, error) {
	token, err := o.client.RefreshToken(ctx, refreshToken)
	if err != nil {
		return nil, err
	}

	return &memory.SpotifyOAuthToken{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		Scope:        token.Scope,
		ExpiresAt:    token.ExpiresAt,
	}, nil
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}

func writeInternalError(w http.ResponseWriter, err error) {
	writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
}

func writeDomainError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, usecase.ErrInvalidInput):
		writeError(w, http.StatusBadRequest, err.Error())
	case errors.Is(err, usecase.ErrNotFound):
		writeError(w, http.StatusNotFound, err.Error())
	case errors.Is(err, usecase.ErrPermissionDenied):
		writeError(w, http.StatusForbidden, err.Error())
	default:
		writeInternalError(w, err)
	}
}
