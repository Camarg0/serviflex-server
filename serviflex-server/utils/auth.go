package utils

import (
	"context"
	"net/http"
	"servico-api/config"
	"servico-api/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Login realiza autenticação simples via Firestore
// @Summary Login de usuário
// @Description Autentica um usuário com email e senha
// @Tags Autenticação
// @Accept json
// @Produce json
// @Param credenciais body models.Login true "Email e Senha"
// @Success 200 {object} models.Usuario
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /login [post]
func Login(c *gin.Context) {
	var credenciais struct {
		Email string `json:"email"`
		Senha string `json:"senha"`
	}
	if err := c.ShouldBindJSON(&credenciais); err != nil {
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

	colecoes := []string{"clientes", "profissionais", "admin"}
	for _, colecao := range colecoes {
		docs, err := client.Collection(colecao).
			Where("email", "==", credenciais.Email).
			Where("senha", "==", credenciais.Senha).
			Limit(1).Documents(ctx).GetAll()

		if err == nil && len(docs) > 0 {
			var usuario models.Usuario
			if err := docs[0].DataTo(&usuario); err == nil {
				usuario.Tipo = colecao
				c.JSON(http.StatusOK, gin.H{
					"mensagem": "Login realizado com sucesso",
					"usuario":  usuario,
				})
				return
			}
		}
	}

	c.JSON(http.StatusUnauthorized, gin.H{"error": "Email ou senha inválidos"})
}

// CadastrarUsuario cria um novo cliente ou profissional
// @Summary Cadastro de usuário
// @Description Cadastra novo cliente ou profissional no Firestore
// @Tags Autenticação
// @Accept json
// @Produce json
// @Param usuario body models.Usuario true "Dados do usuário"
// @Success 201 {object} models.Usuario
// @Failure 400 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /cadastro [post]
func CadastrarUsuario(c *gin.Context) {
	var novoUsuario models.Usuario

	if err := c.ShouldBindJSON(&novoUsuario); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos"})
		return
	}

	if novoUsuario.Tipo != "clientes" && novoUsuario.Tipo != "profissionais" && novoUsuario.Tipo != "admin" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Tipo inválido (use 'clientes', 'admin' ou 'profissionais')"})
		return
	}

	ctx := context.Background()
	client, err := config.App.Firestore(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao conectar com Firestore"})
		return
	}
	defer client.Close()

	exists, err := client.Collection(novoUsuario.Tipo).
		Where("email", "==", novoUsuario.Email).
		Limit(1).Documents(ctx).GetAll()
	if err == nil && len(exists) > 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "Email já cadastrado"})
		return
	}

	novoUsuario.ID = uuid.New().String()

	_, err = client.Collection(novoUsuario.Tipo).Doc(novoUsuario.ID).Set(ctx, novoUsuario)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao salvar usuário"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"mensagem": "Usuário cadastrado com sucesso", "usuario": novoUsuario})
}
