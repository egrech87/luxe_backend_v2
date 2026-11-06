package daoOperator

type Operator string

const (
	Equals        Operator = "="
	NotEquals     Operator = "!="
	Greater       Operator = ">"
	Smaller       Operator = "<"
	GreaterEquals Operator = ">="
	SmallerEquals Operator = "<="
	Is            Operator = "IS"
	IsNot         Operator = "IS NOT"
	Between       Operator = "BETWEEN"
	Like          Operator = "LIKE"
	In            Operator = "IN"
)
