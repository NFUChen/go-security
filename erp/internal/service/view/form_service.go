package view

import (
	"context"
	. "go-security/erp/internal/repository"
	"go-security/erp/internal/service"
	baseApp "go-security/security/service"
	"slices"
)

type FormService struct {
	UserService                 *baseApp.UserService
	ProfileService              *service.ProfileService
	NotificationApproachService *service.NotificationApproachService
}

func NewFormService(
	profileService *service.ProfileService,
	userService *baseApp.UserService,
	notificationApproachService *service.NotificationApproachService,
) *FormService {
	return &FormService{
		ProfileService:              profileService,
		UserService:                 userService,
		NotificationApproachService: notificationApproachService,
	}
}

type Form struct {
	Fields []*FormField `json:"fields"`
}

func (form *Form) AsMap(excludeReadOnly bool) map[string]any {
	formMap := make(map[string]any)
	for _, field := range form.Fields {
		if excludeReadOnly && field.ReadOnly {
			continue
		}

		formMap[field.Key] = field
	}
	return formMap
}

type FieldType string

const (
	FieldTypeText     FieldType = "text"
	FieldTypeTextArea FieldType = "textarea"
	FieldTypeFile     FieldType = "file"
	FieldTypeCheckbox FieldType = "checkbox"
	FieldTypeCombobox FieldType = "combobox"
)

var OptionTypes = []FieldType{
	FieldTypeCheckbox,
	FieldTypeCombobox,
}

type FieldOption struct {
	Label           string `json:"label"`
	IsDisabled      bool   `json:"is_disabled"`
	DisabledMessage string `json:"disabled_message"`
	IsChecked       bool   `json:"checked"`
}

type FormField struct {
	Key      string         `json:"key"`
	Label    string         `json:"label"`
	Type     FieldType      `json:"type"`
	ReadOnly bool           `json:"read_only"`
	Required bool           `json:"required"`
	Options  []*FieldOption `json:"options,omitempty"`
	Value    any            `json:"value,omitempty"`

	Method string `json:"method"`
	URL    string `json:"url"`
}

func (service *FormService) GetUserProfileFormTemplate() *Form {
	approaches := service.NotificationApproachService.GetAllNotificationTypes()
	options := []*FieldOption{}
	for _, approach := range approaches {
		option := FieldOption{Label: string(approach)}
		options = append(options, &option)
	}
	fields := []*FormField{
		{
			Key:      "full_name",
			Label:    "使用者名稱",
			Type:     FieldTypeText,
			ReadOnly: true,
			Required: true,
		},
		{
			Key:      "phone_number",
			Label:    "使用者手機",
			Type:     FieldTypeText,
			ReadOnly: false,
			Required: true,
		},
		{
			Key:      "address",
			Label:    "地址",
			Type:     FieldTypeText,
			ReadOnly: false,
			Required: true,
		},
		{
			Key:      "notification_approaches",
			Label:    "通知方式",
			Type:     FieldTypeCheckbox,
			ReadOnly: false,
			Required: true,
			Options:  options,
		},
		{
			Key:      "user_description",
			Label:    "使用者描述",
			Type:     FieldTypeTextArea,
			ReadOnly: false,
			Required: true,
			Options:  options,
		},
		{
			Key:      "profile_picture_url",
			Label:    "個人照片",
			Type:     FieldTypeFile,
			ReadOnly: false,
			Required: false, // this can be uploaded later
		},
	}

	return &Form{Fields: fields}
}

func (service *FormService) GetUserProfileForm(ctx context.Context, userID uint) (*Form, error) {
	user, err := service.UserService.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	profile, _ := service.ProfileService.FindProfileByUserId(ctx, userID)
	notificationEnabledMapLookup := make(map[NotificationType]bool)

	var phoneNumber string
	var address string
	var profilePictureURL string
	if profile != nil {
		phoneNumber = profile.PhoneNumber
		address = profile.Address
		if len(profile.ProfilePictureObjectName) == 0 {
			profilePictureURL = profile.ProfilePictureObjectName
		}
		for _, approach := range profile.NotificationApproaches {
			notificationEnabledMapLookup[approach.Name] = approach.Enabled
		}

	}

	checkOptions := []*FieldOption{}
	for _, currentType := range service.NotificationApproachService.GetAllNotificationTypes() {
		isDisabled := true
		disabledMessage := ""
		switch currentType {
		case NotificationTypeEmail:
			if user.IsVerified {
				isDisabled = false
				break
			}
			disabledMessage = "使用者未驗證Email"
		case NotificationTypeSMS:
			if profile != nil && profile.IsPhoneNumberVerified {
				isDisabled = false
				break
			}
			disabledMessage = "使用者未驗證手機號碼"
		case NotificationTypeLineMessage:
			if baseApp.PlatformType(user.Platform.Name) == baseApp.PlatformLine {
				isDisabled = false
				break
			}
			disabledMessage = "使用者登入平台不是Line"
		}
		enabled, ok := notificationEnabledMapLookup[currentType]
		isChecked := ok && enabled
		checkOption := FieldOption{
			Label:           string(currentType),
			IsChecked:       isChecked,
			IsDisabled:      isDisabled,
			DisabledMessage: disabledMessage,
		}

		checkOptions = append(checkOptions, &checkOption)
	}

	populatedValues := map[string]any{
		"full_name":               user.Name,
		"phone_number":            phoneNumber,
		"address":                 address,
		"notification_approaches": checkOptions,
		"profile_picture_url":     profilePictureURL,
	}

	form := service.GetUserProfileFormTemplate()
	// TODO: also need to populate notification options value once we have notification approach service setup
	for idx := range form.Fields {
		formField := form.Fields[idx]
		value, ok := populatedValues[formField.Key]
		if !ok {
			continue
		}

		isOptionValue := slices.Contains(OptionTypes, formField.Type)
		if isOptionValue {
			formField.Options = checkOptions
		} else {
			formField.Value = value
		}
	}

	return form, nil
}

func (service *FormService) GetProductFormTemplate() *Form {
	fields := []*FormField{
		{
			Key:      "name",
			Label:    "產品名稱",
			Type:     FieldTypeText,
			ReadOnly: false,
			Required: true,
		},
		{
			Key:      "price",
			Label:    "價格",
			Type:     FieldTypeText,
			ReadOnly: false,
			Required: true,
		},
		{
			Key:      "description",
			Label:    "產品描述",
			Type:     FieldTypeTextArea,
			ReadOnly: false,
			Required: true,
		},
	}

	return &Form{Fields: fields}
}

func (service *FormService) GetUserFormTemplate() *Form {
	fields := []*FormField{
		{
			Key:      "name",
			Label:    "名稱",
			Type:     FieldTypeText,
			ReadOnly: false,
			Required: true,
		},
		{
			Key:      "email",
			Label:    "電子郵件",
			Type:     FieldTypeText,
			ReadOnly: false,
			Required: true,
		},
		{
			Key:      "password",
			Label:    "密碼",
			Type:     FieldTypeText,
			ReadOnly: false,
			Required: true,
		},
	}

	return &Form{Fields: fields}
}
