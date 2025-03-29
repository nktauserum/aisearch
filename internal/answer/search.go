package answer

import (
	"context"
	"time"

	"github.com/nktauserum/aisearch/shared"
)

func Search(ctx context.Context, search_info SearchInfo) ([]shared.Website, error) {
	ctx_timeout, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()

	content, err := ExtractInfo(ctx_timeout, search_info.Queries...)
	if err != nil {
		return nil, err
	}

	websites := make([]shared.Website, len(content))
	for site := range content {
		websites = append(websites, site)
	}

	return websites, nil
}
