package controllers

import (
	"context"
	"fmt"
	"net/http"
	"servico-api/config"
	"servico-api/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// CriarEstabelecimento cria um novo estabelecimento
// @Summary Criar estabelecimento
// @Tags Estabelecimentos
// @Accept json
// @Produce json
// @Param estabelecimento body models.EstabelecimentoInput true "Dados do estabelecimento"
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /estabelecimentos [post]
func CriarEstabelecimento(c *gin.Context) {
	var input models.EstabelecimentoInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos"})
		return
	}

	estabID := uuid.New().String()
	uid := c.GetString("uid") // você pode usar middleware para injetar o uid do auth

	estab := models.Estabelecimento{
		Nome:           input.Nome,
		Descricao:      input.Descricao,
		FotoURL:        input.FotoURL,
		Categoria:      input.Categoria,
		Localizacao:    input.Localizacao,
		CriadoEm:       time.Now(),
		ResponsavelUID: uid,
	}

	ctx := context.Background()
	client, err := config.App.Firestore(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao conectar ao Firestore"})
		return
	}
	defer client.Close()

	_, err = client.Collection("estabelecimentos").Doc(estabID).Set(ctx, estab)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao criar estabelecimento"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": estabID})
}

// EditarEstabelecimento atualiza os dados de um estabelecimento
// @Summary Editar estabelecimento
// @Tags Estabelecimentos
// @Accept json
// @Produce json
// @Param id path string true "ID do estabelecimento"
// @Param estabelecimento body models.EstabelecimentoInput true "Dados atualizados"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /estabelecimentos/{id} [put]
func EditarEstabelecimento(c *gin.Context) {
	id := c.Param("id")

	var input models.EstabelecimentoInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos"})
		return
	}

	ctx := context.Background()
	client, err := config.App.Firestore(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao conectar ao Firestore"})
		return
	}
	defer client.Close()

	// Garante que o estabelecimento existe
	docRef := client.Collection("estabelecimentos").Doc(id)
	_, err = docRef.Get(ctx)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Estabelecimento não encontrado"})
		return
	}

	update := models.Estabelecimento{
		Nome:           input.Nome,
		Descricao:      input.Descricao,
		FotoURL:        input.FotoURL,
		Categoria:      input.Categoria,
		Localizacao:    input.Localizacao,
		CriadoEm:       time.Now(),
		ResponsavelUID: c.GetString("uid"), // opcionalmente, pode manter o antigo
	}

	if _, err := docRef.Set(ctx, update); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao atualizar"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Estabelecimento atualizado"})
}

// ListarEstabelecimentos retorna todos os estabelecimentos
// @Summary Listar estabelecimentos
// @Tags Estabelecimentos
// @Produce json
// @Success 200 {array} models.Estabelecimento
// @Router /estabelecimentos [get]
func ListarEstabelecimentos(c *gin.Context) {
	ctx := context.Background()
	client, err := config.App.Firestore(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao conectar"})
		return
	}
	defer client.Close()

	docs, err := client.Collection("estabelecimentos").Documents(ctx).GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar dados"})
		return
	}

	var estabelecimentos []map[string]interface{}
	for _, doc := range docs {
		var e models.Estabelecimento
		if err := doc.DataTo(&e); err == nil {
			// Criar um mapa que inclui o ID do documento
			estabelecimentoComID := map[string]interface{}{
				"id":             doc.Ref.ID,
				"nome":           e.Nome,
				"descricao":      e.Descricao,
				"fotoURL":        e.FotoURL,
				"categoria":      e.Categoria,
				"localizacao":    e.Localizacao,
				"criadoEm":       e.CriadoEm,
				"responsavelUid": e.ResponsavelUID,
			}
			estabelecimentos = append(estabelecimentos, estabelecimentoComID)
		}
	}

	c.JSON(http.StatusOK, estabelecimentos)
}

// BuscarEstabelecimentoPorID retorna um estabelecimento por ID
// @Summary Buscar estabelecimento por ID
// @Tags Estabelecimentos
// @Produce json
// @Param id path string true "ID do estabelecimento"
// @Success 200 {object} models.Estabelecimento
// @Failure 404 {object} map[string]string
// @Router /estabelecimentos/{id} [get]
func BuscarEstabelecimentoPorID(c *gin.Context) {
	id := c.Param("id")
	ctx := context.Background()
	client, err := config.App.Firestore(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao conectar"})
		return
	}
	defer client.Close()

	doc, err := client.Collection("estabelecimentos").Doc(id).Get(ctx)
	if err != nil || !doc.Exists() {
		c.JSON(http.StatusNotFound, gin.H{"error": "Estabelecimento não encontrado"})
		return
	}

	var e models.Estabelecimento
	doc.DataTo(&e)

	// Retornar com ID incluído
	estabelecimentoComID := map[string]interface{}{
		"id":             doc.Ref.ID,
		"nome":           e.Nome,
		"descricao":      e.Descricao,
		"fotoURL":        e.FotoURL,
		"categoria":      e.Categoria,
		"localizacao":    e.Localizacao,
		"criadoEm":       e.CriadoEm,
		"responsavelUid": e.ResponsavelUID,
	}

	c.JSON(http.StatusOK, estabelecimentoComID)
}
func RelatorioFaturamentoEstabelecimento(c *gin.Context) {
	estabID := c.Param("id")

	ctx := context.Background()
	client, err := config.App.Firestore(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao conectar com Firestore"})
		return
	}
	defer client.Close()

	// Buscar todos os agendamentos do estabelecimento
	agendamentos, err := client.Collection("agendamentos").
		Where("estabelecimento_id", "==", estabID).
		Documents(ctx).GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar agendamentos"})
		return
	}

	if len(agendamentos) == 0 {
		c.JSON(http.StatusOK, gin.H{"total_faturado": 0, "quantidade_agendamentos": 0})
		return
	}

	// Buscar todos os procedimentos do sistema (ou filtrar por estabelecimento se tiver esse vínculo)
	procDocs, err := client.Collection("procedimentos").Documents(ctx).GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar procedimentos"})
		return
	}

	// Mapa com chave combinando profissional_id + nome do procedimento
	procMap := make(map[string]float64)
	for _, doc := range procDocs {
		var p models.Procedimento
		if err := doc.DataTo(&p); err == nil {
			key := fmt.Sprintf("%s|%s", p.ProfissionalID, p.Nome)
			procMap[key] = p.Preco
		}
	}

	total := 0.0
	for _, ag := range agendamentos {
		profID := ag.Data()["profissional_id"].(string)
		nomeProc := ag.Data()["procedimento"].(string)
		key := fmt.Sprintf("%s|%s", profID, nomeProc)
		if preco, ok := procMap[key]; ok {
			total += preco
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"estabelecimento_id":      estabID,
		"quantidade_agendamentos": len(agendamentos),
		"total_faturado":          total,
	})
}
func RelatorioAvaliacoesPorEstabelecimento(c *gin.Context) {
	estabID := c.Param("id")

	ctx := context.Background()
	client, err := config.App.Firestore(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao conectar com Firestore"})
		return
	}
	defer client.Close()

	// Buscar todas as avaliações do estabelecimento
	docs, err := client.Collection("avaliacoes").
		Where("estabelecimento_id", "==", estabID).
		Documents(ctx).GetAll()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar avaliações"})
		return
	}

	if len(docs) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"estabelecimento_id":    estabID,
			"quantidade_avaliacoes": 0,
			"media_nota":            0,
			"avaliacoes":            []models.Avaliacao{},
		})
		return
	}

	var avaliacoes []models.Avaliacao
	var total float64 = 0

	for _, doc := range docs {
		var a models.Avaliacao
		if err := doc.DataTo(&a); err == nil {
			total += float64(a.Nota)
			avaliacoes = append(avaliacoes, a)
		}
	}

	media := total / float64(len(avaliacoes))

	c.JSON(http.StatusOK, gin.H{
		"estabelecimento_id":    estabID,
		"quantidade_avaliacoes": len(avaliacoes),
		"media_nota":            media,
		"avaliacoes":            avaliacoes,
	})
}
func RelatorioAgendamentosPorMesEstabelecimento(c *gin.Context) {
	estabID := c.Param("id")

	ctx := context.Background()
	client, err := config.App.Firestore(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao conectar com Firestore"})
		return
	}
	defer client.Close()

	// Últimos 12 meses
	agora := time.Now()
	inicio := agora.AddDate(0, -11, 0)
	inicio = time.Date(inicio.Year(), inicio.Month(), 1, 0, 0, 0, 0, time.Local)

	// Buscar agendamentos do estabelecimento
	docs, err := client.Collection("agendamentos").
		Where("estabelecimento_id", "==", estabID).
		Where("data_hora", ">=", inicio).
		Where("data_hora", "<=", agora).
		Documents(ctx).GetAll()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar agendamentos"})
		return
	}

	// Inicializa contagem
	contagem := make(map[string]int)
	for i := 0; i < 12; i++ {
		mes := inicio.AddDate(0, i, 0)
		key := fmt.Sprintf("%04d-%02d", mes.Year(), mes.Month())
		contagem[key] = 0
	}

	// Conta agendamentos
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
		"estabelecimento_id":   estabID,
		"agendamentos_por_mes": contagem,
	})
}
