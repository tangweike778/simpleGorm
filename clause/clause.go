package clause

// Interface clause interface
type Interface interface {
	Name() string
	Build(Builder)
	MergeClause(*Clause)
}

type Writer interface {
	WriteByte(byte) error
	WriteString(string) (int, error)
}

// Builder builder interface
type Builder interface {
	Writer
	WriteQuoted(field interface{})
	AddVar(Writer, ...interface{})
	AddError(error) error
}

// Clause
type Clause struct {
	Name       string // WHERE
	Expression Expression
	//Builder    ClauseBuilder
}

// Build build clause
func (c Clause) Build(builder Builder) {
	if c.Expression != nil {
		if c.Name != "" {
			builder.WriteString(c.Name)
			builder.WriteByte(' ')
		}

		c.Expression.Build(builder)
	}
}

// Table quote with name
type Table struct {
	Name string
	Raw  bool
}

type Column struct {
	Table string
	Name  string
	Raw   bool
}
