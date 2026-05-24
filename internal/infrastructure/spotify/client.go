package spotify

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	"demo-tmp/internal/usecase"
)

const (
	ProviderName   = "spotify"
	defaultBaseURL = "https://api.spotify.com/v1"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type Client struct {
	baseURL    *url.URL
	httpClient HTTPClient
}

func NewClient(httpClient HTTPClient) (*Client, error) {
	return NewClientWithBaseURL(httpClient, defaultBaseURL)
}

func NewClientWithBaseURL(httpClient HTTPClient, rawBaseURL string) (*Client, error) {
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 10 * time.Second}
	}

	baseURL, err := url.Parse(rawBaseURL)
	if err != nil {
		return nil, err
	}

	return &Client{
		baseURL:    baseURL,
		httpClient: httpClient,
	}, nil
}

func (c *Client) ListCollections(ctx context.Context, accessToken string) ([]usecase.MusicCollection, error) {
	if strings.TrimSpace(accessToken) == "" {
		return nil, usecase.ErrInvalidInput
	}

	var collections []usecase.MusicCollection
	nextPath := "/me/playlists?limit=50"

	for nextPath != "" {
		var page playlistsPage
		if err := c.get(ctx, accessToken, nextPath, &page); err != nil {
			return nil, err
		}

		for _, item := range page.Items {
			collections = append(collections, usecase.MusicCollection{
				Provider:    ProviderName,
				ID:          item.ID,
				Name:        item.Name,
				Description: item.Description,
				ArtworkURL:  firstImageURL(item.Images),
				ExternalURL: item.ExternalURLs.Spotify,
			})
		}

		nextPath = relativePathFromAbsolute(c.baseURL, page.Next)
	}

	return collections, nil
}

func (c *Client) ListCollectionTracks(
	ctx context.Context,
	accessToken string,
	collectionID string,
) ([]usecase.MusicTrack, error) {
	if strings.TrimSpace(accessToken) == "" || strings.TrimSpace(collectionID) == "" {
		return nil, usecase.ErrInvalidInput
	}

	var tracks []usecase.MusicTrack
	nextPath := fmt.Sprintf("/playlists/%s/items?limit=100&additional_types=track", url.PathEscape(collectionID))

	for nextPath != "" {
		var page playlistItemsPage
		if err := c.get(ctx, accessToken, nextPath, &page); err != nil {
			return nil, err
		}

		for _, item := range page.Items {
			if item.Track == nil || item.Track.Type != "track" || item.Track.IsLocal {
				continue
			}
			tracks = append(tracks, toMusicTrack(*item.Track))
		}

		nextPath = relativePathFromAbsolute(c.baseURL, page.Next)
	}

	return tracks, nil
}

func (c *Client) GetTrack(ctx context.Context, accessToken string, trackID string) (*usecase.MusicTrack, error) {
	if strings.TrimSpace(accessToken) == "" || strings.TrimSpace(trackID) == "" {
		return nil, usecase.ErrInvalidInput
	}

	var track spotifyTrack
	if err := c.get(ctx, accessToken, fmt.Sprintf("/tracks/%s", url.PathEscape(trackID)), &track); err != nil {
		return nil, err
	}

	result := toMusicTrack(track)
	return &result, nil
}

func (c *Client) get(ctx context.Context, accessToken string, rawPath string, dest any) error {
	requestURL, err := c.resolve(rawPath)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return usecase.ErrPermissionDenied
	}
	if resp.StatusCode == http.StatusNotFound {
		return usecase.ErrNotFound
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		return fmt.Errorf("spotify api error: status=%d body=%s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	if err := json.NewDecoder(resp.Body).Decode(dest); err != nil {
		return err
	}

	return nil
}

func (c *Client) resolve(rawPath string) (string, error) {
	if rawPath == "" {
		return "", errors.New("empty path")
	}

	u, err := url.Parse(rawPath)
	if err != nil {
		return "", err
	}
	if u.IsAbs() {
		return u.String(), nil
	}

	base := *c.baseURL
	base.Path = path.Join(c.baseURL.Path, u.Path)
	base.RawQuery = u.RawQuery
	return base.String(), nil
}

func relativePathFromAbsolute(base *url.URL, raw string) string {
	if strings.TrimSpace(raw) == "" {
		return ""
	}

	nextURL, err := url.Parse(raw)
	if err != nil {
		return ""
	}
	if !nextURL.IsAbs() {
		return raw
	}

	if !sameHost(base, nextURL) {
		return raw
	}

	if nextURL.RawQuery == "" {
		return nextURL.Path
	}

	return nextURL.Path + "?" + nextURL.RawQuery
}

func sameHost(a *url.URL, b *url.URL) bool {
	return a != nil && b != nil && strings.EqualFold(a.Scheme, b.Scheme) && strings.EqualFold(a.Host, b.Host)
}

func toMusicTrack(track spotifyTrack) usecase.MusicTrack {
	return usecase.MusicTrack{
		Provider:    ProviderName,
		ID:          track.ID,
		Title:       track.Name,
		ArtistName:  joinArtistNames(track.Artists),
		PreviewURL:  track.PreviewURL,
		ExternalURL: track.ExternalURLs.Spotify,
		ArtworkURL:  firstImageURL(track.Album.Images),
	}
}

func joinArtistNames(artists []spotifyArtist) string {
	names := make([]string, 0, len(artists))
	for _, artist := range artists {
		if strings.TrimSpace(artist.Name) == "" {
			continue
		}
		names = append(names, artist.Name)
	}

	return strings.Join(names, ", ")
}

func firstImageURL(images []spotifyImage) string {
	if len(images) == 0 {
		return ""
	}
	return images[0].URL
}

type playlistsPage struct {
	Items []spotifyPlaylist `json:"items"`
	Next  string            `json:"next"`
}

type playlistItemsPage struct {
	Items []spotifyPlaylistItem `json:"items"`
	Next  string                `json:"next"`
}

type spotifyPlaylist struct {
	ID           string         `json:"id"`
	Name         string         `json:"name"`
	Description  string         `json:"description"`
	Images       []spotifyImage `json:"images"`
	ExternalURLs spotifyURLMap  `json:"external_urls"`
}

type spotifyPlaylistItem struct {
	Track *spotifyTrack `json:"track"`
}

type spotifyTrack struct {
	ID           string          `json:"id"`
	Name         string          `json:"name"`
	Type         string          `json:"type"`
	IsLocal      bool            `json:"is_local"`
	PreviewURL   string          `json:"preview_url"`
	Artists      []spotifyArtist `json:"artists"`
	Album        spotifyAlbum    `json:"album"`
	ExternalURLs spotifyURLMap   `json:"external_urls"`
}

type spotifyArtist struct {
	Name string `json:"name"`
}

type spotifyAlbum struct {
	Images []spotifyImage `json:"images"`
}

type spotifyImage struct {
	URL string `json:"url"`
}

type spotifyURLMap struct {
	Spotify string `json:"spotify"`
}
