#!/bin/bash

# Instala o Swaggo se ainda não estiver instalado
if ! command -v swag &> /dev/null
then
    echo "Instalando Swaggo..."
    go install github.com/swaggo/swag/cmd/swag@latest
fi

# Gera a documentação Swagger
echo "Gerando documentação Swagger..."
swag init -g cmd/api/main.go -o ./docs

echo "Documentação Swagger gerada com sucesso!"
echo "Acesse http://localhost:8080/swagger/index.html após iniciar o servidor"
