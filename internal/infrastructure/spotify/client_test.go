package spotify

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"testing"
)

type roundTripFunc func(req *http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func newTestHTTPClient(fn roundTripFunc) *http.Client {
	return &http.Client{Transport: fn}
}

func jsonResponse(body string) *http.Response {
	return &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Type": []string{"application/json"}},
		Body:       io.NopCloser(bytes.NewBufferString(body)),
	}
}

func TestClientListCollections(t *testing.T) {
	httpClient := newTestHTTPClient(func(r *http.Request) (*http.Response, error) {
		if r.URL.Path != "/v1/me/playlists" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if got := r.Header.Get("Authorization"); got != "Bearer token" {
			t.Fatalf("unexpected auth header: %s", got)
		}

		return jsonResponse(`{
			"items": [
				{
					"id": "pl1",
					"name": "Favorites",
					"description": "my playlist",
					"external_urls": {"spotify": "https://open.spotify.com/playlist/pl1"},
					"images": [{"url": "https://image/1"}]
				}
			],
			"next": ""
		}`), nil
	})

	client, err := NewClientWithBaseURL(httpClient, "https://api.spotify.test/v1")
	if err != nil {
		t.Fatal(err)
	}

	got, err := client.ListCollections(context.Background(), "token")
	if err != nil {
		t.Fatal(err)
	}

	if len(got) != 1 {
		t.Fatalf("expected 1 collection, got %d", len(got))
	}
	if got[0].Provider != ProviderName || got[0].ID != "pl1" || got[0].ArtworkURL != "https://image/1" {
		t.Fatalf("unexpected collection: %+v", got[0])
	}
}

func TestClientListCollectionTracksSkipsNonTrackItems(t *testing.T) {
	httpClient := newTestHTTPClient(func(r *http.Request) (*http.Response, error) {
		if r.URL.Path != "/v1/playlists/pl1/items" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}

		return jsonResponse(`{
			"items": [
				{"track": null},
				{"track": {"type": "episode", "id": "ep1", "name": "ep"}},
				{"track": {
					"type": "track",
					"id": "tr1",
					"name": "Song",
					"preview_url": "https://preview",
					"is_local": false,
					"artists": [{"name": "A1"}, {"name": "A2"}],
					"album": {"images": [{"url": "https://art"}]},
					"external_urls": {"spotify": "https://open.spotify.com/track/tr1"}
				}}
			],
			"next": ""
		}`), nil
	})

	client, err := NewClientWithBaseURL(httpClient, "https://api.spotify.test/v1")
	if err != nil {
		t.Fatal(err)
	}

	got, err := client.ListCollectionTracks(context.Background(), "token", "pl1")
	if err != nil {
		t.Fatal(err)
	}

	if len(got) != 1 {
		t.Fatalf("expected 1 track, got %d", len(got))
	}
	if got[0].ArtistName != "A1, A2" || got[0].ExternalURL != "https://open.spotify.com/track/tr1" {
		t.Fatalf("unexpected track: %+v", got[0])
	}
}

func TestClientGetTrack(t *testing.T) {
	httpClient := newTestHTTPClient(func(r *http.Request) (*http.Response, error) {
		if r.URL.Path != "/v1/tracks/tr1" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}

		return jsonResponse(`{
			"id": "tr1",
			"name": "Song",
			"type": "track",
			"preview_url": "",
			"is_local": false,
			"artists": [{"name": "Artist"}],
			"album": {"images": [{"url": "https://art"}]},
			"external_urls": {"spotify": "https://open.spotify.com/track/tr1"}
		}`), nil
	})

	client, err := NewClientWithBaseURL(httpClient, "https://api.spotify.test/v1")
	if err != nil {
		t.Fatal(err)
	}

	got, err := client.GetTrack(context.Background(), "token", "tr1")
	if err != nil {
		t.Fatal(err)
	}

	if got == nil || got.Provider != ProviderName || got.Title != "Song" {
		t.Fatalf("unexpected track: %+v", got)
	}
}
