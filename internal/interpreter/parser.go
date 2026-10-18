package interpreter

type Parser struct {
	lexer 		*Lexer
	handlers	map[TokenType]func(*Parser) (Node, error)
	trees   	[]Node

	r					int
	tokens		[]Token
}

func NewParser() *Parser {
	parser := &Parser{
		lexer: NewLexer(),
		handlers: map[TokenType]func(*Parser) (Node, error) {
			INSERT: newInsetNode,
			// SELECT: newSelectNode,
		},
	}
	return parser
}

func (p *Parser) ParseQuery(query string) error {
	p.lexer.Write(query)
	for {
		t, err := p.lexer.Next()
		if err != nil {
			return err
		}

		p.tokens = append(p.tokens, t)

		if t.Type == EOF {
			break
		}
	}

	for {
		if p.currentToken().Type == EOF {
			break
		}

		handler, ok := p.handlers[p.currentToken().Type]
		if !ok {
			return NewError(
				UnknownBaseNodeType,
				p.currentToken().LineNumber,
				p.currentToken().Start,
				p.currentToken().Text,
				)
		}

		node, err := handler(p)
		if err != nil {
			return err
		}

		p.trees = append(p.trees, node)
	}

	return nil
}

func (p *Parser) currentToken() Token {
	if p.r >= len(p.tokens) {
		return Token{}
	}

	return p.tokens[p.r]
}

func (p *Parser) peakToken() Token {
	if (p.r + 1) >= len(p.tokens) {
		return Token{}
	}

	return p.tokens[p.r + 1]
}

// returns true if there are more tokens to read, false if EOF
func (p *Parser) nextToken() error {
	if p.r >= len(p.tokens) - 1 {
		return NewError(
			UnexpectedEOF,
			p.currentToken().LineNumber,
			p.currentToken().Start,
			p.currentToken().Text,
			)

	}
	p.r = p.r + 1

	return nil
}
