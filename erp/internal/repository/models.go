package repository

import (
	"encoding/json"
	. "go-security/security/repository"
	"time"
)

type PricingPolicy struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	Name        string `gorm:"type:varchar(100);not null;unique" json:"name"` // Policy name
	Description string `gorm:"type:text" json:"description"`                  // Optional description

	PolicyPrices []PolicyPrice `gorm:"foreignKey:PolicyID" json:"policy_prices"` // Prices for specific products
	CreatedAt    time.Time     `json:"created_at"`
	UpdatedAt    time.Time     `json:"updated_at"`
	DeletedAt    *time.Time    `json:"deleted_at"`
}

type PolicyPrice struct {
	ID        uint       `gorm:"primaryKey" json:"id"`
	PolicyID  uint       `gorm:"not null;uniqueIndex:policy_product_idx" json:"policy_id"`  // Foreign key linking to PricingPolicy
	ProductID uint       `gorm:"not null;uniqueIndex:policy_product_idx" json:"product_id"` // Foreign key linking to Product
	Product   Product    `gorm:"foreignKey:ProductID; references:ID" json:"-"`              // Product relationship
	Price     int        `gorm:"type:int;not null" json:"price"`                            // Price for the product under this policy
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

type OrderNotification struct {
	NotifiedOrderState OrderState `gorm:"not null"  json:"notified_order_state"`

	ID        uint       `gorm:"primaryKey" json:"id"` // Auto-increment primary key
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

// TODO: NotificatonApproachService for creating all notification approaches for a user
type NotificationApproach struct {
	ID     uint `gorm:"primaryKey" json:"id"`
	UserID uint `gorm:"not null" json:"user_id"` // Foreign key linking to User.ID
	User   User `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"-"`

	Name    NotificationType `gorm:"type:varchar(50);not null" json:"approach"`
	Enabled bool             `gorm:"default:true" json:"enabled"`
}

type UserProfile struct {
	ID     uint `gorm:"primaryKey" json:"id"`
	UserID uint `gorm:"not null; unique" json:"user_id"` // Foreign key linking to User.ID
	User   User `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"-"`

	NotificationApproaches []NotificationApproach `gorm:"foreignKey:UserID" json:"notification_approaches"`
	PhoneNumber            string                 `gorm:"type:varchar(20)" json:"phone_number"` // for SMS
	IsPhoneNumberVerified  bool                   `gorm:"default:false" json:"is_phone_number_verified"`

	PricingPolicyID uint          `json:"pricing_policy_id"`                                              // Foreign key linking to PricingPolicy
	PricingPolicy   PricingPolicy `gorm:"foreignKey:PricingPolicyID;references:ID" json:"pricing_policy"` // Many-to-many relationship

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`

	UserDescription          string `gorm:"type:text" json:"user_description"`
	Address                  string `gorm:"type:text" json:"address"`
	ProfilePictureObjectName string `gorm:"type:text" json:"profile_picture_object_name"`
}

func (profile *UserProfile) AllNotificationTypes() []NotificationType {
	var types []NotificationType
	for _, approach := range profile.NotificationApproaches {
		types = append(types, approach.Name)
	}
	return types
}

type CustomerOrder struct {
	UserID     uint       `gorm:"not null" json:"user_id"` // Foreign key linking to User.ID
	OrderState OrderState `gorm:"type:varchar(50); not null" json:"order_state"`
	OrderDate  time.Time  `gorm:"type:date"`
	User       User       `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"-"` // Relationship to User
	Products   []Product  `gorm:"many2many:order_products;" json:"products"`

	ID        uint       `gorm:"primaryKey" json:"id"` // Auto-increment primary key
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

func (order *CustomerOrder) AddProduct(product *Product) {
	order.Products = append(order.Products, *product)
}

func (order *CustomerOrder) RemoveProduct(product *Product) {
	for idx, p := range order.Products {
		if p.ID != product.ID {
			continue
		}
		order.Products = append(order.Products[:idx], order.Products[idx+1:]...)
	}
}

func (order *CustomerOrder) AsJson() (string, error) {
	_json, err := json.Marshal(order)
	if err != nil {
		return "", err
	}

	var _map map[string]any
	if err := json.Unmarshal(_json, &_map); err != nil {
		return "", err
	}

	updatedJson, err := json.Marshal(_map)
	if err != nil {
		return "", err
	}
	return string(updatedJson), nil
}

type Product struct {
	ID          uint    `gorm:"primaryKey" json:"id"` // Auto-increment primary key
	Name        string  `gorm:"type:varchar(100); not null" json:"name"`
	Description string  `gorm:"type:text" json:"description"`
	PictureURL  *string `gorm:"type:text" json:"picture_url"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}
