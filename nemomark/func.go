package nemomark

import "strings"

var Markdown_Handlers = map[string]func(MarkdownFucntion) string{
	"bold":      bold,
	"italic":    italic,
	"cancel":    strikehrough,
	"underline": underline,
	"link":      link,
}

func RenderPlain(input []string) string {
	str := strings.Join(input, "")
	return str
}

func bold(input MarkdownFucntion) string {
	str := strings.Join(input.Context, "")
	return "<strong>" + str + "</strong>"
}

func italic(input MarkdownFucntion) string {
	str := strings.Join(input.Context, "")
	return "<em>" + str + "</em>"
}

func strikehrough(input MarkdownFucntion) string {
	str := strings.Join(input.Context, "")
	return "<del>" + str + "</del>"
}

func underline(input MarkdownFucntion) string {
	str := strings.Join(input.Context, "")
	return "<u>" + str + "</u>"
}

func link(input MarkdownFucntion) string {
	str := strings.Join(input.Context, "")
	link, isexist := input.Args["url"]

	if !isexist {
		link = ""
	}

	return `<a href="` + link + `">` + str + "</a>"
}
