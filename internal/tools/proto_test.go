package tools

import (
	"context"
	"testing"
	"time"
)

func TestGetContent(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	GetWebResourceContent(ctx, "https://www.bbc.com/russian/articles/cg704gpj0yxo")
}
