package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

const defaultConfigPath = "config.local.json"

type config struct {
	Server  serverConfig  `json:"server"`
	Spotify spotifyConfig `json:"spotify"`
}

type serverConfig struct {
	Addr string `json:"addr"`
}

type spotifyConfig struct {
	ClientID     string   `json:"client_id"`
	ClientSecret string   `json:"client_secret"`
	RedirectURI  string   `json:"redirect_uri"`
	Scopes       []string `json:"scopes"`
}

func loadConfig() (*config, error) {
	cfg := defaultConfig()

	if path := configPath(); path != "" {
		fileCfg, err := loadConfigFile(path)
		if err != nil {
			return nil, err
		}
		mergeConfig(cfg, fileCfg)
	}

	overrideWithEnv(cfg)

	if cfg.Spotify.ClientID == "" || cfg.Spotify.ClientSecret == "" || cfg.Spotify.RedirectURI == "" {
		return nil, errors.New("spotify config is required: set config.local.json or SPOTIFY_CLIENT_ID, SPOTIFY_CLIENT_SECRET, SPOTIFY_REDIRECT_URI")
	}

	return cfg, nil
}

func defaultConfig() *config {
	return &config{
		Server: serverConfig{
			Addr: ":8080",
		},
		Spotify: spotifyConfig{
			Scopes: []string{
				"playlist-read-private",
				"playlist-read-collaborative",
			},
		},
	}
}

func configPath() string {
	if path := os.Getenv("APP_CONFIG_FILE"); path != "" {
		return path
	}

	if _, err := os.Stat(defaultConfigPath); err == nil {
		return defaultConfigPath
	}

	return ""
}

func loadConfigFile(path string) (*config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config file %s: %w", path, err)
	}

	var cfg config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parse config file %s: %w", path, err)
	}

	return &cfg, nil
}

func mergeConfig(dst *config, src *config) {
	if src == nil {
		return
	}
	if src.Server.Addr != "" {
		dst.Server.Addr = src.Server.Addr
	}
	if src.Spotify.ClientID != "" {
		dst.Spotify.ClientID = src.Spotify.ClientID
	}
	if src.Spotify.ClientSecret != "" {
		dst.Spotify.ClientSecret = src.Spotify.ClientSecret
	}
	if src.Spotify.RedirectURI != "" {
		dst.Spotify.RedirectURI = src.Spotify.RedirectURI
	}
	if len(src.Spotify.Scopes) > 0 {
		dst.Spotify.Scopes = append([]string{}, src.Spotify.Scopes...)
	}
}

func overrideWithEnv(cfg *config) {
	if value := os.Getenv("SERVER_ADDR"); value != "" {
		cfg.Server.Addr = value
	}
	if value := os.Getenv("SPOTIFY_CLIENT_ID"); value != "" {
		cfg.Spotify.ClientID = value
	}
	if value := os.Getenv("SPOTIFY_CLIENT_SECRET"); value != "" {
		cfg.Spotify.ClientSecret = value
	}
	if value := os.Getenv("SPOTIFY_REDIRECT_URI"); value != "" {
		cfg.Spotify.RedirectURI = value
	}
}
