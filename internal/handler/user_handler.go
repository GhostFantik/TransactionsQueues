package handler

import (
	"TransactionsQueues/pkg"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"net/http"
)

func (h Handler) CreateUser() func(c *gin.Context) {
	return func(c *gin.Context) {
		var input struct {
			Name string `json:"name"`
		}
		if err := c.ShouldBindJSON(&input); err != nil {
			log.Err(err).Stack().Msg("ошибка парсинга данных JSON")
			return
		}
		user, err := h.repositories.CreateUser(input.Name)
		if err != nil {
			ErrorResponse(c, err)
			return
		}
		c.JSON(http.StatusCreated, user)
		return
	}
}

func (h Handler) GetUser() func(c *gin.Context) {
	return func(c *gin.Context) {
		name, ok := c.GetQuery("name")
		if !ok {
			ErrorResponse(c, pkg.NewError(http.StatusBadRequest, "bad request"))
			return
		}
		user, err := h.repositories.GetUserByName(name)
		if err != nil {
			ErrorResponse(c, err)
			return
		}
		c.JSON(http.StatusOK, user)
	}
}
