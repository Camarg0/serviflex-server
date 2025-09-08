package main

import (
	"servico-api/config"
	"servico-api/routes"
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors" // Importando o pacote de CORS
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	// Configuração do Firebase (não relacionada ao CORS, mas necessária para o seu backend)
	config.InitFirebase()

	r := gin.Default()

	// Configuração do CORS
	r.Use(cors.Default()) // Permite todas as origens. Pode ser ajustado para mais restrições

	// Caso queira permitir apenas uma origem específica (como o seu frontend):
	// r.Use(cors.New(cors.Config{
	// 	AllowOrigins:     []string{"http://localhost:5173"},
	// 	AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
	// 	AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
	// 	AllowCredentials: true,
	// }))

	// Configuração das rotas
	routes.SetupRoutes(r)

	// Configuração do Swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Inicia o servidor na porta 8080
	r.Run(":8080")
}
