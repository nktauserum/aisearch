package answer

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/nktauserum/aisearch/pkg/ai/models"
	"github.com/nktauserum/aisearch/pkg/parser"
	"github.com/nktauserum/aisearch/shared"
)

func ExtractInfo(ctx context.Context, queries ...string) ([]shared.Website, error) {
	// Создаем контекст с timeout для всей операции
	ctxTimeout, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()

	log.Println("Ищем в интернете")
	urls, err := DoSearchQueries(queries)
	if err != nil {
		log.Printf("Ошибка DoSearchQueries: %v", err)
		return nil, err
	}

	log.Printf("Получено %d URL", len(urls))

	// Канал для сбора результатов. Размер буфера равен количеству URL, чтобы избежать блокировок.
	resultsCh := make(chan shared.Website, len(urls))
	// Канал для сбора ошибок (необязательный, если требуется обработка ошибок)
	errorCh := make(chan error, len(urls))

	var wg sync.WaitGroup

	// Для каждого URL запускаем горутину.
	for _, url := range urls {
		// Чтобы переменная url не была переопределена в замыкании
		url := url
		wg.Add(1)
		go func() {
			defer wg.Done()

			// Вызываем функцию получения контента с использованием ctxTimeout.
			siteInfo, err := parser.GetContent(ctxTimeout, url)
			if err != nil {
				log.Printf("Ошибка получения контента для %s: %v", url, err)
				// Можно отправлять ошибку на канал, если требуется дальнейшая обработка
				errorCh <- err
				return
			}

			log.Printf("Сайт %s обработан", url)
			resultsCh <- siteInfo
		}()
	}

	// Ожидаем завершения всех горутин в отдельной горутине,
	// после чего закрываем канал результатов.
	go func() {
		wg.Wait()
		close(resultsCh)
		close(errorCh)
	}()

	// Собираем результаты до окончания работы или таймаута.
	var websites []shared.Website
collectLoop:
	for {
		select {
		case site, ok := <-resultsCh:
			if !ok {
				// Канал закрыт – собираем все результаты закончены.
				break collectLoop
			}
			websites = append(websites, site)
		case <-ctxTimeout.Done():
			log.Println("Таймаут выполнения, возвращаем уже полученные результаты")
			break collectLoop
		}
	}

	// Если требуется возврат ошибки при наличии ошибок, можно объединить их
	// или вернуть ошибку, если ни одного сайта не удалось получить.
	// Здесь возвращаем только результаты.
	return websites, nil
}

// Даёт ответ на запрос по переданному контенту
func Research(ctx context.Context, conversation *models.Conversation, content string) (string, error) {
	// Get summary from AI
	summary, err := conversation.Continue(ctx, models.Message{Text: content})
	if err != nil {
		log.Println(err)
		return "", err
	}

	// Write summary to file

	return summary, nil
}
