package storage

import (
	"context"
	"log/slog"

	"be/pkg/config"
	"be/pkg/storage/db"
)

type Storer interface {
	StoreCloser
	TitleAboutMeGetter
}

type TitleAboutMeGetter interface {
	// GetTitle fetches the page title
	GetTitle(int32) (string, error)
	// GetAboutMe fetches about me content
	GetAboutMe(int32) (string, error)
}

type StoreCloser interface {
	Close() error
}

type Storage struct {
	db Storer
}

func NewStorage(ctx context.Context, log *slog.Logger, conf config.Config) (*Storage, error) {
	d, err := db.NewDb(ctx, log, conf)
	if err != nil {
		return nil, err
	}

	return &Storage{
		db: d,
	}, nil
}

func NewMockStorage() (*Storage, error) {
	return &Storage{
		db: db.NewMockDb(),
	}, nil
}

func (s *Storage) DB() Storer {
	return s.db
}

func (s *Storage) Close() error {
	return s.DB().Close()
}
