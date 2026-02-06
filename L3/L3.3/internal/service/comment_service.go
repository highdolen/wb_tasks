package service

import (
	"comment/internal/models"
	"comment/internal/storage"
)

type CommentService struct {
	store *storage.MemoryStorage
}

func NewCommentService(store *storage.MemoryStorage) *CommentService {
	return &CommentService{store: store}
}

func (s *CommentService) Create(text string, parentID *int64) *models.Comment {
	return s.store.Add(text, parentID)
}

func (s *CommentService) Delete(id int64) {
	s.store.DeleteTree(id)
}

func (s *CommentService) GetByID(id int64) (*models.Comment, bool) {
	return s.store.GetByID(id)
}

func (s *CommentService) GetChildren(parentID int64) []*models.Comment {
	return s.store.GetChildren(parentID)
}

func (s *CommentService) GetRoots(limit, offset int) []*models.Comment {
	return s.store.GetRoots(limit, offset)
}

func (s *CommentService) Search(q string) []*models.Comment {
	return s.store.Search(q)
}
