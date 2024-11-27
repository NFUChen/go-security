package repository

import (
	"fmt"
	"go-security/security"
	"net/mail"
	"time"
)

type PostgresDataSourceConfig struct {
	Host         string `yaml:"host"`
	Port         int    `yaml:"port"`
	User         string `yaml:"user"`
	Password     string `yaml:"password"`
	DatabaseName string `yaml:"db_name"`
}

func (config *PostgresDataSourceConfig) AsDSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.Host, config.Port, config.User, config.Password, config.DatabaseName,
	)
}

type RoleIndex uint

type UserRole struct {
	Name      string `gorm:"type:varchar(50);not null" json:"name"`
	RoleIndex uint   `gorm:"type:int;unique;not null" json:"role_index"`
	Users     []User `gorm:"foreignKey:RoleID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"users"`

	ID        uint       `gorm:"primaryKey" json:"id"` // Auto-increment primary key
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

type User struct {
	Name       string   `gorm:"type:varchar(100);not null" json:"name"`
	Email      string   `gorm:"type:varchar(100);unique;not null" json:"email"`
	Password   string   `gorm:"type:varchar(255);not null" json:"-"` // Excluded from JSON responses
	RoleID     uint     `gorm:"not null" json:"role_id"`             // Foreign key
	IsVerified bool     `gorm:"default:false" json:"is_verified"`
	Role       UserRole `gorm:"foreignKey:RoleID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"role"`
	PlatformID uint     `gorm:"not null" json:"platform_id"` // Foreign key
	Platform   Platform `gorm:"foreignKey:PlatformID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"platform"`
	ExternalID *string  `gorm:"type:varchar(100);unique" json:"external_id"`

	ID        uint       `gorm:"primaryKey" json:"id"` // Auto-increment primary key
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

type Platform struct {
	Name string `gorm:"type:varchar(100);not null;unique" json:"name"`

	ID        uint       `gorm:"primaryKey" json:"id"` // Auto-increment primary key
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

func (user *User) Validate() error {
	if len(user.Name) == 0 {
		return security.UserNameNotAllowed
	}
	_, err := mail.ParseAddress(user.Email)
	if len(user.Email) == 0 || err != nil {
		return security.UserEmailNotAllowed
	}
	if len(user.Password) == 0 {
		return security.UserPasswordNotAllowed
	}
	return nil
}
