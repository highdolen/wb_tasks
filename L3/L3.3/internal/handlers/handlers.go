package handlers

import (
	"net/http"
	"strconv"

	"comment/internal/models"
	"comment/internal/service"

	"github.com/wb-go/wbf/ginext"
)

type CommentHandler struct {
	svc *service.CommentService
}

// NewCommentHandler - конструктор CommentHandler
func NewCommentHandler(s *service.CommentService) *CommentHandler {
	return &CommentHandler{svc: s}
}

// Register - регистрация роутов
func (h *CommentHandler) Register(r *ginext.Engine) {
	r.POST("/comments", h.Create)
	r.GET("/comments", h.GetComments)
	r.GET("/comments/search", h.Search)
	r.DELETE("/comments/:id", h.Delete)
}

type createReq struct {
	Text     string `json:"text"`
	ParentID *int64 `json:"parent_id"`
}

// Create - handler для созданиея комментария
func (h *CommentHandler) Create(c *ginext.Context) {
	var req createReq
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ginext.H{"error": err.Error()})
		return
	}

	comment := h.svc.Create(req.Text, req.ParentID)
	if comment == nil {
		c.JSON(http.StatusBadRequest, ginext.H{"error": "invalid parent"})
		return
	}

	c.JSON(http.StatusOK, comment)
}

// GetComments - handler для получения комментария
func (h *CommentHandler) GetComments(c *ginext.Context) {
	parentStr := c.Query("parent")
	limitStr := c.DefaultQuery("limit", "5")
	pageStr := c.DefaultQuery("page", "1")

	limit, _ := strconv.Atoi(limitStr)
	page, _ := strconv.Atoi(pageStr)
	offset := (page - 1) * limit

	if parentStr == "" {
		roots := h.svc.GetRoots(limit, offset)
		var result []*models.Comment
		for _, r := range roots {
			r.Children = h.svc.GetChildren(r.ID)
			result = append(result, r)
		}
		c.JSON(http.StatusOK, result)
		return
	}

	// поддерево конкретного комментария
	id, err := strconv.ParseInt(parentStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ginext.H{"error": "invalid parent id"})
		return
	}

	children := h.svc.GetChildren(id)
	c.JSON(http.StatusOK, children)
}

// Search - handler для поиска комментария
func (h *CommentHandler) Search(c *ginext.Context) {
	q := c.Query("q")
	if q == "" {
		c.JSON(http.StatusBadRequest, ginext.H{"error": "empty query"})
		return
	}
	result := h.svc.Search(q)
	c.JSON(http.StatusOK, result)
}

// Delete - удаление комментария
func (h *CommentHandler) Delete(c *ginext.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ginext.H{"error": "invalid id"})
		return
	}

	h.svc.Delete(id)
	c.JSON(http.StatusOK, ginext.H{"deleted": id})
}
