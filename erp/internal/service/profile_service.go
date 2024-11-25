package service

import . "go-security/erp/internal/repository"

type ProfileService struct {
	profileRepository IProfileRepository
}

func NewProfileService(profileRepository IProfileRepository) *ProfileService {
	return &ProfileService{
		profileRepository: profileRepository,
	}
}

func (service *ProfileService) FindProfileByCustomerId(customerId uint) (*CustomerProfile, error) {
	return service.profileRepository.FindProfileByCustomerId(customerId)
}
