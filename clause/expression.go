package clause

// Expression 表达式接口
type Expression interface {
	Build(builder Builder)
}
