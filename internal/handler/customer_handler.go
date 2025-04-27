package handler

import (
	"net/http"
	"strconv"

	"github.com/wandermaia/customer-api/internal/domain/model"
	"github.com/wandermaia/customer-api/internal/domain/service"

	"github.com/gin-gonic/gin"
)

type CustomerHandler struct {
	service service.CustomerService
}

// NewCustomerHandler cria uma nova instância do handler de clientes
func NewCustomerHandler(service service.CustomerService) *CustomerHandler {
	return &CustomerHandler{
		service: service,
	}
}

// RegisterRoutes registra as rotas do handler no router do Gin
func (h *CustomerHandler) RegisterRoutes(router *gin.Engine) {
	customers := router.Group("/api/customers")
	{
		customers.POST("", h.CreateCustomer)
		customers.GET("", h.GetAllCustomers)
		customers.GET("/count", h.CountCustomers)
		customers.GET("/:id", h.GetCustomerByID)
		customers.GET("/search", h.GetCustomersByName)
		customers.PUT("/:id", h.UpdateCustomer)
		customers.DELETE("/:id", h.DeleteCustomer)
	}
}

// CreateCustomer cria um novo cliente
// @Summary Criar um novo cliente
// @Description Cria um novo cliente com os dados fornecidos
// @Tags customers
// @Accept json
// @Produce json
// @Param customer body model.Customer true "Dados do cliente"
// @Success 201 {object} model.Customer
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /customers [post]
func (h *CustomerHandler) CreateCustomer(c *gin.Context) {
	var customer model.Customer
	if err := c.ShouldBindJSON(&customer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos"})
		return
	}

	if err := h.service.CreateCustomer(c.Request.Context(), &customer); err != nil {
		if err == service.ErrInvalidCustomer {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao criar cliente"})
		return
	}

	c.JSON(http.StatusCreated, customer)
}

// GetCustomerByID busca um cliente pelo ID
// @Summary Buscar cliente por ID
// @Description Retorna os dados de um cliente específico com base no ID
// @Tags customers
// @Accept json
// @Produce json
// @Param id path int true "ID do Cliente"
// @Success 200 {object} model.Customer
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /customers/{id} [get]
func (h *CustomerHandler) GetCustomerByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	customer, err := h.service.GetCustomerByID(c.Request.Context(), uint(id))
	if err != nil {
		if err == service.ErrCustomerNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar cliente"})
		return
	}

	c.JSON(http.StatusOK, customer)
}

// GetAllCustomers retorna todos os clientes
// @Summary Listar todos os clientes
// @Description Retorna uma lista com todos os clientes cadastrados
// @Tags customers
// @Accept json
// @Produce json
// @Success 200 {array} model.Customer
// @Failure 500 {object} map[string]string
// @Router /customers [get]
func (h *CustomerHandler) GetAllCustomers(c *gin.Context) {
	customers, err := h.service.GetAllCustomers(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar clientes"})
		return
	}

	c.JSON(http.StatusOK, customers)
}

// GetCustomersByName busca clientes pelo nome
// @Summary Buscar clientes por nome
// @Description Retorna uma lista de clientes que correspondem ao nome fornecido
// @Tags customers
// @Accept json
// @Produce json
// @Param name query string true "Nome do cliente"
// @Success 200 {array} model.Customer
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /customers/search [get]
func (h *CustomerHandler) GetCustomersByName(c *gin.Context) {
	name := c.Query("name")
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nome não fornecido"})
		return
	}

	customers, err := h.service.GetCustomersByName(c.Request.Context(), name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar clientes"})
		return
	}

	c.JSON(http.StatusOK, customers)
}

// UpdateCustomer atualiza um cliente existente
// @Summary Atualizar cliente
// @Description Atualiza os dados de um cliente existente
// @Tags customers
// @Accept json
// @Produce json
// @Param id path int true "ID do Cliente"
// @Param customer body model.Customer true "Dados atualizados do cliente"
// @Success 200 {object} model.Customer
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /customers/{id} [put]
func (h *CustomerHandler) UpdateCustomer(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var customer model.Customer
	if err := c.ShouldBindJSON(&customer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos"})
		return
	}

	customer.ID = uint(id)

	if err := h.service.UpdateCustomer(c.Request.Context(), &customer); err != nil {
		if err == service.ErrInvalidCustomer {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err == service.ErrCustomerNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao atualizar cliente"})
		return
	}

	c.JSON(http.StatusOK, customer)
}

// DeleteCustomer remove um cliente pelo ID
// @Summary Excluir cliente
// @Description Remove um cliente do sistema com base no ID
// @Tags customers
// @Accept json
// @Produce json
// @Param id path int true "ID do Cliente"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /customers/{id} [delete]
func (h *CustomerHandler) DeleteCustomer(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	if err := h.service.DeleteCustomer(c.Request.Context(), uint(id)); err != nil {
		if err == service.ErrCustomerNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao excluir cliente"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// CountCustomers retorna o número total de clientes
// @Summary Contar clientes
// @Description Retorna o número total de clientes cadastrados no sistema
// @Tags customers
// @Accept json
// @Produce json
// @Success 200 {object} utils.CountResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /customers/count [get]
func (h *CustomerHandler) CountCustomers(c *gin.Context) {
	count, err := h.service.CountCustomers(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao contar clientes"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"count": count})
}
