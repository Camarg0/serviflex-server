package models

import "time"

type Profissional struct {
	ID                string    `json:"id,omitempty" firestore:"id,omitempty"`
	Nome              string    `json:"nome" firestore:"nome"`
	Email             string    `json:"email" firestore:"email"`
	Senha             string    `json:"senha" firestore:"senha"`
	ImagemURL         string    `json:"imagem_url,omitempty" firestore:"imagem_url,omitempty"`
	Telefone          string    `json:"telefone,omitempty" firestore:"telefone,omitempty"`
	EstabelecimentoID string    `json:"estabelecimentoId,omitempty" firestore:"estabelecimentoId,omitempty"`
	CriadoEm          time.Time `json:"criadoEm" firestore:"criadoEm"`
}

type ProfissionalEstabelecimento struct {
	UID          string    `firestore:"uid"`
	Nome         string    `firestore:"nome"`
	Status       string    `firestore:"status"` // "ativo", "pendente"
	AdicionadoEm time.Time `firestore:"adicionadoEm"`
}
