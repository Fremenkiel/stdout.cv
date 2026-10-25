package interpreter

import (
	"bytes"
)

type Node interface {
	Table()					Node
	Columns()				[]string
	String()				string
	ReturnsValue()	bool
}

type baseNode struct {
	Token Token
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
			INVALID_INSERT_STATEMENT,
			parser.currentToken().LineNumber,
			parser.currentToken().End,
			parser.currentToken().Text,
			)
	}

	var node InsertNode
	node.Token =  parser.currentToken()

	if err := parser.nextToken(); err != nil {
		return nil, err
	}

	if parser.currentToken().Type != INTO {
		return nil, NewError(
			INVALID_INSERT_STATEMENT,
			parser.currentToken().LineNumber,
			parser.currentToken().End,
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
			INVALID_INSERT_STATEMENT,
			parser.currentToken().LineNumber,
			parser.currentToken().End,
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

func (i *InsertNode) Columns() []string {
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

func (i *InsertNode) ReturnsValue() bool {
	return false
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
			INVALID_SELECT_STATEMENT,
			parser.currentToken().LineNumber,
			parser.currentToken().End,
			parser.currentToken().Text,
			)
	}

	var node SelectNode
	node.Token = parser.currentToken()

	if err := parser.nextToken(); err != nil {
		return nil, err
	}

	for {
		var column Node
		if !parser.currentToken().isField() {
			break
		}

		fieldNode, err := newFieldNode(parser)
		if err != nil {
			return nil, err
		}
		column = fieldNode 

		node.columns = append(node.columns, column)
	}

	if parser.currentToken().Type != FROM {
		if err := consumeEndChar(parser); err != nil {
			return nil, err
		}
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
	if table := s.from; table != nil {
		return table.Table()
	}

	return nil
}

func (s *SelectNode) Columns() []string {
	var columns []string

	for i := range s.columns {
		columns = append(columns, s.columns[i].Columns()...)
	}

	return columns
}

func (s *SelectNode) String() string {
	var buf bytes.Buffer

	buf.WriteString("SELECT")

	if len(s.columns) > 0 {
		buf.WriteString(" ")
	}

	for i := range s.columns {
		buf.WriteString(s.columns[i].String())

		if i < len(s.columns) - 1 {
			buf.WriteString(", ")
		}
	}

	if s.from == nil {
		buf.WriteString(";")
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

func (s *SelectNode) ReturnsValue() bool {
	return true
}

// UPDATE
type UpdateNode struct {
	baseNode
	table				Node
	set					[]Node
	where				Node
}

var _ Node = (*UpdateNode)(nil)

func newUpdateNode(parser *Parser) (Node, error) {
	if parser.currentToken().Type != UPDATE {
		return nil, NewError(
			INVALID_UPDATE_STATEMENT,
			parser.currentToken().LineNumber,
			parser.currentToken().End,
			parser.currentToken().Text,
			)
	}

	var node UpdateNode
	node.Token = parser.currentToken()

	if err := parser.nextToken(); err != nil {
		return nil, err
	}

	if parser.currentToken().Type != FIELD {
		return nil, NewError(
			INVALID_UPDATE_STATEMENT,
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
			INVALID_UPDATE_STATEMENT,
			parser.currentToken().LineNumber,
			parser.currentToken().End,
			string(SET),
			)
	}

	// consume set tokes
	if err := parser.nextToken(); err != nil {
		return nil, err
	}

	for {
		if parser.currentToken().Type == EOF {
			return nil, NewError(
				INVALID_END_OF_STATEMENT,
				parser.currentToken().LineNumber,
				parser.currentToken().End,
				";",
				)
		}

		if parser.currentToken().Type == SEMICOL {
			break
		}
		
		if parser.currentToken().Type == WHERE {
			break
		}

		set, err := newOperatorNode(parser)
		if err != nil {
			return nil, err
		}

		node.set = append(node.set, set)
	}

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
func (u *UpdateNode) Columns() []string {
	return nil
}

func (u *UpdateNode) String() string {
	var buf bytes.Buffer

	buf.WriteString("UPDATE")

	if u.table != nil {
		buf.WriteString(" ")
		buf.WriteString(u.table.String())
	}

	buf.WriteString(" SET")

	if len(u.set) > 0 {
		buf.WriteString(" ")
	}

	for i := range u.set {
		buf.WriteString(u.set[i].String())

		if i < len(u.set) - 1 {
			buf.WriteString(", ")
		}
	}

	if u.where != nil {
		buf.WriteString(" ")
		buf.WriteString(u.where.String())
	}

	buf.WriteString(";")

	return buf.String()
}

func (u *UpdateNode) ReturnsValue() bool {
	return false
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
			INVALID_DELETE_STATEMENT,
			parser.currentToken().LineNumber,
			parser.currentToken().End,
			parser.currentToken().Text,
			)
	}

	var node DeleteNode
	node.Token = parser.currentToken()

	if err := parser.nextToken(); err != nil {
		return nil, err
	}

	if parser.currentToken().Type != FROM {
		return nil, NewError(
			INVALID_DELETE_STATEMENT,
			parser.currentToken().LineNumber,
			parser.currentToken().End,
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

func (d *DeleteNode) Columns() []string {
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

func (d *DeleteNode) ReturnsValue() bool {
	return false
}

type FieldNode struct {
	baseNode
}

var _ Node = (*FieldNode)(nil)

func newFieldNode(parser *Parser) (Node, error) {
	if !parser.currentToken().isField() {
		return nil, NewError(
			INVALID_FIELD,
			parser.currentToken().LineNumber,
			parser.currentToken().End,
			parser.currentToken().Text,
			)
	}

	var node Node

	switch parser.currentToken().Type {
	case FUNC: 
		funcNode, err := newFuncNode(parser)
		if err != nil {
			return nil, err
		}
		node = funcNode
	default:
		fieldNode := FieldNode{}
		fieldNode.Token = parser.currentToken()

		if err := parser.nextToken(); err != nil {
			return nil, err
		}

		node = &fieldNode
	}

	if parser.currentToken().Type == AS {
		castNode, err := newCastNode(parser, node)
		if err != nil {
			return nil, err
		}

		return castNode, nil
	}

	return node, nil
}

func (f *FieldNode) Table() Node {
	return f
}

func (f *FieldNode) Columns() []string {
	return []string{f.String()}
}

func (f *FieldNode) String() string {
	if f.Token.Type == STRING {
		return "'" + f.Token.Text + "'"
	}

	return f.Token.Text
}

func (f *FieldNode) ReturnsValue() bool {
	return false
}

type ExpressionNode struct {
	baseNode
	precedence	int

	prefix,
	suffix 	Node
}

var _ Node = (*ExpressionNode)(nil)

func newExpressionNode(parser *Parser) (Node, error) {
	var node ExpressionNode

	node.Token = parser.currentToken()

	if err := parser.nextToken(); err != nil {
		return nil, err
	}

	precedence, ok := precedenceMap[node.Token.Type]
	if !ok {
		return nil, NewError(
			INVALID_TYPE_FOR_PRECEDENCE,
			node.Token.LineNumber,
			node.Token.End,
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
func (e *ExpressionNode) Columns() []string {
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

func (e *ExpressionNode) ReturnsValue() bool {
	return false
}

type OperatorNode struct {
	infix				[]baseNode
	prefix,
	suffix			Node
	precedence	int
}

var _ Node = (*OperatorNode)(nil)

func newOperatorNode(parser *Parser) (Node, error) {
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

		if !parser.currentToken().isOperator() {
			return nil, NewError(
				INVALID_OPERATOR_TOKEN,
				parser.currentToken().LineNumber,
				parser.currentToken().End,
				parser.currentToken().Text,
				)
		}

		infix := baseNode{
			Token: parser.currentToken(),
		}

		node.infix = append(node.infix, infix)

		precedence, ok := precedenceMap[infix.Token.Type]
		if !ok {
			return nil, NewError(
				INVALID_TYPE_FOR_PRECEDENCE,
				infix.Token.LineNumber,
				infix.Token.End,
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

func (o *OperatorNode) Columns() []string {
	return nil
}

func (o *OperatorNode) String() string {
	var buf bytes.Buffer

	if prefix := o.prefix; prefix != nil {
		if str := prefix.String(); len(str) > 0 {
			buf.WriteString(str)
		}
		
		if len(o.infix) != 0 {
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

func (o *OperatorNode) ReturnsValue() bool {
	return false
}

type CastNode struct {
	baseNode
	prefix,
	suffix		Node
}

var _ Node = (*CastNode)(nil)

func newCastNode(parser *Parser, prefix Node) (Node, error) {
	var node CastNode

	if parser.currentToken().Type != AS {
		return nil, NewError(
			INVALID_BASE_NODE_TYPE,
			parser.currentToken().LineNumber,
			parser.currentToken().End,
			parser.currentToken().Text,
			)
	}

	node.prefix = prefix

	node.Token = parser.currentToken()

	// consume AS token
	if err := parser.nextToken(); err != nil {
		return nil, err
	}
	
	suffix, err := newFieldNode(parser)
	if err != nil {
		return nil, err
	}

	node.suffix = suffix 

	return &node, nil
}

func (c *CastNode) Columns() []string {
	return []string{c.suffix.String()}
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

func (c *CastNode) ReturnsValue() bool {
	return false
}

// FUNC
type FuncNode struct {
	baseNode

	infix	[]Node
	suffix Node
}

var _ Node = (*FuncNode)(nil)

func newFuncNode(parser *Parser) (Node, error) {
	var node FuncNode

	if parser.currentToken().Type != FUNC {
		return nil, NewError(
			INVALID_BASE_NODE_TYPE,
			parser.currentToken().LineNumber,
			parser.currentToken().End,
			parser.currentToken().Text,
			)
	}

	node.Token = parser.currentToken()

	if err := parser.nextToken(); err != nil {
		return nil, err
	}

	if parser.currentToken().Type != L_PARAN {
		return nil, NewError(
			INVALID_FUNC_CHAR,
			parser.currentToken().LineNumber,
			parser.currentToken().End,
			"(",
			)
	}

	// consume left paran
	if err := parser.nextToken(); err != nil {
		return nil, err
	}
	
	for {
		if !parser.currentToken().isField() {
			break
		}

		infix, err := newFieldNode(parser)
		if err != nil {
			return nil, err
		}

		node.infix = append(node.infix, infix)
	}

	if parser.currentToken().Type != R_PARAN {
		return nil, NewError(
			INVALID_FUNC_CHAR,
			parser.currentToken().LineNumber,
			parser.currentToken().End,
			")",
			)
	}

	// consume right paran
	if err := parser.nextToken(); err != nil {
		return nil, err
	}

	return &node, nil
}

func (f *FuncNode) Columns() []string {
	return []string{f.String()}
}

func (f *FuncNode) Table() Node {
	return nil
}

func (f *FuncNode) String() string {
	var buf bytes.Buffer

	buf.WriteString(f.Token.Text)
	buf.WriteString("(")

	for i := range f.infix {
		buf.WriteString(f.infix[i].String())

		if i < (len(f.infix) - 1) {
			buf.WriteString(", ")
		}
	}

	buf.WriteString(")")

	return buf.String()
}

func (f *FuncNode) ReturnsValue() bool {
	return false
}

// list
type ListNode struct {
	nodes []Node
}

var _ Node = (*ListNode)(nil)

func newListNode(parser *Parser) (Node, error) {
	var list ListNode
	if parser.currentToken().Type != L_PARAN {
		return nil, NewError(
			INVALID_LIST_FORMAT,
			parser.currentToken().LineNumber,
			parser.currentToken().End,
			"(",
			)
	}

	// consume l paran
	if err := parser.nextToken(); err != nil {
		return nil, err
	}

	for {
		if parser.currentToken().Type == EOF {
			return nil, NewError(
				INVALID_LIST_FORMAT,
				parser.currentToken().LineNumber,
				parser.currentToken().End,
				")",
				)
		}

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

func (l *ListNode) Columns() []string {
	var columns []string

	for i := range l.nodes {
		columns = append(columns, l.nodes[i].String())
	}

	return columns
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

func (l *ListNode) ReturnsValue() bool {
	return false
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
			INVALID_BASE_NODE_TYPE,
			parser.currentToken().LineNumber,
			parser.currentToken().End,
			parser.currentToken().Text,
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

func (w *WhereNode) Columns() []string {
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

func (w *WhereNode) ReturnsValue() bool {
	return false
}

func consumeEndChar(parser *Parser) error {
	if parser.currentToken().Type != SEMICOL {
		return NewError(
			INVALID_END_OF_STATEMENT,
			parser.currentToken().LineNumber,
			parser.currentToken().End,
			";",
			)
	}

	// consume semicol token
	if err := parser.nextToken(); err != nil {
		return err
	}

	return nil
}
