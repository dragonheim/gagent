package picol

import (
	"unicode"
	"unicode/utf8"
)

/*
 * pt_ESC is @TODO
 */
const (
	pt_ESC = iota
	pt_STR
	pt_CMD
	pt_VAR
	pt_SEP
	pt_EOL
	pt_EOF
)

type parserStruct struct {
	text              string
	p, start, end, ln int
	insidequote       int
	Type              int
}

func InitParser(text string) *parserStruct {
	return &parserStruct{text, 0, 0, 0, len(text), 0, pt_EOL}
}

func (p *parserStruct) next() {
	_, w := utf8.DecodeRuneInString(p.text[p.p:])
	p.p += w
	p.ln -= w
}

func (p *parserStruct) current() rune {
	r, _ := utf8.DecodeRuneInString(p.text[p.p:])
	return r
}

func (p *parserStruct) token() (t string) {
	defer recover()
	return p.text[p.start:p.end]
}

func (p *parserStruct) parseSep() string {
	p.start = p.p
	for ; p.p < len(p.text); p.next() {
		if !unicode.IsSpace(p.current()) {
			break
		}
	}
	p.end = p.p
	p.Type = pt_SEP
	return p.token()
}

func (p *parserStruct) parseEol() string {
	p.start = p.p

	for ; p.p < len(p.text); p.next() {
		if p.current() == ';' || unicode.IsSpace(p.current()) {
			// pass
		} else {
			break
		}
	}

	p.end = p.p
	p.Type = pt_EOL
	return p.token()
}

func (p *parserStruct) parseCommand() string {
	level, blevel := 1, 0
	p.next() // skip
	p.start = p.p
Loop:
	for {
		switch {
		case p.ln == 0:
			break Loop
		case p.current() == '[' && blevel == 0:
			level++
		case p.current() == ']' && blevel == 0:
			level--
			if level == 0 {
				break Loop
			}
		case p.current() == '\\':
			p.next()
		case p.current() == '{':
			blevel++
		case p.current() == '}' && blevel != 0:
			blevel--
		}
		p.next()
	}
	p.end = p.p
	p.Type = pt_CMD
	if p.p < len(p.text) && p.current() == ']' {
		p.next()
	}
	return p.token()
}

func (p *parserStruct) parseVar() string {
	p.next() // skip the $
	p.start = p.p

	if p.current() == '{' {
		p.Type = pt_VAR
		return p.parseBrace()
	}

	for p.p < len(p.text) {
		c := p.current()
		if unicode.IsLetter(c) || ('0' <= c && c <= '9') || c == '_' {
			p.next()
			continue
		}
		break
	}

	if p.start == p.p { // It's just a single char string "$"
		p.start = p.p - 1
		p.end = p.p
		p.Type = pt_STR
	} else {
		p.end = p.p
		p.Type = pt_VAR
	}
	return p.token()
}

func (p *parserStruct) parseBrace() string {
	level := 1
	p.next() // skip
	p.start = p.p

Loop:
	for p.p < len(p.text) {
		c := p.current()
		switch {
		case p.ln >= 2 && c == '\\':
			p.next()
		case p.ln == 0 || c == '}':
			level--
			if level == 0 || p.ln == 0 {
				break Loop
			}
		case c == '{':
			level++
		}
		p.next()
	}
	p.end = p.p
	if p.ln != 0 { // Skip final closed brace
		p.next()
	}
	return p.token()
}

func (p *parserStruct) parseString() string {
	newword := p.Type == pt_SEP || p.Type == pt_EOL || p.Type == pt_STR

	if c := p.current(); newword && c == '{' {
		p.Type = pt_STR
		return p.parseBrace()
	} else if newword && c == '"' {
		p.insidequote = 1
		p.next() // skip
	}

	p.start = p.p

Loop:
	for ; p.ln != 0; p.next() {
		switch p.current() {
		case '\\':
			if p.ln >= 2 {
				p.next()
			}
		case '$', '[':
			break Loop
		case '"':
			if p.insidequote != 0 {
				p.end = p.p
				p.Type = pt_ESC
				p.next()
				p.insidequote = 0
				return p.token()
			}
		}
		if p.current() == ';' || unicode.IsSpace(p.current()) {
			if p.insidequote == 0 {
				break Loop
			}
		}
	}

	p.end = p.p
	p.Type = pt_ESC
	return p.token()
}

func (p *parserStruct) parseComment() string {
	for p.ln != 0 && p.current() != '\n' {
		p.next()
	}
	return p.token()
}

func (p *parserStruct) GetToken() string {
	for {
		if p.ln == 0 {
			if p.Type != pt_EOL && p.Type != pt_EOF {
				p.Type = pt_EOL
			} else {
				p.Type = pt_EOF
			}
			return p.token()
		}

		switch p.current() {
		case ' ', '\t', '\r':
			if p.insidequote != 0 {
				return p.parseString()
			}
			return p.parseSep()
		case '\n', ';':
			if p.insidequote != 0 {
				return p.parseString()
			}
			return p.parseEol()
		case '[':
			return p.parseCommand()
		case '$':
			return p.parseVar()
		case '#':
			if p.Type == pt_EOL {
				p.parseComment()
				continue
			}
			return p.parseString()
		default:
			return p.parseString()
		}
	}
	/*	return p.token() /* unreached */
}
