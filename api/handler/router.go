package handler

import (
	"dna-test/service"
	"github.com/gin-gonic/gin"
)

func BindRoutes(service service.Service, g *gin.Engine) {
	h := NewHandler(service)

	// Swagger can be added here
	//g.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	g.GET("/healthcheck", h.HealthCheck)

	account := g.Group("/account")
	{
		account.POST("/", h.AccountHandler.CreateAccount)
		account.PATCH("/freeze/:id", h.AccountHandler.FreezeAccount)
		account.GET("/:id", h.AccountHandler.GetAccount)
	}
}
