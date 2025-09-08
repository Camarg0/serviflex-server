package models

import "time"

type Agendamento struct {
	ID                string    `firestore:"id"`
	ClienteID         string    `firestore:"clienteId"`
	ProfissionalID    string    `firestore:"profissionalId"`
	EstabelecimentoID string    `firestore:"estabelecimentoId"`
	Procedimento      string    `firestore:"procedimento"` // nome do procedimento
	DataHora          time.Time `firestore:"dataHora"`
}
