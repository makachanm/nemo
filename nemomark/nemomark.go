package nemomark

import "nemo/nemomark/core"

type Nemomark struct {
	Lexer    Lexer
	Parser   Parser
	NLexer   core.NLexer
	NParser  core.NParser
	Renderer core.Renderer
	isLegacy bool
}

func NewNemomark(isLegacy bool) *Nemomark {
	return &Nemomark{
		Lexer:    NewLexer(),
		Parser:   NewParser(),
		NLexer:   core.NewNLexer(),
		NParser:  core.NewNParser(),
		Renderer: core.NewRenderer(),
		isLegacy: isLegacy,
	}
}

func (n *Nemomark) Mark(input string) string {
	if n.isLegacy {
		lexed := n.Lexer.Tokenize(input, core.TokenMap)
		parsed := n.Parser.Parse(&lexed)
		result := n.Renderer.Render(parsed)

		return result
	} else {
		lexed := n.NLexer.Toknize(input, core.NTokenMap)
		parsed := n.NParser.Parse(lexed)
		result := n.Renderer.Render(parsed)

		return result
	}
}
