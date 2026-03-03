package service

import (
	"context"
	"io"
	"time"

	"salesTracker/internal/repository"
)

type ExportService struct {
	repo *repository.ItemRepository
}

func NewExportService(repo *repository.ItemRepository) *ExportService {
	return &ExportService{repo: repo}
}

// Принимаем time.Time, увеличиваем до конца дня и передаём в репозиторий
func (s *ExportService) ExportCSV(ctx context.Context, from, to time.Time, w io.Writer) error {
	// Включаем весь день
	to = to.Add(23*time.Hour + 59*time.Minute + 59*time.Second)
	return s.repo.ExportCSV(ctx, from, to, w)
}
