package repository

type SecurityModelProvider struct{}

func NewSecurityModelProvider() *SecurityModelProvider {
	return &SecurityModelProvider{}
}

func (provider *SecurityModelProvider) ProvideModels() []any {
	return []any{
		&User{},
	}
}
