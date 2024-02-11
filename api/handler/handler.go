package handler

import (
	"dna-test/service"
	"github.com/gin-gonic/gin"
)

type Handlers struct {
	AccountHandler *AccountHandler
}

func NewHandler(service service.Service) *Handlers {
	return &Handlers{
		AccountHandler: NewAccountHandler(service),
	}
}

func (s *Handlers) HealthCheck(c *gin.Context) {
	c.JSON(200, "OK")
	return
}
