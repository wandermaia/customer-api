# API RESTful em Golang com Gin Framework para Gestão de Clientes

Repositório para o desafio final do bootcamp de Arquitetura de Sofltware, presente na pós graduação Arquitetura de Software e Inteligência Artificial ministrado pela XP Educação.



## Estrutura de Pastas (Padrão MVC)

```bash
/customer-api
  /cmd
    /api
      main.go              # Ponto de entrada da aplicação
  /internal
    /config
      config.go            # Configurações da aplicação
    /domain
      /model
        customer.go        # Modelo de dados do cliente
      /repository
        customer_repo.go   # Interface do repositório
        postgres_repo.go   # Implementação do repositório com PostgreSQL
      /service
        customer_service.go # Lógica de negócios
    /handler
      customer_handler.go  # Controladores da API REST
    /middleware
      logging.go           # Middleware para logging
      auth.go              # Middleware para autenticação (opcional)
    /utils
      error.go            # Utilitários para tratamento de erros
  /pkg
    /database
      postgres.go         # Conexão com o banco de dados
  /docs
    swagger.yaml          # Documentação da API (opcional)
  docker-compose.yml      # Configuração Docker
  Dockerfile              # Para containerização
  go.mod                  # Gerenciamento de dependências
  go.sum                  # Checksum das dependências
  README.md               # Documentação
```

## Explicação da Estrutura de Pastas e Componentes

### 1. Padrão MVC no Contexto da API

- **Model (Modelo)**: Representado pelos arquivos na pasta `/internal/domain/model` - define a estrutura de dados do cliente e suas validações.
- **View**: Em APIs RESTful, a "View" é substituída pela resposta JSON enviada ao cliente.
- **Controller**: Representado pelos handlers na pasta `/internal/handler` - recebe as requisições HTTP, chama os serviços adequados e formata as respostas.

### 2. Componentes Principais

- **Repositório**: Responsável pela persistência de dados, encapsula todas as operações de banco de dados.
- **Serviço**: Contém a lógica de negócios, validações e regras da aplicação.
- **Handler**: Gerencia as requisições HTTP, rotas e formata as respostas.
- **Middleware**: Funções executadas antes ou depois dos handlers para tarefas como logging e autenticação.
- **Config**: Configurações globais da aplicação.
- **Utils**: Funções utilitárias reutilizáveis.

go mod init github.com/wandermaia/customer-api