// Adicione ao arquivo internal/utils/error.go

package utils

// ErrorResponse modelo para respostas de erro
// @Description Modelo padrão para respostas de erro da API
type ErrorResponse struct {
	Error string `json:"error" example:"Mensagem de erro"`
}

// CountResponse modelo para resposta de contagem
// @Description Modelo para resposta de contagem de registros
type CountResponse struct {
	Count int64 `json:"count" example:"42"`
}
