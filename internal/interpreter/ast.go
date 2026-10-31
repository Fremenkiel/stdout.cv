package interpreter

import (
	"bytes"
)

type Node interface {
	Table()					Node
	Columns()				[]string
	String()				string
}

type baseToken struct {
	token Token
}

type BaseNode interface {
	Node
	Token()					Token
	ReturnsValue() 	bool
	AltersSchema()	bool
	AltersRows()		bool
}

type baseNode struct {
	baseToken
}

func (b *baseNode) ReturnsValue() bool {
	return b.Token().Type.returnsValue()
}

func (b *baseNode) AltersSchema() bool {
	return b.Token().Type.altersSchema()
}

func (b *baseNode) AltersRows() bool {
	return b.Token().Type.altersRows()
}

func (b *baseNode) Token() Token {
	return b.token
}

// INSERT
type InsertNode struct {
	baseNode

	columns Node
	into    Node
	values 	Node
}

var _ BaseNode = (*InsertNode)(nil)

func newInsetNode(parser *Parser) (BaseNode, error) {
	if parser.currentToken().Type != INSERT {
		return nil, NewError(
			INVALID_INSERT_STATEMENT,
			parser.currentToken().LineNumber,
			parser.currentToken().End,
			parser.currentToken().Text,
			)
	}

	node := InsertNode{
		token:  parser.currentToken(),
	}

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

	buf.WriteString(i.token.Text)
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

var _ BaseNode = (*SelectNode)(nil)

func newSelectNode(parser *Parser) (BaseNode, error) {
	if parser.currentToken().Type != SELECT {
		return nil, NewError(
			INVALID_SELECT_STATEMENT,
			parser.currentToken().LineNumber,
			parser.currentToken().End,
			parser.currentToken().Text,
			)
	}

	node := SelectNode{
		token: parser.currentToken(),
	}

	// consume select token
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

		fieldNode, err := newFieldNode(parser)
		if err != nil {
			return nil, err
		}

		node.columns = append(node.columns, fieldNode)

		if parser.currentToken().Type != COMMA {
			if parser.currentToken().isField() {
				return nil, NewError(
					INVALID_SELECT_STATEMENT,
					parser.currentToken().LineNumber,
					parser.currentToken().End,
					parser.currentToken().Text,
					)
			}
			break
		}

		// consume comma token
		if err := parser.nextToken(); err != nil {
			return nil, err
		}
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

// UPDATE
type UpdateNode struct {
	baseNode

	table				Node
	set					[]Node
	where				Node
}

var _ BaseNode = (*UpdateNode)(nil)

func newUpdateNode(parser *Parser) (BaseNode, error) {
	if parser.currentToken().Type != UPDATE {
		return nil, NewError(
			INVALID_UPDATE_STATEMENT,
			parser.currentToken().LineNumber,
			parser.currentToken().End,
			parser.currentToken().Text,
			)
	}

	node := UpdateNode{
		token: parser.currentToken(),
	}

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

		set, err := newOperatorNode(parser)
		if err != nil {
			return nil, err
		}

		node.set = append(node.set, set)

		if parser.currentToken().Type != COMMA {
			if parser.currentToken().isField() {
				return nil, NewError(
					INVALID_UPDATE_STATEMENT,
					parser.currentToken().LineNumber,
					parser.currentToken().End,
					parser.currentToken().Text,
					)
			}
			break
		}

		// consume comma token
		if err := parser.nextToken(); err != nil {
			return nil, err
		}
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

// DELETE
type DeleteNode struct {
	baseNode

	from    Node
	where		Node
}

var _ BaseNode = (*DeleteNode)(nil)

func newDeleteNode(parser *Parser) (BaseNode, error) {
	if parser.currentToken().Type != DELETE {
		return nil, NewError(
			INVALID_DELETE_STATEMENT,
			parser.currentToken().LineNumber,
			parser.currentToken().End,
			parser.currentToken().Text,
			)
	}

	node := DeleteNode{
		token: parser.currentToken(),
	}

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

type DropNode struct {
	baseNode

	modifier	[]Node
	table			Node
}

var _ BaseNode = (*DropNode)(nil)

func newDropNode(parser *Parser) (BaseNode, error) {
	if parser.currentToken().Type != DROP {
		return nil, NewError(
			INVALID_DROP_STATEMENT,
			parser.currentToken().LineNumber,
			parser.currentToken().End,
			parser.currentToken().Text,
			)
	}

	node := DropNode{
		token: parser.currentToken(),
	}

	// consume DROP
	if err := parser.nextToken(); err != nil {
		return nil, err
	}

	if parser.currentToken().Type != TABLE {
		return nil, NewError(
			INVALID_DROP_STATEMENT,
			parser.currentToken().LineNumber,
			parser.currentToken().End,
			parser.currentToken().Text,
			)
	}

	// consume TABLE
	if err := parser.nextToken(); err != nil {
		return nil, err
	}

	if parser.currentToken().Type == IF {
		if parser.peakToken().Type != EXISTS {
			return nil, NewError(
				INVALID_DROP_STATEMENT,
				parser.currentToken().LineNumber,
				parser.currentToken().End,
				parser.currentToken().Text,
				)
		}

		modifier, err := newOperatorNode(parser)
		if err != nil {
			return nil, err
		}

		node.modifier = append(node.modifier, modifier)
	}

	table, err := newFieldNode(parser)
	if err != nil {
		return nil, err
	}

	node.table = table

	if err := consumeEndChar(parser); err != nil {
		return nil, err
	}

	return &node, nil
}

func (d *DropNode) Table() Node {
	return d.table.Table()
}

func (d *DropNode) Columns() []string {
	return nil
}

func (d *DropNode) String() string {
	var buf bytes.Buffer

	buf.WriteString("DROP TABLE ")

	for i := range d.modifier {
		buf.WriteString(d.modifier[i].String())
		buf.WriteString(" ")
	}

	buf.WriteString(d.table.String())
	buf.WriteString(";")

	return buf.String()
}

type CreateNode struct {
	baseNode

	table 	Node
	columns	[]Node
}

var _ BaseNode = (*CreateNode)(nil)

func newCreateNode(parser *Parser) (BaseNode, error) {
	if parser.currentToken().Type != CREATE {
		return nil, NewError(
			INVALID_CREATE_STATEMENT,
			parser.currentToken().LineNumber,
			parser.currentToken().End,
			parser.currentToken().Text,
			)
	}

	node := CreateNode{
		token: parser.currentToken(),
	}

	// consume CREATE
	if err := parser.nextToken(); err != nil {
		return nil, err
	}

	if parser.currentToken().Type != TABLE {
		return nil, NewError(
			INVALID_CREATE_STATEMENT,
			parser.currentToken().LineNumber,
			parser.currentToken().End,
			parser.currentToken().Text,
			)
	}

	// consume TABLE
	if err := parser.nextToken(); err != nil {
		return nil, err
	}

	table, err := newFieldNode(parser)
	if err != nil {
		return nil, err
	}

	node.table = table

	if parser.currentToken().Type != L_PARAN {
		return nil, NewError(
			INVALID_CREATE_STATEMENT,
			parser.currentToken().LineNumber,
			parser.currentToken().End,
			parser.currentToken().Text,
			)
	}

	// consume L PARAN
	if err := parser.nextToken(); err != nil {
		return nil, err
	}

	for {
		if parser.currentToken().Type == EOF {
			return nil, NewError(
				INVALID_CREATE_STATEMENT,
				parser.currentToken().LineNumber,
				parser.currentToken().End,
				")",
				)
		}
		
		if parser.currentToken().Type == R_PARAN {
			break
		}

		definition, err := newDefinitionNode(parser)
		if err != nil {
			return nil, err
		}

		node.columns = append(node.columns, definition)

		if parser.currentToken().Type == COMMA {
			// consume COMMA
			if err := parser.nextToken(); err != nil {
				return nil, err
			}
		}
	}

	// consume R PARAN
	if err := parser.nextToken(); err != nil {
		return nil, err
	}

	if err := consumeEndChar(parser); err != nil {
		return nil, err
	}

	return &node, nil
}

func (c *CreateNode) Table() Node {
	return c.table.Table()
}

func (c *CreateNode) Columns() []string {
	var slice = make([]string, 0, len(c.columns))

	for i := range c.columns {
		slice = append(slice, c.columns[i].Columns()...)
	}

	return slice
}

func (c *CreateNode) String() string {
	var buf bytes.Buffer

	buf.WriteString("CREATE TABLE ")

	buf.WriteString(c.table.String())
	buf.WriteString("(")

	for i := range c.columns {
		buf.WriteString(c.columns[i].String())

		if i < len(c.columns) - 1 {
			buf.WriteString(", ")
		}
	}

	buf.WriteString(");")

	return buf.String()
}

type AlterNode struct {
	baseNode

	table 			Node
	operator		Node
	definition	Node

	renamesTable bool
}

var _ BaseNode = (*AlterNode)(nil)

func newAlterNode(parser *Parser) (BaseNode, error) {
	if parser.currentToken().Type != ALTER {
		return nil, NewError(
			INVALID_ALTER_STATEMENT,
			parser.currentToken().LineNumber,
			parser.currentToken().End,
			parser.currentToken().Text,
			)
	}

	node := AlterNode{
		token: parser.currentToken(),
	}

	// consume ALTER
	if err := parser.nextToken(); err != nil {
		return nil, err
	}

	if parser.currentToken().Type != TABLE {
		return nil, NewError(
			INVALID_ALTER_STATEMENT,
			parser.currentToken().LineNumber,
			parser.currentToken().End,
			parser.currentToken().Text,
			)
	}

	// consume TABLE
	if err := parser.nextToken(); err != nil {
		return nil, err
	}

	table, err := newFieldNode(parser)
	if err != nil {
		return nil, err
	}

	node.table = table

	actionMap, ok := alterMap[parser.currentToken().Type]
	if !ok {
		return nil, NewError(
			INVALID_ALTER_STATEMENT,
			parser.currentToken().LineNumber,
			parser.currentToken().End,
			parser.currentToken().Text,
			)
	}

	if _, ok := actionMap[parser.peakToken().Type]; !ok {
		return nil, NewError(
			INVALID_ALTER_STATEMENT,
			parser.currentToken().LineNumber,
			parser.currentToken().End,
			parser.currentToken().Text,
			)
	}

	node.renamesTable = parser.currentToken().Type == RENAME &&
		parser.peakToken().Type == TO

	operator, err := newOperatorNode(parser)
	if err != nil {
		return nil, err
	}

	node.operator = operator

	definition, err := newDefinitionNode(parser)
	if err != nil {
		return nil, err
	}

	node.definition = definition 

	if err := consumeEndChar(parser); err != nil {
		return nil, err
	}

	return &node, nil
}

func (a *AlterNode) Table() Node {
	return a.table.Table()
}

func (a *AlterNode) Columns() []string {
	return nil
}

func (a *AlterNode) String() string {
	var buf bytes.Buffer

	buf.WriteString("ALTER TABLE ")
	buf.WriteString(a.table.String())
	buf.WriteString(" ")

	buf.WriteString(a.operator.String())
	buf.WriteString(" ")
	buf.WriteString(a.definition.String())

	buf.WriteString(";")

	return buf.String()
}

func (a *AlterNode) NewTableName() string {
	if a.renamesTable {
		return a.definition.String()
	}

	return ""
}

// column and table
type FieldNode struct {
	baseToken
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

	switch {
	case parser.peakToken().Type == L_PARAN && !parser.currentIsTable(): 
		funcNode, err := newFuncNode(parser)
		if err != nil {
			return nil, err
		}
		node = funcNode
	default:
		fieldNode := FieldNode{
			token: parser.currentToken(),
		}

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
	if f.token.Type == STRING {
		return "'" + f.token.Text + "'"
	}

	return f.token.Text
}

type DefinitionNode struct {
	baseToken

	definition	Node
}

var _ Node = (*DefinitionNode)(nil)

func newDefinitionNode(parser *Parser) (Node, error) {
	if parser.currentToken().Type != FIELD {
		return nil, NewError(
			INVALID_COLUMN_DEFINITION,
			parser.currentToken().LineNumber,
			parser.currentToken().End,
			parser.currentToken().Text,
			)
	}

	node := DefinitionNode{
		token: parser.currentToken(),
	}

	if err := parser.nextToken(); err != nil {
		return nil, err
	}

	definition, err := newOperatorNode(parser)
	if err != nil {
		return nil, err
	}

	node.definition = definition

	return &node, nil
}

func (c *DefinitionNode) Table() Node {
	return nil
}

func (c *DefinitionNode) Columns() []string {
	return []string{c.token.Text}
}

func (c *DefinitionNode) String() string {
	var buf bytes.Buffer

	buf.WriteString(c.token.Text)

	if definition := c.definition.String(); len(definition) > 0 {
		buf.WriteString(" ")
		buf.WriteString(definition)
	}

	return buf.String()
}

type ExpressionNode struct {
	baseToken
	precedence	int

	prefix,
	suffix 	Node
}

var _ Node = (*ExpressionNode)(nil)

func newExpressionNode(parser *Parser) (Node, error) {
	node := ExpressionNode{
		token: parser.currentToken(),
	}

	if err := parser.nextToken(); err != nil {
		return nil, err
	}

	precedence, ok := precedenceMap[node.token.Type]
	if !ok {
		return nil, NewError(
			INVALID_TYPE_FOR_PRECEDENCE,
			node.token.LineNumber,
			node.token.End,
			node.token.Text,
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

	buf.WriteString(e.token.Text)

	if suffix := e.suffix; suffix != nil {
		if str := suffix.String(); len(str) > 0 {
			buf.WriteString(" ")
			buf.WriteString(str)
		}
	}

	return buf.String()
}

type OperatorNode struct {
	infix				[]baseToken
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

	outerloop:
	for {
		switch {
		case parser.currentToken().Type == EOF:
			return &node, nil
		case parser.currentToken().Type == SEMICOL:
			return &node, nil
		case parser.currentToken().Type == COMMA:
			return &node, nil
		case parser.currentToken().Type == R_PARAN:
			return &node, nil
		case parser.currentToken().isField():
			break outerloop
		case !parser.currentToken().isOperator():
			return nil, NewError(
				INVALID_OPERATOR_TOKEN,
				parser.currentToken().LineNumber,
				parser.currentToken().End,
				parser.currentToken().Text,
				)
		}

		infix := baseToken{
			token: parser.currentToken(),
		}

		node.infix = append(node.infix, infix)

		precedence, ok := precedenceMap[infix.token.Type]
		if !ok {
			return nil, NewError(
				INVALID_TYPE_FOR_PRECEDENCE,
				infix.token.LineNumber,
				infix.token.End,
				infix.token.Text,
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
		buf.WriteString(o.infix[i].token.Text)

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
	baseToken
	prefix,
	suffix		Node
}

var _ Node = (*CastNode)(nil)

func newCastNode(parser *Parser, prefix Node) (Node, error) {
	if parser.currentToken().Type != AS {
		return nil, NewError(
			INVALID_BASE_NODE_TYPE,
			parser.currentToken().LineNumber,
			parser.currentToken().End,
			parser.currentToken().Text,
			)
	}

	node := CastNode{
		prefix: prefix,
		token: parser.currentToken(),
	}

	// consume AS token
	if err := parser.nextToken(); err != nil {
		return nil, err
	}

	if parser.currentToken().Type != FIELD {
		return nil, NewError(
			INVALID_CAST_FORMAT,
			parser.currentToken().LineNumber,
			parser.currentToken().End,
			parser.currentToken().Text,
			)
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
	buf.WriteString(c.token.Text)
	buf.WriteString(" ")
	buf.WriteString(c.suffix.String())

	return buf.String()
}

// FUNC
type FuncNode struct {
	baseToken

	infix	[]Node
	suffix Node
}

var _ Node = (*FuncNode)(nil)

func newFuncNode(parser *Parser) (Node, error) {
	if parser.currentToken().Type != FIELD {
		return nil, NewError(
			INVALID_BASE_NODE_TYPE,
			parser.currentToken().LineNumber,
			parser.currentToken().End,
			parser.currentToken().Text,
			)
	}

	node := FuncNode{
		token: parser.currentToken(),
	}

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

	buf.WriteString(f.token.Text)
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

		valueNode, err := newFieldNode(parser)
		if err != nil {
			return nil, err
		}

		list.nodes = append(list.nodes, valueNode)

		if parser.currentToken().Type != COMMA {
			if parser.currentToken().isField() {
				return nil, NewError(
					INVALID_LIST_FORMAT,
					parser.currentToken().LineNumber,
					parser.currentToken().End,
					parser.currentToken().Text,
					)
			}
			break
		}

		// consume comma
		if err := parser.nextToken(); err != nil {
			return nil, err
		}
	}

	if parser.currentToken().Type != R_PARAN {
		return nil, NewError(
			INVALID_LIST_FORMAT,
			parser.currentToken().LineNumber,
			parser.currentToken().End,
			")",
			)
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

// WHERE
type WhereNode struct {
	baseToken
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

	node := WhereNode{
		token: parser.currentToken(),
	}

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
