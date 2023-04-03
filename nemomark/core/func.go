package nmcore

import (
	"strings"
)

var MarkdownHandlers = map[string]func(MarkdownFucntion) string{
	"bold":      bold,
	"italic":    italic,
	"cancel":    strikehrough,
	"underline": underline,
	"link":      link,
	"image":     image,
	"video":     video,
	"code":      code,
	"head":      head,
	"small":     small,
}

func RenderPlain(input []string) string {
	str := strings.Join(input, "")
	str = strings.ReplaceAll(str, "\n", "<br />")

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

func image(input MarkdownFucntion) string {
	altstr := strings.Join(input.Context, "")
	src, isexist := input.Args["url"]

	if !isexist {
		return "alt: " + altstr
	}

	imgtag := `<img src="` + src + `" class="content-image" alt="` + altstr + `">`

	return `<a href="` + src + `">` + imgtag + `</a>`
}

func video(input MarkdownFucntion) string {
	src, isexist := input.Args["url"]

	if !isexist {
		return ""
	}

	vidtag := `<video src="` + src + `" class="content-video" controls></video>`

	return `<a href="` + src + `">` + vidtag + `</a>`
}

func code(input MarkdownFucntion) string {
	str := strings.Join(input.Context, "")
	str = strings.ReplaceAll(str, "\n", "<br />")

	return `<div class="code"> <code>` + str + `</code> </div>`
}

func head(input MarkdownFucntion) string {
	str := strings.Join(input.Context, "")
	return "<h2>" + str + "</h2><hr />"
}

func small(input MarkdownFucntion) string {
	str := strings.Join(input.Context, "")
	return "<h6>" + str + "</h6>"
}
