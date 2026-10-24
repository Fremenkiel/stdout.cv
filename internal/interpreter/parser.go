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
			INSERT:	newInsetNode,
			SELECT:	newSelectNode,
			UPDATE: newUpdateNode,
			DELETE:	newDeleteNode,
		},
		r: -1,
	}
	return parser
}

func (p *Parser) ParseQuery(query string) error {
	p.trees = make([]Node, 0, 1)
	p.lexer.Write(query)
	p.r = 0

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
				UnsupportedQuery,
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

func (p *Parser) GetExecureQueries() []Node {
	var mod int = 0
	if rq := p.GetReturnQuery(); rq != nil {
		mod = 1
	}

	if (len(p.trees) - mod) < 0 {
		return nil
	}

	return p.trees[0:len(p.trees) - mod]
}

func (p *Parser) GetReturnQuery() Node {
	if len(p.trees) == 0 {
		return nil
	}

	if rq := p.trees[len(p.trees) - 1]; rq.ReturnsValue() {
		return rq
	}

	return nil
}

func (p *Parser) expressionParse(basePrecedence int) (Node, error) {
	left, err := p.newParseNode()
	if err != nil {
		return nil, err
	}

	for {
		operator := p.currentToken()

		operatorPrecedence, ok := precedenceMap[operator.Type]
		if !ok {
			return nil, NewError(
				InvalidTypeForPrecedence,
				operator.LineNumber,
				operator.Start,
				operator.Text,
				)
		}

		if operatorPrecedence < basePrecedence {
			break
		}

		if err := p.nextToken(); err != nil {
			return nil, err
		}

		right, err := p.expressionParse(operatorPrecedence)
		if err != nil {
			return nil, err
		}

		node := ExpressionNode{
			precedence: operatorPrecedence,
			prefix: left,
			suffix: right,
		}
		node.Token = operator

		left = &node
	}

	return left, nil
}

func (p *Parser) newParseNode() (Node, error) {
	if p.currentToken().isField() {
		return newOperatorNode(p)
	}

	return newExpressionNode(p)
}

func (p *Parser) currentToken() Token {
	if p.r >= len(p.tokens) {
		return Token{}
	}

	return p.tokens[p.r]
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
