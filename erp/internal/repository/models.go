package repository

import (
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

type NotificationApproach struct {
	ID     uint `gorm:"primaryKey" json:"id"`
	UserID uint `gorm:"not null" json:"user_id"` // Foreign key linking to User.ID
	User   User `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"-"`

	Name    NotificationType `gorm:"type:varchar(50);not null" json:"approach"`
	Enabled bool             `gorm:"default:false" json:"enabled"`
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

func (profile *UserProfile) HasProfilePicture() bool {
	return len(profile.ProfilePictureObjectName) != 0
}

func (profile *UserProfile) AllNotificationTypes() []NotificationType {
	var types []NotificationType
	for _, approach := range profile.NotificationApproaches {
		types = append(types, approach.Name)
	}
	return types
}

type Order struct {
	UserID     uint        `gorm:"not null" json:"user_id"` // Foreign key linking to User.ID
	OrderState OrderState  `gorm:"type:varchar(50); not null" json:"order_state"`
	OrderDate  time.Time   `gorm:"type:date"`
	User       User        `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"-"` // Relationship to User
	OrderItems []OrderItem `gorm:"foreignKey:OrderID" json:"order_items"`                                                   // Relationship to OrderItem

	ID        uint       `gorm:"primaryKey" json:"id"` // Auto-increment primary key
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

type OrderItem struct {
	OrderID      uint    `gorm:"not null" json:"order_id"`   // Foreign key linking to Order.ID
	ProductID    uint    `gorm:"not null" json:"product_id"` // Foreign key linking to Product.ID
	Product      Product `gorm:"foreignKey:ProductID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"product"`
	Quantity     uint    `gorm:"type:int;not null" json:"quantity"`
	PricePerUnit int     `gorm:"type:int;not null" json:"price_per_unit"`
	TotalPrice   int     `gorm:"type:int;not null" json:"total_price"`
}

type ProductCategory struct {
	ID          uint      `gorm:"primaryKey" json:"id"` // Auto-increment primary key
	Name        string    `gorm:"type:varchar(100);not null;unique" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	Products    []Product `gorm:"foreignKey:CategoryID" json:"products"` // One-to-Many relationship

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `gorm:"index" json:"deleted_at"`
}

type Product struct {
	ID                       uint   `gorm:"primaryKey" json:"id"` // Auto-increment primary key
	Name                     string `gorm:"type:varchar(100);not null;unique" json:"name"`
	Description              string `gorm:"type:text" json:"description"`
	ProfilePictureObjectName string `gorm:"type:text" json:"profile_picture_object_name"`

	CategoryID uint            `json:"category_id"` // Foreign key
	Category   ProductCategory `gorm:"foreignKey:CategoryID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;default:null" json:"category"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `gorm:"index" json:"deleted_at"`

	ProfilePictureURL string `gorm:"-" json:"profile_picture_url"`
	Cost              uint   `gorm:"type:int;not null" json:"cost"`
}

func (product *Product) HasProfilePicture() bool {
	return len(product.ProfilePictureObjectName) != 0
}
