package paste

import (
	"context"

	"github.com/gabrielolivrp/pastebin-api/pkg/database"
)

type PasteRepository interface {
	Create(ctx context.Context, user Paste) error
	GetByID(ctx context.Context, id string) (Paste, error)
}

type pasteRepository struct {
	dbClient database.Client
}

func NewPasteRepository(db database.Client) PasteRepository {
	return &pasteRepository{
		dbClient: db,
	}
}

func (r *pasteRepository) Create(ctx context.Context, paste Paste) error {
	if err := r.dbClient.DB().WithContext(ctx).Create(&paste).Error; err != nil {
		return err
	}
	return nil
}

func (r *pasteRepository) GetByID(ctx context.Context, id string) (Paste, error) {
	var paste Paste
	if err := r.dbClient.DB().WithContext(ctx).Where("id = ?", id).First(&paste).Error; err != nil {
		return Paste{}, err
	}
	return paste, nil
}
