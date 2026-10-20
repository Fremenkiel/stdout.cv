package interpreter

import (
	"bytes"
	"fmt"
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
	if !parser.currentToken().isField() {
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

func (f *FieldNode) Columns() []Node {
	return []Node{f}
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

	node.Token = parser.currentToken()

	if err := parser.nextToken(); err != nil {
		return nil, err
	}

	precedence, ok := precedenceMap[node.Token.Type]
	if !ok {
		return nil, NewError(
			InvalidTypeForPrecedence,
			node.Token.LineNumber,
			node.Token.Start,
			node.Token.Text,
			)
	}

	node.precedence = precedence

	return &node, nil
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

	if prefix := e.prefix; prefix != nil {
		if str := prefix.String(); len(str) > 0 {
			buf.WriteString(str)
			buf.WriteString(" ")
		}
	}

	buf.WriteString(e.Token.Text)


	if suffix := e.suffix; suffix != nil {
		if str := suffix.String(); len(str) > 0 {
			buf.WriteString(" ")
			buf.WriteString(str)
		}
	}

	return buf.String()
}

type OperatorNode struct {
	infix				[]baseNode
	prefix,
	suffix			Node
	precedence	int
}

var _ Node = (*OperatorNode)(nil)

func newOperatorNode(parser *Parser) (*OperatorNode, error) {
	var node OperatorNode
	node.precedence = 0

	if parser.currentToken().isField() {
		prefixNode, err := newFieldNode(parser)
		if err != nil {
			return nil, err
		}

		node.prefix = prefixNode
	}

	for {
		if parser.currentToken().Type == EOF {
			return &node, nil
		}

		if parser.currentToken().isField() {
			break
		}

		infix := baseNode{
			Token: parser.currentToken(),
		}

		node.infix = append(node.infix, infix)

		precedence, ok := precedenceMap[infix.Token.Type]
		if !ok {
			return nil, NewError(
				InvalidTypeForPrecedence,
				infix.Token.LineNumber,
				infix.Token.Start,
				infix.Token.Text,
				)
		}

		if precedence > node.precedence {
			node.precedence = precedence
		}

		if err := parser.nextToken(); err != nil {
			return nil, err
		}
	}

	// suffix is only relevant if prefix is also set
	if node.prefix != nil {
		suffixNode, err := newFieldNode(parser)
		if err != nil {
			return nil, err
		}

		node.suffix = suffixNode
	}

	return &node, nil
}

func (o *OperatorNode) Table() Node {
	return nil
}

func (o *OperatorNode) Columns() []Node {
	return nil
}

func (o *OperatorNode) String() string {
	var buf bytes.Buffer

	if prefix := o.prefix; prefix != nil {
		if str := prefix.String(); len(str) > 0 {
			buf.WriteString(str)
			buf.WriteString(" ")
		}
	}

	for i := range o.infix {
		buf.WriteString(o.infix[i].Token.Text)

		if i < len(o.infix) - 1 {
			buf.WriteString(" ")
		}
	}

	if suffix := o.suffix; suffix != nil {
		if str := suffix.String(); len(str) > 0 {
			buf.WriteString(" ")
			buf.WriteString(str)
		}
	}

	return buf.String()
}

type CastNode struct {
	baseNode
	prefix,
	suffix		Node
}

var _ Node = (*CastNode)(nil)

func newCastNode(parser *Parser) (*CastNode, error) {
	var node CastNode

	if parser.peakToken().Type != AS {
		return nil, NewError(
			InvalidCastType,
			parser.peakToken().LineNumber,
			parser.peakToken().Start,
			parser.peakToken().Text,
			)
	}

	prefix, err := newFieldNode(parser)
	if err != nil {
		return nil, err
	}

	node.prefix = prefix

	node.Token = parser.currentToken()

	if err := parser.nextToken(); err != nil {
		return nil, err
	}
	
	suffix, err := newFieldNode(parser)
	if err != nil {
		return nil, err
	}

	node.suffix =suffix 

	return &node, nil
}

func (c *CastNode) Columns() []Node {
	return []Node{c.prefix}
}

func (c *CastNode) Table() Node {
	return c.prefix
}

func (c *CastNode) String() string {
	var buf bytes.Buffer

	buf.WriteString(c.prefix.String())
	buf.WriteString(" ")
	buf.WriteString(c.Token.Text)
	buf.WriteString(" ")
	buf.WriteString(c.suffix.String())

	return buf.String()
}

// list
type ListNode struct {
	nodes []Node
}

var _ Node = (*ListNode)(nil)

func newListNode(parser *Parser) (*ListNode, error) {
	var list ListNode
	if parser.currentToken().Type != L_PARAN {
		return nil, NewError(
			InvalidListFormat,
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
	if parser.currentToken().Type != INSERT {
		return nil, NewError(
			InvalidBaseNodeType,
			parser.currentToken().LineNumber,
			parser.currentToken().Start,
			fmt.Sprintf("expected %s got %s", INSERT, parser.currentToken().Type),
			)
	}

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
			string(INTO),
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
			string(VALUES),
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

func newSelectNode(parser *Parser) (Node, error) {
	if parser.currentToken().Type != SELECT {
		return nil, NewError(
			InvalidBaseNodeType,
			parser.currentToken().LineNumber,
			parser.currentToken().Start,
			fmt.Sprintf("expected %s got %s", SELECT, parser.currentToken().Type),
			)
	}

	var node SelectNode
	node.Token = parser.currentToken()

	if err := parser.nextToken(); err != nil {
		return nil, err
	}

	for {
		var column Node
		if parser.currentToken().Type != FIELD {
			break
		}

		if parser.peakToken().Type == AS {
			response, err := newCastNode(parser)
			if err != nil {
				return nil, err
			}
			column = response 
		} else {
			response, err := newFieldNode(parser)
			if err != nil {
				return nil, err
			}
			column = response 
		}

		node.columns = append(node.columns, column)
	}

	if parser.currentToken().Type != FROM {
		return &node, nil
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
		where, err := newWhereNode(parser)
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
	var columns []Node

	for i := range s.columns {
		columns = append(columns, s.columns[i].Columns()...)
	}

	return columns
}

func (s *SelectNode) String() string {
	var buf bytes.Buffer

	buf.WriteString("SELECT ")

	for i := range s.columns {
		buf.WriteString(s.columns[i].String())

		if i < len(s.columns) - 1 {
			buf.WriteString(", ")
		}
	}

	if s.from == nil {
		return buf.String()
	}

	buf.WriteString(" ")
	buf.WriteString("FROM ")
	buf.WriteString(s.from.String())

	if s.where != nil {
		buf.WriteString(" ")
		buf.WriteString(s.where.String())
	}

	buf.WriteString(";")

	return buf.String()
}

// UPDATE
type UpdateNode struct {
	baseNode
	table				Node
	set					Node
	where				Node
}

var _ Node = (*UpdateNode)(nil)

func newUpdateNode(parser *Parser) (Node, error) {
	if parser.currentToken().Type != UPDATE {
		return nil, NewError(
			InvalidBaseNodeType,
			parser.currentToken().LineNumber,
			parser.currentToken().Start,
			fmt.Sprintf("expected %s got %s", UPDATE, parser.currentToken().Type),
			)
	}

	var node UpdateNode
	node.Token = parser.currentToken()

	if err := parser.nextToken(); err != nil {
		return nil, err
	}

	if parser.currentToken().Type != FIELD {
		return nil, NewError(
			InvalidUpdateStatement,
			parser.currentToken().LineNumber,
			parser.currentToken().Start,
			"table",
			)
	}

	table, err := newFieldNode(parser)
	if err != nil {
		return nil, err
	}

	node.table = table

	if parser.currentToken().Type != SET {
		return nil, NewError(
			InvalidUpdateStatement,
			parser.currentToken().LineNumber,
			parser.currentToken().Start,
			string(SET),
			)
	}

	// consume set tokes
	if err := parser.nextToken(); err != nil {
		return nil, err
	}

	set, err := parser.expressionParse(0)
	if err != nil {
		return nil, err
	}

	node.set = set

	if parser.currentToken().Type == WHERE {
		where, err := newWhereNode(parser)
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

func (u *UpdateNode) Table() Node {
	return u.table
}

// update do not return data for our use
func (u *UpdateNode) Columns() []Node {
	return nil
}

func (u *UpdateNode) String() string {
	var buf bytes.Buffer

	buf.WriteString("UPDATE ")
	buf.WriteString(u.table.String())
	buf.WriteString(" SET ")

	buf.WriteString(u.set.String())

	if u.where != nil {
		buf.WriteString(" ")
		buf.WriteString(u.where.String())
	}

	buf.WriteString(";")

	return buf.String()
}

// DELETE
type DeleteNode struct {
	baseNode
	from    Node
	where		Node
}

var _ Node = (*DeleteNode)(nil)

func newDeleteNode(parser *Parser) (Node, error) {
	if parser.currentToken().Type != DELETE {
		return nil, NewError(
			InvalidBaseNodeType,
			parser.currentToken().LineNumber,
			parser.currentToken().Start,
			fmt.Sprintf("expected %s got %s", DELETE, parser.currentToken().Type),
			)
	}

	var node DeleteNode
	node.Token = parser.currentToken()

	if err := parser.nextToken(); err != nil {
		return nil, err
	}

	if parser.currentToken().Type != FROM {
		return nil, NewError(
			InvalidDeleteStatement,
			parser.currentToken().LineNumber,
			parser.currentToken().Start,
			string(FROM),
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
		where, err := newWhereNode(parser)
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

func (d *DeleteNode) Table() Node {
	return d.from
}

func (d *DeleteNode) Columns() []Node {
	return nil
}

func (d *DeleteNode) String() string {
	var buf bytes.Buffer

	buf.WriteString("DELETE FROM ")
	buf.WriteString(d.from.String())

	if d.where != nil {
		buf.WriteString(" ")
		buf.WriteString(d.where.String())
	}

	buf.WriteString(";")

	return buf.String()
}

// WHERE
type WhereNode struct {
	baseNode
	prefix			Node
	expression	Node
}

var _ Node = (*WhereNode)(nil)

// do not consume end char as where can't be a stand alone statement
func newWhereNode(parser *Parser) (Node, error) {
	if parser.currentToken().Type != WHERE {
		return nil, NewError(
			InvalidBaseNodeType,
			parser.currentToken().LineNumber,
			parser.currentToken().Start,
			fmt.Sprintf("expected %s got %s", WHERE, parser.currentToken().Type),
			)
	}

	var node WhereNode
	node.Token = parser.currentToken()

	if err := parser.nextToken(); err != nil {
		return nil, err
	}

	if !parser.currentToken().isField() {
		prefix, err := newOperatorNode(parser)
		if err != nil {
			return nil, err
		}

		node.prefix = prefix
	}
	
	expression, err := parser.expressionParse(0)
	if err != nil {
		return nil, err
	}

	node.expression = expression

	return &node, nil
}

func (w *WhereNode) Table() Node {
	return nil
}

func (w *WhereNode) Columns() []Node {
	return nil
}

func (w *WhereNode) String() string {
	var buf bytes.Buffer

	buf.WriteString("WHERE ")

	if w.prefix != nil {
		buf.WriteString(w.prefix.String())
		buf.WriteString(" ")
	}
	
	buf.WriteString(w.expression.String())

	return buf.String()
}

func consumeEndChar(parser *Parser) error {
	if parser.currentToken().Type != SEMICOL {
		return NewError(
			InvalidEndOfStatement,
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
