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
	prefix,
	suffix 	Node
}

var _ Node = (*ExpressionNode)(nil)

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

func newInsetNode(parser *Parser) (Node, error) {
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

	if parser.currentToken().Type != SEMICOL {
		return nil, NewError(
			InvalidInsertStatement,
			parser.currentToken().LineNumber,
			parser.currentToken().Start,
			";",
		)
	}

	if err := parser.nextToken(); err != nil {
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
}

var _ Node = (*SelectNode)(nil)

func (s *SelectNode) Table() Node {
	return s.from
}

func (s *SelectNode) Columns() []Node {
	return s.columns
}

func (s* SelectNode) String() string {
	return ""
}
