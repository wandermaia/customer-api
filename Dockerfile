FROM golang:1.24.2-alpine AS builder

WORKDIR /app

# Instala o Swaggo
RUN go install github.com/swaggo/swag/cmd/swag@latest

# Copia arquivos de dependência
COPY go.mod go.sum ./
RUN go mod download

# Copia o restante do código
COPY . .

# Gera a documentação Swagger
RUN swag init -g cmd/api/main.go -o ./docs

# Comenta as linhas geradas no swagger que impedem a compilção.
# É um erro que veio no swagg e essa é uma solução de contorno.
#
RUN sed -i 's/LeftDelim/\/\/LeftDelim/g' docs/docs.go && sed -i 's/RightDelim/\/\/LeftDelim/g' docs/docs.go

# Compila a aplicação
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o customer-api ./cmd/api

# Imagem final
FROM alpine:latest  

WORKDIR /root/

# Copia o binário compilado e a documentação
COPY --from=builder /app/customer-api .
COPY --from=builder /app/docs ./docs

# Expõe a porta
EXPOSE 8080

# Comando para executar a aplicação
CMD ["./customer-api"]
