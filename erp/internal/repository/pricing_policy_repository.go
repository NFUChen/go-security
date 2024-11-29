package repository

import (
	"context"
	"gorm.io/gorm"
)

type IPricingPolicyRepository interface {
	FindPolicyByName(ctx context.Context, name string) (*PricingPolicy, error)
	FindPolicyByID(ctx context.Context, ID uint) (*PricingPolicy, error)
	AddPolicy(ctx context.Context, policy *PricingPolicy) error
	FindAllPolicies(ctx context.Context) ([]*PricingPolicy, error)
}

type PricingPolicyRepository struct {
	Engine *gorm.DB
}

func (repo PricingPolicyRepository) FindPolicyByName(ctx context.Context, name string) (*PricingPolicy, error) {
	policy := PricingPolicy{}
	tx := repo.createPreloadQuery(ctx).Where("name = ?", name).First(&policy)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &policy, nil
}

func NewPricingPolicyRepository(engine *gorm.DB) *PricingPolicyRepository {
	return &PricingPolicyRepository{
		Engine: engine,
	}
}
func (repo PricingPolicyRepository) FindPolicyByID(ctx context.Context, ID uint) (*PricingPolicy, error) {
	policy := PricingPolicy{}
	tx := repo.createPreloadQuery(ctx).Where("id = ?", ID).First(&policy)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &policy, nil
}

func (repo PricingPolicyRepository) AddPolicy(ctx context.Context, policy *PricingPolicy) error {
	tx := repo.Engine.WithContext(ctx).Create(policy)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

func (repo PricingPolicyRepository) FindAllPolicies(ctx context.Context) ([]*PricingPolicy, error) {
	policies := []*PricingPolicy{}
	tx := repo.createPreloadQuery(ctx).Find(&policies)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return policies, nil
}

func (repo PricingPolicyRepository) createPreloadQuery(ctx context.Context) *gorm.DB {
	return repo.Engine.WithContext(ctx)
}
