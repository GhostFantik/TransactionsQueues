package handler

import (
	"TransactionsQueues/internal/repository"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	repositories *repository.Repository
}

func NewHandler(repositories *repository.Repository) *Handler {
	return &Handler{repositories: repositories}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.Default()

	api := router.Group("/api")
	{
		user := api.Group("/user")
		{
			user.POST("/create", h.CreateUser())
		}
		balance := api.Group("/balance")
		{
			balance.POST("/put", h.PutBalance())
			balance.POST("/pop", h.PopBalance())
		}
	}
	return router
}
