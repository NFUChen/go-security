package service

import . "go-security/erp/internal/repository"

type PricingPolicyService struct {
	PricingPolicyRepository IPricingPolicyRepository
}

func (service *PricingPolicyService) GetPolicyByID(ID uint) (*PricingPolicy, error) {
	return service.PricingPolicyRepository.FindPolicyByID(ID)
}

func NewPricingPolicyService(repo IPricingPolicyRepository) *PricingPolicyService {
	return &PricingPolicyService{
		PricingPolicyRepository: repo,
	}
}
