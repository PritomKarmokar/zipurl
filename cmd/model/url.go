package model

import "time"

type URL struct {
	ID          uint64    `gorm:"primaryKey;size:26"`
	URL         string    `gorm:"type:text;not null"`
	HashedToken string    `gorm:"type:text;not null"`
	CreatedAt   time.Time `gorm:"type:timestamptz;not null;default:now()"`
	UpdatedAt   time.Time `gorm:"type:timestamptz;not null;default:now()"`
}

func (u *URL) TableName() string {
	return "urls"
}
