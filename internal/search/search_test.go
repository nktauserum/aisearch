package search_test

import (
	"testing"

	"github.com/nktauserum/aisearch/internal/answer"
	"github.com/nktauserum/aisearch/internal/search"
)

// func TestDucksearch(t *testing.T) {
// 	result, err := search.("nezhegol")
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	t.Log(result)
// }

func TestDucksearch(t *testing.T) {
	result, err := search.SearchTavily("nezhegol")
	if err != nil {
		t.Fatal(err)
	}

	t.Log(result)
}

func TestDoSearchQueries(t *testing.T) {
	queries := []string{"Как приготовить борщ?", "Что представляет собой Чикагская школа?", "В каком году произошла отмена крепостного права?"}
	results, err := answer.DoSearchQueries(queries)
	if err != nil {
		t.Fatal(err)
	}

	if len(results) == 0 {
		t.Error("expected non-empty results")
	}

	t.Log(results)
}
