package picol_test

import (
	"testing"

	"github.com/dragonheim/gagent/pkg/picol"
)

func TestParser(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected []int
	}{
		{
			"Simple test",
			"set x 10\nincr x",
			[]int{
				picol.ParserTokenSTR,
				picol.ParserTokenSEP,
				picol.ParserTokenSTR,
				picol.ParserTokenSEP,
				picol.ParserTokenSTR,
				picol.ParserTokenEOL,
				picol.ParserTokenSTR,
				picol.ParserTokenSEP,
				picol.ParserTokenSTR,
				picol.ParserTokenEOL,
				picol.ParserTokenEOF,
			},
		},
		{
			"Variable and command test",
			"set x $y\nputs [expr $x * 2]",
			[]int{
				picol.ParserTokenSTR,
				picol.ParserTokenSEP,
				picol.ParserTokenSTR,
				picol.ParserTokenSEP,
				picol.ParserTokenVAR,
				picol.ParserTokenEOL,
				picol.ParserTokenSTR,
				picol.ParserTokenSEP,
				picol.ParserTokenCMD,
				picol.ParserTokenEOL,
				picol.ParserTokenEOF,
			},
		},
		{
			"Braces and quotes test",
			`set x {"Hello World"}`,
			[]int{
				picol.ParserTokenSTR,
				picol.ParserTokenSEP,
				picol.ParserTokenSTR,
				picol.ParserTokenSEP,
				picol.ParserTokenSTR,
				picol.ParserTokenEOL,
				picol.ParserTokenEOF,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			parser := picol.InitParser(tc.input)

			for _, expectedType := range tc.expected {
				token := parser.GetToken()
				if parser.Type != expectedType {
					t.Errorf("Expected token type %d, got %d", expectedType, parser.Type)
				}
				if parser.Type == picol.ParserTokenEOF {
					break
				}
			}
		})
	}
}
