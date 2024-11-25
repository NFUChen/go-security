package repository

import (
	"encoding/json"
	. "go-security/security/repository"
	"gorm.io/gorm"
	"time"
)

type OrderNotification struct {
	NotifiedOrderState OrderState `gorm:"not null"  json:"notified_order_state"`

	ID        uint       `gorm:"primaryKey" json:"id"` // Auto-increment primary key
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

type CustomerProfile struct {
	CustomerID             uint                   `gorm:"not null" json:"customer_id"`                                                                 // Foreign key linking to User.ID
	Customer               User                   `gorm:"foreignKey:CustomerID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"-"` // Relationship to User
	NotificationApproaches []NotificationApproach `gorm:"not null; type:json" json:"notification_approach"`
	PhoneNumber            string                 `gorm:"type:varchar(20)" json:"phone_number"`   // for SMS
	EmailAddress           string                 `gorm:"type:varchar(100)" json:"email_address"` // for Email
	LineID                 string                 `gorm:"type:varchar(100)" json:"line_id"`       // for LINE

	ID        uint       `gorm:"primaryKey" json:"id"` // Auto-increment primary key
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

type CustomerOrder struct {
	CustomerID uint       `gorm:"not null" json:"customer_id"` // Foreign key linking to User.ID
	OrderState OrderState `gorm:"type:varchar(50); not null" json:"order_state"`
	OrderDate  time.Time  `gorm:"type:date"`
	Customer   User       `gorm:"foreignKey:CustomerID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"-"` // Relationship to User
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

func (order *CustomerOrder) TotalAmount() int {
	amount := 0
	for _, product := range order.Products {
		amount += product.Price
	}
	return amount
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

	_map["total_amount"] = order.TotalAmount()
	updatedJson, err := json.Marshal(_map)
	if err != nil {
		return "", err
	}
	return string(updatedJson), nil
}

type Product struct {
	gorm.Model
	ID    uint   `gorm:"primaryKey" json:"id"` // Auto-increment primary key
	Name  string `gorm:"type:varchar(100); not null" json:"name"`
	Price int    `gorm:"type:int; not null" json:"amount"`
}
