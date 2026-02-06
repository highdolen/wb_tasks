package storage

import (
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"comment/internal/models"
)

type MemoryStorage struct {
	mu       sync.RWMutex
	comments map[int64]*models.Comment
	nextID   int64
}

// NewMemoryStorage - констуктор MemoryStorage
func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		comments: make(map[int64]*models.Comment),
		nextID:   1,
	}
}

// Add - добавления комментария в in-memory
func (s *MemoryStorage) Add(text string, parentID *int64) *models.Comment {
	s.mu.Lock()
	defer s.mu.Unlock()

	id := s.nextID
	s.nextID++

	var path string
	var level int
	if parentID == nil {
		path = fmt.Sprintf("%06d", id)
		level = 0
	} else {
		parent, ok := s.comments[*parentID]
		if !ok {
			return nil
		}
		path = parent.Path + "." + fmt.Sprintf("%06d", id)
		level = parent.Level + 1
	}

	c := &models.Comment{
		ID:        id,
		ParentID:  parentID,
		Text:      text,
		CreatedAt: time.Now(),
		Path:      path,
		Level:     level,
	}

	s.comments[id] = c
	return c
}

// GetByID - получение комментария по айди
func (s *MemoryStorage) GetByID(id int64) (*models.Comment, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	c, ok := s.comments[id]
	return c, ok
}

// GetChildren - получение дочерних комментариев
func (s *MemoryStorage) GetChildren(parentID int64) []*models.Comment {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var children []*models.Comment
	for _, c := range s.comments {
		if c.ParentID != nil && *c.ParentID == parentID {
			children = append(children, c)
		}
	}

	sort.Slice(children, func(i, j int) bool {
		return children[i].CreatedAt.Before(children[j].CreatedAt)
	})

	return children
}

// GetRoots - получение только корневых комментариев(для пагинации)
func (s *MemoryStorage) GetRoots(limit, offset int) []*models.Comment {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var roots []*models.Comment
	for _, c := range s.comments {
		if c.ParentID == nil {
			roots = append(roots, c)
		}
	}

	sort.Slice(roots, func(i, j int) bool {
		return roots[i].CreatedAt.Before(roots[j].CreatedAt)
	})

	if offset >= len(roots) {
		return []*models.Comment{}
	}

	end := offset + limit
	if end > len(roots) {
		end = len(roots)
	}

	return roots[offset:end]
}

// DeleteTree - удаление комменатриев(дерева)
func (s *MemoryStorage) DeleteTree(id int64) {
	s.mu.Lock()
	defer s.mu.Unlock()

	root, ok := s.comments[id]
	if !ok {
		return
	}

	prefix := root.Path
	for cid, c := range s.comments {
		if strings.HasPrefix(c.Path, prefix) {
			delete(s.comments, cid)
		}
	}
}

// Search - поиск комментариев по ключевым словам
func (s *MemoryStorage) Search(q string) []*models.Comment {
	s.mu.RLock()
	defer s.mu.RUnlock()
	q = strings.ToLower(q)

	var result []*models.Comment
	for _, c := range s.comments {
		if strings.Contains(strings.ToLower(c.Text), q) {
			result = append(result, c)
		}
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].CreatedAt.Before(result[j].CreatedAt)
	})
	return result
}
