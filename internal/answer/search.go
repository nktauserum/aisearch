package answer

import (
	"context"
	"log"
	"time"

	"github.com/nktauserum/aisearch/shared"
)

func Search(ctx context.Context, search_info SearchInfo) ([]shared.Website, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	searchStart := time.Now()
	content, err := ExtractInfo(ctx, search_info.Queries...)
	if err != nil {
		return nil, err
	}
	log.Printf("Search completed: %v ms", time.Since(searchStart).Milliseconds())

	return content, nil
}
