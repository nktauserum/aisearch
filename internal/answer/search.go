package answer

import (
	"context"
	"time"

	"github.com/nktauserum/aisearch/shared"
)

func Search(ctx context.Context, search_info SearchInfo) ([]shared.Website, error) {
	// Устанавливаем таймаут для операции поиска
	ctxTimeout, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()

	// Получаем канал с результатами
	websites, err := ExtractInfo(ctxTimeout, search_info.Queries...)
	if err != nil {
		return nil, err
	}

	return websites, nil
}
