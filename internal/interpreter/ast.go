package interpreter

import (
	"bytes"
)

type Node interface {
	Table() Node
	Columns() []Node
	String() string
}

type baseNode struct {
	Token Token
}

type FieldNode struct {
	baseNode
}

var _ Node = (*FieldNode)(nil)

func newFieldNode(parser *Parser) (*FieldNode, error) {
	if parser.currentToken().Type != FIELD &&
	parser.currentToken().Type != STRING &&
	parser.currentToken().Type != NUMBER {
		return nil, NewError(
			InvalidField,
			parser.currentToken().LineNumber,
			parser.currentToken().Start,
			parser.currentToken().Text,
			)
	}

	var node FieldNode
	node.Token = parser.currentToken()

	if err := parser.nextToken(); err != nil {
		return nil, err
	}

	return &node, nil
}

func (f *FieldNode) Table() Node {
	return f
}

// is handled by the parent node
func (f *FieldNode) Columns() []Node {
	return nil
}

func (f *FieldNode) String() string {
	if f.Token.Type == STRING {
		return "'" + f.Token.Text + "'"
	}

	return f.Token.Text
}

type ExpressionNode struct {
	baseNode
	precedence	int

	prefix,
	suffix 	Node
}

var _ Node = (*ExpressionNode)(nil)

func newExpressionNode(parser *Parser) (*ExpressionNode, error) {
	var node ExpressionNode

	if parser.currentToken().Type == FIELD {
		prefixNode, err := newFieldNode(parser)
		if err != nil {
			return nil, err
		}

		node.prefix = prefixNode

		node.Token = parser.currentToken()

		if err := parser.nextToken(); err != nil {
			return nil, err
		}

		suffixNode, err := newFieldNode(parser)
		if err != nil {
			return nil, err
		}

		node.suffix = suffixNode
	} else {
		node.Token = parser.currentToken()
	}

	precedence, ok := precedenceMap[node.Token.Type]
	if !ok {
		return nil, NewError(
			InvalidTypeForExpressionPrecedence,
			node.Token.LineNumber,
			node.Token.Start,
			node.Token.Text,
			)
	}

	node.precedence = precedence

	if err := parser.nextToken(); err != nil {
		return nil, err
	}

	return &node, nil
}

func prattParse(parser *Parser, basePrecedence int) (*ExpressionNode, error) {
	left, err := newExpressionNode(parser)
	if err != nil {
		return nil, err
	}

	for {
		operator := parser.peakToken()
		if operator.Type == EOF {
			break
		}

		operatorPrecedence, ok := precedenceMap[operator.Type]
		if !ok {
			return nil, NewError(
				InvalidTypeForExpressionPrecedence,
				operator.LineNumber,
				operator.Start,
				operator.Text,
				)
		}

		if err := parser.nextToken(); err != nil {
			return nil, err
		}
		right, err := prattParse(parser, operatorPrecedence)
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

// expressions can't hold the table node
func (e *ExpressionNode) Table() Node {
	return nil
}

// expressions can't hold returned columns
func (e *ExpressionNode) Columns() []Node {
	return nil
}

func (e *ExpressionNode) String() string {
	var buf bytes.Buffer

	buf.WriteString(e.prefix.String())
	buf.WriteString(" ")
	buf.WriteString(e.Token.Text)
	buf.WriteString(" ")
	buf.WriteString(e.suffix.String())

	return buf.String()
}

type OperatorNode struct {
	prefix,
	infix,
	suffix			Node
	precedence	int
}

var _ Node = (*OperatorNode)(nil)

func newOperatorNode() (*OperatorNode, error) {
	return &OperatorNode{}, nil
}

func (o *OperatorNode) Table() Node {
	return nil
}

func (o *OperatorNode) Columns() []Node {
	return nil
}

func (o *OperatorNode) String() string {
	return ""
}

type ListNode struct {
	nodes []Node
}

var _ Node = (*ListNode)(nil)

func newListNode(parser *Parser) (*ListNode, error) {
	var list ListNode
	if parser.currentToken().Type != L_PARAN {
		return nil, NewError(
			InvalidInsertStatement,
			parser.currentToken().LineNumber,
			parser.currentToken().Start,
			"(",
			)
	}

	// consume l paran
	if err := parser.nextToken(); err != nil {
		return nil, err
	}

	for {
		if parser.currentToken().Type == R_PARAN {
			break
		}

		valueNode, err := newFieldNode(parser)
		if err != nil {
			return nil, err
		}

		list.nodes = append(list.nodes, valueNode)
	}

	// consume r paran
	if err := parser.nextToken(); err != nil {
		return nil, err
	}

	return &list, nil
}

// list can't hold the table node
func (l *ListNode) Table() Node {
	return nil
}

func (l *ListNode) Columns() []Node {
	return l.nodes
}

func (l *ListNode) String() string {
	var buf bytes.Buffer

	buf.WriteString("(")

	for i := range l.nodes {
		buf.WriteString(l.nodes[i].String())

		if i < len(l.nodes) - 1 {
			buf.WriteString(", ")
		}
	}

	buf.WriteString(")")

	return buf.String()
}


// INSERT
type InsertNode struct {
	baseNode

	columns Node
	into    Node
	values 	Node
}

var _ Node = (*InsertNode)(nil)

func newInsetNode(parser *Parser) (*InsertNode, error) {
	var node InsertNode
	node.Token =  parser.currentToken()

	if err := parser.nextToken(); err != nil {
		return nil, err
	}

	if parser.currentToken().Type != INTO {
		return nil, NewError(
			InvalidInsertStatement,
			parser.currentToken().LineNumber,
			parser.currentToken().Start,
			"INTO",
			)
	}

	// consume into token
	if err := parser.nextToken(); err != nil {
		return nil, err
	}

	intoNode, err := newFieldNode(parser)
	if err != nil {
		return nil, err
	}

	node.into = intoNode

	columns, err := newListNode(parser)
	if err != nil {
		return nil, err
	}

	node.columns = columns

	if parser.currentToken().Type != VALUES {
		return nil, NewError(
			InvalidInsertStatement,
			parser.currentToken().LineNumber,
			parser.currentToken().Start,
			"VALUES",
			)
	}

	// consume values tokes
	if err := parser.nextToken(); err != nil {
		return nil, err
	}

	values, err := newListNode(parser)
	if err != nil {
		return nil, err
	}

	node.values = values

	if err := consumeEndChar(parser); err != nil {
		return nil, err
	}

	return &node, nil
}

func (i *InsertNode) Table() Node {
	return i.into
}

func (i *InsertNode) Columns() []Node {
	return i.columns.Columns()
}

func (i *InsertNode) String() string {
	var buf bytes.Buffer

	buf.WriteString(i.Token.Text)
	buf.WriteString(" INTO ")
	buf.WriteString(i.into.String())
	buf.WriteString(" ")
	buf.WriteString(i.columns.String())
	buf.WriteString(" VALUES ")
	buf.WriteString(i.values.String())
	buf.WriteString(";")

	return buf.String()
}

// SELECT
type SelectNode struct {
	baseNode
	columns []Node
	from    Node
	where		Node
}

var _ Node = (*SelectNode)(nil)

func newSelectNode(parser *Parser) (*SelectNode, error) {
	var node SelectNode
	node.Token = parser.currentToken()

	if err := parser.nextToken(); err != nil {
		return nil, err
	}

	for {
		if parser.currentToken().Type != FIELD {
			break
		}

		column, err := newFieldNode(parser)
		if err != nil {
			return nil, err
		}

		node.columns = append(node.columns, column)
	}

	if parser.currentToken().Type != FROM {
		return nil, NewError(
			InvalidInsertStatement,
			parser.currentToken().LineNumber,
			parser.currentToken().Start,
			"FROM",
			)
	}

	// consume from tokes
	if err := parser.nextToken(); err != nil {
		return nil, err
	}

	from, err := newFieldNode(parser)
	if err != nil {
		return nil, err
	}

	node.from = from

	if parser.currentToken().Type == WHERE {
		where, err := newExpressionNode(parser)
		if err != nil {
			return nil, err
		}

		node.where = where
	}

	if err := consumeEndChar(parser); err != nil {
		return nil, err
	}

	return &node, nil
}

func (s *SelectNode) Table() Node {
	return s.from
}

func (s *SelectNode) Columns() []Node {
	return s.columns
}

func (s* SelectNode) String() string {
	return ""
}

func consumeEndChar(parser *Parser) error {
	if parser.currentToken().Type != SEMICOL {
		return NewError(
			InvalidInsertStatement,
			parser.currentToken().LineNumber,
			parser.currentToken().Start,
			";",
			)
	}

	// consume semicol token
	if err := parser.nextToken(); err != nil {
		return err
	}

	return nil
}
