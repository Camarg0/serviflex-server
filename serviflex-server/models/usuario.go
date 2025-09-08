package models

import "time"

type Usuario struct {
	ID                string    `json:"id" firestore:"id"`
	Nome              string    `json:"nome" firestore:"nome"`
	Email             string    `json:"email" firestore:"email"`
	Senha             string    `json:"senha" firestore:"senha"`
	Tipo              string    `json:"tipo,omitempty"`
	Telefone          string    `json:"telefone" firestore:"telefone"`
	FotoURL           string    `json:"fotoUrl" firestore:"fotoURL"`
	EstabelecimentoID string    `json:"estabelecimentoId" firestore:"estabelecimentoId,omitempty"` // só se for profissional
	CriadoEm          time.Time `json:"criadoEm" firestore:"criadoEm"`
}

type Login struct {
	Email string `json:"email"`
	Senha string `json:"senha"`
}
