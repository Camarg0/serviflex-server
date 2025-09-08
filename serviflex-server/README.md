# Servico API

API para gerenciamento de agendamentos entre clientes e profissionais, como um "Ifood de serviços".

## Tecnologias Utilizadas

- Go (Golang)
- Gin (framework web)
- Firebase Firestore (banco de dados)
- Swagger (documentação)
- Postman (testes de API)

## Estrutura de Pastas

servico-api/  
├── config/             # Configuração do Firebase  
├── controllers/        # Handlers da API  
├── models/             # Modelos de dados  
├── routes/             # Organização das rotas  
├── utils/              # Funções auxiliares (ex: login)  
├── main.go             # Entry point  
├── go.mod / go.sum     # Dependências  
└── docs/               # Gerado pelo swag init

## Como Rodar o Projeto

go mod tidy  
export GOOGLE_APPLICATION_CREDENTIALS="caminho/para/sua/credencial.json"  
go run main.go

## Gerar Documentação Swagger

go install github.com/swaggo/swag/cmd/swag@latest  
swag init  
Acesse: http://localhost:8080/swagger/index.html

## Endpoints Organizados

### Autenticação

- POST /api/login – Login simples  
- POST /api/cadastro – Cadastro de cliente ou profissional

### Cliente

- GET /api/estabelecimentos  
- POST /api/agendamentos  
- GET /api/agendamentos/cliente/:id

### Profissional

- GET /api/agendamentos/profissional/:id  
- GET /api/horarios/:id

### Horários

- POST /api/horarios

### Procedimentos

- POST /api/procedimentos  
- GET /api/procedimentos/:id  
- PUT /api/procedimentos/:id  
- DELETE /api/procedimentos/:id

### Upload

- PUT /api/upload/{tipo}/{id} (tipo: profissional ou procedimento)

## Importar no Postman

Coleção Postman gerada no formato JSON (ver próximo bloco).

## Notas

- Sem autenticação JWT, uso apenas de login manual.  
- Firestore precisa de índices compostos para certas queries.  
- Suporte a imagens via URL salva no Firestore.
