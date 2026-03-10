package audit

import "time"

// History - структура записи истории изменений товара
type History struct {
	ID        int       `json:"id"`
	ItemID    int       `json:"item_id"`
	Action    string    `json:"action"`
	OldValue  string    `json:"old_value"`
	NewValue  string    `json:"new_value"`
	ChangedBy string    `json:"changed_by"`
	ChangedAt time.Time `json:"changed_at"`
}
