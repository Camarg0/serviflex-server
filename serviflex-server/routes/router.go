package routes

import (
	"servico-api/controllers"
	"servico-api/utils"

	"github.com/gin-gonic/gin"
)

// SetupRoutes configura todas as rotas principais da API
func SetupRoutes(router *gin.Engine) {
	api := router.Group("/api")

	SetupAuthRoutes(api)
	SetupClienteRoutes(api)
	SetupProfissionalRoutes(api)
	SetupProcedimentoRoutes(api)
	SetupHorarioRoutes(api)
	SetupAgendamentoRoutes(api)
	SetupUploadRoutes(api)
	SetupEstabelecimentoRoutes(api)
	SetupAdminRoutes(api) // ADICIONE ESTA LINHA

}

func SetupAuthRoutes(rg *gin.RouterGroup) {
	rg.POST("/login", utils.Login)
	rg.POST("/cadastro", utils.CadastrarUsuario)
}

func SetupClienteRoutes(rg *gin.RouterGroup) {
	rg.GET("/agendamentos/cliente/:id", controllers.ListarAgendamentosPorCliente)
	rg.PUT("/usuarios/:id", controllers.EditarUsuario) // Agora requer ?tipo=clientes|profissionais|admin
}

func SetupEstabelecimentoRoutes(rg *gin.RouterGroup) {
	rg.POST("/estabelecimentos", controllers.CriarEstabelecimento)
	rg.PUT("/estabelecimentos/:id", controllers.EditarEstabelecimento)
	rg.GET("/estabelecimentos", controllers.ListarEstabelecimentos)
	rg.GET("/estabelecimentos/:id", controllers.BuscarEstabelecimentoPorID)
	rg.GET("/relatorios/estabelecimento/faturamento/:id", controllers.RelatorioFaturamentoEstabelecimento)
	rg.GET("/relatorios/avaliacoes/estabelecimento/:id", controllers.RelatorioAvaliacoesPorEstabelecimento)
	rg.GET("/relatorios/agendamentos/estabelecimento/:id", controllers.RelatorioAgendamentosPorMesEstabelecimento)
	rg.POST("/estabelecimentos/profissionais/convidar", controllers.ConvidarProfissional)
	rg.POST("/estabelecimentos/profissionais/notificacao/:id", controllers.AceitarOuRecusarConvite)
	rg.DELETE("/estabelecimentos/:estId/profissionais/:profId", controllers.RemoverProfissional)
	rg.GET("/estabelecimentos/:id/profissionais", controllers.ListarProfissionaisDoEstabelecimento)
}

func SetupProfissionalRoutes(rg *gin.RouterGroup) {
	rg.GET("/agendamentos/profissional/:id", controllers.ListarAgendamentosPorProfissional)
	rg.GET("/horarios/:id", controllers.ListarHorariosPorProfissional)
	rg.GET("/profissionais/:uid/convites-pendentes", controllers.ListarConvitesPendentes)
	rg.GET("/profissionais", controllers.ListarProfissionais)         // NOVA ROTA
	rg.GET("/profissionais/:uid", controllers.BuscarProfissionalPorID) // NOVA ROTA
	rg.GET("/relatorios/profissional/faturamento/:id", controllers.RelatorioFaturamentoProfissional)
	rg.GET("/relatorios/avaliacoes/profissional/:id", controllers.RelatorioAvaliacoesPorProfissional)
	rg.GET("/relatorios/agendamentos/profissional/:id", controllers.RelatorioAgendamentosPorMesProfissional)

}

func SetupProcedimentoRoutes(rg *gin.RouterGroup) {
	rg.POST("/procedimentos", controllers.CriarProcedimento)
	rg.GET("/procedimentos/:id", controllers.ListarProcedimentosPorProfissional)
	rg.PUT("/procedimentos/:id", controllers.AtualizarProcedimento)
	rg.DELETE("/procedimentos/:id", controllers.DeletarProcedimento)
}

func SetupHorarioRoutes(rg *gin.RouterGroup) {
	rg.POST("/horarios", controllers.CriarHorario)
	rg.PUT("/horarios/:id", controllers.EditarHorario)   // NOVA ROTA
	rg.DELETE("/horarios/:id", controllers.ExcluirHorario) // NOVA ROTA
}

func SetupUploadRoutes(rg *gin.RouterGroup) {
	rg.PUT("/upload/:tipo/:id", controllers.SetImagemURL) // tipo = profissional | procedimento
}

func SetupAgendamentoRoutes(rg *gin.RouterGroup) {
	rg.POST("/agendamentos", controllers.AgendarHorario)
}

func SetupAdminRoutes(rg *gin.RouterGroup) {
	rg.GET("/admins", controllers.ListarAdmins)
	rg.GET("/admins/:id", controllers.BuscarAdminPorID)
	rg.POST("/admins", controllers.CriarAdmin)
	rg.PUT("/admins/:id", controllers.EditarAdmin)
	rg.DELETE("/admins/:id", controllers.ExcluirAdmin)
}
