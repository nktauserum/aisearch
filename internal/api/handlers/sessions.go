package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nktauserum/aisearch/pkg/ai/client"
)

func SessionListHandler(c *gin.Context) {
	// Объявляем контекст в десять секунд
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Берём список всех текущих сессий
	memory := client.GetMemory()
	list := memory.GetConversationList(ctx)

	c.JSON(http.StatusOK, list)
}
