package repository

import (
	"context"
	"gorm.io/gorm"
	"time"
)

type IOrderRepository interface {
	FindOrdersByCustomerIDAndDate(ctx context.Context, customerID uint, datetime time.Time) ([]*CustomerOrder, error)
	FindOrderByID(ctx context.Context, orderID uint) (*CustomerOrder, error)
	UpdateOrderState(ctx context.Context, orderID uint, state OrderState) error
	CreateOrder(ctx context.Context, order *CustomerOrder) error
}

type OrderRepository struct {
	Engine *gorm.DB
}

func (repo *OrderRepository) FindOrderByID(ctx context.Context, orderID uint) (*CustomerOrder, error) {
	//TODO implement me
	panic("implement me")
}

func (repo *OrderRepository) UpdateOrderState(ctx context.Context, orderID uint, state OrderState) error {
	//TODO implement me
	panic("implement me")
}

func NewOrderRepository(engine *gorm.DB) *OrderRepository {
	return &OrderRepository{Engine: engine}
}

func (repo *OrderRepository) FindOrdersByCustomerIDAndDate(ctx context.Context, customerID uint, datetime time.Time) ([]*CustomerOrder, error) {
	orders := []*CustomerOrder{}
	tx := repo.Engine.WithContext(ctx).Where("customer_id = ? AND order_date = ?", customerID, datetime.Format(time.DateOnly)).Find(&orders)
	return orders, tx.Error
}

func (repo *OrderRepository) CreateOrder(ctx context.Context, order *CustomerOrder) error {
	tx := repo.Engine.WithContext(ctx).Create(order)
	return tx.Error
}
