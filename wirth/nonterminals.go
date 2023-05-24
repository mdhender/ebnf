package wirth

type Syntax struct {
	Start       *Token
	Productions map[string]*Production
}

type Production struct {
	Name       *Token
	Expression *Expression
}

type Expression struct {
	Terms []*Term
}

type Term struct {
	Factors []*Factor
}

type Factor struct {
	NonTerminal *Token
	Terminal    *Token
	Group       *Expression
	Option      *Expression
	Repetition  *Expression
}
