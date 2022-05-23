package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"net/http"
)

func (h Handler) PutBalance() func(c *gin.Context) {
	return func(c *gin.Context) {
		var input struct {
			UserId uuid.UUID `json:"user_id"`
			Money  int       `json:"money"`
		}
		if err := c.ShouldBindJSON(&input); err != nil {
			log.Err(err).Stack().Msg("ошибка парсинга данных JSON")
			return
		}
		err := h.repositories.PutBalance(input.UserId, input.Money)
		if err != nil {
			ErrorResponse(c, err)
			return
		}
		c.JSON(http.StatusCreated, "success")
		return
	}
}

func (h Handler) PopBalance() func(c *gin.Context) {
	return func(c *gin.Context) {
		var input struct {
			UserId uuid.UUID `json:"user_id"`
			Money  int       `json:"money"`
		}
		if err := c.ShouldBindJSON(&input); err != nil {
			log.Err(err).Stack().Msg("ошибка парсинга данных JSON")
			return
		}
		err := h.repositories.PopBalance(input.UserId, input.Money)
		if err != nil {
			ErrorResponse(c, err)
			return
		}
		c.JSON(http.StatusCreated, "success")
		return
	}
}
