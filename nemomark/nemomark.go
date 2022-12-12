package nemomark

type Nemomark struct {
	Lexer    Lexer
	Parser   Parser
	Renderer Renderer
}

func NewNemomark() *Nemomark {
	return &Nemomark{
		Lexer:    NewLexer(),
		Parser:   NewParser(),
		Renderer: NewRenderer(),
	}
}

func (n *Nemomark) Mark(input string) string {
	lexed := n.Lexer.Tokenize(input, TokenMap)
	parsed := n.Parser.Parse(&lexed)
	result := n.Renderer.Render(parsed)

	return result
}
