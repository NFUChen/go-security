package repository

import (
	"context"
	"go-security/erp/internal"
	"gorm.io/gorm"
)

type IPricingPolicyRepository interface {
	FindPolicyByName(ctx context.Context, name string) (*PricingPolicy, error)
	FindPolicyByID(ctx context.Context, ID uint) (*PricingPolicy, error)
	AddPolicy(ctx context.Context, policy *PricingPolicy) error
	AddPolicyPrice(ctx context.Context, policyID uint, policyPrice *PolicyPrice) error
	DeletePricingPolicy(ctx context.Context, policyID uint) error
	DeletePolicyPrice(ctx context.Context, policyPriceID uint) error
	FindAllPolicies(ctx context.Context) ([]*PricingPolicy, error)
}

type PricingPolicyRepository struct {
	Engine *gorm.DB
}

func (repo *PricingPolicyRepository) DeletePricingPolicy(ctx context.Context, policyID uint) error {
	policy := PricingPolicy{}
	if tx := repo.Engine.WithContext(ctx).First(&policy, policyID); tx.Error != nil {
		return tx.Error
	}

	if tx := repo.Engine.WithContext(ctx).Delete(&policy); tx.Error != nil {
		return tx.Error
	}
	return nil
}

func NewPricingPolicyRepository(engine *gorm.DB) *PricingPolicyRepository {
	return &PricingPolicyRepository{
		Engine: engine,
	}
}

func (repo *PricingPolicyRepository) AddPolicyPrice(ctx context.Context, policyID uint, policyPrice *PolicyPrice) error {
	policy, err := repo.FindPolicyByID(ctx, policyID)
	if err != nil {
		return internal.PricingPolicyNotFound
	}
	policy.PolicyPrices = append(policy.PolicyPrices, *policyPrice)
	tx := repo.Engine.WithContext(ctx).Save(policy)
	return tx.Error
}

func (repo *PricingPolicyRepository) DeletePolicyPrice(ctx context.Context, policyPriceID uint) error {
	policyPrice := PolicyPrice{}
	if tx := repo.Engine.WithContext(ctx).First(&policyPrice, policyPriceID); tx.Error != nil {
		return tx.Error
	}

	if tx := repo.Engine.WithContext(ctx).Delete(&policyPrice); tx.Error != nil {
		return tx.Error
	}
	return nil
}

func (repo *PricingPolicyRepository) FindPolicyByName(ctx context.Context, name string) (*PricingPolicy, error) {
	policy := PricingPolicy{}
	tx := repo.createPreloadQuery(ctx).Where("name = ?", name).First(&policy)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &policy, nil
}

func (repo *PricingPolicyRepository) FindPolicyByID(ctx context.Context, ID uint) (*PricingPolicy, error) {
	policy := PricingPolicy{}
	tx := repo.createPreloadQuery(ctx).Where("id = ?", ID).First(&policy)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &policy, nil
}

func (repo *PricingPolicyRepository) AddPolicy(ctx context.Context, policy *PricingPolicy) error {
	tx := repo.Engine.WithContext(ctx).Create(policy)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

func (repo *PricingPolicyRepository) FindAllPolicies(ctx context.Context) ([]*PricingPolicy, error) {
	policies := []*PricingPolicy{}
	tx := repo.createPreloadQuery(ctx).Find(&policies)
	if tx.Error != nil {
		return nil, tx.Error
	}

	for _, policy := range policies {
		defaultPolicyPrices := []PolicyPrice{}
		if policy.PolicyPrices == nil {
			policy.PolicyPrices = defaultPolicyPrices
		}
	}
	return policies, nil
}

func (repo *PricingPolicyRepository) createPreloadQuery(ctx context.Context) *gorm.DB {
	return repo.Engine.WithContext(ctx)
}
