package models

import "time"

type Estabelecimento struct {
	Nome           string    `firestore:"nome"`
	Descricao      string    `firestore:"descricao"`
	FotoURL        string    `firestore:"fotoURL"`
	Categoria      string    `firestore:"categoria"`
	Localizacao    Endereco  `firestore:"localizacao"`
	CriadoEm       time.Time `firestore:"criadoEm"`
	ResponsavelUID string    `firestore:"responsavelUid"`
}

type Endereco struct {
	Endereco string `firestore:"endereco"`
	Cidade   string `firestore:"cidade"`
	UF       string `firestore:"uf"`
}

type EstabelecimentoInput struct {
	Nome        string   `json:"nome" binding:"required"`
	Descricao   string   `json:"descricao"`
	FotoURL     string   `json:"fotoURL"`
	Categoria   string   `json:"categoria"`
	Localizacao Endereco `json:"localizacao" binding:"required"`
}

// VinculoProfissionalInput representa a requisição para adicionar um profissional a um estabelecimento.
type VinculoProfissionalInput struct {
	EstabelecimentoID string `json:"estabelecimento_id" binding:"required"` // ID do estabelecimento
	ProfissionalUID   string `json:"profissional_uid" binding:"required"`   // UID do profissional a ser convidado
}

// RespostaConviteInput representa a resposta de um profissional a um convite
type RespostaConviteInput struct {
	ConviteID string `json:"convite_id" binding:"required"` // ID do convite no Firestore
	Aceito    bool   `json:"aceito"`                        // true para aceitar, false para recusar
}
