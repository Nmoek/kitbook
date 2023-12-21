package repository

import (
	"context"
	"kitbook/internal/domain"
)

type HistoryRepository interface {
	AddRecord(ctx context.Context, record domain.HistoryRecord) error
}
