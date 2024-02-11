package handler

import (
	"dna-test/models"
	"dna-test/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

// Additional tools like logger can be added
type AccountHandler struct {
	Service service.Service
}

func NewAccountHandler(service service.Service) *AccountHandler {
	return &AccountHandler{service}
}

func (h *AccountHandler) CreateAccount(c *gin.Context) {
	ctx := c.Request.Context()

	var input models.Account
	err := c.ShouldBind(&input)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	err = h.Service.CreateAccount(ctx, &input)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusCreated, input)
}

func (h *AccountHandler) FreezeAccount(c *gin.Context) {
	ctx := c.Request.Context()

	id := c.Param("id")

	err := h.Service.FreezeAccount(ctx, id)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, "done")
}

func (h *AccountHandler) GetAccount(c *gin.Context) {
	ctx := c.Request.Context()

	id := c.Param("id")

	account, err := h.Service.GetAccount(ctx, id)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, account)

}
