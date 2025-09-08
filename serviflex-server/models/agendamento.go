package models

import "time"

type Agendamento struct {
    ID                string    `firestore:"id" json:"id"`
    ClienteID         string    `firestore:"clienteId" json:"cliente_id"`
    ProfissionalID    string    `firestore:"profissionalId" json:"profissional_id"`
    EstabelecimentoID string    `firestore:"estabelecimentoId" json:"estabelecimento_id"`
    Procedimento      string    `firestore:"procedimento" json:"procedimento"`
    DataHora          time.Time `firestore:"dataHora" json:"data_hora"`
}
