package models

import "time"

type Admin struct {
	ID       string    `json:"id" firestore:"id"`
	Nome     string    `json:"nome" firestore:"nome"`
	Email    string    `json:"email" firestore:"email"`
	Senha    string    `json:"senha" firestore:"senha"`
	CriadoEm time.Time `json:"criadoEm" firestore:"criadoEm"`
	Tipo	 string    `json:"tipo" firestore:"tipo"`
}

