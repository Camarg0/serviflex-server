package models

import "time"

type Notificacao struct {
	ParaUID           string    `firestore:"paraUid"`
	Tipo              string    `firestore:"tipo"` // "convite_estabelecimento"
	Mensagem          string    `firestore:"mensagem"`
	EstabelecimentoID string    `firestore:"estabelecimentoId"`
	Respondido        bool      `firestore:"respondido"`
	Resposta          *string   `firestore:"resposta"` // nil, "aceito", "recusado"
	CriadoEm          time.Time `firestore:"criadoEm"`
}
