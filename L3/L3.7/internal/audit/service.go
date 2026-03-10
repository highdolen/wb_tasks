package audit

import (
	"context"
	"encoding/json"
)

// Service - сервис для работы с историей изменений товаров
type Service struct {
	repo *Repository
}

// NewService - создание нового сервиса истории
func NewService(repo *Repository) *Service {
	return &Service{
		repo: repo,
	}
}

// GetHistory - получение истории изменений товара
func (s *Service) GetHistory(ctx context.Context, itemID int) ([]History, error) {
	return s.repo.GetHistory(ctx, itemID)
}

// FilterHistory - фильтрация истории изменений по параметрам
func (s *Service) FilterHistory(
	ctx context.Context,
	user, action, from, to string,
) ([]History, error) {

	return s.repo.FilterHistory(ctx, user, action, from, to)
}

// GetDiff - получение различий между старой и новой версией товара
func (s *Service) GetDiff(oldStr, newStr string) (map[string]interface{}, error) {

	var oldData map[string]interface{}
	var newData map[string]interface{}

	if err := json.Unmarshal([]byte(oldStr), &oldData); err != nil {
		return nil, err
	}
	if err := json.Unmarshal([]byte(newStr), &newData); err != nil {
		return nil, err
	}

	diff := make(map[string]interface{})

	for key, newVal := range newData {

		oldVal := oldData[key]

		if oldVal != newVal {
			diff[key] = []interface{}{oldVal, newVal}
		}
	}

	return diff, nil
}
