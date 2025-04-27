package service_test // Use _test package convention

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/wandermaia/customer-api/internal/domain/model"
	mock_repository "github.com/wandermaia/customer-api/internal/domain/repository/mock" // Import the generated mock
	"github.com/wandermaia/customer-api/internal/domain/service"
)

// Helper para configurar o mock e o serviço para cada teste
func setup(t *testing.T) (context.Context, service.CustomerService, *mock_repository.MockCustomerRepository) {
	ctrl := gomock.NewController(t)
	// Não precisamos mais do DeferFinish explícito com gomock > v1.6
	// defer ctrl.Finish() // Descomente se usar versão antiga do gomock ou preferir ser explícito

	mockRepo := mock_repository.NewMockCustomerRepository(ctrl)
	customerService := service.NewCustomerService(mockRepo)
	ctx := context.Background() // Usar um contexto básico para os testes

	return ctx, customerService, mockRepo
}

func TestCustomerService_CreateCustomer(t *testing.T) {
	ctx, customerService, mockRepo := setup(t)

	t.Run("Success", func(t *testing.T) {
		customer := &model.Customer{Name: "Valid User", Email: "valid@example.com"}

		// Expectativa: O método Create do repositório será chamado com o customer
		// e deve retornar nil (sem erro).
		mockRepo.EXPECT().Create(ctx, customer).Return(nil).Times(1)

		err := customerService.CreateCustomer(ctx, customer)

		assert.NoError(t, err)
	})

	t.Run("Validation Error", func(t *testing.T) {
		// Cliente inválido (sem nome, por exemplo, assumindo que Validate() verifica isso)
		customer := &model.Customer{Email: "invalid@example.com"}

		// Expectativa: O método Create do repositório NÃO deve ser chamado.
		mockRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Times(0)

		err := customerService.CreateCustomer(ctx, customer)

		// Verifica se o erro retornado é o erro de validação esperado.
		assert.Error(t, err)
		assert.Equal(t, service.ErrInvalidCustomer, err)
	})

	t.Run("Repository Error", func(t *testing.T) {
		customer := &model.Customer{Name: "Valid User", Email: "valid@example.com"}
		repoErr := errors.New("database connection refused")

		// Expectativa: O método Create será chamado, mas retornará um erro.
		mockRepo.EXPECT().Create(ctx, customer).Return(repoErr).Times(1)

		err := customerService.CreateCustomer(ctx, customer)

		// Verifica se o erro retornado é o erro de operação de banco de dados esperado.
		assert.Error(t, err)
		assert.Equal(t, service.ErrDatabaseOperation, err)
		// Opcional: verificar se o erro original não foi perdido (se necessário)
		// assert.True(t, errors.Is(err, service.ErrDatabaseOperation)) // Go 1.13+
	})
}

func TestCustomerService_GetCustomerByID(t *testing.T) {
	ctx, customerService, mockRepo := setup(t)
	testID := uint(1)

	t.Run("Success", func(t *testing.T) {
		expectedCustomer := &model.Customer{ID: testID, Name: "Found User", Email: "found@example.com"}

		// Expectativa: GetByID será chamado com o ID correto e retornará o cliente.
		mockRepo.EXPECT().GetByID(ctx, testID).Return(expectedCustomer, nil).Times(1)

		customer, err := customerService.GetCustomerByID(ctx, testID)

		assert.NoError(t, err)
		assert.Equal(t, expectedCustomer, customer)
	})

	t.Run("Not Found", func(t *testing.T) {
		// Simula o erro que o repositório retornaria se não encontrasse (pode variar)
		// Usaremos um erro genérico aqui, pois o serviço deve traduzi-lo.
		repoErr := errors.New("record not found") // ou gorm.ErrRecordNotFound, etc.

		// Expectativa: GetByID será chamado, mas retornará um erro.
		mockRepo.EXPECT().GetByID(ctx, testID).Return(nil, repoErr).Times(1)

		customer, err := customerService.GetCustomerByID(ctx, testID)

		// Verifica se o erro retornado é o ErrCustomerNotFound esperado.
		assert.Error(t, err)
		assert.Equal(t, service.ErrCustomerNotFound, err)
		assert.Nil(t, customer) // Garante que nenhum cliente foi retornado.
	})
}

func TestCustomerService_GetAllCustomers(t *testing.T) {
	ctx, customerService, mockRepo := setup(t)

	t.Run("Success", func(t *testing.T) {
		expectedCustomers := []*model.Customer{
			{ID: 1, Name: "User One", Email: "one@example.com"},
			{ID: 2, Name: "User Two", Email: "two@example.com"},
		}

		// Expectativa: GetAll será chamado e retornará a lista de clientes.
		mockRepo.EXPECT().GetAll(ctx).Return(expectedCustomers, nil).Times(1)

		customers, err := customerService.GetAllCustomers(ctx)

		assert.NoError(t, err)
		assert.Equal(t, expectedCustomers, customers)
	})

	t.Run("Repository Error", func(t *testing.T) {
		repoErr := errors.New("failed to fetch all")

		// Expectativa: GetAll será chamado e retornará um erro.
		mockRepo.EXPECT().GetAll(ctx).Return(nil, repoErr).Times(1)

		customers, err := customerService.GetAllCustomers(ctx)

		// Verifica se o erro original do repositório foi repassado.
		assert.Error(t, err)
		assert.Equal(t, repoErr, err) // O serviço repassa o erro diretamente aqui.
		assert.Nil(t, customers)
	})
}

func TestCustomerService_GetCustomersByName(t *testing.T) {
	ctx, customerService, mockRepo := setup(t)
	searchName := "Test"

	t.Run("Success", func(t *testing.T) {
		expectedCustomers := []*model.Customer{
			{ID: 1, Name: "Test User 1", Email: "test1@example.com"},
			{ID: 3, Name: "Test User 3", Email: "test3@example.com"},
		}

		// Expectativa: GetByName será chamado com o nome correto.
		mockRepo.EXPECT().GetByName(ctx, searchName).Return(expectedCustomers, nil).Times(1)

		customers, err := customerService.GetCustomersByName(ctx, searchName)

		assert.NoError(t, err)
		assert.Equal(t, expectedCustomers, customers)
	})

	t.Run("Empty Name Error", func(t *testing.T) {
		// Expectativa: GetByName NÃO deve ser chamado.
		mockRepo.EXPECT().GetByName(gomock.Any(), gomock.Any()).Times(0)

		customers, err := customerService.GetCustomersByName(ctx, "") // Chama com nome vazio

		// Verifica se o erro retornado é o de validação.
		assert.Error(t, err)
		assert.Equal(t, service.ErrInvalidCustomer, err)
		assert.Nil(t, customers)
	})

	t.Run("Repository Error", func(t *testing.T) {
		repoErr := errors.New("search failed")

		// Expectativa: GetByName será chamado e retornará um erro.
		mockRepo.EXPECT().GetByName(ctx, searchName).Return(nil, repoErr).Times(1)

		customers, err := customerService.GetCustomersByName(ctx, searchName)

		// Verifica se o erro original do repositório foi repassado.
		assert.Error(t, err)
		assert.Equal(t, repoErr, err)
		assert.Nil(t, customers)
	})
}

func TestCustomerService_UpdateCustomer(t *testing.T) {
	ctx, customerService, mockRepo := setup(t)
	testID := uint(1)
	now := time.Now()
	customerToUpdate := &model.Customer{ID: testID, Name: "Updated Name", Email: "updated@example.com", UpdatedAt: now}

	t.Run("Success", func(t *testing.T) {
		existingCustomer := &model.Customer{ID: testID, Name: "Old Name", Email: "old@example.com"}

		// Expectativa 1: GetByID será chamado para verificar a existência.
		mockRepo.EXPECT().GetByID(ctx, testID).Return(existingCustomer, nil).Times(1)
		// Expectativa 2: Update será chamado com os dados atualizados.
		mockRepo.EXPECT().Update(ctx, customerToUpdate).Return(nil).Times(1)

		err := customerService.UpdateCustomer(ctx, customerToUpdate)

		assert.NoError(t, err)
	})

	t.Run("Validation Error", func(t *testing.T) {
		invalidCustomer := &model.Customer{ID: testID, Name: ""} // Nome inválido

		// Expectativa: Nem GetByID nem Update devem ser chamados.
		mockRepo.EXPECT().GetByID(gomock.Any(), gomock.Any()).Times(0)
		mockRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Times(0)

		err := customerService.UpdateCustomer(ctx, invalidCustomer)

		assert.Error(t, err)
		assert.Equal(t, service.ErrInvalidCustomer, err)
	})

	t.Run("Not Found Error", func(t *testing.T) {
		repoErr := errors.New("record not found")

		// Expectativa: GetByID será chamado e retornará erro.
		mockRepo.EXPECT().GetByID(ctx, testID).Return(nil, repoErr).Times(1)
		// Expectativa: Update NÃO deve ser chamado.
		mockRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Times(0)

		err := customerService.UpdateCustomer(ctx, customerToUpdate)

		assert.Error(t, err)
		assert.Equal(t, service.ErrCustomerNotFound, err)
	})

	t.Run("Repository Update Error", func(t *testing.T) {
		existingCustomer := &model.Customer{ID: testID, Name: "Old Name", Email: "old@example.com"}
		repoErr := errors.New("database update failed")

		// Expectativa 1: GetByID será chamado e retornará sucesso.
		mockRepo.EXPECT().GetByID(ctx, testID).Return(existingCustomer, nil).Times(1)
		// Expectativa 2: Update será chamado, mas retornará erro.
		mockRepo.EXPECT().Update(ctx, customerToUpdate).Return(repoErr).Times(1)

		err := customerService.UpdateCustomer(ctx, customerToUpdate)

		assert.Error(t, err)
		assert.Equal(t, service.ErrDatabaseOperation, err)
	})
}

func TestCustomerService_DeleteCustomer(t *testing.T) {
	ctx, customerService, mockRepo := setup(t)
	testID := uint(1)

	t.Run("Success", func(t *testing.T) {
		existingCustomer := &model.Customer{ID: testID, Name: "To Delete", Email: "delete@example.com"}

		// Expectativa 1: GetByID será chamado para verificar existência.
		mockRepo.EXPECT().GetByID(ctx, testID).Return(existingCustomer, nil).Times(1)
		// Expectativa 2: Delete será chamado com o ID correto.
		mockRepo.EXPECT().Delete(ctx, testID).Return(nil).Times(1)

		err := customerService.DeleteCustomer(ctx, testID)

		assert.NoError(t, err)
	})

	t.Run("Not Found Error", func(t *testing.T) {
		repoErr := errors.New("record not found")

		// Expectativa: GetByID será chamado e retornará erro.
		mockRepo.EXPECT().GetByID(ctx, testID).Return(nil, repoErr).Times(1)
		// Expectativa: Delete NÃO deve ser chamado.
		mockRepo.EXPECT().Delete(gomock.Any(), gomock.Any()).Times(0)

		err := customerService.DeleteCustomer(ctx, testID)

		assert.Error(t, err)
		assert.Equal(t, service.ErrCustomerNotFound, err)
	})

	t.Run("Repository Delete Error", func(t *testing.T) {
		existingCustomer := &model.Customer{ID: testID, Name: "To Delete", Email: "delete@example.com"}
		repoErr := errors.New("database delete failed")

		// Expectativa 1: GetByID será chamado e retornará sucesso.
		mockRepo.EXPECT().GetByID(ctx, testID).Return(existingCustomer, nil).Times(1)
		// Expectativa 2: Delete será chamado, mas retornará erro.
		mockRepo.EXPECT().Delete(ctx, testID).Return(repoErr).Times(1)

		err := customerService.DeleteCustomer(ctx, testID)

		assert.Error(t, err)
		assert.Equal(t, service.ErrDatabaseOperation, err)
	})
}

func TestCustomerService_CountCustomers(t *testing.T) {
	ctx, customerService, mockRepo := setup(t)

	t.Run("Success", func(t *testing.T) {
		expectedCount := int64(42)

		// Expectativa: Count será chamado e retornará a contagem.
		mockRepo.EXPECT().Count(ctx).Return(expectedCount, nil).Times(1)

		count, err := customerService.CountCustomers(ctx)

		assert.NoError(t, err)
		assert.Equal(t, expectedCount, count)
	})

	t.Run("Repository Error", func(t *testing.T) {
		repoErr := errors.New("failed to count")

		// Expectativa: Count será chamado e retornará um erro.
		mockRepo.EXPECT().Count(ctx).Return(int64(0), repoErr).Times(1)

		count, err := customerService.CountCustomers(ctx)

		// Verifica se o erro original do repositório foi repassado.
		assert.Error(t, err)
		assert.Equal(t, repoErr, err)
		assert.Equal(t, int64(0), count) // Verifica se a contagem retornada é 0 em caso de erro.
	})
}
