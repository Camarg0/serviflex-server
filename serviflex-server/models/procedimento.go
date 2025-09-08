package models

type Procedimento struct {
	ID             string  `json:"id,omitempty" firestore:"id,omitempty"`
	ProfissionalID string  `json:"profissional_id" firestore:"profissional_id"`
	Nome           string  `json:"nome" firestore:"nome"`
	Descricao      string  `json:"descricao" firestore:"descricao"`
	Preco          float64 `json:"preco" firestore:"preco"`
	DuracaoMin     int     `json:"duracao_min" firestore:"duracao_min"`
	ImagemURL      string  `json:"imagem_url,omitempty" firestore:"imagem_url,omitempty"`
}
