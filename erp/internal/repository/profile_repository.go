package repository

type IProfileRepository interface {
	FindProfileByCustomerId(customerId uint) (*CustomerProfile, error)
}

type ProfileRepository struct{}

func (repo ProfileRepository) FindProfileByCustomerId(customerId uint) (*CustomerProfile, error) {
	//TODO implement me
	panic("implement me")
}

func NewProfileRepository() *ProfileRepository {
	return &ProfileRepository{}
}
