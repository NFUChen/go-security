package repository

type CoreModelProvider struct{}

func (provider *CoreModelProvider) ProvideModels() []any {
	return []any{
		&Order{},
		&OrderItem{},
		&UserProfile{},
		&OrderNotification{},
		&Product{},
		&ProductCategory{},
		&PricingPolicy{},
		&PolicyPrice{},
		&NotificationApproach{},
	}
}

func NewCoreModelProvider() *CoreModelProvider {
	return &CoreModelProvider{}
}
