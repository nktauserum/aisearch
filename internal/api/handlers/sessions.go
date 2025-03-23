package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	mw "github.com/nktauserum/aisearch/internal/api/middleware"
	"github.com/nktauserum/aisearch/pkg/ai/client"
)

func SessionListHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// Объявляем контекст в десять секунд
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Берём список всех текущих сессий
	memory := client.GetMemory()
	list := memory.GetConversationList(ctx)

	// Переводим в json и отдаём
	result, err := json.Marshal(&list)
	if err != nil {
		mw.ErrorHandler(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(result)
}
