package interpreter

type Parser struct {
	lexer 		*Lexer
	handlers	map[TokenType]func(*Parser) (BaseNode, error)

	r					int
	tokens		[]Token
}

func NewParser() *Parser {
	parser := &Parser{
		lexer: NewLexer(),
		handlers: map[TokenType]func(*Parser) (BaseNode, error) {
			INSERT:	newInsetNode,
			SELECT:	newSelectNode,
			UPDATE: newUpdateNode,
			DELETE:	newDeleteNode,
			DROP: newDropNode,
			ALTER: newAlterNode,
			CREATE: newCreateNode,
		},
		r: -1,
	}
	return parser
}

func (p *Parser) ParseSql(sql string) ([]BaseNode, error) {
	trees := make([]BaseNode, 0, 3)
	p.lexer.Write(sql)
	p.r = 0

	for {
		t, err := p.lexer.Next()
		if err != nil {
			return nil, err
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
			return nil, NewError(
				UNSUPPORTED_QUERY,
				p.currentToken().LineNumber,
				p.currentToken().End,
				p.currentToken().Text,
				)
		}

		node, err := handler(p)
		if err != nil {
			return nil, err
		}

		trees = append(trees, node)
	}

	return trees, nil
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
				INVALID_TYPE_FOR_PRECEDENCE,
				operator.LineNumber,
				operator.End,
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
			token: operator,
		}

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

func (p *Parser) peakToken() Token {
	if (p.r + 1) >= len(p.tokens) {
		return Token{}
	}

	return p.tokens[p.r + 1]
}

func (p *Parser) currentIsTable() bool {
	if p.r == 0 {
		return false
	}

	if p.tokens[p.r - 1].isTablePrefix() {
		return true
	}

	return false
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
			UNEXPECTED_EOF,
			p.currentToken().LineNumber,
			p.currentToken().End,
			p.currentToken().Text,
			)

	}
	p.r = p.r + 1

	return nil
}
