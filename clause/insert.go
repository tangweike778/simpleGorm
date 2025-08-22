package clause

const (
	CurrentTable string = "~~~ct~~~"
)

var (
	currentTable = Table{Name: CurrentTable}
)

type Insert struct {
	Table Table
}

func (insert Insert) Name() string {
	return "INSERT"
}

func (insert Insert) Build(builder Builder) {
	builder.WriteString("INTO ")
	if insert.Table.Name != "" {
		builder.WriteQuoted(insert.Table)
	} else {
		builder.WriteQuoted(currentTable)
	}
}

func (insert Insert) MergeClause(clause *Clause) {
	if v, ok := clause.Expression.(Insert); ok {
		if insert.Table.Name == "" {
			insert.Table = v.Table
		}
	}
	clause.Expression = insert
}
