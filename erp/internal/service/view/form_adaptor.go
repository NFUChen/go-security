package view

import (
	"context"
	"fmt"
	. "go-security/erp/internal/repository"
	"go-security/erp/internal/service"
)

type FormAdaptor struct {
	ProductService *service.ProductService
}

func NewFormAdaptor(productService *service.ProductService) *FormAdaptor {
	return &FormAdaptor{
		ProductService: productService,
	}
}

func (service *FormAdaptor) FormToProduct(ctx context.Context, form *Form) (*Product, error) {
	if form == nil {
		return nil, fmt.Errorf("form is nil")
	}
	product := &Product{}

	for _, field := range form.Fields {
		switch field.Key {
		case "name":
			if value, ok := field.Value.(string); ok {
				product.Name = value
			} else {
				return nil, fmt.Errorf("invalid type for name")
			}
		case "category_id":
			label, ok := field.Value.(string)
			if !ok {
				return nil, fmt.Errorf("invalid type for category_id")
			}
			category, err := service.ProductService.FindCategoryByName(ctx, label)
			if err != nil {
				return nil, err
			}
			product.CategoryID = category.ID

		case "description":
			if value, ok := field.Value.(string); ok {
				product.Description = value
			} else {
				return nil, fmt.Errorf("invalid type for description")
			}
		}
	}
	return product, nil
}

func (service *FormAdaptor) FormToProductCategory(form *Form) (*ProductCategory, error) {
	if form == nil {
		return nil, fmt.Errorf("form is nil")
	}
	category := &ProductCategory{}

	for _, field := range form.Fields {
		switch field.Key {
		case "name":
			if value, ok := field.Value.(string); ok {
				category.Name = value
			} else {
				return nil, fmt.Errorf("invalid type for name")
			}
		case "description":
			if value, ok := field.Value.(string); ok {
				category.Description = value
			} else {
				return nil, fmt.Errorf("invalid type for description")
			}
		}
	}
	return category, nil
}

func (service *FormAdaptor) FormToUserProfile(userID uint, form *Form) (*UserProfile, error) {
	if form == nil {
		return nil, fmt.Errorf("form is nil")
	}
	userProfile := &UserProfile{}

	for _, field := range form.Fields {
		switch field.Key {
		case "full_name":
			// Assuming this maps to the User model rather than UserProfile directly.
			// No direct mapping here for UserProfile.

		case "phone_number":
			if value, ok := field.Value.(string); ok {
				userProfile.PhoneNumber = value
			} else {
				return nil, fmt.Errorf("invalid type for phone_number")
			}

		case "address":
			if value, ok := field.Value.(string); ok {
				userProfile.Address = value
			} else {
				return nil, fmt.Errorf("invalid type for address")
			}

		case "notification_approaches":
			approaches := []NotificationApproach{}
			for _, opt := range field.Options {
				approach := NotificationApproach{
					UserID:  userID,
					Name:    NotificationType(opt.Label),
					Enabled: opt.IsChecked,
				}
				approaches = append(approaches, approach)
			}
			userProfile.NotificationApproaches = approaches
		case "user_description":
			if value, ok := field.Value.(string); ok {
				userProfile.UserDescription = value
			} else {
				return nil, fmt.Errorf("invalid type for user_description")
			}

		}
	}

	return userProfile, nil
}
