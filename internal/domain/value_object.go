package domain

import (
	"strings"
)

type Track struct {
	title       string
	artistName  string
	previewURL  string
	externalURL string
	artworkURL  string
}

func NewTrack(
	title string,
	artistName string,
	previewURL string,
	externalURL string,
	artworkURL string,
) (Track, error) {
	if isBlank(title) || isBlank(artistName) {
		return Track{}, ErrInvalidTrack
	}

	return Track{
		title:       strings.TrimSpace(title),
		artistName:  strings.TrimSpace(artistName),
		previewURL:  strings.TrimSpace(previewURL),
		externalURL: strings.TrimSpace(externalURL),
		artworkURL:  strings.TrimSpace(artworkURL),
	}, nil
}

func (t Track) Title() string {
	return t.title
}

func (t Track) ArtistName() string {
	return t.artistName
}

func (t Track) PreviewURL() string {
	return t.previewURL
}

func (t Track) ExternalURL() string {
	return t.externalURL
}

func (t Track) ArtworkURL() string {
	return t.artworkURL
}

type EncounterSource string

const (
	EncounterSourceQR       EncounterSource = "qr"
	EncounterSourceSameTime EncounterSource = "same_time"
	EncounterSourceLocation EncounterSource = "location"
	EncounterSourceBLE      EncounterSource = "ble"
)

func (s EncounterSource) IsValid() bool {
	switch s {
	case EncounterSourceQR,
		EncounterSourceSameTime,
		EncounterSourceLocation,
		EncounterSourceBLE:
		return true
	default:
		return false
	}
}

type ReactionType string

const (
	ReactionTypeLiked     ReactionType = "liked"
	ReactionTypeSaved     ReactionType = "saved"
	ReactionTypeSurprised ReactionType = "surprised"
	ReactionTypeListened  ReactionType = "listened"
)

func (t ReactionType) IsValid() bool {
	switch t {
	case ReactionTypeLiked,
		ReactionTypeSaved,
		ReactionTypeSurprised,
		ReactionTypeListened:
		return true
	default:
		return false
	}
}

func NormalizeTags(tags []string) []string {
	result := make([]string, 0, len(tags))
	seen := map[string]struct{}{}

	for _, tag := range tags {
		t := strings.TrimSpace(tag)
		if t == "" {
			continue
		}

		if _, ok := seen[t]; ok {
			continue
		}

		seen[t] = struct{}{}
		result = append(result, t)
	}

	return result
}
