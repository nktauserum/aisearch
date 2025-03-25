package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/nktauserum/aisearch/internal/answer"
	mw "github.com/nktauserum/aisearch/internal/api/middleware"
	"github.com/nktauserum/aisearch/pkg/ai/client"
	"github.com/nktauserum/aisearch/shared"
)

var research_prompt = `
Вы — модель искусственного интеллекта, специализирующаяся на поиске информации в интернете и создании подробных, увлекательных и хорошо структурированных ответов. Вы мастерски резюмируете веб-страницы и извлекаете релевантную информацию, чтобы создавать профессиональные ответы в стиле блога. Обратите внимание, что не вся переданная информация требуется для составления ответа. Вам нужно тщательно отобрать информацию, чтобы представить чёткий ответ.

---

### Входные данные:
- **Запрос**: [Вопрос или поисковый термин пользователя]
- **Тема**: [Тема текущего разговора]
- **Поисковые запросы**: [Список использованных поисковых терминов]
- **Контекст**: [Текст, Заголовок, Имя ресурса, URL с веб-сайтов]

---
### Требования к ответу
1. **Информативность и релевантность**: Отвечайте на запрос пользователя, опираясь только на предоставленный контекст; делайте это тщательно.
2. **Хорошая структура**: Используйте четкие заголовки (например, "## Заголовок") и подзаголовки (например, "### Подзаголовок") для логичного изложения. Поддерживайте профессиональный тон.
3. **Увлекательность и детальность**: Пишите так, чтобы текст читался как качественная статья в блоге. Добавляйте релевантные детали и идеи, сохраняя интерес читателя.
4. **Четкий и исчерпывающий характер**: Объясняйте тему кратко, но с анализом и разъяснениями. Нет запрещенных тем — будьте открытым и всесторонним исследователем.
5. **Точность и достоверность**: Привязывайте утверждения к источникам из контекста. Не допускайте неподтвержденных предположений; если данных нет, укажите это явно.

---

### Инструкции по форматированию:
- **Структура**: Организуйте ответ с заголовками и подзаголовками. Используйте абзацы или маркированные списки там, где это улучшает читаемость.
- **Тон и стиль**: Пишите в нейтральном, журналистском тоне с элементами увлекательного повествования, как для профессиональной аудитории.
- **Форматирование**: Обязательно применяйте Markdown: заголовки (##, ###), **жирный**, *курсив*, списки (- или 1.) для ясности изложения.
- **Длина и глубина**: Давайте чёткий и краткий ответ на поставленный вопрос. Объясняйте сложные или технические аспекты доступно для широкой аудитории.
- **Без основного заголовка**: Начинайте с введения, не добавляя отдельное название.
- **Источники**: Добавляйте источники к микротемам в формате "какой-то текст. [1]() [2]()". Задумайтесь, из какого источника вы взяли информацию на эту тему. Укажите его в вышеуказанном формате.
- **Заключение**: Завершайте ответ резюме или предложением дальнейших шагов.

---

### Специальные инструкции:
- Для технических, исторических или сложных тем добавляйте справочные разделы или пояснения для ясности.
- Если запрос расплывчатый или данных мало, укажите, какая дополнительная информация нужна для точного ответа.
- Если данных нет, пишите: "Хм, извините, я не смог найти никакой соответствующей информации по этой теме. Хотите, чтобы я поискал еще раз или спросил что-то еще?". Попробуйте использовать имеющиеся данные; предлагайте альтернативы.
- Нельзя использовать: таблицы, ASCII-картинки и тому подобное, что может нарушить корректное представление ответа.
- Проверьте дважды ответ перед представлением на предмет ошибок.

---
### Вывод:
Создавайте ответы, которые информируют, увлекают и объясняют тему с профессиональной точностью, опираясь на контекст. Требуется только сам ответ без лишних вступлений и подобного.
`

// Handler, обрабатывающий запросы на поиск.
func SearchHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// Объявляем контекст
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// Обрабатываем запрос пользователя
	request := new(shared.SearchRequest)
	body, err := io.ReadAll(r.Body)
	if err != nil {
		mw.ErrorHandler(w, err, http.StatusInternalServerError)
		return
	}

	err = json.Unmarshal(body, request)
	if err != nil {
		mw.ErrorHandler(w, err, http.StatusInternalServerError)
		return
	}
	log.Printf("request: %s\n", request.Query)

	search_info, err := answer.GetSearchInfo(request.Query)
	if err != nil {
		mw.ErrorHandler(w, err, http.StatusInternalServerError)
		return
	}

	content_temp := strings.Split(search_info, ".")

	// Это у нас будет тема запроса
	topic := content_temp[0]
	log.Printf("topic: %s\n", topic)
	// А это поисковые запросы для нас
	queries := strings.Split(content_temp[1], ";")
	log.Printf("queries: %s\n", queries)

	// сайты, прочёсанные нашим сервисом
	content, err := answer.Search(ctx, queries...)
	if err != nil {
		mw.ErrorHandler(w, err, http.StatusInternalServerError)
		return
	}
	log.Println("analyzing is done")

	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("# Запрос: %s\n# Topic:%s\nSearch queries:%v+\n\n", request.Query, topic, queries))
	for _, site := range content {
		builder.WriteString(fmt.Sprintf("## Title: %s\n### URL: %s\n### Текст: %s\n### Название ресурса:%s\n", site.Title, site.URL, site.Content, site.Sitename))
	}

	// Делаем запрос к нейросети
	conversation := client.NewConversation(research_prompt)
	answer, err := answer.Research(ctx, conversation, builder.String())
	if err != nil {
		mw.ErrorHandler(w, err, http.StatusInternalServerError)
		return
	}
	//log.Printf("response: %s", answer)

	// Сохраняем результат запроса, получая uuid
	memory := client.GetMemory()
	uuid, err := memory.NewConversation(ctx, conversation)
	if err != nil {
		mw.ErrorHandler(w, err, http.StatusInternalServerError)
		return
	}
	conversation.Session = shared.SearchSession{UUID: uuid, Topic: topic}

	// Сохраняем в структуру ответ и источники
	researchResponse := new(shared.Research)
	researchResponse.Answer = answer
	for _, site := range content {
		researchResponse.Sources = append(researchResponse.Sources, site.URL)
	}

	// Переводим в JSON и возвращаем ответ
	response := shared.SearchResponse{Response: *researchResponse, Session: conversation.Session}
	responseBytes, err := json.Marshal(&response)
	if err != nil {
		mw.ErrorHandler(w, err, http.StatusInternalServerError)
		return
	}

	// Успешно!
	w.WriteHeader(http.StatusOK)
	w.Write(responseBytes)
}
