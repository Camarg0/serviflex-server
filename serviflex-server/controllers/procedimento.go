package controllers

import (
	"context"
	"net/http"
	"servico-api/config"
	"servico-api/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// CriarProcedimento adiciona novo procedimento para um profissional
// @Summary Criar procedimento
// @Tags Procedimentos
// @Accept json
// @Produce json
// @Param procedimento body models.Procedimento true "Procedimento"
// @Success 201 {object} models.Procedimento
// @Router /procedimentos [post]
func CriarProcedimento(c *gin.Context) {
	var proc models.Procedimento
	if err := c.ShouldBindJSON(&proc); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos"})
		return
	}

	proc.ID = uuid.New().String()

	ctx := context.Background()
	client, err := config.App.Firestore(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao conectar ao Firestore"})
		return
	}
	defer client.Close()

	_, err = client.Collection("procedimentos").Doc(proc.ID).Set(ctx, proc)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao salvar procedimento"})
		return
	}

	c.JSON(http.StatusCreated, proc)
}

// ListarProcedimentosPorProfissional lista todos os procedimentos de um profissional
// @Summary Listar procedimentos por profissional
// @Tags Procedimentos
// @Produce json
// @Param id path string true "ID do profissional"
// @Success 200 {array} models.Procedimento
// @Router /procedimentos/{id} [get]
func ListarProcedimentosPorProfissional(c *gin.Context) {
	profID := c.Param("id")
	ctx := context.Background()
	client, err := config.App.Firestore(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro na conexão com Firestore"})
		return
	}
	defer client.Close()

	docs, err := client.Collection("procedimentos").
		Where("profissional_id", "==", profID).
		Documents(ctx).GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar procedimentos"})
		return
	}

	var procedimentos []models.Procedimento
	for _, doc := range docs {
		var proc models.Procedimento
		if err := doc.DataTo(&proc); err == nil {
			procedimentos = append(procedimentos, proc)
		}
	}

	c.JSON(http.StatusOK, procedimentos)
}

// AtualizarProcedimento atualiza dados de um procedimento
// @Summary Atualizar procedimento
// @Tags Procedimentos
// @Accept json
// @Produce json
// @Param id path string true "ID do procedimento"
// @Param procedimento body models.Procedimento true "Procedimento atualizado"
// @Success 200 {object} map[string]string
// @Router /procedimentos/{id} [put]
func AtualizarProcedimento(c *gin.Context) {
	id := c.Param("id")
	var proc models.Procedimento
	if err := c.ShouldBindJSON(&proc); err != nil {
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

	_, err = client.Collection("procedimentos").Doc(id).Set(ctx, proc)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao atualizar procedimento"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"mensagem": "Procedimento atualizado"})
}

// DeletarProcedimento remove um procedimento
// @Summary Deletar procedimento
// @Tags Procedimentos
// @Produce json
// @Param id path string true "ID do procedimento"
// @Success 200 {object} map[string]string
// @Router /procedimentos/{id} [delete]
func DeletarProcedimento(c *gin.Context) {
	id := c.Param("id")
	ctx := context.Background()
	client, err := config.App.Firestore(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro na conexão com Firestore"})
		return
	}
	defer client.Close()

	_, err = client.Collection("procedimentos").Doc(id).Delete(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao deletar procedimento"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"mensagem": "Procedimento removido com sucesso"})
}
