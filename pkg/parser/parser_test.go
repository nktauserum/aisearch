package parser

import (
	"context"
	"testing"
	"time"
)

func TestMakeReadable(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	raw_content, err := GetContent(ctx, "https://pkg.go.dev/crypto/tls")
	if err != nil {
		t.Fatal(err)
	}

	t.Log(raw_content.Content)
}
