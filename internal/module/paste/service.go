package paste

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/gabrielolivrp/pastebin-api/pkg/cache"
	"github.com/gabrielolivrp/pastebin-api/pkg/database"
	"github.com/google/uuid"
)

var (
	ErrPasteNotFound = errors.New("paste not found")
)

type PasteService interface {
	Create(ctx context.Context, params CreatePasteParams) (Paste, error)
	GetByID(ctx context.Context, id string) (Paste, error)
}

type pasteService struct {
	repo        PasteRepository
	cacheClient cache.Client
}

func NewPasteService(repo PasteRepository, cacheClient cache.Client) PasteService {
	return &pasteService{
		repo:        repo,
		cacheClient: cacheClient,
	}
}

type CreatePasteParams struct {
	Content string
	Title   string
	Lang    string
}

type CreatePasteResult struct {
	ID        uuid.UUID  `json:"id"`
	Title     string     `json:"title"`
	Content   string     `json:"content"`
	Lang      string     `json:"lang"`
	ExpiresAt *time.Time `json:"expires_at"`
	CreatedAt time.Time  `json:"created_at"`
}

func (s *pasteService) Create(ctx context.Context, params CreatePasteParams) (Paste, error) {
	expiresAt := calculateExpiration()
	paste := Paste{
		ID:        uuid.New(),
		Content:   params.Content,
		Title:     params.Title,
		Lang:      params.Lang,
		CreatedAt: time.Now(),
		ExpiresAt: expiresAt,
	}

	if err := s.repo.Create(ctx, paste); err != nil {
		return Paste{}, err
	}
	return paste, nil
}

type GetPasteByIDResult struct {
	ID        uuid.UUID  `json:"id"`
	Title     string     `json:"title"`
	Content   string     `json:"content"`
	Lang      string     `json:"lang"`
	ExpiresAt *time.Time `json:"expires_at"`
	CreatedAt time.Time  `json:"created_at"`
}

func (s *pasteService) GetByID(ctx context.Context, id string) (Paste, error) {
	exists, err := s.cacheClient.Has(ctx, id)
	if err != nil {
		return Paste{}, err
	}

	if exists {
		data, err := s.cacheClient.Get(ctx, id)
		if err != nil {
			return Paste{}, err
		}

		var paste Paste
		if err := json.Unmarshal([]byte(data), &paste); err != nil {
			return Paste{}, err
		}
		return paste, nil
	}

	paste, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if database.ErrNotFound(err) {
			return Paste{}, ErrPasteNotFound
		}
		return Paste{}, err
	}

	data, err := json.Marshal(paste)
	if err == nil {
		err = s.cacheClient.Set(ctx, id, data, 10*time.Minute)
		if err != nil {
			return Paste{}, err
		}
	}

	return paste, nil
}

func calculateExpiration() *time.Time {
	expiresAt := time.Now().Add(24 * time.Hour)
	return &expiresAt
}
