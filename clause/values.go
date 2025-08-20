package clause

type Values struct {
	Columns []Column
	Values  [][]interface{}
}

func (Values) Name() string { return "VALUES" }

func (values Values) Build(builder Builder) {
	if len(values.Columns) > 0 {
		builder.WriteByte('(')
		for idx, column := range values.Columns {
			if idx > 0 {
				builder.WriteByte(',')
			}
			builder.WriteQuoted(column)
		}
		builder.WriteByte(')')
		for idx, value := range values.Values {
			if idx > 0 {
				builder.WriteByte(',')
			}
			builder.WriteByte('(')
			builder.AddVar(builder, value...)
			builder.WriteByte(')')
		}
	} else {
		builder.WriteString("DEFAULT VALUES")
	}
}

func (values Values) MergeClause(clause *Clause) {
	clause.Name = ""
	clause.Expression = values
}
