package service

import (
	"context"
	"errors"

	"github.com/wandermaia/customer-api/internal/domain/model"
	"github.com/wandermaia/customer-api/internal/domain/repository"
)

var (
	ErrInvalidCustomer   = errors.New("dados do cliente inválidos")
	ErrCustomerNotFound  = errors.New("cliente não encontrado")
	ErrDatabaseOperation = errors.New("erro na operação do banco de dados")
)

// CustomerService define as operações de serviço para clientes
type CustomerService interface {
	CreateCustomer(ctx context.Context, customer *model.Customer) error
	GetCustomerByID(ctx context.Context, id uint) (*model.Customer, error)
	GetAllCustomers(ctx context.Context) ([]*model.Customer, error)
	GetCustomersByName(ctx context.Context, name string) ([]*model.Customer, error)
	UpdateCustomer(ctx context.Context, customer *model.Customer) error
	DeleteCustomer(ctx context.Context, id uint) error
	CountCustomers(ctx context.Context) (int64, error)
}

type customerService struct {
	repo repository.CustomerRepository
}

// NewCustomerService cria uma nova instância do serviço de clientes
func NewCustomerService(repo repository.CustomerRepository) CustomerService {
	return &customerService{
		repo: repo,
	}
}

// CreateCustomer cria um novo cliente
func (s *customerService) CreateCustomer(ctx context.Context, customer *model.Customer) error {
	if err := customer.Validate(); err != nil {
		return ErrInvalidCustomer
	}

	if err := s.repo.Create(ctx, customer); err != nil {
		return ErrDatabaseOperation
	}

	return nil
}

// GetCustomerByID busca um cliente pelo ID
func (s *customerService) GetCustomerByID(ctx context.Context, id uint) (*model.Customer, error) {
	customer, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, ErrCustomerNotFound
	}

	return customer, nil
}

// GetAllCustomers retorna todos os clientes
func (s *customerService) GetAllCustomers(ctx context.Context) ([]*model.Customer, error) {
	return s.repo.GetAll(ctx)
}

// GetCustomersByName busca clientes pelo nome
func (s *customerService) GetCustomersByName(ctx context.Context, name string) ([]*model.Customer, error) {
	if name == "" {
		return nil, ErrInvalidCustomer
	}

	return s.repo.GetByName(ctx, name)
}

// UpdateCustomer atualiza um cliente existente
func (s *customerService) UpdateCustomer(ctx context.Context, customer *model.Customer) error {
	if err := customer.Validate(); err != nil {
		return ErrInvalidCustomer
	}

	// Verifica se o cliente existe
	_, err := s.repo.GetByID(ctx, customer.ID)
	if err != nil {
		return ErrCustomerNotFound
	}

	if err := s.repo.Update(ctx, customer); err != nil {
		return ErrDatabaseOperation
	}

	return nil
}

// DeleteCustomer remove um cliente pelo ID
func (s *customerService) DeleteCustomer(ctx context.Context, id uint) error {
	// Verifica se o cliente existe
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return ErrCustomerNotFound
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return ErrDatabaseOperation
	}

	return nil
}

// CountCustomers retorna o número total de clientes
func (s *customerService) CountCustomers(ctx context.Context) (int64, error) {
	return s.repo.Count(ctx)
}
