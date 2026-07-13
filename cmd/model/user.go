package model

import (
	"time"

	"github.com/PritomKarmokar/zipurl/cmd/utils"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Role string
type UserStatus string

const (
	UserRole  Role = "user"
	AdminRole Role = "admin"

	StatusActive   UserStatus = "active"
	StatusInActive UserStatus = "inactive"
	StatusBlocked  UserStatus = "blocked"
)

type User struct {
	ID        string     `gorm:"primaryKey;size:26" json:"id"`
	UserName  string     `gorm:"column:username;type:varchar(50);not null" json:"username"`
	FirstName string     `gorm:"column:first_name;type:varchar(25);not null" json:"first_name"`
	LastName  string     `gorm:"column:last_name;type:varchar(25);not null" json:"last_name"`
	Email     string     `gorm:"column:email;type:varchar(100);not null" json:"email"`
	Password  string     `gorm:"column:password;type:text" json:"-"`
	Status    UserStatus `gorm:"column:status;type:varchar(10);not null;default:active" json:"status"`
	Role      Role       `gorm:"column:role;type:varchar(10);not null" json:"role"`

	CreatedAt  time.Time  `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time  `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
	DeletedAt  *time.Time `gorm:"column:deleted_at" json:"deleted_at,omitempty"`
	LastLogin  *time.Time `gorm:"column:last_login" json:"last_login,omitempty"`
	DateJoined time.Time  `gorm:"column:date_joined,autoUpdateTime" json:"date_joined"`
}

func (u *User) TableName() string {
	return "users"
}

func (u *User) SetPassword(password string) error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashed)
	return nil
}

func (u *User) CheckPassword(providedPassword string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(providedPassword)) == nil
}

func (u *User) UpdateLastLoginTime(db *gorm.DB) error {
	now := time.Now()
	result := db.Model(&User{}).
		Where("id = ? AND status = ?", u.ID, StatusActive).
		Update("last_login", now)

	if result.Error != nil {
		return result.Error
	}
	u.LastLogin = &now
	return nil
}

func (u *User) BeforeCreate(db *gorm.DB) error {
	if u.ID == "" {
		u.ID = utils.GenerateULID()
	}
	return nil
}
