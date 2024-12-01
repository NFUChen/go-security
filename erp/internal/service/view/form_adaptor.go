package view

import (
	"fmt"
	. "go-security/erp/internal/repository"
)

type FormAdaptor struct{}

func NewFormAdaptor() *FormAdaptor {
	return &FormAdaptor{}
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
