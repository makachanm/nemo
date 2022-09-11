package nemomark

type Nemomark struct {
	lexer    Lexer
	parser   Parser
	renderer Renderer
}

func MakeNemomark() Nemomark {
	return Nemomark{
		lexer:    NewLexer(),
		parser:   NewParser(),
		renderer: NewRenderer(),
	}
}

func (n *Nemomark) Mark(input string) string {
	lexed := n.lexer.Tokenize(input, TokenMap)
	parsed := n.parser.Parse(&lexed)
	result := n.renderer.render(parsed)

	return result
}
