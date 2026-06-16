package router

import (
	"infinite-canvas-server/handler"
	"infinite-canvas-server/middleware"
	"infinite-canvas-server/service"

	"github.com/gin-gonic/gin"
)

func Setup(r *gin.Engine, authService *service.AuthService, authHandler *handler.AuthHandler, userHandler *handler.UserHandler, creditHandler *handler.CreditHandler) {
	r.Use(middleware.Cors())

	api := r.Group("/api")

	api.POST("/auth/register", authHandler.Register)
	api.POST("/auth/login", authHandler.Login)

	auth := api.Group("")
	auth.Use(middleware.AuthRequired(authService))
	{
		auth.GET("/auth/me", authHandler.Me)

		auth.GET("/credits/balance", creditHandler.GetBalance)
		auth.GET("/credits/transactions", creditHandler.GetTransactions)
		auth.GET("/credits/estimate", creditHandler.EstimateCost)

		admin := auth.Group("")
		admin.Use(middleware.AdminRequired())
		{
			admin.GET("/users", userHandler.List)
			admin.GET("/credits/pricing", creditHandler.ListPricing)
			admin.POST("/credits/pricing", creditHandler.SavePricing)
			admin.DELETE("/credits/pricing/:id", creditHandler.DeletePricing)
		}
	}
}
