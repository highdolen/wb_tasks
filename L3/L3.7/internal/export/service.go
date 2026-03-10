package export

import (
	"context"
	"encoding/csv"
	"io"
	"strconv"
	"time"
)

// Service - сервис для экспорта истории изменений
type Service struct {
	repo *Repository
}

// NewService - создание сервиса экспорта
func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// ExportHistoryCSV - экспорт истории изменений в CSV файл
func (s *Service) ExportHistoryCSV(ctx context.Context, writer io.Writer) error {

	history, err := s.repo.GetAllHistory(ctx)
	if err != nil {
		return err
	}

	csvWriter := csv.NewWriter(writer)

	err = csvWriter.Write([]string{
		"id",
		"item_id",
		"action",
		"old_value",
		"new_value",
		"changed_by",
		"changed_at",
	})
	if err != nil {
		return err
	}

	for _, h := range history {

		row := []string{
			strconv.Itoa(h.ID),
			strconv.Itoa(h.ItemID),
			h.Action,
			h.OldValue,
			h.NewValue,
			h.ChangedBy,
			h.ChangedAt.Format(time.RFC3339),
		}

		err := csvWriter.Write(row)
		if err != nil {
			return err
		}
	}

	csvWriter.Flush()

	return csvWriter.Error()
}
