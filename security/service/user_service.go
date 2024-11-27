package service

import (
	"context"
	"go-security/security"
	. "go-security/security/repository"
)

const (
	RoleSuperAdmin  string = "super_admin"
	RoleAdmin       string = "admin"
	RoleGuest       string = "guest"
	RoleBlockedUser string = "blocked_user"
)

var BuiltinRoles = []UserRole{
	{Name: RoleSuperAdmin, RoleIndex: 1000},
	{Name: RoleAdmin, RoleIndex: 500},
	{Name: RoleGuest, RoleIndex: 1},
	{Name: RoleBlockedUser, RoleIndex: 0},
}

type PlatformType string

const (
	PlatformSelf   PlatformType = "Self"
	PlatformGoogle PlatformType = "Google"
	PlatformLine   PlatformType = "LINE"
)

var BuiltinPlatforms = []Platform{
	{Name: string(PlatformSelf)},
	{Name: string(PlatformGoogle)},
	{Name: string(PlatformLine)},
}

type UserService struct {
	UserRepository IUserRepository
}

func (service *UserService) PostConstruct() {
	service.addBuiltinRoles()
	service.addBuiltinPlatforms()
}

func (service *UserService) GetUserPlatform(ctx context.Context, userID uint) (*Platform, error) {
	user, err := service.UserRepository.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return service.UserRepository.FindPlatformByID(ctx, user.PlatformID)
}

func (service *UserService) addBuiltinPlatforms() {
	for _, platform := range BuiltinPlatforms {
		_ = service.UserRepository.AddPlatform(context.Background(), &platform)
	}
}

func (service *UserService) addBuiltinRoles() {
	for _, role := range BuiltinRoles {
		_ = service.AddRole(context.Background(), &role)
	}
}

func NewUserService(userRepository IUserRepository) *UserService {
	return &UserService{
		UserRepository: userRepository,
	}
}

func (service *UserService) AddPlatform(ctx context.Context, platform *Platform) error {
	return service.UserRepository.AddPlatform(ctx, platform)
}

func (service *UserService) AddRole(ctx context.Context, role *UserRole) error {
	return service.UserRepository.AddRole(ctx, role)
}

func (service *UserService) FindAllRoles(ctx context.Context) ([]*UserRole, error) {
	return service.UserRepository.FindAllRoles(ctx)
}

func (service *UserService) FindRoleByName(ctx context.Context, name string) (*UserRole, error) {
	if len(name) == 0 {
		return nil, security.UserRoleNotAllowed
	}
	return service.UserRepository.FindRoleByName(ctx, name)
}
func (service *UserService) FindPlatformByName(ctx context.Context, name PlatformType) (*Platform, error) {
	if len(name) == 0 {
		return nil, security.UserPlatformEmpty
	}
	return service.UserRepository.FindPlatformByName(ctx, string(name))
}

func (service *UserService) FindAllUsers(ctx context.Context) ([]*User, error) {
	return service.UserRepository.FindAll(ctx)
}

func (service *UserService) FindUserByID(ctx context.Context, id uint) (*User, error) {
	user, err := service.UserRepository.FindByID(ctx, id)
	if err != nil {
		return nil, security.UserNotFound
	}
	return user, nil
}

func (service *UserService) FindUserByEmail(ctx context.Context, email string) (*User, error) {
	user, err := service.UserRepository.FindByEmail(ctx, email)
	if err != nil {
		return nil, security.UserNotFound
	}
	return user, nil
}
func (service *UserService) FindUserByUserName(ctx context.Context, name string) (*User, error) {
	user, err := service.UserRepository.FindByUserName(ctx, name)
	if err != nil {
		return nil, security.UserNotFound
	}
	return user, nil
}

func (service *UserService) SaveUser(ctx context.Context, user *User) (*User, error) {
	err := service.UserRepository.Save(ctx, user)
	return user, err
}

func (service *UserService) ResetUserPassword(ctx context.Context, user *User, password string) error {
	return service.UserRepository.UpdateUserPassword(ctx, user, password)
}

func (service *UserService) ActivateUser(ctx context.Context, user *User) error {
	if user.IsVerified {
		return nil
	}
	return service.UserRepository.ActivateUser(ctx, user)
}
