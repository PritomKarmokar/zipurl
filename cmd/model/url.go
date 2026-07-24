package model

import "time"

const DefaultClickCount int64 = 10e6

type URL struct {
	ID          string  `gorm:"primaryKey;size:26"`
	UserID      *string `gorm:"size:26;index"`
	User        *User   `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE"`
	URL         string  `gorm:"type:text;not null"`
	HashedToken string  `gorm:"type:text;not null"`

	ExpiresAt  *time.Time `gorm:"type:timestamptz"`
	ClickCount int64      `gorm:"not null;default:0"`
	MaxClicks  *int64     `gorm:"type:bigint"`
	CreatedAt  time.Time  `gorm:"type:timestamptz;not null;default:now()"`
	UpdatedAt  time.Time  `gorm:"type:timestamptz;not null;default:now()"`
}

func (u *URL) TableName() string {
	return "urls"
}
