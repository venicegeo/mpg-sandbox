package dsl

// shunting yard implementation adapted from https://github.com/mgenware/go-shunting-yard

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"text/scanner"
)

//===========================================================================
var arrayTypeRegexp *regexp.Regexp

func init() {
	arrayTypeRegexp = regexp.MustCompile(`^\[(\d+)\]`)
}

//===========================================================================

type TokenId int

const (
	TokenInvalid TokenId = iota
	TokenEquals
	TokenNotEquals
	TokenGreaterThan
	TokenGreaterOrEqualThan
	TokenLessThan
	TokenLessOrEqualThan
	TokenAdd
	TokenSubtract
	TokenMultiply
	TokenDivide
	TokenExponent
	TokenMod
	TokenBitwiseOr
	TokenBitwiseAnd
	TokenLogicalAnd
	TokenLogicalOr
	TokenLeftParen
	TokenRightParen
	TokenLeftBracket
	TokenRightBracket
	TokenSymbol // 16
	TokenNumber
	TokenTypeSlice
	TokenTypeArray
	TokenTypeMap
)

type Token struct {
	Line   int
	Column int
	Text   string
	Id     TokenId
	Value  interface{}
}

func (t *Token) String() string {
	s := fmt.Sprintf("[%d:%d] id=%d text=\"%s\"", t.Line, t.Column, t.Id, t.Text)
	if t.Value != nil {
		s += fmt.Sprintf(" value=<%v>", t.Value)
	}
	return s
}

func convertId(r rune) TokenId {
	switch r {
	case -2:
		return TokenSymbol
	case -3:
		return TokenNumber
	case 40:
		return TokenLeftParen
	case 41:
		return TokenRightParen
	case 42:
		return TokenMultiply
	case 43:
		return TokenAdd
	case 60:
		return TokenLessThan
	case 62:
		return TokenGreaterThan
	case 91:
		return TokenLeftBracket
	case 93:
		return TokenRightBracket
	case 124:
		return TokenBitwiseOr
	default:
		return TokenInvalid
	}
}

//===========================================================================

type Scanner struct{}

func (s *Scanner) Scan(str string) ([]Token, error) {
	tokens, err := s.scan(str)
	if err != nil {
		return nil, err
	}

	tokens, err = peepholer(tokens)
	if err != nil {
		return nil, err
	}

	return tokens, nil
}

//===========================================================================

func (s *Scanner) scan(input string) ([]Token, error) {

	tokens := []Token{}

	var sx scanner.Scanner
	sx.Init(strings.NewReader(input))

	var tok rune

	for {
		tok = sx.Scan()
		if tok == scanner.EOF {
			break
		}

		id := convertId(tok)
		if id == TokenInvalid {
			return nil, fmt.Errorf("Unknown token type %v (%s)", tok, sx.TokenText())
		}
		token := Token{
			Line:   sx.Pos().Line,
			Column: sx.Pos().Column,
			Text:   sx.TokenText(),
			Id:     id,
		}
		//log.Printf("TOK: %s", token.String())
		tokens = append(tokens, token)
	}

	return tokens, nil
}

func peepholer(tokens []Token) ([]Token, error) {

	result := []Token{}
	push := func(t Token) {
		result = append(result, t)
	}

	// combine these two tokens:
	//   ||, &&
	//   []
	// combine these three tokens:
	//   [map]
	//   [int]

	// the last index in use
	last := len(tokens) - 1

	for i := 0; i <= last; i += 0 { // TODO
		atLeastTwoLeft := i <= last-1
		atLeastThreeLeft := i <= last-2

		if atLeastThreeLeft {
			s := tokens[i].Text + tokens[i+1].Text + tokens[i+2].Text
			arrayMatch, arrayLen := matchArrayTypePrefix(s)

			switch {
			case s == "[map]":
				t := Token{
					Line:   tokens[i].Line,
					Column: tokens[i].Column,
					Text:   s,
					Id:     TokenTypeMap,
				}
				push(t)
				i += 3
				continue
			case arrayMatch: // [123]
				t := Token{
					Line:   tokens[i].Line,
					Column: tokens[i].Column,
					Text:   s,
					Id:     TokenTypeArray,
					Value:  arrayLen,
				}
				push(t)
				i += 3
				continue
			}
		}

		if atLeastTwoLeft {
			s := tokens[i].Text + tokens[i+1].Text

			switch s {
			case "||":
				t := Token{
					Line:   tokens[i].Line,
					Column: tokens[i].Column,
					Text:   s,
					Id:     TokenLogicalOr,
				}
				push(t)
				i += 2
				continue
			case "&&":
				t := Token{
					Line:   tokens[i].Line,
					Column: tokens[i].Column,
					Text:   s,
					Id:     TokenLogicalAnd,
				}
				push(t)
				i += 2
				continue
			case "[]":
				t := Token{
					Line:   tokens[i].Line,
					Column: tokens[i].Column,
					Text:   s,
					Id:     TokenTypeSlice,
				}
				push(t)
				i += 2
				continue
			}
		}

		// no peephole match
		push(tokens[i])
		i++
	}

	return result, nil
}

func matchArrayTypePrefix(s string) (bool, int) {
	ok := arrayTypeRegexp.Match([]byte(s))
	if !ok {
		return false, -1
	}
	sub := arrayTypeRegexp.FindSubmatch([]byte(s))

	siz, err := strconv.Atoi(string(sub[1]))
	if err != nil || siz < 1 {
		panic(err)
	}
	return true, siz
}