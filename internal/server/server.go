package server

import (
	"github.com/gin-gonic/gin"

	"github.com/koba-masa/equivalent_exchange_v2/internal/handler"
)

func New() *gin.Engine {
	engine := gin.New()
	engine.Use(gin.Logger(), gin.Recovery())
	engine.GET("/health", handler.Health)
	return engine
}
