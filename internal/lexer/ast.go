package lexer

type Node interface {
	TableName() Node
	Columns() []Node
}

type baseNode struct {
	Token	Token
}

type Statemant struct {
	baseNode
	nodes	[]Node
}
var _ Node = (*Statemant)(nil)

func (s *Statemant) TableName() Node {
	switch s.Token.Type {
	case UPDATE:
		return s.nodes[0]
	}

	for i := range s.nodes {
		if node := s.nodes[i].TableName(); node != nil {
			return node
		}
	}

	return nil
}

func (s *Statemant) Columns() []Node {
	switch s.Token.Type {
	case SELECT:
		return s.nodes[:len(s.nodes) - 1]
	}

	for i := range s.nodes {
		if node := s.nodes[i].Columns(); node != nil {
			return node
		}
	}

	return nil
}

type Expression struct {
	baseNode
	nodes []Node
}
var _ Node = (*Expression)(nil)

func (e *Expression) TableName() Node {
	switch e.Token.Type {
	case INTO:
		return e.nodes[0]
	case FROM:
		return e.nodes[0]
	}

	for i := range e.nodes {
		if node := e.nodes[i].TableName(); node != nil {
			return node
		}
	}

	return nil
}

func (e *Expression) Columns() []Node {
	switch e.Token.Type {
	case INTO:
		return e.nodes[1].Columns()
	case SET:
		return e.nodes[:len(e.nodes) - 1]
	}

	return nil
}

type Field struct {
	baseNode
	node	Node
}
var _ Node = (*Field)(nil)

func (f *Field) TableName() Node {
	return f.node.TableName()
}

func (f *Field) Columns() []Node {
	return f.node.Columns()
}

type Operator struct {
	baseNode
	nodes	[]Node
}
var _ Node = (*Operator)(nil)

func (o *Operator) TableName() Node {
	for i := range o.nodes {
		if node := o.nodes[i].TableName(); node != nil {
			return node
		}
	}

	return nil
}

func (o *Operator) Columns() []Node {
	for i := range o.nodes {
		if node := o.nodes[i].Columns(); node != nil {
			return node
		}
	}

	return nil
}

