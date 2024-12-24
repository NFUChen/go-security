package service

import (
	"context"
	"github.com/rs/zerolog/log"
	"go-security/erp/internal"
	. "go-security/erp/internal/repository"
)

type PricingPolicyService struct {
	DefaultPolicyName       string
	PricingPolicyRepository IPricingPolicyRepository
}

func NewPricingPolicyService(repo IPricingPolicyRepository) *PricingPolicyService {
	return &PricingPolicyService{
		DefaultPolicyName:       "Default",
		PricingPolicyRepository: repo,
	}
}

func (service *PricingPolicyService) PostConstruct() {
	policy, err := service.CreateDefaultPricingPolicy(context.Background())
	if err != nil {
		log.Fatal().Msgf("Failed to create default pricing policy: %v", err)
	}
	log.Info().Msgf("Created default pricing policy with ID: %d", policy.ID)
}

func (service *PricingPolicyService) GetPolicyByID(ctx context.Context, ID uint) (*PricingPolicy, error) {
	return service.PricingPolicyRepository.FindPolicyByID(ctx, ID)
}

func (service *PricingPolicyService) GetDefaultPolicy(ctx context.Context) (*PricingPolicy, error) {
	return service.PricingPolicyRepository.FindPolicyByName(ctx, service.DefaultPolicyName)
}

func (service *PricingPolicyService) CreateDefaultPricingPolicy(ctx context.Context) (*PricingPolicy, error) {
	if policy, err := service.PricingPolicyRepository.FindPolicyByName(ctx, service.DefaultPolicyName); err == nil {
		log.Info().Msgf("Default pricing policy already exists, skipping creation")
		return policy, nil
	}

	policy := &PricingPolicy{
		Name:         service.DefaultPolicyName,
		Description:  "Default pricing policy",
		PolicyPrices: []PolicyPrice{},
	}
	err := service.PricingPolicyRepository.AddPolicy(ctx, policy)
	if err != nil {
		return nil, err
	}
	return policy, nil
}

func (service *PricingPolicyService) GetAllPolicies(ctx context.Context) ([]*PricingPolicy, error) {
	return service.PricingPolicyRepository.FindAllPolicies(ctx)
}

func (service *PricingPolicyService) CreateNewPricingPolicy(ctx context.Context, name string, description string) error {
	policy := &PricingPolicy{
		Name:        name,
		Description: description,
	}

	policyFound, err := service.PricingPolicyRepository.FindPolicyByName(ctx, name)
	if policyFound != nil && err == nil {
		return internal.PricingPolicyAlreadyExists
	}

	return service.PricingPolicyRepository.AddPolicy(ctx, policy)
}

func (service *PricingPolicyService) AddPolicyPriceToPolicy(ctx context.Context, policyId uint, policyPrice *PolicyPrice) error {
	policy, err := service.PricingPolicyRepository.FindPolicyByID(ctx, policyId)
	if err != nil {
		return internal.PricingPolicyNotFound
	}
	return service.PricingPolicyRepository.AddPolicyPrice(ctx, policy.ID, policyPrice)
}
