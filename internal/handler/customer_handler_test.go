// Package handler_test contém os testes de unidade para os handlers HTTP da API de clientes.
// Estes testes utilizam mocks para isolar a camada de handler da camada de serviço,
// garantindo que apenas a lógica do handler (parsing de requisições, chamadas de serviço,
// formatação de respostas e tratamento de erros) seja testada.
package handler_test // Use _test package convention

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/wandermaia/customer-api/internal/domain/model"
	"github.com/wandermaia/customer-api/internal/domain/service"
	mock_service "github.com/wandermaia/customer-api/internal/domain/service/mock" // Import the generated mock
	"github.com/wandermaia/customer-api/internal/handler"
)

// setupTestRouter cria e configura um novo Gin engine em modo de teste.
// Ele instancia o CustomerHandler com o mock do serviço fornecido,
// registra as rotas do handler nesse engine e retorna o engine e um ResponseRecorder
// para capturar as respostas das requisições simuladas.
func setupTestRouter(t *testing.T, mockService *mock_service.MockCustomerService) (*gin.Engine, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)                                  // Garante que o Gin opere em modo de teste (menos verbose)
	router := gin.New()                                        // Cria um novo engine Gin sem middlewares padrão
	recorder := httptest.NewRecorder()                         // Cria um gravador para capturar a resposta HTTP
	customerHandler := handler.NewCustomerHandler(mockService) // Cria o handler com o serviço mockado
	customerHandler.RegisterRoutes(router)                     // Registra as rotas da API no router
	return router, recorder
}

// performRequest simula o envio de uma requisição HTTP para o router de teste.
// Cria uma nova requisição HTTP com o método, caminho e corpo (opcional) fornecidos.
// Se um corpo for fornecido, define o cabeçalho Content-Type como application/json.
// Executa a requisição contra o router, e a resposta é capturada pelo recorder.
func performRequest(router *gin.Engine, recorder *httptest.ResponseRecorder, method, path string, body []byte) {
	// Cria a requisição HTTP. Ignora o erro pois os parâmetros são controlados pelo teste.
	req, _ := http.NewRequest(method, path, bytes.NewBuffer(body))
	if body != nil {
		// Define o tipo de conteúdo se houver corpo na requisição
		req.Header.Set("Content-Type", "application/json")
	}
	// Envia a requisição para o router, que a direciona para o handler apropriado.
	// A resposta do handler é escrita no recorder.
	router.ServeHTTP(recorder, req)
}

// TestCustomerHandler_CreateCustomer testa o endpoint POST /api/customers.
// Ele verifica os cenários de sucesso na criação de um cliente, falha por JSON inválido,
// falha devido a erro de validação retornado pelo serviço e falha por erro interno do serviço.
func TestCustomerHandler_CreateCustomer(t *testing.T) {
	// Cria o controlador do Gomock para gerenciar o ciclo de vida dos mocks.
	mockCtrl := gomock.NewController(t)
	// Garante que as expectativas dos mocks foram atendidas no final do teste.
	defer mockCtrl.Finish()

	// Cria o mock do CustomerService e o router de teste.
	mockService := mock_service.NewMockCustomerService(mockCtrl)
	router, recorder := setupTestRouter(t, mockService)

	// ctx := context.Background() // Contexto não usado diretamente nas chamadas do mock aqui, gomock.Any() é usado.

	// Subteste para o cenário de sucesso.
	t.Run("Success", func(t *testing.T) {
		// Define os dados de entrada e a saída esperada.
		customerInput := model.Customer{Name: "Test User", Email: "test@example.com"}
		customerOutput := model.Customer{ID: 1, Name: "Test User", Email: "test@example.com"}
		customerJSON, _ := json.Marshal(customerInput) // Converte a entrada para JSON.

		// Define a expectativa: o método CreateCustomer do serviço deve ser chamado uma vez.
		mockService.EXPECT().
			CreateCustomer(gomock.Any(), gomock.Any()). // Aceita qualquer contexto e ponteiro para Customer.
			DoAndReturn(func(ctx context.Context, c *model.Customer) error {
				// Simula a lógica do serviço: atribui um ID ao cliente criado.
				c.ID = customerOutput.ID
				// Verifica se os dados recebidos pelo serviço correspondem à entrada.
				assert.Equal(t, customerInput.Name, c.Name)
				assert.Equal(t, customerInput.Email, c.Email)
				return nil // Retorna nil para simular sucesso.
			}).Times(1) // Espera que seja chamado exatamente uma vez.

		// Reseta o recorder para este subteste específico.
		recorder = httptest.NewRecorder()
		// Executa a requisição POST.
		performRequest(router, recorder, http.MethodPost, "/api/customers", customerJSON)

		// Verifica o código de status HTTP (esperado: 201 Created).
		assert.Equal(t, http.StatusCreated, recorder.Code)

		// Verifica o corpo da resposta.
		var createdCustomer model.Customer
		err := json.Unmarshal(recorder.Body.Bytes(), &createdCustomer) // Decodifica o JSON da resposta.
		assert.NoError(t, err)                                         // Garante que não houve erro na decodificação.
		assert.Equal(t, customerOutput, createdCustomer)               // Compara o cliente retornado com o esperado.
	})

	// Subteste para o cenário de JSON inválido na requisição.
	t.Run("Invalid JSON", func(t *testing.T) {
		recorder = httptest.NewRecorder()
		// Envia um corpo JSON malformado.
		performRequest(router, recorder, http.MethodPost, "/api/customers", []byte(`{"invalid json`))

		// Verifica se o status é 400 Bad Request.
		assert.Equal(t, http.StatusBadRequest, recorder.Code)
		// Verifica se a mensagem de erro esperada está no corpo da resposta.
		assert.Contains(t, recorder.Body.String(), "Dados inválidos")
	})

	// Subteste para o cenário onde o serviço retorna um erro de validação.
	t.Run("Service Validation Error", func(t *testing.T) {
		// Prepara uma entrada que causaria um erro de validação no serviço (ex: email faltando).
		customerInput := model.Customer{Name: "Bad User"}
		customerJSON, _ := json.Marshal(customerInput)

		// Define a expectativa: CreateCustomer será chamado, mas retornará um erro de validação.
		mockService.EXPECT().
			CreateCustomer(gomock.Any(), gomock.Any()).
			Return(service.ErrInvalidCustomer). // Simula o serviço retornando o erro específico.
			Times(1)

		recorder = httptest.NewRecorder()
		performRequest(router, recorder, http.MethodPost, "/api/customers", customerJSON)

		// Verifica se o status é 400 Bad Request.
		assert.Equal(t, http.StatusBadRequest, recorder.Code)
		// Verifica se a mensagem de erro do serviço está na resposta.
		assert.Contains(t, recorder.Body.String(), service.ErrInvalidCustomer.Error())
	})

	// Subteste para o cenário onde o serviço retorna um erro interno genérico.
	t.Run("Service Internal Error", func(t *testing.T) {
		customerInput := model.Customer{Name: "Test User", Email: "test@example.com"}
		customerJSON, _ := json.Marshal(customerInput)
		serviceErr := errors.New("database connection failed") // Erro genérico simulado.

		// Define a expectativa: CreateCustomer será chamado e retornará um erro interno.
		mockService.EXPECT().
			CreateCustomer(gomock.Any(), gomock.Any()).
			Return(serviceErr). // Simula o serviço retornando o erro genérico.
			Times(1)

		recorder = httptest.NewRecorder()
		performRequest(router, recorder, http.MethodPost, "/api/customers", customerJSON)

		// Verifica se o status é 500 Internal Server Error.
		assert.Equal(t, http.StatusInternalServerError, recorder.Code)
		// Verifica se a mensagem de erro genérica do handler está na resposta.
		assert.Contains(t, recorder.Body.String(), "Erro ao criar cliente")
	})
}

// TestCustomerHandler_GetCustomerByID testa o endpoint GET /api/customers/{id}.
// Verifica os cenários de sucesso ao buscar um cliente, falha por ID inválido na URL,
// falha quando o cliente não é encontrado (Not Found) e falha por erro interno do serviço.
func TestCustomerHandler_GetCustomerByID(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockService := mock_service.NewMockCustomerService(mockCtrl)
	router, recorder := setupTestRouter(t, mockService)

	// Define um ID de teste.
	testID := uint(123)
	testIDStr := strconv.FormatUint(uint64(testID), 10) // Converte o ID para string para usar na URL.

	// Subteste para o cenário de sucesso.
	t.Run("Success", func(t *testing.T) {
		// Define o cliente esperado que o serviço retornará.
		expectedCustomer := model.Customer{ID: testID, Name: "Found User", Email: "found@example.com"}

		// Define a expectativa: GetCustomerByID será chamado com o ID correto e retornará o cliente.
		mockService.EXPECT().
			GetCustomerByID(gomock.Any(), testID). // Espera o ID específico.
			Return(&expectedCustomer, nil).        // Retorna o cliente encontrado (como ponteiro) e nenhum erro.
			Times(1)

		recorder = httptest.NewRecorder()
		// Executa a requisição GET para o ID específico.
		performRequest(router, recorder, http.MethodGet, "/api/customers/"+testIDStr, nil)

		// Verifica o status 200 OK.
		assert.Equal(t, http.StatusOK, recorder.Code)

		// Verifica o corpo da resposta.
		var foundCustomer model.Customer
		err := json.Unmarshal(recorder.Body.Bytes(), &foundCustomer)
		assert.NoError(t, err)
		assert.Equal(t, expectedCustomer, foundCustomer) // Compara o cliente da resposta com o esperado.
	})

	// Subteste para o cenário de formato de ID inválido na URL.
	t.Run("Invalid ID Format", func(t *testing.T) {
		recorder = httptest.NewRecorder()
		// Executa a requisição com um ID não numérico.
		performRequest(router, recorder, http.MethodGet, "/api/customers/abc", nil)

		// Verifica o status 400 Bad Request.
		assert.Equal(t, http.StatusBadRequest, recorder.Code)
		// Verifica a mensagem de erro específica.
		assert.Contains(t, recorder.Body.String(), "ID inválido")
	})

	// Subteste para o cenário onde o cliente não é encontrado pelo serviço.
	t.Run("Not Found", func(t *testing.T) {
		// Define a expectativa: GetCustomerByID será chamado, mas retornará ErrCustomerNotFound.
		mockService.EXPECT().
			GetCustomerByID(gomock.Any(), testID).
			Return(nil, service.ErrCustomerNotFound). // Retorna nil para cliente e o erro específico.
			Times(1)

		recorder = httptest.NewRecorder()
		performRequest(router, recorder, http.MethodGet, "/api/customers/"+testIDStr, nil)

		// Verifica o status 404 Not Found.
		assert.Equal(t, http.StatusNotFound, recorder.Code)
		// Verifica se a mensagem de erro do serviço está na resposta.
		assert.Contains(t, recorder.Body.String(), service.ErrCustomerNotFound.Error())
	})

	// Subteste para o cenário de erro interno do serviço.
	t.Run("Service Internal Error", func(t *testing.T) {
		serviceErr := errors.New("database query failed") // Erro genérico simulado.

		// Define a expectativa: GetCustomerByID será chamado e retornará um erro interno.
		mockService.EXPECT().
			GetCustomerByID(gomock.Any(), testID).
			Return(nil, serviceErr). // Retorna nil para cliente e o erro genérico.
			Times(1)

		recorder = httptest.NewRecorder()
		performRequest(router, recorder, http.MethodGet, "/api/customers/"+testIDStr, nil)

		// Verifica o status 500 Internal Server Error.
		assert.Equal(t, http.StatusInternalServerError, recorder.Code)
		// Verifica a mensagem de erro genérica do handler.
		assert.Contains(t, recorder.Body.String(), "Erro ao buscar cliente")
	})
}

// TestCustomerHandler_GetAllCustomers testa o endpoint GET /api/customers.
// Verifica o cenário de sucesso ao buscar todos os clientes e o cenário de erro interno do serviço.
func TestCustomerHandler_GetAllCustomers(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockService := mock_service.NewMockCustomerService(mockCtrl)
	router, recorder := setupTestRouter(t, mockService)

	// Subteste para o cenário de sucesso.
	t.Run("Success", func(t *testing.T) {
		// Define a lista esperada de clientes (como slice de ponteiros, conforme a assinatura do serviço).
		expectedCustomers := []*model.Customer{
			{ID: 1, Name: "User One", Email: "one@example.com"},
			{ID: 2, Name: "User Two", Email: "two@example.com"},
		}

		// Define a expectativa: GetAllCustomers será chamado e retornará a lista de clientes.
		mockService.EXPECT().
			GetAllCustomers(gomock.Any()).
			Return(expectedCustomers, nil). // Retorna o slice de ponteiros e nenhum erro.
			Times(1)

		recorder = httptest.NewRecorder()
		performRequest(router, recorder, http.MethodGet, "/api/customers", nil)

		// Verifica o status 200 OK.
		assert.Equal(t, http.StatusOK, recorder.Code)

		// Verifica o corpo da resposta.
		// Assume que o handler serializa para um slice de valores no JSON.
		var customers []model.Customer
		err := json.Unmarshal(recorder.Body.Bytes(), &customers)
		assert.NoError(t, err)

		// Cria um slice de valores esperado para comparação com o resultado do JSON.
		expectedValues := []model.Customer{
			{ID: 1, Name: "User One", Email: "one@example.com"},
			{ID: 2, Name: "User Two", Email: "two@example.com"},
		}
		assert.Equal(t, expectedValues, customers) // Compara a lista da resposta com os valores esperados.
	})

	// Subteste para o cenário de erro interno do serviço.
	t.Run("Service Internal Error", func(t *testing.T) {
		serviceErr := errors.New("failed to fetch all") // Erro genérico simulado.

		// Define a expectativa: GetAllCustomers será chamado e retornará um erro interno.
		mockService.EXPECT().
			GetAllCustomers(gomock.Any()).
			Return(nil, serviceErr). // Retorna nil para a lista e o erro genérico.
			Times(1)

		recorder = httptest.NewRecorder()
		performRequest(router, recorder, http.MethodGet, "/api/customers", nil)

		// Verifica o status 500 Internal Server Error.
		assert.Equal(t, http.StatusInternalServerError, recorder.Code)
		// Verifica a mensagem de erro genérica do handler.
		assert.Contains(t, recorder.Body.String(), "Erro ao buscar clientes")
	})
}

// TestCustomerHandler_GetCustomersByName testa o endpoint GET /api/customers/search?name={name}.
// Verifica os cenários de sucesso na busca por nome, falha se o nome não for fornecido
// e falha por erro interno do serviço.
func TestCustomerHandler_GetCustomersByName(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockService := mock_service.NewMockCustomerService(mockCtrl)
	router, recorder := setupTestRouter(t, mockService)

	searchName := "Test" // Nome a ser buscado.

	// Subteste para o cenário de sucesso.
	t.Run("Success", func(t *testing.T) {
		// Define a lista esperada de clientes encontrados (como slice de ponteiros).
		expectedCustomers := []*model.Customer{
			{ID: 1, Name: "Test User 1", Email: "test1@example.com"},
			{ID: 3, Name: "Test User 3", Email: "test3@example.com"},
		}

		// Define a expectativa: GetCustomersByName será chamado com o nome correto e retornará a lista.
		mockService.EXPECT().
			GetCustomersByName(gomock.Any(), searchName). // Espera o nome específico.
			Return(expectedCustomers, nil).               // Retorna a lista encontrada e nenhum erro.
			Times(1)

		recorder = httptest.NewRecorder()
		// Executa a requisição GET com o parâmetro de query 'name'.
		performRequest(router, recorder, http.MethodGet, "/api/customers/search?name="+searchName, nil)

		// Verifica o status 200 OK.
		assert.Equal(t, http.StatusOK, recorder.Code)

		// Verifica o corpo da resposta. Assume que o handler retorna ponteiros ou valores que podem ser comparados diretamente.
		// Ajuste o tipo de 'customers' se necessário (para []model.Customer).
		var customers []*model.Customer
		err := json.Unmarshal(recorder.Body.Bytes(), &customers)
		assert.NoError(t, err)
		// testify/assert.Equal pode comparar slices de ponteiros por valor dos elementos apontados.
		assert.Equal(t, expectedCustomers, customers)
	})

	// Subteste para o cenário onde o parâmetro 'name' não é fornecido.
	t.Run("Name Not Provided", func(t *testing.T) {
		recorder = httptest.NewRecorder()
		// Executa a requisição sem o parâmetro 'name'.
		performRequest(router, recorder, http.MethodGet, "/api/customers/search", nil)

		// Verifica o status 400 Bad Request.
		assert.Equal(t, http.StatusBadRequest, recorder.Code)
		// Verifica a mensagem de erro específica.
		assert.Contains(t, recorder.Body.String(), "Nome não fornecido")
	})

	// Subteste para o cenário de erro interno do serviço.
	t.Run("Service Internal Error", func(t *testing.T) {
		serviceErr := errors.New("search failed") // Erro genérico simulado.

		// Define a expectativa: GetCustomersByName será chamado e retornará um erro interno.
		mockService.EXPECT().
			GetCustomersByName(gomock.Any(), searchName).
			Return(nil, serviceErr). // Retorna nil para a lista e o erro genérico.
			Times(1)

		recorder = httptest.NewRecorder()
		performRequest(router, recorder, http.MethodGet, "/api/customers/search?name="+searchName, nil)

		// Verifica o status 500 Internal Server Error.
		assert.Equal(t, http.StatusInternalServerError, recorder.Code)
		// Verifica a mensagem de erro genérica do handler.
		assert.Contains(t, recorder.Body.String(), "Erro ao buscar clientes")
	})
}

// TestCustomerHandler_UpdateCustomer testa o endpoint PUT /api/customers/{id}.
// Verifica os cenários de sucesso na atualização, falha por ID inválido, JSON inválido,
// erro de validação do serviço, erro de cliente não encontrado pelo serviço e erro interno do serviço.
func TestCustomerHandler_UpdateCustomer(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockService := mock_service.NewMockCustomerService(mockCtrl)
	router, recorder := setupTestRouter(t, mockService)

	testID := uint(456)
	testIDStr := strconv.FormatUint(uint64(testID), 10)

	// Subteste para o cenário de sucesso.
	t.Run("Success", func(t *testing.T) {
		// Define os dados de atualização e o resultado esperado.
		customerUpdateInput := model.Customer{Name: "Updated Name", Email: "updated@example.com"}
		customerUpdateInputJSON, _ := json.Marshal(customerUpdateInput)
		// O cliente esperado na resposta deve ter o ID da URL e os dados do corpo.
		expectedUpdatedCustomer := model.Customer{ID: testID, Name: "Updated Name", Email: "updated@example.com"}

		// Define a expectativa: UpdateCustomer será chamado.
		mockService.EXPECT().
			UpdateCustomer(gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, c *model.Customer) error {
				// Verifica se o ID da URL foi corretamente atribuído ao objeto Customer passado para o serviço.
				assert.Equal(t, testID, c.ID)
				// Verifica se os dados do corpo foram corretamente passados.
				assert.Equal(t, customerUpdateInput.Name, c.Name)
				assert.Equal(t, customerUpdateInput.Email, c.Email)
				// Simula o sucesso da atualização no serviço (retorna nil).
				// O serviço pode ou não modificar 'c' ou retornar o objeto atualizado;
				// aqui assumimos que o handler retornará o objeto 'c' como ele foi passado (com ID setado).
				return nil
			}).Times(1)

		recorder = httptest.NewRecorder()
		// Executa a requisição PUT com o ID e o corpo JSON.
		performRequest(router, recorder, http.MethodPut, "/api/customers/"+testIDStr, customerUpdateInputJSON)

		// Verifica o status 200 OK.
		assert.Equal(t, http.StatusOK, recorder.Code)

		// Verifica o corpo da resposta.
		var updatedCustomer model.Customer
		err := json.Unmarshal(recorder.Body.Bytes(), &updatedCustomer)
		assert.NoError(t, err)
		// Compara o cliente retornado com o esperado.
		assert.Equal(t, expectedUpdatedCustomer, updatedCustomer)
	})

	// Subteste para o cenário de formato de ID inválido na URL.
	t.Run("Invalid ID Format", func(t *testing.T) {
		customerUpdateInput := model.Customer{Name: "Updated Name", Email: "updated@example.com"}
		customerUpdateInputJSON, _ := json.Marshal(customerUpdateInput)

		recorder = httptest.NewRecorder()
		// Executa a requisição com um ID não numérico.
		performRequest(router, recorder, http.MethodPut, "/api/customers/xyz", customerUpdateInputJSON)

		// Verifica o status 400 Bad Request.
		assert.Equal(t, http.StatusBadRequest, recorder.Code)
		assert.Contains(t, recorder.Body.String(), "ID inválido")
	})

	// Subteste para o cenário de JSON inválido no corpo da requisição.
	t.Run("Invalid JSON Body", func(t *testing.T) {
		recorder = httptest.NewRecorder()
		// Envia um corpo JSON malformado.
		performRequest(router, recorder, http.MethodPut, "/api/customers/"+testIDStr, []byte(`{"invalid`))

		// Verifica o status 400 Bad Request.
		assert.Equal(t, http.StatusBadRequest, recorder.Code)
		assert.Contains(t, recorder.Body.String(), "Dados inválidos")
	})

	// Subteste para o cenário onde o serviço retorna um erro de validação.
	t.Run("Service Validation Error", func(t *testing.T) {
		// Prepara uma entrada inválida (ex: nome vazio).
		customerUpdateInput := model.Customer{Name: ""}
		customerUpdateInputJSON, _ := json.Marshal(customerUpdateInput)

		// Define a expectativa: UpdateCustomer será chamado e retornará erro de validação.
		mockService.EXPECT().
			UpdateCustomer(gomock.Any(), gomock.Any()).
			Return(service.ErrInvalidCustomer).
			Times(1)

		recorder = httptest.NewRecorder()
		performRequest(router, recorder, http.MethodPut, "/api/customers/"+testIDStr, customerUpdateInputJSON)

		// Verifica o status 400 Bad Request.
		assert.Equal(t, http.StatusBadRequest, recorder.Code)
		assert.Contains(t, recorder.Body.String(), service.ErrInvalidCustomer.Error())
	})

	// Subteste para o cenário onde o serviço não encontra o cliente para atualizar.
	t.Run("Service Not Found Error", func(t *testing.T) {
		customerUpdateInput := model.Customer{Name: "Updated Name", Email: "updated@example.com"}
		customerUpdateInputJSON, _ := json.Marshal(customerUpdateInput)

		// Define a expectativa: UpdateCustomer será chamado e retornará ErrCustomerNotFound.
		mockService.EXPECT().
			UpdateCustomer(gomock.Any(), gomock.Any()).
			Return(service.ErrCustomerNotFound).
			Times(1)

		recorder = httptest.NewRecorder()
		performRequest(router, recorder, http.MethodPut, "/api/customers/"+testIDStr, customerUpdateInputJSON)

		// Verifica o status 404 Not Found.
		assert.Equal(t, http.StatusNotFound, recorder.Code)
		assert.Contains(t, recorder.Body.String(), service.ErrCustomerNotFound.Error())
	})

	// Subteste para o cenário de erro interno do serviço.
	t.Run("Service Internal Error", func(t *testing.T) {
		customerUpdateInput := model.Customer{Name: "Updated Name", Email: "updated@example.com"}
		customerUpdateInputJSON, _ := json.Marshal(customerUpdateInput)
		serviceErr := errors.New("database update failed") // Erro genérico simulado.

		// Define a expectativa: UpdateCustomer será chamado e retornará um erro interno.
		mockService.EXPECT().
			UpdateCustomer(gomock.Any(), gomock.Any()).
			Return(serviceErr).
			Times(1)

		recorder = httptest.NewRecorder()
		performRequest(router, recorder, http.MethodPut, "/api/customers/"+testIDStr, customerUpdateInputJSON)

		// Verifica o status 500 Internal Server Error.
		assert.Equal(t, http.StatusInternalServerError, recorder.Code)
		assert.Contains(t, recorder.Body.String(), "Erro ao atualizar cliente")
	})
}

// TestCustomerHandler_DeleteCustomer testa o endpoint DELETE /api/customers/{id}.
// Verifica os cenários de sucesso na exclusão, falha por ID inválido,
// falha quando o cliente não é encontrado e falha por erro interno do serviço.
func TestCustomerHandler_DeleteCustomer(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockService := mock_service.NewMockCustomerService(mockCtrl)
	router, recorder := setupTestRouter(t, mockService)

	testID := uint(789)
	testIDStr := strconv.FormatUint(uint64(testID), 10)

	// Subteste para o cenário de sucesso.
	t.Run("Success", func(t *testing.T) {
		// Define a expectativa: DeleteCustomer será chamado com o ID correto e retornará nil (sucesso).
		mockService.EXPECT().
			DeleteCustomer(gomock.Any(), testID).
			Return(nil). // Simula exclusão bem-sucedida.
			Times(1)

		recorder = httptest.NewRecorder()
		// Executa a requisição DELETE.
		performRequest(router, recorder, http.MethodDelete, "/api/customers/"+testIDStr, nil)

		// Verifica o status 204 No Content.
		assert.Equal(t, http.StatusNoContent, recorder.Code)
		// Verifica que o corpo da resposta está vazio.
		assert.Empty(t, recorder.Body.String())
	})

	// Subteste para o cenário de formato de ID inválido na URL.
	t.Run("Invalid ID Format", func(t *testing.T) {
		recorder = httptest.NewRecorder()
		// Executa a requisição com um ID não numérico.
		performRequest(router, recorder, http.MethodDelete, "/api/customers/invalid-id", nil)

		// Verifica o status 400 Bad Request.
		assert.Equal(t, http.StatusBadRequest, recorder.Code)
		assert.Contains(t, recorder.Body.String(), "ID inválido")
	})

	// Subteste para o cenário onde o cliente a ser excluído não é encontrado.
	t.Run("Not Found", func(t *testing.T) {
		// Define a expectativa: DeleteCustomer será chamado e retornará ErrCustomerNotFound.
		mockService.EXPECT().
			DeleteCustomer(gomock.Any(), testID).
			Return(service.ErrCustomerNotFound). // Simula cliente não encontrado.
			Times(1)

		recorder = httptest.NewRecorder()
		performRequest(router, recorder, http.MethodDelete, "/api/customers/"+testIDStr, nil)

		// Verifica o status 404 Not Found.
		assert.Equal(t, http.StatusNotFound, recorder.Code)
		assert.Contains(t, recorder.Body.String(), service.ErrCustomerNotFound.Error())
	})

	// Subteste para o cenário de erro interno do serviço.
	t.Run("Service Internal Error", func(t *testing.T) {
		serviceErr := errors.New("database delete failed") // Erro genérico simulado.

		// Define a expectativa: DeleteCustomer será chamado e retornará um erro interno.
		mockService.EXPECT().
			DeleteCustomer(gomock.Any(), testID).
			Return(serviceErr). // Simula erro interno.
			Times(1)

		recorder = httptest.NewRecorder()
		performRequest(router, recorder, http.MethodDelete, "/api/customers/"+testIDStr, nil)

		// Verifica o status 500 Internal Server Error.
		assert.Equal(t, http.StatusInternalServerError, recorder.Code)
		assert.Contains(t, recorder.Body.String(), "Erro ao excluir cliente")
	})
}

// TestCustomerHandler_CountCustomers testa o endpoint GET /api/customers/count.
// Verifica o cenário de sucesso ao contar os clientes e o cenário de erro interno do serviço.
func TestCustomerHandler_CountCustomers(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockService := mock_service.NewMockCustomerService(mockCtrl)
	router, recorder := setupTestRouter(t, mockService)

	// Subteste para o cenário de sucesso.
	t.Run("Success", func(t *testing.T) {
		expectedCount := int64(42) // Contagem esperada.

		// Define a expectativa: CountCustomers será chamado e retornará a contagem.
		mockService.EXPECT().
			CountCustomers(gomock.Any()).
			Return(expectedCount, nil). // Retorna a contagem e nenhum erro.
			Times(1)

		recorder = httptest.NewRecorder()
		performRequest(router, recorder, http.MethodGet, "/api/customers/count", nil)

		// Verifica o status 200 OK.
		assert.Equal(t, http.StatusOK, recorder.Code)

		// Verifica o corpo da resposta (espera-se um JSON como {"count": 42}).
		var responseBody map[string]int64
		err := json.Unmarshal(recorder.Body.Bytes(), &responseBody)
		assert.NoError(t, err)
		// Verifica se a chave "count" existe e tem o valor esperado.
		assert.Equal(t, expectedCount, responseBody["count"])
	})

	// Subteste para o cenário de erro interno do serviço.
	t.Run("Service Internal Error", func(t *testing.T) {
		serviceErr := errors.New("failed to count") // Erro genérico simulado.

		// Define a expectativa: CountCustomers será chamado e retornará um erro interno.
		mockService.EXPECT().
			CountCustomers(gomock.Any()).
			Return(int64(0), serviceErr). // Retorna 0 para contagem e o erro.
			Times(1)

		recorder = httptest.NewRecorder()
		performRequest(router, recorder, http.MethodGet, "/api/customers/count", nil)

		// Verifica o status 500 Internal Server Error.
		assert.Equal(t, http.StatusInternalServerError, recorder.Code)
		// Verifica a mensagem de erro genérica do handler.
		assert.Contains(t, recorder.Body.String(), "Erro ao contar clientes")
	})
}
