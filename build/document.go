package build

import (
	"bytes"
	"nemo/nemomark"
	"strconv"
	"text/template"
)

type MarkdownFucntion nemomark.MarkdownFucntion

type TimeStamp struct {
	Year  int
	Month int
	Day   int
	Hour  int
	Min   int
}

type Document_Meta struct {
	Title     string
	Timestamp TimeStamp
	Summary   string
	Path      string
}

type Document struct {
	Meta    Document_Meta
	Head    Header
	Foot    Footer
	Nav     Nav
	Content string
}

type Header struct {
	IsNotIndex bool
	Header     string
}

type Footer struct {
	IsNotIndex bool
	Footer     string
}

type Nav struct {
	IsNotIndex bool
	Navbar     string
}

type IndexPage struct {
	Indexs []Document_Meta
	Head   Header
	Foot   Footer
	Nav    Nav
}

func MakeMetaData() Document_Meta {
	return Document_Meta{}
}

func NewDocument() Document {
	return Document{}
}

func NewIndexData() IndexPage {
	return IndexPage{}
}

func MakeTimeStamp(year int, month int, day int, hour int, min int) TimeStamp {
	return TimeStamp{
		Year:  year,
		Month: month,
		Day:   day,
		Hour:  hour,
		Min:   min,
	}
}

func (t *TimeStamp) StampSize() int {
	return t.Year + t.Month + t.Day + t.Hour + t.Min
}

func (d *Document) ParseMeta(input string) {
	var metadata Document_Meta = MakeMetaData()

	lexer := nemomark.NewLexer()
	parser := nemomark.NewParser()

	lexed := lexer.Tokenize(input, nemomark.TokenMap)
	parse_result := parser.Parse(&lexed)

	for _, ctx := range parse_result.Child {
		switch ctx.Func_context.Fucntion_name {
		case "title":
			metadata.Title = ctx.Func_context.Context[0]

		case "summary":
			metadata.Summary = ctx.Func_context.Context[0]

		case "timestamp":
			year, yerr := strconv.Atoi(ctx.Func_context.Args["year"])
			if yerr != nil {
				year = 0
			}
			month, merr := strconv.Atoi(ctx.Func_context.Args["month"])
			if merr != nil {
				month = 0
			}
			date, derr := strconv.Atoi(ctx.Func_context.Args["day"])
			if derr != nil {
				date = 0
			}
			hour, derr := strconv.Atoi(ctx.Func_context.Args["hour"])
			if derr != nil {
				hour = 0
			}
			min, derr := strconv.Atoi(ctx.Func_context.Args["min"])
			if derr != nil {
				min = 0
			}

			stamp := MakeTimeStamp(year, month, date, hour, min)
			metadata.Timestamp = stamp
		}
	}

	d.Meta = metadata
}

func BuildHeader(skin Skin, Ctx interface{}) Header {
	t, err := template.ParseFiles(skin.Info.Paths.Header)

	if err != nil {
		panic(err)
	}

	var writer bytes.Buffer
	t.Execute(&writer, Ctx)

	var head Header
	head.Header = writer.String()

	return head
}

func BuildFooter(skin Skin, Ctx interface{}) Footer {
	t, err := template.ParseFiles(skin.Info.Paths.Footer)

	if err != nil {
		panic(err)
	}

	var writer bytes.Buffer
	t.Execute(&writer, Ctx)

	var foot Footer
	foot.Footer = writer.String()

	return foot
}

func BuildNav(skin Skin, Ctx interface{}) Nav {
	t, err := template.ParseFiles(skin.Info.Paths.Nav)

	if err != nil {
		panic(err)
	}

	var writer bytes.Buffer
	t.Execute(&writer, Ctx)

	var nav Nav
	nav.Navbar = writer.String()

	return nav
}
