package wirth

func Parse(input []byte) (*Syntax, error) {
	p := &parser{tokens: Scan(input)}
	p.eof = p.tokens[len(p.tokens)-1]

	syntax := p.parse()
	return syntax, nil
}

type parser struct {
	tokens []*Token // all the tokens in the input
	eof    *Token   // last token in the input
	peek   *Token   // one token look-ahead
}

// parser parses a grammar file.
func (p *parser) parse() (syntax *Syntax) {
	syntax = &Syntax{Productions: make(map[string]*Production)}

	return syntax
}
