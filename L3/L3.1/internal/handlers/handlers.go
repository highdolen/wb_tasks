package handlers

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"delayed_notifier/internal/models"
	"delayed_notifier/internal/rabbitmq"
	"delayed_notifier/internal/storage"
)

type NotificationHandlers struct {
	store    storage.NotificationStorage
	producer *rabbitmq.Producer
}

func NewNotificationHandlers(store storage.NotificationStorage, producer *rabbitmq.Producer) *NotificationHandlers {
	return &NotificationHandlers{
		store:    store,
		producer: producer,
	}
}

// POST /notify — создаёт уведомление
func (h *NotificationHandlers) CreateNotification(c *gin.Context) {
	var req models.Notification

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	req.ID = uuid.New().String()

	req.Status = models.StatusPending

	req.Meta = models.NotificationMeta{
		Attempt:    0,
		MaxAttempt: 5,
		CreatedAt:  time.Now(),
	}

	if req.SendAt.IsZero() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "send_at is required"})
		return
	}

	if err := h.store.Save(req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to save notification",
		})
		return
	}

	if err := h.producer.PublishNotification(c.Request.Context(), req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to enqueue notification"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": req.ID})
}

// GET /notify/:id — получение статуса
func (h *NotificationHandlers) GetNotification(c *gin.Context) {
	id := c.Param("id")

	n, ok := h.store.Get(id)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "notification not found"})
		return
	}

	c.JSON(http.StatusOK, n)
}

// DELETE /notify/:id — отмена уведомления
func (h *NotificationHandlers) DeleteNotification(c *gin.Context) {
	id := c.Param("id")

	log.Printf("handlers: try cancel id=%s", id)

	n, ok := h.store.Get(id)
	if !ok {
		log.Printf("handlers: delete id=%s not found in storage", id)
		c.JSON(http.StatusNotFound, gin.H{"error": "notification not found"})
		return
	}

	log.Printf(
		"handlers: delete id=%s current_status=%s",
		n.ID,
		n.Status,
	)

	n.Status = models.StatusCancelled

	if err := h.store.Save(n); err != nil {
		log.Printf("handlers: save failed id=%s err=%v", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to cancel"})
		return
	}

	log.Printf("handlers: id=%s status set to CANCELLED", id)

	c.JSON(http.StatusOK, gin.H{"status": "cancelled"})
}
