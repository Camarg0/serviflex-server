package seed

import (
	"context"
	"log"
	"servico-api/config"
	"servico-api/models"
	"time"

	"github.com/google/uuid"
)

func SeedFirestore() {
	ctx := context.Background()
	client, err := config.App.Firestore(ctx)
	if err != nil {
		log.Fatalf("Erro ao conectar ao Firestore: %v", err)
	}
	defer client.Close()

	// Limpa dados antigos para facilitar testes (opcional)
	collections := []string{"estabelecimentos", "profissionais", "clientes", "horarios", "procedimentos", "agendamentos"}
	for _, col := range collections {
		docs, _ := client.Collection(col).Documents(ctx).GetAll()
		for _, doc := range docs {
			_, _ = doc.Ref.Delete(ctx)
		}
	}

	// Cria um estabelecimento
	estabID := uuid.New().String()
	estabelecimento := models.Estabelecimento{
		Nome:           "Studio da Beleza",
		Descricao:      "Espaço completo de estética e bem-estar",
		FotoURL:        "https://exemplo.com/foto.jpg",
		Categoria:      "Beleza",
		CriadoEm:       time.Now(),
		ResponsavelUID: "uid-admin-123",
		Localizacao: models.Endereco{
			Endereco: "Rua das Flores, 123",
			Cidade:   "Uberlândia",
			UF:       "MG",
		},
	}
	_, _ = client.Collection("estabelecimentos").Doc(estabID).Set(ctx, estabelecimento)

	// Cria um profissional vinculado a esse estabelecimento
	profID := uuid.New().String()
	profissional := models.Profissional{
		ID:                profID,
		Nome:              "Maria Silva",
		Email:             "maria@exemplo.com",
		Senha:             "123456",
		ImagemURL:         "https://exemplo.com/maria.jpg",
		Telefone:          "(34) 99999-9999",
		EstabelecimentoID: estabID,
		CriadoEm:          time.Now(),
	}
	_, _ = client.Collection("profissionais").Doc(profID).Set(ctx, profissional)

	// Cria um cliente para testes
	clienteID := uuid.New().String()
	cliente := models.Cliente{
		ID:       clienteID,
		Nome:     "João Cliente",
		Email:    "joao@cliente.com",
		Senha:    "123456",
		Telefone: "(34) 98888-7777",
		FotoURL:  "https://exemplo.com/joao.jpg",
		CriadoEm: time.Now(),
	}
	_, _ = client.Collection("clientes").Doc(clienteID).Set(ctx, cliente)

	log.Println("✅ Seed executado com sucesso!")
}
