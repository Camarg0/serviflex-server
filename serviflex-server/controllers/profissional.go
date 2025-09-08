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
)

// CriarHorario adiciona horários de atendimento para um profissional
// @Summary Criar horários
// @Tags Horários
// @Accept json
// @Produce json
// @Param horario body models.HorarioInput true "Horário de atendimento"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Router /horarios [post]
func CriarHorario(c *gin.Context) {
	var input models.HorarioInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos"})
		return
	}

	ctx := context.Background()
	client, err := config.App.Firestore(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro na conexão com Firestore"})
		return
	}
	defer client.Close()

	existentes, err := client.Collection("horarios").
		Where("profissional_id", "==", input.ProfissionalID).
		Documents(ctx).GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao verificar horários existentes"})
		return
	}

	diasJaCadastrados := make(map[string]bool)
	for _, doc := range existentes {
		data := doc.Data()
		if dia, ok := data["dia_semana"].(string); ok {
			diasJaCadastrados[dia] = true
		}
	}

	var criados []models.Horario
	var ignorados []string

	for _, dia := range input.DiasSemana {
		if diasJaCadastrados[dia] {
			ignorados = append(ignorados, dia)
			continue
		}

		horario := models.Horario{
			ID:             uuid.New().String(),
			ProfissionalID: input.ProfissionalID,
			DiaSemana:      dia,
			HoraInicio:     input.HoraInicio,
			HoraFim:        input.HoraFim,
			Disponivel:     true,
		}

		_, err := client.Collection("horarios").Doc(horario.ID).Set(ctx, horario)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao salvar horário"})
			return
		}

		criados = append(criados, horario)
	}

	c.JSON(http.StatusCreated, gin.H{
		"criados":   criados,
		"ignorados": ignorados,
	})
}

// ListarHorariosPorProfissional lista os horários de atendimento do profissional
// @Summary Listar horários por profissional
// @Tags Horários
// @Produce json
// @Param id path string true "ID do profissional"
// @Success 200 {array} models.Horario
// @Router /horarios/{id} [get]
func ListarHorariosPorProfissional(c *gin.Context) {
	profissionalID := c.Param("id")
	ctx := context.Background()
	client, err := config.App.Firestore(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro na conexão com Firestore"})
		return
	}
	defer client.Close()

	docs, err := client.Collection("horarios").Where("profissional_id", "==", profissionalID).Documents(ctx).GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar horários"})
		return
	}

	var horarios []models.Horario
	for _, doc := range docs {
		var h models.Horario
		if err := doc.DataTo(&h); err == nil {
			horarios = append(horarios, h)
		}
	}

	c.JSON(http.StatusOK, horarios)
}

// ListarAgendamentosPorProfissional lista os agendamentos de um profissional
// @Summary Listar agendamentos por profissional
// @Tags Agendamentos
// @Produce json
// @Param id path string true "ID do profissional"
// @Success 200 {array} models.Agendamento
// @Router /agendamentos/profissional/{id} [get]
func ListarAgendamentosPorProfissional(c *gin.Context) {
	profissionalID := c.Param("id")
	ctx := context.Background()
	client, err := config.App.Firestore(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro na conexão com Firestore"})
		return
	}
	defer client.Close()

	agDocs, err := client.Collection("agendamentos").
		Where("profissional_id", "==", profissionalID).
		OrderBy("data_hora", firestore.Asc).
		Documents(ctx).GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar agendamentos"})
		return
	}

	type AgendamentoComCliente struct {
		models.Agendamento
		ClienteNome string `json:"cliente_nome"`
	}

	var agendamentos []AgendamentoComCliente

	for _, doc := range agDocs {
		var ag models.Agendamento
		if err := doc.DataTo(&ag); err != nil {
			continue
		}

		clienteDoc, err := client.Collection("clientes").Doc(ag.ClienteID).Get(ctx)
		nomeCliente := ""
		if err == nil && clienteDoc.Exists() {
			dadosCliente := clienteDoc.Data()
			if nome, ok := dadosCliente["nome"].(string); ok {
				nomeCliente = nome
			}
		}

		agendamentos = append(agendamentos, AgendamentoComCliente{
			Agendamento: ag,
			ClienteNome: nomeCliente,
		})
	}

	c.JSON(http.StatusOK, agendamentos)
}

// ConvidarProfissional adiciona um profissional a um estabelecimento com notificação pendente
// @Summary Convidar profissional
// @Tags ProfissionaisEstabelecimento
// @Accept json
// @Produce json
// @Param dados body models.VinculoProfissionalInput true "IDs do profissional e estabelecimento"
// @Success 201 {object} map[string]string
// @Router /estabelecimentos/profissionais/convidar [post]
func ConvidarProfissional(c *gin.Context) {
	var input models.VinculoProfissionalInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos"})
		return
	}

	notificacao := models.Notificacao{
		ParaUID:           input.ProfissionalUID,
		Tipo:              "convite_estabelecimento",
		Mensagem:          "Você foi convidado para o estabelecimento",
		EstabelecimentoID: input.EstabelecimentoID,
		Respondido:        false,
		Resposta:          nil,
		CriadoEm:          time.Now(),
	}

	ctx := context.Background()
	client, err := config.App.Firestore(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro no Firestore"})
		return
	}
	defer client.Close()

	_, err = client.Collection("notificacoes").Doc(uuid.New().String()).Set(ctx, notificacao)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao enviar notificação"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Convite enviado"})
}

// AceitarOuRecusarConvite atualiza o status da notificação e vincula o profissional se aceito
// @Summary Aceitar ou recusar convite
// @Tags ProfissionaisEstabelecimento
// @Accept json
// @Produce json
// @Param id path string true "ID da notificação"
// @Param body body map[string]string true "resposta: aceito | recusado"
// @Success 200 {object} map[string]string
// @Router /estabelecimentos/profissionais/notificacao/{id} [post]
func AceitarOuRecusarConvite(c *gin.Context) {
	notifID := c.Param("id")
	var body map[string]string
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Requisição inválida"})
		return
	}

	resposta := body["resposta"]
	if resposta != "aceito" && resposta != "recusado" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Resposta inválida"})
		return
	}

	ctx := context.Background()
	client, err := config.App.Firestore(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro no Firestore"})
		return
	}
	defer client.Close()

	doc, err := client.Collection("notificacoes").Doc(notifID).Get(ctx)
	if err != nil || !doc.Exists() {
		c.JSON(http.StatusNotFound, gin.H{"error": "Notificação não encontrada"})
		return
	}

	var notif models.Notificacao
	doc.DataTo(&notif)

	client.Collection("notificacoes").Doc(notifID).Update(ctx, []firestore.Update{
		{Path: "respondido", Value: true},
		{Path: "resposta", Value: resposta},
	})

	if resposta == "aceito" {
		vinculo := models.ProfissionalEstabelecimento{
			UID:          notif.ParaUID,
			Nome:         "",
			Status:       "ativo",
			AdicionadoEm: time.Now(),
		}
		client.Collection("estabelecimentos").Doc(notif.EstabelecimentoID).
			Collection("profissionais").Doc(notif.ParaUID).Set(ctx, vinculo)

		client.Collection("profissionais").Doc(notif.ParaUID).Update(ctx, []firestore.Update{
			{Path: "estabelecimentoId", Value: notif.EstabelecimentoID},
		})
	}

	c.JSON(http.StatusOK, gin.H{"message": "Resposta registrada com sucesso"})
}

// ListarConvitesPendentes lista convites pendentes para um profissional
// @Summary Listar convites pendentes
// @Tags ProfissionaisEstabelecimento
// @Produce json
// @Param uid path string true "UID do profissional"
// @Success 200 {array} models.Notificacao
// @Router /profissionais/{uid}/convites-pendentes [get]
func ListarConvitesPendentes(c *gin.Context) {
    uid := c.Param("uid")
    ctx := context.Background()
    client, err := config.App.Firestore(ctx)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao conectar"})
        return
    }
    defer client.Close()

    docs, err := client.Collection("notificacoes").
        Where("paraUid", "==", uid).
        Where("tipo", "==", "conviteestabelecimento").
        Where("respondido", "==", false).
        Documents(ctx).GetAll()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar convites"})
        return
    }

    // Struct de resposta incluindo o ID da notificação
    type NotificacaoComID struct {
        ID                string      json:"id"
        models.Notificacao
    }

    var convites []NotificacaoComID
    for , doc := range docs {
        var n models.Notificacao
        if err := doc.DataTo(&n); err == nil {
            convites = append(convites, NotificacaoComID{
                ID:           doc.Ref.ID,
                Notificacao:  n,
            })
        }
    }

    c.JSON(http.StatusOK, convites)
}

// RemoverProfissional remove o vínculo de um profissional com o estabelecimento
// @Summary Remover profissional
// @Tags ProfissionaisEstabelecimento
// @Produce json
// @Param estId path string true "ID do estabelecimento"
// @Param profId path string true "ID do profissional"
// @Success 200 {object} map[string]string
// @Router /estabelecimentos/{estId}/profissionais/{profId} [delete]
func RemoverProfissional(c *gin.Context) {
	estID := c.Param("estId")
	profID := c.Param("profId")

	ctx := context.Background()
	client, err := config.App.Firestore(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro no Firestore"})
		return
	}
	defer client.Close()

	client.Collection("estabelecimentos").Doc(estID).Collection("profissionais").Doc(profID).Delete(ctx)

	client.Collection("profissionais").Doc(profID).Update(ctx, []firestore.Update{
		{Path: "estabelecimentoId", Value: ""},
	})

	c.JSON(http.StatusOK, gin.H{"message": "Profissional removido com sucesso"})
}

// ListarProfissionaisDoEstabelecimento lista os profissionais vinculados a um estabelecimento
// @Summary Listar profissionais do estabelecimento
// @Tags ProfissionaisEstabelecimento
// @Produce json
// @Param id path string true "ID do estabelecimento"
// @Success 200 {array} models.ProfissionalEstabelecimento
// @Router /estabelecimentos/{id}/profissionais [get]
func ListarProfissionaisDoEstabelecimento(c *gin.Context) {
	estID := c.Param("id")

	ctx := context.Background()
	client, err := config.App.Firestore(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao conectar"})
		return
	}
	defer client.Close()

	docs, err := client.Collection("estabelecimentos").Doc(estID).Collection("profissionais").Documents(ctx).GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar profissionais"})
		return
	}

	var lista []models.ProfissionalEstabelecimento
	for _, doc := range docs {
		var p models.ProfissionalEstabelecimento
		if err := doc.DataTo(&p); err == nil {
			lista = append(lista, p)
		}
	}

	c.JSON(http.StatusOK, lista)
}

func RelatorioFaturamentoProfissional(c *gin.Context) {
	profID := c.Param("id")

	ctx := context.Background()
	client, err := config.App.Firestore(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro na conexão com Firestore"})
		return
	}
	defer client.Close()

	// Buscar agendamentos do profissional
	agendamentos, err := client.Collection("agendamentos").
		Where("profissional_id", "==", profID).
		Documents(ctx).GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar agendamentos"})
		return
	}

	if len(agendamentos) == 0 {
		c.JSON(http.StatusOK, gin.H{"total_faturado": 0, "quantidade_agendamentos": 0})
		return
	}

	// Buscar todos os procedimentos do profissional
	procDocs, err := client.Collection("procedimentos").
		Where("profissional_id", "==", profID).
		Documents(ctx).GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar procedimentos"})
		return
	}

	precos := make(map[string]float64)
	for _, doc := range procDocs {
		var p models.Procedimento
		if err := doc.DataTo(&p); err == nil {
			precos[p.Nome] = p.Preco
		}
	}

	total := 0.0
	for _, ag := range agendamentos {
		nomeProc := ag.Data()["procedimento"].(string)
		if preco, ok := precos[nomeProc]; ok {
			total += preco
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"profissional_id":         profID,
		"quantidade_agendamentos": len(agendamentos),
		"total_faturado":          total,
	})
}
func RelatorioAvaliacoesPorProfissional(c *gin.Context) {
	profID := c.Param("id")

	ctx := context.Background()
	client, err := config.App.Firestore(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao conectar com Firestore"})
		return
	}
	defer client.Close()

	// Buscar todas as avaliações do profissional
	docs, err := client.Collection("avaliacoes").
		Where("profissional_id", "==", profID).
		Documents(ctx).GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar avaliações"})
		return
	}

	if len(docs) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"profissional_id":       profID,
			"quantidade_avaliacoes": 0,
			"media_nota":            0,
			"avaliacoes":            []models.Avaliacao{},
		})
		return
	}

	totalNotas := 0.0
	var avaliacoes []models.Avaliacao

	for _, doc := range docs {
		var a models.Avaliacao
		if err := doc.DataTo(&a); err == nil {
			totalNotas += float64(a.Nota)
			avaliacoes = append(avaliacoes, a)
		}
	}

	media := totalNotas / float64(len(avaliacoes))

	c.JSON(http.StatusOK, gin.H{
		"profissional_id":       profID,
		"quantidade_avaliacoes": len(avaliacoes),
		"media_nota":            media,
		"avaliacoes":            avaliacoes,
	})
}
func RelatorioAgendamentosPorMesProfissional(c *gin.Context) {
	profID := c.Param("id")

	ctx := context.Background()
	client, err := config.App.Firestore(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao conectar com Firestore"})
		return
	}
	defer client.Close()

	// Definir intervalo: últimos 12 meses
	agora := time.Now()
	inicio := agora.AddDate(0, -11, 0) // 11 meses atrás
	inicio = time.Date(inicio.Year(), inicio.Month(), 1, 0, 0, 0, 0, time.Local)

	// Buscar agendamentos do profissional dentro do intervalo
	docs, err := client.Collection("agendamentos").
		Where("profissional_id", "==", profID).
		Where("data_hora", ">=", inicio).
		Where("data_hora", "<=", agora).
		Documents(ctx).GetAll()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar agendamentos"})
		return
	}

	// Inicializar contagem mensal
	contagem := make(map[string]int)
	for i := 0; i < 12; i++ {
		mes := inicio.AddDate(0, i, 0)
		key := fmt.Sprintf("%04d-%02d", mes.Year(), mes.Month())
		contagem[key] = 0
	}

	// Contar agendamentos por mês
	for _, doc := range docs {
		var ag models.Agendamento
		if err := doc.DataTo(&ag); err == nil {
			data := ag.DataHora.In(time.Local)
			key := fmt.Sprintf("%04d-%02d", data.Year(), data.Month())
			if _, ok := contagem[key]; ok {
				contagem[key]++
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"profissional_id":      profID,
		"agendamentos_por_mes": contagem,
	})
}
