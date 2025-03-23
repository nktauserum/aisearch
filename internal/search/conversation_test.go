package search_test

import (
	"testing"

	"github.com/nktauserum/aisearch/internal/search"
)

func TestDialog(t *testing.T) {
	messages := make(chan string)
	go func() {
		messages <- "Что такое либертарианство?"
		messages <- "Как оно развивалось?"
		messages <- "Состояние в России?"
		close(messages)
	}()

	search.Dialog(messages)

}
