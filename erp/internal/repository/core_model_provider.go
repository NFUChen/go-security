package repository

type CoreModelProvider struct{}

func (provider *CoreModelProvider) ProvideModels() []any {
	return []any{
		&CustomerOrder{},
		&UserProfile{},
		&OrderNotification{},
		&Product{},
		&PricingPolicy{},
		&PolicyPrice{},
	}
}

func NewCoreModelProvider() *CoreModelProvider {
	return &CoreModelProvider{}
}
