package models

type Horario struct {
	ID             string `json:"id,omitempty" firestore:"id,omitempty"`
	ProfissionalID string `json:"profissional_id" firestore:"profissional_id"`
	DiaSemana      string `json:"dia_semana" firestore:"dia_semana"`
	HoraInicio     string `json:"hora_inicio" firestore:"hora_inicio"`
	HoraFim        string `json:"hora_fim" firestore:"hora_fim"`
	Disponivel     bool   `json:"disponivel" firestore:"disponivel"`
}
type HorarioInput struct {
	ProfissionalID string   `json:"profissional_id"`
	DiasSemana     []string `json:"dias_semana"`
	HoraInicio     string   `json:"hora_inicio"`
	HoraFim        string   `json:"hora_fim"`
}
