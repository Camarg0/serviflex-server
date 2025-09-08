package controllers

import (
	"context"
	"fmt"
	"net/http"
	"servico-api/config"
	"servico-api/models"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
) // adicione no topo se ainda não tiver

// AgendarHorario cria um agendamento entre cliente e profissional
// @Summary Agendar horário
// @Tags Cliente
// @Accept json
// @Produce json
// @Param agendamento body models.Agendamento true "Dados do agendamento"
// @Success 201 {object} models.Agendamento
// @Router /agendamentos [post]
func AgendarHorario(c *gin.Context) {
	var agendamento models.Agendamento
	if err := c.ShouldBindJSON(&agendamento); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos"})
		return
	}

	ctx := context.Background()
	client, err := config.App.Firestore(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao conectar com Firestore"})
		return
	}
	defer client.Close()

	// Busca procedimento
	procDocs, err := client.Collection("procedimentos").
		Where("profissional_id", "==", agendamento.ProfissionalID).
		Where("nome", "==", agendamento.Procedimento).
		Limit(1).Documents(ctx).GetAll()

	if err != nil || len(procDocs) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Procedimento inválido"})
		return
	}

	duracaoMin := int(procDocs[0].Data()["duracao_min"].(int64))
	duracao := time.Duration(duracaoMin) * time.Minute

	// Corrigir o fuso e descobrir dia da semana
	localDateTime := agendamento.DataHora.In(time.Local)
	diaSemana := diaDaSemana(localDateTime.Weekday())
	fmt.Println("Dia da semana calculado:", diaSemana)

	// Buscar horários do profissional no Firestore
	horariosSnap, err := client.Collection("horarios").
		Where("profissional_id", "==", agendamento.ProfissionalID).
		Where("dia_semana", "==", diaSemana).
		Documents(ctx).GetAll()

	if err != nil {
		fmt.Println("Erro ao buscar horários:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar horários"})
		return
	}

	if len(horariosSnap) == 0 {
		fmt.Println("Nenhum horário encontrado para", agendamento.ProfissionalID, "no dia", diaSemana)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Profissional não trabalha nesse dia"})
		return
	}

	fmt.Println("Horários encontrados:")
	for _, h := range horariosSnap {
		fmt.Println(h.Data())
	}

	// Validar se o horário está dentro do expediente
	agHora := agendamento.DataHora.In(time.Local)
	inicioAg := time.Date(0, 1, 1, agHora.Hour(), agHora.Minute(), 0, 0, time.UTC)
	fimAg := inicioAg.Add(duracao)

	valid := false
	for _, h := range horariosSnap {
		data := h.Data()
		hIniStr := data["hora_inicio"].(string)
		hFimStr := data["hora_fim"].(string)

		hIni, _ := time.Parse("15:04", hIniStr)
		hFim, _ := time.Parse("15:04", hFimStr)

		if (inicioAg.Equal(hIni) || inicioAg.After(hIni)) && (fimAg.Equal(hFim) || fimAg.Before(hFim)) {
			valid = true
			break
		}
	}

	if !valid {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Horário não está dentro do expediente do profissional"})
		return
	}

	// Criar agendamento
	agendamento.ID = uuid.New().String()
	client.Collection("agendamentos").Doc(agendamento.ID).Set(ctx, agendamento)
	c.JSON(http.StatusCreated, agendamento)
}

// ListarAgendamentosPorCliente retorna todos os agendamentos do cliente ordenados por data
// @Summary Listar agendamentos do cliente
// @Tags Agendamentos
// @Produce json
// @Param id path string true "ID do cliente"
// @Success 200 {array} models.Agendamento
// @Router /agendamentos/cliente/{id} [get]
func ListarAgendamentosPorCliente(c *gin.Context) {
	clienteID := c.Param("id")
	ctx := context.Background()
	client, err := config.App.Firestore(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro na conexão com Firestore"})
		return
	}
	defer client.Close()

	docs, err := client.Collection("agendamentos").
		Where("cliente_id", "==", clienteID).
		OrderBy("data_hora", firestore.Asc).
		Documents(ctx).GetAll()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar agendamentos"})
		return
	}

	var agendamentos []models.Agendamento
	for _, doc := range docs {
		var ag models.Agendamento
		if err := doc.DataTo(&ag); err == nil {
			agendamentos = append(agendamentos, ag)
		}
	}

	c.JSON(http.StatusOK, agendamentos)
}

// EditarUsuario permite que qualquer tipo de usuário atualize seus dados básicos
// @Summary Editar perfil do usuário
// @Tags Usuários
// @Accept json
// @Produce json
// @Param id path string true "ID do usuário"
// @Param tipo query string true "Tipo do usuário (clientes, profissionais, admin)"
// @Param usuario body interface{} true "Dados atualizados do usuário"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /usuarios/{id} [put]
func EditarUsuario(c *gin.Context) {
	id := c.Param("id")
	tipo := c.Query("tipo") // Recebe o tipo como query parameter

	if tipo == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parâmetro 'tipo' é obrigatório (clientes, profissionais, admin)"})
		return
	}

	if tipo != "clientes" && tipo != "profissionais" && tipo != "admin" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Tipo inválido. Use: clientes, profissionais ou admin"})
		return
	}

	ctx := context.Background()
	client, err := config.App.Firestore(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao conectar com Firestore"})
		return
	}
	defer client.Close()

	// Verifica se o usuário existe na collection específica
	docRef := client.Collection(tipo).Doc(id)
	snap, err := docRef.Get(ctx)
	if err != nil || !snap.Exists() {
		c.JSON(http.StatusNotFound, gin.H{"error": "Usuário não encontrado na collection " + tipo})
		return
	}

	// Processa a atualização baseada no tipo
	switch tipo {
	case "clientes":
		var input models.Cliente
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos para cliente"})
			return
		}
		input.ID = id // Mantém o ID original
		if _, err := docRef.Set(ctx, input); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao atualizar cliente"})
			return
		}

	case "profissionais":
		var input models.Profissional
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos para profissional"})
			return
		}
		input.ID = id // Mantém o ID original
		if _, err := docRef.Set(ctx, input); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao atualizar profissional"})
			return
		}

	case "admin":
		var input models.Admin
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos para admin"})
			return
		}
		input.ID = id // Mantém o ID original
		if _, err := docRef.Set(ctx, input); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao atualizar admin"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Usuário atualizado com sucesso"})
}

func diaDaSemana(weekday time.Weekday) string {
	switch weekday {
	case time.Monday:
		return "Segunda"
	case time.Tuesday:
		return "Terça"
	case time.Wednesday:
		return "Quarta"
	case time.Thursday:
		return "Quinta"
	case time.Friday:
		return "Sexta"
	case time.Saturday:
		return "Sábado"
	case time.Sunday:
		return "Domingo"
	default:
		return ""
	}
}
