package repository

import (
	"context"
	"gorm.io/gorm"
	"time"
)

type IOrderRepository interface {
	FindOrdersByCustomerIDAndDate(ctx context.Context, customerID uint, datetime time.Time) ([]*Order, error)
	FindOrderByID(ctx context.Context, orderID uint) (*Order, error)
	UpdateOrderState(ctx context.Context, orderID uint, state OrderState) error
	CreateOrder(ctx context.Context, order *Order) error
}

type OrderRepository struct {
	Engine *gorm.DB
}

func (repo *OrderRepository) FindOrderByID(ctx context.Context, orderID uint) (*Order, error) {
	order := Order{}
	tx := repo.Engine.WithContext(ctx).Where("id = ?", orderID).First(&order)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &order, nil
}

func (repo *OrderRepository) UpdateOrderState(ctx context.Context, orderID uint, state OrderState) error {
	tx := repo.Engine.WithContext(ctx).Model(&Order{}).Where("id = ?", orderID).Update("state", state)
	return tx.Error
}

func NewOrderRepository(engine *gorm.DB) *OrderRepository {
	return &OrderRepository{Engine: engine}
}

func (repo *OrderRepository) FindOrdersByCustomerIDAndDate(ctx context.Context, customerID uint, datetime time.Time) ([]*Order, error) {
	orders := []*Order{}
	tx := repo.Engine.WithContext(ctx).Where("customer_id = ? AND order_date = ?", customerID, datetime.Format(time.DateOnly)).Find(&orders)
	return orders, tx.Error
}

func (repo *OrderRepository) CreateOrder(ctx context.Context, order *Order) error {
	tx := repo.Engine.WithContext(ctx).Create(order)
	return tx.Error
}
