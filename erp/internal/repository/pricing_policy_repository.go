package repository

import "gorm.io/gorm"

type IPricingPolicyRepository interface {
	FindPolicyByID(ID uint) (*PricingPolicy, error)
}

type PricingPolicyRepository struct {
	Engine *gorm.DB
}

func (repo PricingPolicyRepository) FindPolicyByID(ID uint) (*PricingPolicy, error) {
	policy := PricingPolicy{}
	tx := repo.Engine.Where("id = ?", ID).First(&policy)
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
