package models

import "time"

type Avaliacao struct {
	ProfissionalID    string    `firestore:"profissionalId"`
	ClienteID         string    `firestore:"clienteId"`
	EstabelecimentoID string    `firestore:"estabelecimentoId"`
	Nota              float64   `firestore:"nota"`
	Comentario        string    `firestore:"comentario"`
	Data              time.Time `firestore:"data"`
}
