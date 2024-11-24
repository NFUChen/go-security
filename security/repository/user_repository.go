package repository

import (
	"context"
	"gorm.io/gorm"
)

type IUserRepository interface {
	FindAll(ctx context.Context) ([]*User, error)
	FindByID(ctx context.Context, id uint) (*User, error)
	Save(ctx context.Context, user *User) error
	DeleteByID(ctx context.Context, id uint) error
	FindByEmail(ctx context.Context, email string) (*User, error)
	FindByUserName(ctx context.Context, name string) (*User, error)
	AddRole(ctx context.Context, role *UserRole) error
	FindAllRoles(ctx context.Context) ([]*UserRole, error)
	FindRoleByName(ctx context.Context, name string) (*UserRole, error)
	UpdateUserPassword(ctx context.Context, user *User, password string) error
	ActivateUser(ctx context.Context, user *User) error
}

type UserRepository struct {
	Engine *gorm.DB
}

func (repo *UserRepository) ActivateUser(ctx context.Context, user *User) error {
	return repo.Engine.WithContext(ctx).Model(user).Update("is_verified", true).Error
}

func (repo *UserRepository) AddRole(ctx context.Context, role *UserRole) error {
	return repo.Engine.WithContext(ctx).Create(role).Error
}

func (repo *UserRepository) FindAllRoles(ctx context.Context) ([]*UserRole, error) {
	var roles []*UserRole
	err := repo.Engine.WithContext(ctx).Find(&roles).Error
	return roles, err
}

func (repo *UserRepository) FindRoleByName(ctx context.Context, name string) (*UserRole, error) {
	var role UserRole
	err := repo.Engine.WithContext(ctx).First(&role, "name = ?", name).Error
	return &role, err
}

func (repo *UserRepository) FindByEmail(ctx context.Context, email string) (*User, error) {
	var user User
	tx := repo.Engine.WithContext(ctx).Preload("Role").First(&user, "email = ?", email)
	return &user, tx.Error
}

func (repo *UserRepository) FindByUserName(ctx context.Context, name string) (*User, error) {
	var user User
	tx := repo.Engine.WithContext(ctx).Preload("Role").First(&user, "name = ?", name)
	return &user, tx.Error
}

func NewUserRepository(engine *gorm.DB) *UserRepository {
	return &UserRepository{
		Engine: engine,
	}
}

func (repo *UserRepository) FindAll(ctx context.Context) ([]*User, error) {
	var users []*User
	err := repo.Engine.WithContext(ctx).Find(&users).Error
	return users, err
}

func (repo *UserRepository) FindByID(ctx context.Context, id uint) (*User, error) {
	var user User
	err := repo.Engine.WithContext(ctx).First(&user, id).Error
	return &user, err
}

func (repo *UserRepository) Save(ctx context.Context, user *User) error {
	return repo.Engine.WithContext(ctx).Save(user).Error
}

func (repo *UserRepository) DeleteByID(ctx context.Context, id uint) error {
	return repo.Engine.WithContext(ctx).Delete(&User{}, id).Error
}

func (repo *UserRepository) UpdateUserPassword(ctx context.Context, user *User, password string) error {
	return repo.Engine.WithContext(ctx).Model(user).Update("password", password).Error
}
