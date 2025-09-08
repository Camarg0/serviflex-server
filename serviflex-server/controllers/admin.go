package controllers

import (
	"context"
	"net/http"
	"servico-api/config"
	"servico-api/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ListarAdmins retorna todos os admins
func ListarAdmins(c *gin.Context) {
	ctx := context.Background()
	client, err := config.App.Firestore(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao conectar"})
		return
	}
	defer client.Close()

	docs, err := client.Collection("admin").Documents(ctx).GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar admins"})
		return
	}

	var admins []models.Admin
	for _, doc := range docs {
		var a models.Admin
		if err := doc.DataTo(&a); err == nil {
			admins = append(admins, a)
		}
	}
	c.JSON(http.StatusOK, admins)
}

// BuscarAdminPorID retorna um admin por ID
func BuscarAdminPorID(c *gin.Context) {
	id := c.Param("id")
	ctx := context.Background()
	client, err := config.App.Firestore(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao conectar"})
		return
	}
	defer client.Close()

	doc, err := client.Collection("admin").Doc(id).Get(ctx)
	if err != nil || !doc.Exists() {
		c.JSON(http.StatusNotFound, gin.H{"error": "Admin não encontrado"})
		return
	}

	var a models.Admin
	doc.DataTo(&a)
	c.JSON(http.StatusOK, a)
}

// CriarAdmin cria um novo admin
func CriarAdmin(c *gin.Context) {
	var input models.Admin
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos"})
		return
	}
	input.ID = uuid.New().String()
	input.Tipo = "admin"

	ctx := context.Background()
	client, err := config.App.Firestore(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao conectar"})
		return
	}
	defer client.Close()

	_, err = client.Collection("admin").Doc(input.ID).Set(ctx, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao criar admin"})
		return
	}
	c.JSON(http.StatusCreated, input)
}

// EditarAdmin atualiza os dados de um admin
func EditarAdmin(c *gin.Context) {
	id := c.Param("id")
	var input models.Admin
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos"})
		return
	}
	input.ID = id
	input.Tipo = "admin"

	ctx := context.Background()
	client, err := config.App.Firestore(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao conectar"})
		return
	}
	defer client.Close()

	_, err = client.Collection("admin").Doc(id).Set(ctx, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao atualizar admin"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Admin atualizado com sucesso"})
}

// ExcluirAdmin remove um admin
func ExcluirAdmin(c *gin.Context) {
	id := c.Param("id")
	ctx := context.Background()
	client, err := config.App.Firestore(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao conectar"})
		return
	}
	defer client.Close()

	_, err = client.Collection("admin").Doc(id).Delete(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao excluir admin"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Admin excluído com sucesso"})
}
