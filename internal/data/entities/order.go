package entities

import (
	"server/internal/core/shared/utils/cryptoutil"
	"server/internal/core/shared/utils/stringutil"
	"time"

	"gorm.io/gorm"
)

type OrderStatusEnum string

type OrderPriorityEnum string

const (
	OrderStatusNew       OrderStatusEnum = "new"
	OrderStatusPicking   OrderStatusEnum = "picking"
	OrderStatusPacked    OrderStatusEnum = "packed"
	OrderStatusShipped   OrderStatusEnum = "shipped"
	OrderStatusBackorder OrderStatusEnum = "backorder"
	OrderStatusCancelled OrderStatusEnum = "cancelled"
)

const (
	OrderPriorityLow    OrderPriorityEnum = "low"
	OrderPriorityNormal OrderPriorityEnum = "normal"
	OrderPriorityHigh   OrderPriorityEnum = "high"
	OrderPriorityUrgent OrderPriorityEnum = "urgent"
)

type Order struct {
	ID           string            `gorm:"type:varchar(150);primaryKey;unique" json:"id"`
	LocationID   string            `gorm:"type:varchar(150);not null;index;uniqueIndex:idx_orders_location_number" json:"location_id"`
	OrderNumber  string            `gorm:"type:varchar(100);not null;uniqueIndex:idx_orders_location_number" json:"order_number"`
	CustomerName string            `gorm:"type:varchar(200);not null" json:"customer_name"`
	Channel      string            `gorm:"type:varchar(120);not null;default:'manual'" json:"channel"`
	Status       OrderStatusEnum   `gorm:"type:varchar(25);not null;default:'new';check(status IN ('new','picking','packed','shipped','backorder','cancelled'))" json:"status"`
	Priority     OrderPriorityEnum `gorm:"type:varchar(25);not null;default:'normal';check(priority IN ('low','normal','high','urgent'))" json:"priority"`
	ItemsCount   uint64            `gorm:"not null;default:0" json:"items_count"`
	TotalAmount  float64           `gorm:"type:double;not null;default:0" json:"total_amount"`
	DueAt        *time.Time        `gorm:"null" json:"due_at,omitempty"`
	Notes        *string           `gorm:"type:text;null" json:"notes,omitempty"`
	CreatedAt    time.Time         `json:"created_at"`
	UpdatedAt    time.Time         `json:"updated_at"`
	Location     *Location         `gorm:"foreignKey:LocationID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"location,omitempty"`
}

func (e *Order) TableName() string {
	return "orders"
}

func (e *Order) BeforeCreate(tx *gorm.DB) error {
	if stringutil.IsEmpty(e.ID) {
		e.ID = cryptoutil.UID()
	}
	return nil
}
