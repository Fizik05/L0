package handler

import (
	"github.com/Fizik05/L0/pkg/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()
	router.LoadHTMLGlob("templates/*.html")

	api := router.Group("/api")
	{
		api.GET("/:order_uid", h.getOrder)
	}

	return router
}
