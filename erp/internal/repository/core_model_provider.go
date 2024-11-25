package repository

type CoreModelProvider struct{}

func (provider *CoreModelProvider) ProvideModels() []any {
	return []any{
		&CustomerOrder{},
		&CustomerProfile{},
		&OrderNotification{},
		&Product{},
	}
}

func NewCoreModelProvider() *CoreModelProvider {
	return &CoreModelProvider{}
}
