package shared

import (
	"fmt"
	"log"
	"os"
)

type ParseMode interface {
	Design() string
	Sources() string
}

type FormatMD struct {
	design  string
	sources string
}

func NewFormatMD() *FormatMD {
	absPath, err := os.Getwd()
	if err != nil {
		log.Println(err)
		return nil
	}

	info, err := os.ReadFile(fmt.Sprintf("%s/prompt/format/markdown_v2.md", absPath))
	if err != nil {
		log.Println(err)
		return nil
	}

	design := fmt.Sprintf(`
	Вы можете использовать жирный, курсивный, подчеркнутый, зачеркнутый текст, спойлеры, блочные цитаты, а также встроенные ссылки и предварительно отформатированный код. Нельзя использовать что-то, помимо этого.
	Используйте указанные инструкции по экранированию символов при форматировании.

	%s
	`, string(info))

	sources := "[число 1](url), [число 2](url)"

	return &FormatMD{design: design, sources: sources}
}

func (f *FormatMD) Design() string {
	return f.design
}

func (f *FormatMD) Sources() string {
	return f.sources
}

type FormatHTML struct {
	design  string
	sources string
}

func NewFormatHTML() *FormatHTML {
	absPath, err := os.Getwd()
	if err != nil {
		log.Println(err)
		return nil
	}

	info, err := os.ReadFile(fmt.Sprintf("%s/prompt/format/html.md", absPath))
	if err != nil {
		log.Println(err)
		return nil
	}

	design := fmt.Sprintf(`
	Используйте HTML-теги по указанным инструкциям. Нельзя использовать теги <p>, <br>, <ul>, <li> и вообще любые, помимо указанных.
	Вы можете использовать жирный, курсивный, подчеркнутый, зачеркнутый текст, спойлеры, блочные цитаты, а также встроенные ссылки и предварительно отформатированный код. Нельзя использовать что-то, помимо этого.
	При составлении ответа экранируйте указанные символы. 
	Используйте корректный способ представления ссылок. Всегда указывайте URL с помощью атрибута href.
	Заголовки делайте в виде выделенного жирным шрифтом текста с помощью тега <b>
	Не указывайте, что форматируете в HTML. Требуется лишь сам ответ с допустимыми тегами, не блок кода.
	Недопустимо markdown-форматирование.

	%s
	`, string(info))

	sources := "<a href=\"http://www.example.com/\">[1]</a>, <a href=\"http://www.example2.com/\">[2]</a>"

	return &FormatHTML{design: design, sources: sources}
}

func (f *FormatHTML) Design() string {
	return f.design
}

func (f *FormatHTML) Sources() string {
	return f.sources
}

type FormatXML struct {
	design  string
	sources string
}

func NewFormatXML() *FormatXML {
	absPath, err := os.Getwd()
	if err != nil {
		log.Println(err)
		return nil
	}

	info, err := os.ReadFile(fmt.Sprintf("%s/prompt/format/xml.md", absPath))
	if err != nil {
		log.Println(err)
		return nil
	}

	design := string(info)

	sources := `
	Текст, нуждающийся в подтверждении. <link url="http://www.example.com/">[1]</link> <link url="http://www.example2.com/">[2]</link>

	Не указывайте источники в конце ответа отдельным абзацем. 
	`

	return &FormatXML{design: design, sources: sources}
}

func (f *FormatXML) Design() string {
	return f.design
}

func (f *FormatXML) Sources() string {
	return f.sources
}
