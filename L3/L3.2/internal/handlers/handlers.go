package handlers

import (
	"net/http"

	"shortener/internal/repository"

	"github.com/wb-go/wbf/ginext"
)

type Handler struct {
	repo *repository.Repository
}

func NewRepository(repo *repository.Repository) *Handler {
	return &Handler{repo: repo}
}

type shortenRequest struct {
	URL        string  `json:"url" binding:"required"`
	CustomCode *string `json:"custom_code"`
}

type shortenResponse struct {
	Code string `json:"code"`
}

// POST /shorten
func (h *Handler) CreateShortLink(c *ginext.Context) {
	var req shortenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ginext.H{
			"error": "invalid request body",
		})
		return
	}

	code, err := h.repo.CreateShortLink(
		c.Request.Context(),
		req.URL,
		req.CustomCode,
	)
	if err != nil {
		if err == repository.ErrCodeAlreadyExists {
			c.JSON(http.StatusConflict, ginext.H{
				"error": "short code already exists",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, ginext.H{
			"error": "internal error",
		})
		return
	}

	c.JSON(http.StatusOK, shortenResponse{Code: code})
}

// GET /s/:code
func (h *Handler) Redirect(c *ginext.Context) {
	code := c.Param("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, ginext.H{
			"error": "code is required",
		})
		return
	}

	originalURL, err := h.repo.GetOriginalURL(
		c.Request.Context(),
		code,
	)
	if err != nil {
		if err == repository.ErrNotFound {
			c.JSON(http.StatusNotFound, ginext.H{
				"error": "link not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, ginext.H{
			"error": "internal error",
		})
		return
	}

	// аналитика (ошибка не ломает редирект)
	userAgent := c.GetHeader("User-Agent")
	_ = h.repo.SaveVisit(c.Request.Context(), code, userAgent)

	c.Redirect(http.StatusFound, originalURL)
}

// GET /analytics/:code
func (h *Handler) Analytics(c *ginext.Context) {
	code := c.Param("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, ginext.H{
			"error": "code is required",
		})
		return
	}

	stats, err := h.repo.GetAnalytics(
		c.Request.Context(),
		code,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ginext.H{
			"error": "internal error",
		})
		return
	}

	c.JSON(http.StatusOK, stats)
}
