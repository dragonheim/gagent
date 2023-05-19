package picol

import (
	"unicode"
	"unicode/utf8"
)

// Define parser token types
const (
	ParserTokenESC = iota
	ParserTokenSTR
	ParserTokenCMD
	ParserTokenVAR
	ParserTokenSEP
	ParserTokenEOL
	ParserTokenEOF
)

// parserStruct represents the parser state
type parserStruct struct {
	text              string
	p, start, end, ln int
	insidequote       int
	Type              int
}

// initParser initializes a new parserStruct instance
func InitParser(text string) *parserStruct {
	return &parserStruct{text: text, ln: len(text), Type: ParserTokenEOL}
}

// next advances the parser position by one rune
func (p *parserStruct) next() {
	_, w := utf8.DecodeRuneInString(p.text[p.p:])
	p.p += w
	p.ln -= w
}

// current returns the current rune at the parser position
func (p *parserStruct) current() rune {
	r, _ := utf8.DecodeRuneInString(p.text[p.p:])
	return r
}

// token returns the current token text between start and end positions
func (p *parserStruct) token() (t string) {
	defer recover()
	return p.text[p.start:p.end]
}

// parseSep parses whitespace separators
func (p *parserStruct) parseSep() string {
	p.start = p.p
	for ; p.p < len(p.text); p.next() {
		if !unicode.IsSpace(p.current()) {
			break
		}
	}
	p.end = p.p
	p.Type = ParserTokenSEP
	return p.token()
}

// parseEol parses end of line and comments
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
	p.Type = ParserTokenEOL
	return p.token()
}

// parseCommand parses a command within brackets
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
	p.Type = ParserTokenCMD
	if p.p < len(p.text) && p.current() == ']' {
		p.next()
	}
	return p.token()
}

// parseVar parses a variable reference
func (p *parserStruct) parseVar() string {
	p.next() // skip the $
	p.start = p.p

	if p.current() == '{' {
		p.Type = ParserTokenVAR
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
		p.Type = ParserTokenSTR
	} else {
		p.end = p.p
		p.Type = ParserTokenVAR
	}
	return p.token()
}

// parseBrace parses a brace-enclosed string
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

// parseString parses a string with or without quotes
func (p *parserStruct) parseString() string {
	newword := p.Type == ParserTokenSEP || p.Type == ParserTokenEOL || p.Type == ParserTokenSTR

	if c := p.current(); newword && c == '{' {
		p.Type = ParserTokenSTR
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
				p.Type = ParserTokenESC
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
	p.Type = ParserTokenESC
	return p.token()
}

// parseComment skips over comment text
func (p *parserStruct) parseComment() string {
	for p.ln != 0 && p.current() != '\n' {
		p.next()
	}
	return p.token()
}

// GetToken returns the next token from the parser
func (p *parserStruct) GetToken() string {
	for {
		if p.ln == 0 {
			if p.Type != ParserTokenEOL && p.Type != ParserTokenEOF {
				p.Type = ParserTokenEOL
			} else {
				p.Type = ParserTokenEOF
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
			if p.Type == ParserTokenEOL {
				p.parseComment()
				continue
			}
			return p.parseString()
		default:
			return p.parseString()
		}
	}
}
