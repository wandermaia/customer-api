package repository

import (
	"context"

	"github.com/wandermaia/customer-api/internal/domain/model"
)

// CustomerRepository define as operações do repositório de clientes
type CustomerRepository interface {
	Create(ctx context.Context, customer *model.Customer) error
	GetByID(ctx context.Context, id uint) (*model.Customer, error)
	GetAll(ctx context.Context) ([]*model.Customer, error)
	GetByName(ctx context.Context, name string) ([]*model.Customer, error)
	Update(ctx context.Context, customer *model.Customer) error
	Delete(ctx context.Context, id uint) error
	Count(ctx context.Context) (int64, error)
}
