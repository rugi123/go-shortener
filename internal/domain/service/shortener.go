package service

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/rugi123/go-shortener/internal/domain/model"
)

const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

type Storage interface {
	SaveLink(ctx context.Context, link *model.Link) error
	GetLinkByKey(ctx context.Context, key string) (*model.Link, error)
}

type ShortenerService struct {
	storage   Storage
	keyLength int
}

func NewShortenerService(storage Storage, keyLength int) *ShortenerService {
	return &ShortenerService{
		storage:   storage,
		keyLength: keyLength,
	}
}

func (s *ShortenerService) GenerateKey() string {
	b := make([]byte, s.keyLength)
	for i := range b {
		b[i] = alphabet[rand.Intn(len(alphabet))]
	}
	return string(b)
}

func (s *ShortenerService) ShortenURL(ctx context.Context, original_url string) (string, error) {
	for attemp := 0; attemp < 5; attemp++ {
		key := s.GenerateKey()
		link := &model.Link{
			OriginalURL: original_url,
			ShortKey:    key,
		}
		err := s.storage.SaveLink(ctx, link)
		if err == nil {
			return key, nil
		}
		fmt.Println(err)
	}
	return "", fmt.Errorf("failed to generate unique short URL")
}

func (s *ShortenerService) ExpandURL(ctx context.Context, short_key string) (string, error) {
	link, err := s.storage.GetLinkByKey(ctx, short_key)
	if err != nil {
		return "", err
	}
	return link.OriginalURL, nil
}
