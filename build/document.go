package build

import (
	"bytes"
	"nemo/nemomark"
	"nemo/nemomark/core"
	"strconv"
	"text/template"
)

type MarkdownFucntion core.MarkdownFucntion

type TimeStamp struct {
	Year  int
	Month int
	Day   int
	Hour  int
	Min   int
}

type DocumentMeta struct {
	Title     string
	Timestamp TimeStamp
	Summary   string
	Path      string
	Tags      string
}

type Document struct {
	Meta    DocumentMeta
	Head    Header
	Foot    Footer
	Nav     Nav
	Content string
}

type Header struct {
	IsNotIndex bool
	Header     string
	BlogName   string
	PostName   string
}

type Footer struct {
	IsNotIndex bool
	Footer     string
}

type Nav struct {
	IsNotIndex bool
	Navbar     string
	BlogName   string
}

type IndexPage struct {
	Indexs    []DocumentMeta
	IndexsNum int
	Head      Header
	Foot      Footer
	Nav       Nav
	NextPage  string
	PrevPage  string
}

type AboutPage struct {
	Meta       DocumentMeta
	Head       Header
	Foot       Footer
	Nav        Nav
	Content    string
	BuildInfo  string
	SkinInfo   SkinInfo
	AuthorInfo string
}

type TagsPage struct {
	Tags    map[string][]DocumentMeta
	TagsNum int
	Head    Header
	Foot    Footer
	Nav     Nav
}

func MakeMetaData() DocumentMeta {
	return DocumentMeta{}
}

func NewDocument() Document {
	return Document{}
}

func NewIndexData() IndexPage {
	return IndexPage{}
}

func NewAboutPage() AboutPage {
	return AboutPage{}
}

func NewTagsData() TagsPage {
	return TagsPage{}
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

func (t *TimeStamp) isBiggerStamp(src TimeStamp, cmp TimeStamp) bool {
	if src.Year > cmp.Year {
		return true
	} else if src.Month > cmp.Month {
		return true
	} else if src.Month == cmp.Month {
		if src.Day > cmp.Day {
			return true
		} else if src.Hour > cmp.Hour {
			return true
		} else if src.Min > cmp.Min {
			return true
		}
	}

	return false
}

func (d *Document) ParseMeta(input string) {
	var metadata = MakeMetaData()

	lexer := nemomark.NewLexer()
	parser := nemomark.NewParser()

	lexed := lexer.Tokenize(input, core.TokenMap)
	parseResult := parser.Parse(&lexed)

	for _, ctx := range parseResult.Child {
		switch ctx.FuncContext.FunctionName {
		case "title":
			tagsData := ctx.FuncContext.Context[0]
			if len(tagsData) <= 0 || tagsData == "" {
				metadata.Title = ""
				break
			}

			metadata.Title = tagsData

		case "summary":
			tagsData := ctx.FuncContext.Context[0]
			if len(tagsData) <= 0 || tagsData == "" {
				metadata.Summary = ""
				break
			}

			metadata.Summary = tagsData

		case "timestamp":
			year, yerr := strconv.Atoi(ctx.FuncContext.Args["year"])
			if yerr != nil {
				year = 0
			}
			month, merr := strconv.Atoi(ctx.FuncContext.Args["month"])
			if merr != nil {
				month = 0
			}
			date, derr := strconv.Atoi(ctx.FuncContext.Args["day"])
			if derr != nil {
				date = 0
			}
			hour, derr := strconv.Atoi(ctx.FuncContext.Args["hour"])
			if derr != nil {
				hour = 0
			}
			min, derr := strconv.Atoi(ctx.FuncContext.Args["min"])
			if derr != nil {
				min = 0
			}

			stamp := MakeTimeStamp(year, month, date, hour, min)
			metadata.Timestamp = stamp

		case "tag":
			tagsData := ctx.FuncContext.Context[0]
			if len(tagsData) <= 0 || tagsData == "" {
				metadata.Tags = ""
				break
			}

			metadata.Tags = tagsData
		}
	}

	d.Meta = metadata
}

func BuildHeader(skin Skin, Ctx interface{}) (Header, error) {
	t, err := template.ParseFiles(skin.Info.Paths.Header)

	if err != nil {
		return Header{}, err
	}

	var writer bytes.Buffer
	_ = t.Execute(&writer, Ctx)

	var head Header
	head.Header = writer.String()

	return head, nil
}

func BuildFooter(skin Skin, Ctx interface{}) (Footer, error) {
	t, err := template.ParseFiles(skin.Info.Paths.Footer)

	if err != nil {
		return Footer{}, err
	}

	var writer bytes.Buffer
	_ = t.Execute(&writer, Ctx)

	var foot Footer
	foot.Footer = writer.String()

	return foot, nil
}

func BuildNav(skin Skin, Ctx interface{}) (Nav, error) {
	t, err := template.ParseFiles(skin.Info.Paths.Nav)

	if err != nil {
		return Nav{}, err
	}

	var writer bytes.Buffer
	_ = t.Execute(&writer, Ctx)

	var nav Nav
	nav.Navbar = writer.String()

	return nav, nil
}
