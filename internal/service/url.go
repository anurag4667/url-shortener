package service

import (
	"github.com/anurag4667/url-shortener/internal/database"
	"github.com/anurag4667/url-shortener/internal/short"
)

type URLService struct {
	store *database.MySQLStore
}

func New(store *database.MySQLStore) *URLService {
	return &URLService{store: store}
}

func (s *URLService) Shorten(original string) (string, error) {
	for {
		id, err := short.Generate()
		if err != nil {
			return "", err
		}

		err = s.store.Save(id, original)
		if err == nil {
			return id, nil
		}
		// collision → retry
	}
}

func (s *URLService) Resolve(id string) (string, bool, error) {
	url, ok, err := s.store.Get(id)
	if ok {
		s.store.IncrementClicks(id)
	}
	return url, ok, err
}
