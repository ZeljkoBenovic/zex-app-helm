package backend

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"

	"be/pkg/config"
	"be/pkg/web/types"
)

var ErrNoBackendURL = errors.New("backend url not found")

type Backend struct {
	log *slog.Logger
	url string
}

type UserData struct {
	Title   string
	AboutMe string
}

func NewBackend(log *slog.Logger, conf *config.Config) (*Backend, error) {
	var url string

	if beURL := os.Getenv("BE_URL"); beURL != "" {
		url = beURL
	} else {
		url = conf.Server.BackendURL
	}

	if url == "" && conf.Mode == config.Frontend {
		return nil, ErrNoBackendURL
	}

	return &Backend{
		log: log,
		url: url,
	}, nil
}

func (b *Backend) GetUserContent(id int) (UserData, error) {
	title, err := b.getTitle(id)
	if err != nil {
		return UserData{}, nil
	}

	aboutme, err := b.getAboutme(id)
	if err != nil {
		return UserData{}, nil
	}

	return UserData{
		Title:   title,
		AboutMe: aboutme,
	}, nil
}

func (b *Backend) getTitle(id int) (string, error) {
	title := &types.Response[string]{}
	if err := b.getRequester(fmt.Sprintf("title?id=%d", id), title); err != nil {
		return "", err
	}

	return title.Data, nil
}

func (b *Backend) getAboutme(id int) (string, error) {
	aboutme := &types.Response[string]{}
	if err := b.getRequester(fmt.Sprintf("aboutme?id=%d", id), aboutme); err != nil {
		return "", err
	}

	return aboutme.Data, nil
}

func (b *Backend) getRequester(endpoint string, result any) error {
	raw, err := http.Get(fmt.Sprintf("%s/%s", b.url, endpoint))
	if err != nil {
		return fmt.Errorf("could not make get request: %w", err)
	}

	bytesBody, _ := io.ReadAll(raw.Body)
	if err = json.Unmarshal(bytesBody, result); err != nil {
		return fmt.Errorf("could not unmarshal to json: %w", err)
	}

	return nil
}
