package controllers

import (
	"context"
	"net/http"
	"servico-api/config"

	"cloud.google.com/go/firestore"
	"github.com/gin-gonic/gin"
)

type ImagemInput struct {
	ImagemURL string `json:"imagem_url"`
}

// SetImagemURL atualiza a URL da imagem de um profissional ou procedimento
// @Summary Atualizar imagem do recurso
// @Tags Upload
// @Accept json
// @Produce json
// @Param tipo path string true "Tipo do recurso (profissional ou procedimento)"
// @Param id path string true "ID do recurso"
// @Param imagem body ImagemInput true "URL da imagem"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /upload/{tipo}/{id} [put]
func SetImagemURL(c *gin.Context) {
	tipo := c.Param("tipo")
	id := c.Param("id")

	var input ImagemInput
	if err := c.ShouldBindJSON(&input); err != nil || input.ImagemURL == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "URL da imagem inválida"})
		return
	}

	ctx := context.Background()
	client, err := config.App.Firestore(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao conectar com Firestore"})
		return
	}
	defer client.Close()

	col := ""
	switch tipo {
	case "profissional":
		col = "profissionais"
	case "procedimento":
		col = "procedimentos"
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Tipo inválido"})
		return
	}

	_, err = client.Collection(col).Doc(id).Update(ctx, []firestore.Update{
		{Path: "imagem_url", Value: input.ImagemURL},
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao atualizar imagem no Firestore"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"imagem_url": input.ImagemURL})
}
