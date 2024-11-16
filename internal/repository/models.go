package repository

import (
	"go-security/internal"
	"gorm.io/gorm"
	"net/mail"
	"time"
)

const (
	RoleSuperAdmin  string = "super_admin"
	RoleAdmin       string = "admin"
	RoleGuest       string = "guest"
	RoleBlockedUser string = "blocked_user"
)

type RoleIndex uint

type UserRole struct {
	ID        uint   `gorm:"primaryKey" json:"id"` // Auto-increment primary key
	Name      string `gorm:"type:varchar(50);not null" json:"name"`
	RoleIndex uint   `gorm:"type:int;unique;not null" json:"role_index"`
	Users     []User `gorm:"foreignKey:RoleID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"users"`
}

var BuiltinRoles = []UserRole{
	{Name: RoleSuperAdmin, RoleIndex: 1000},
	{Name: RoleAdmin, RoleIndex: 500},
	{Name: RoleGuest, RoleIndex: 1},
	{Name: RoleBlockedUser, RoleIndex: 0},
}

type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"` // Auto-increment primary key
	Name      string         `gorm:"type:varchar(100);not null" json:"name"`
	Email     string         `gorm:"type:varchar(100);unique;not null" json:"email"`
	Password  string         `gorm:"type:varchar(255);not null" json:"-"` // Excluded from JSON responses
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	RoleID    uint           `gorm:"not null" json:"role_id"` // Foreign key
	Role      UserRole       `gorm:"foreignKey:RoleID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"-"`
}

func (user *User) Validate() error {
	if len(user.Name) == 0 {
		return internal.UserNameNotAllowed
	}
	_, err := mail.ParseAddress(user.Email)
	if len(user.Email) == 0 || err != nil {
		return internal.UserEmailNotAllowed
	}
	if len(user.Password) == 0 {
		return internal.UserPasswordNotAllowed
	}
	return nil
}

func GetAllModels() []any {
	return []any{
		&User{},
		&UserRole{},
	}
}
