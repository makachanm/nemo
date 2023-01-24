package nemomark

type Nemomark struct {
	Lexer    Lexer
	Parser   Parser
	NLexer   NLexer
	NParser  NParser
	Renderer Renderer
	isLegacy bool
}

func NewNemomark(isLegacy bool) *Nemomark {
	return &Nemomark{
		Lexer:    NewLexer(),
		Parser:   NewParser(),
		NLexer:   NewNLexer(),
		NParser:  NewNParser(),
		Renderer: NewRenderer(),
		isLegacy: isLegacy,
	}
}

func (n *Nemomark) Mark(input string) string {
	if n.isLegacy {
		lexed := n.Lexer.Tokenize(input, TokenMap)
		parsed := n.Parser.Parse(&lexed)
		result := n.Renderer.Render(parsed)

		return result
	} else {
		lexed := n.NLexer.Toknize(input, NTokenMap)
		parsed := n.NParser.Parse(lexed)
		result := n.Renderer.Render(parsed)

		return result
	}
}
