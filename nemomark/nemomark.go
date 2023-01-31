package nemomark

import nmcore "nemo/nemomark/core"

type Nemomark struct {
	Lexer    Lexer
	Parser   Parser
	NLexer   nmcore.NLexer
	NParser  nmcore.NParser
	Renderer nmcore.Renderer
	isLegacy bool
}

func NewNemomark(isLegacy bool) *Nemomark {
	return &Nemomark{
		Lexer:    NewLexer(),
		Parser:   NewParser(),
		NLexer:   nmcore.NewNLexer(),
		NParser:  nmcore.NewNParser(),
		Renderer: nmcore.NewRenderer(),
		isLegacy: isLegacy,
	}
}

func (n *Nemomark) Mark(input string) string {
	if n.isLegacy {
		lexed := n.Lexer.Tokenize(input, nmcore.TokenMap)
		parsed := n.Parser.Parse(&lexed)
		result := n.Renderer.Render(parsed)

		return result
	} else {
		lexed := n.NLexer.Toknize(input, nmcore.NTokenMap)
		parsed := n.NParser.Parse(lexed)
		result := n.Renderer.Render(parsed)

		return result
	}
}
