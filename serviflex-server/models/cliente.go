package models

import "time"

type Cliente struct {
	ID       string    `json:"id" firestore:"id"`
	Nome     string    `json:"nome" firestore:"nome"`
	Email    string    `json:"email" firestore:"email"`
	Senha    string    `json:"senha" firestore:"senha"`
	Telefone string    `json:"telefone" firestore:"telefone"`
	FotoURL  string    `json:"fotoUrl" firestore:"fotoURL"`
	CriadoEm time.Time `json:"criadoEm" firestore:"criadoEm"`
}
