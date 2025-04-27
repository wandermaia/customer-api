package repository

import (
	"context"

	"github.com/wandermaia/customer-api/internal/domain/model"

	"gorm.io/gorm"
)

type postgresCustomerRepository struct {
	db *gorm.DB
}

// NewPostgresCustomerRepository cria uma nova instância do repositório PostgreSQL
func NewPostgresCustomerRepository(db *gorm.DB) CustomerRepository {
	return &postgresCustomerRepository{
		db: db,
	}
}

// Create insere um novo cliente no banco de dados
func (r *postgresCustomerRepository) Create(ctx context.Context, customer *model.Customer) error {
	return r.db.WithContext(ctx).Create(customer).Error
}

// GetByID busca um cliente pelo ID
func (r *postgresCustomerRepository) GetByID(ctx context.Context, id uint) (*model.Customer, error) {
	var customer model.Customer
	if err := r.db.WithContext(ctx).First(&customer, id).Error; err != nil {
		return nil, err
	}
	return &customer, nil
}

// GetAll retorna todos os clientes
func (r *postgresCustomerRepository) GetAll(ctx context.Context) ([]*model.Customer, error) {
	var customers []*model.Customer
	if err := r.db.WithContext(ctx).Find(&customers).Error; err != nil {
		return nil, err
	}
	return customers, nil
}

// GetByName busca clientes pelo nome
func (r *postgresCustomerRepository) GetByName(ctx context.Context, name string) ([]*model.Customer, error) {
	var customers []*model.Customer
	if err := r.db.WithContext(ctx).Where("name ILIKE ?", "%"+name+"%").Find(&customers).Error; err != nil {
		return nil, err
	}
	return customers, nil
}

// Update atualiza um cliente existente
func (r *postgresCustomerRepository) Update(ctx context.Context, customer *model.Customer) error {
	return r.db.WithContext(ctx).Save(customer).Error
}

// Delete remove um cliente pelo ID
func (r *postgresCustomerRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.Customer{}, id).Error
}

// Count retorna o número total de clientes
func (r *postgresCustomerRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.Customer{}).Count(&count).Error
	return count, err
}
