package models

import "time"

type OrderDBModel struct {
	Id             int64 `gorm:"id;primaryKey;autoIncrement"`
	UserId         int64 `gorm:"column:user_id;type:int;not null;uniqueIndex;not null"`
	SubscriptionId int64 `gorm:"column:plan_id;type:varchar(255);not null"`
	VoucherId      int64
	PurchasedAt    time.Time  `gorm:"column:purchased_at"`
	CreatedAt      time.Time  `gorm:"column:created_at"`
	DeletedAt      *time.Time `gorm:"column:deleted_at"`
}

func (OrderDBModel) TableName() string {
	return "orders"
}
