package jmespath

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

type astNodeType int

//go:generate stringer -type=astNodeType
const (
	ASTComparator astNodeType = iota
	ASTCurrentNode
	ASTExpRef
	ASTFunctionExpression
	ASTField
	ASTFilterProjection
	ASTFlatten
	ASTIdentity
	ASTIndex
	ASTIndexExpression
	ASTKeyValPair
	ASTLiteral
	ASTMultiSelectHash
	ASTMultiSelectList
	ASTOrExpression
	ASTAndExpression
	ASTNotExpression
	ASTPipe
	ASTProjection
	ASTSubexpression
	ASTSlice
	ASTValueProjection
)

// ASTNode represents the abstract syntax tree of a JMESPath expression.
type ASTNode struct {
	nodeType astNodeType
	value    interface{}
	children []ASTNode
}

func (node ASTNode) String() string {
	return node.PrettyPrint(0)
}

func (node ASTNode) PrettyPrint(indent int) string {
	dotdotdot := ""
	if len(node.children) > 3 {
		dotdotdot = "..."
	}
	firstChildren := node.children
	if len(firstChildren) > 3 {
		firstChildren = firstChildren[:3]
	}
	prettyChildren := make([]string, len(firstChildren))
	for i, child := range firstChildren {
		prettyChildren[i] = child.PrettyPrint(indent + 1)
	}
	displayValue := node.value
	if displayValue == nil {
		displayValue = ""
	}
	return fmt.Sprintf("ASTNode{%s, value: %v, children: [%s%s]}",
		node.nodeType, displayValue,
		strings.Join(prettyChildren, ", "), dotdotdot)
}

type token struct {
	tokType  tokType
	value    interface{}
	position int
	length   int
}

type parser struct {
	tokens      []token
	index       int
	*Lexer
}

// NewParser creates a new JMESPath parser.
func NewParser() *parser {
	p := &parser{}
	return p
}

// Parse parses a JMESPath expression and returns an ASTNode.
func (p *parser) Parse(expression string) (ASTNode, error) {
	nodeType := ASTEmpty
	p.loadTokens(expression)
	if len(p.tokens) == 0 {
		return ASTNode{}, SyntaxError{
			msg:        "Empty expression",
			Expression: expression,
			Offset:     0,
		}
	} else if len(p.tokens) == 1 {
		if p.tokens[0].tokType == tEOF {
			return ASTNode{nodeType: ASTIdentity}, nil
		}
	}
	parsed, err := p.parseExpression(0)
	if err != nil {
		return ASTNode{}, err
	}
	if p.current() != tEOF {
		return ASTNode{}, p.syntaxError(fmt.Sprintf(
			"Unexpected token at the end of the expression: %s", p.lookaheadToken(0).tokType))
	}
	return parsed, nil
}

func (p *parser) loadTokens(expression string) {
	lexer := NewLexer()
	tokens, err := lexer.tokenize(expression)
	if err != nil {
		p.tokens = []token{
			{tokType: tUnknown, value: ""},
		}
		p.index = 0
		return
	}
	p.tokens = tokens
	p.index = 0
}

var ASTEmpty = astNodeType(-1)

func (p *parser) parseExpression(bindingPower int) (ASTNode, error) {
	var err error
	leftToken := p.lookaheadToken(0)
	p.advance()
	leftNode, err := p.nud(leftToken)
	if err != nil {
		return ASTNode{}, err
	}
	currentToken := p.current()
	for bindingPower < bindingPowers[currentToken] {
		p.advance()
		leftNode, err = p.led(currentToken, leftNode)
		if err != nil {
			return ASTNode{}, err
		}
		currentToken = p.current()
	}
	return leftNode, nil
}

func (p *parser) parseSubexpression(bindingPower int) (ASTNode, error) {
	subexpr, err := p.parseExpression(bindingPower)
	if err != nil {
		return ASTNode{}, err
	}
	if p.current() != tEOF {
		return ASTNode{}, p.syntaxError("Unexpected token at the end of the expression")
	}
	return subexpr, nil
}

func (p *parser) nud(token token) (ASTNode, error) {
	switch token.tokType {
	case tUnquotedIdentifier:
		return ASTNode{nodeType: ASTField, value: token.value}, nil
	case tQuotedIdentifier:
		node := ASTNode{nodeType: ASTField, value: token.value}
		if p.current() == tLparen {
			return ASTNode{}, p.syntaxErrorToken("Can't use a quoted identifier as a function name.", token)
		}
		return node, nil
	case tStar:
		leftNode := ASTNode{nodeType: ASTIdentity}
		rightNode, err := p.parseProjectionRHS(bindingPowers[tStar])
		if err != nil {
			return ASTNode{}, err
		}
		return ASTNode{nodeType: ASTValueProjection, value: nil,
			children: []ASTNode{leftNode, rightNode}}, nil
	case tFilter:
		return p.parseFilter(ASTNode{nodeType: ASTIdentity})
	case tLbrace:
		return p.parseMultiSelectHash()
	case tFlatten:
		leftNode := ASTNode{nodeType: ASTFlatten, children: []ASTNode{{nodeType: ASTIdentity}}}
		rightNode, err := p.parseProjectionRHS(bindingPowers[tFlatten])
		if err != nil {
			return ASTNode{}, err
		}
		return ASTNode{nodeType: ASTProjection,
			children: []ASTNode{leftNode, rightNode}}, nil
	case tLbracket:
		tokenType := p.current()
		//  We could have [number], [*], [?], [], [a, b, c]
		if tokenType == tNumber || tokenType == tColon {
			right, err := p.parseIndexExpression()
			if err != nil {
				return ASTNode{}, nil
			}
			return p.projectIfSlice(ASTNode{nodeType: ASTIdentity}, right)
		} else if tokenType == tStar && p.lookahead(1) == tRbracket {
			p.advance()
			p.advance()
			rightNode, err := p.parseProjectionRHS(bindingPowers[tStar])
			if err != nil {
				return ASTNode{}, err
			}
			return ASTNode{nodeType: ASTProjection,
				children: []ASTNode{{nodeType: ASTIdentity}, rightNode}}, nil
		} else {
			return p.parseMultiSelectList()
		}
	case tCurrent:
		return ASTNode{nodeType: ASTCurrentNode}, nil
	case tExpref:
		expression, err := p.parseExpression(bindingPowers[tExpref])
		if err != nil {
			return ASTNode{}, err
		}
		return ASTNode{nodeType: ASTExpRef, children: []ASTNode{expression}}, nil
	case tLparen:
		expression, err := p.parseExpression(0)
		if err != nil {
			return ASTNode{}, err
		}
		if err := p.match(tRparen); err != nil {
			return ASTNode{}, err
		}
		return expression, nil
	case tJSONLiteral:
		return ASTNode{nodeType: ASTLiteral, value: token.value}, nil
	case tStringLiteral:
		return ASTNode{nodeType: ASTLiteral, value: token.value}, nil
	case tNot:
		expression, err := p.parseExpression(bindingPowers[tNot])
		if err != nil {
			return ASTNode{}, err
		}
		return ASTNode{nodeType: ASTNotExpression, children: []ASTNode{expression}}, nil
	case tEOF:
		return ASTNode{}, p.syntaxErrorToken("Incomplete expression", token)
	}
	return ASTNode{}, p.syntaxErrorToken(fmt.Sprintf("Unexpected token: %s", token.tokType), token)
}

func (p *parser) led(tokenType tokType, node ASTNode) (ASTNode, error) {
	switch tokenType {
	case tDot:
		if p.current() != tStar {
			subexpr, err := p.parseSubexpression(bindingPowers[tDot])
			if err != nil {
				return ASTNode{}, err
			}
			return ASTNode{nodeType: ASTSubexpression, children: []ASTNode{node, subexpr}}, nil
		}
		p.advance()
		right, err := p.parseProjectionRHS(bindingPowers[tDot])
		if err != nil {
			return ASTNode{}, err
		}
		return ASTNode{nodeType: ASTValueProjection, children: []ASTNode{node, right}}, nil
	case tPipe:
		right, err := p.parseExpression(bindingPowers[tPipe])
		if err != nil {
			return ASTNode{}, err
		}
		return ASTNode{nodeType: ASTPipe, children: []ASTNode{node, right}}, nil
	case tOr:
		right, err := p.parseExpression(bindingPowers[tOr])
		if err != nil {
			return ASTNode{}, err
		}
		return ASTNode{nodeType: ASTOrExpression, children: []ASTNode{node, right}}, nil
	case tAnd:
		right, err := p.parseExpression(bindingPowers[tAnd])
		if err != nil {
			return ASTNode{}, err
		}
		return ASTNode{nodeType: ASTAndExpression, children: []ASTNode{node, right}}, nil
	case tLbracket:
		tokenType := p.current()
		if tokenType == tNumber || tokenType == tColon {
			right, err := p.parseIndexExpression()
			if err != nil {
				return ASTNode{}, err
			}
			return p.projectIfSlice(node, right)
		}
		// Otherwise a multi-select-list [a, b, c] or [*]
		if tokenType == tStar && p.lookahead(1) == tRbracket {
			p.advance()
			p.advance()
			rightNode, err := p.parseProjectionRHS(bindingPowers[tStar])
			if err != nil {
				return ASTNode{}, err
			}
			return ASTNode{nodeType: ASTProjection, children: []ASTNode{node, rightNode}}, nil
		}
		return p.parseMultiSelectList()
	case tFlatten:
		leftNode := ASTNode{nodeType: ASTFlatten, children: []ASTNode{node}}
		rightNode, err := p.parseProjectionRHS(bindingPowers[tFlatten])
		if err != nil {
			return ASTNode{}, err
		}
		return ASTNode{nodeType: ASTProjection, children: []ASTNode{leftNode, rightNode}}, nil
	case tEQ, tNE, tGT, tGTE, tLT, tLTE:
		right, err := p.parseExpression(bindingPowers[tokenType])
		if err != nil {
			return ASTNode{}, err
		}
		return ASTNode{nodeType: ASTComparator, value: tokenType, children: []ASTNode{node, right}}, nil
	case tFilter:
		return p.parseFilter(node)
	case tStar:
		right, err := p.parseProjectionRHS(bindingPowers[tStar])
		if err != nil {
			return ASTNode{}, err
		}
		return ASTNode{nodeType: ASTValueProjection, children: []ASTNode{node, right}}, nil
	}
	return ASTNode{}, p.syntaxError(fmt.Sprintf("Unexpected token: %s", tokenType))
}

func (p *parser) match(tokenType tokType) error {
	if p.current() == tokenType {
		p.advance()
		return nil
	}
	return p.syntaxError(fmt.Sprintf("Expected %s, received: %s", tokenType, p.lookaheadToken(0).tokType))
}

func (p *parser) parseIndexExpression() (ASTNode, error) {
	if p.lookahead(0) == tColon || p.lookahead(1) == tColon {
		return p.parseSliceExpression()
	}
	index, err := strconv.Atoi(p.lookaheadToken(0).value.(string))
	if err != nil {
		return ASTNode{}, p.syntaxError("Bad token: " + p.lookaheadToken(0).value.(string))
	}
	indexNode := ASTNode{nodeType: ASTIndex, value: index}
	p.advance()
	if err := p.match(tRbracket); err != nil {
		return ASTNode{}, err
	}
	return indexNode, nil
}

func (p *parser) parseSliceExpression() (ASTNode, error) {
	// [start:stop:step]
	// Where any of start, stop, or step are optional
	// They default to start=0, stop=len(array)-1, step=1.
	parts := [3]*int{}
	index := 0
	for p.current() != tRbracket && index < 3 {
		if p.current() == tColon {
			index++
			p.advance()
		} else if p.current() == tNumber {
			intValue, err := strconv.Atoi(p.lookaheadToken(0).value.(string))
			if err != nil {
				return ASTNode{}, err
			}
			parts[index] = &intValue
			p.advance()
		} else {
			return ASTNode{}, p.syntaxError(
				fmt.Sprintf("Unexpected token in slice expression: %s", p.lookaheadToken(0).tokType))
		}
	}
	if err := p.match(tRbracket); err != nil {
		return ASTNode{}, err
	}
	return ASTNode{
		nodeType: ASTSlice,
		value:    parts[:],
	}, nil
}

func (p *parser) parseFilter(node ASTNode) (ASTNode, error) {
	condition, err := p.parseExpression(0)
	if err != nil {
		return ASTNode{}, err
	}
	if err := p.match(tRbracket); err != nil {
		return ASTNode{}, err
	}
	rightNode, err := p.parseProjectionRHS(bindingPowers[tFilter])
	if err != nil {
		return ASTNode{}, err
	}
	return ASTNode{nodeType: ASTFilterProjection,
		children: []ASTNode{node, condition, rightNode}}, nil
}

func (p *parser) parseFunctionExpression(token token) (ASTNode, error) {
	args := []ASTNode{}
	for p.current() != tRparen {
		expression, err := p.parseExpression(0)
		if err != nil {
			return ASTNode{}, err
		}
		if p.current() == tComma {
			if err := p.match(tComma); err != nil {
				return ASTNode{}, err
			}
		}
		args = append(args, expression)
	}
	if err := p.match(tRparen); err != nil {
		return ASTNode{}, err
	}
	return ASTNode{nodeType: ASTFunctionExpression,
		value:    token.value,
		children: args}, nil
}

func (p *parser) parseMultiSelectList() (ASTNode, error) {
	nodes := []ASTNode{}
	for {
		expression, err := p.parseExpression(0)
		if err != nil {
			return ASTNode{}, err
		}
		nodes = append(nodes, expression)
		if p.current() == tRbracket {
			break
		}
		err = p.match(tComma)
		if err != nil {
			return ASTNode{}, err
		}
	}
	p.advance()
	return ASTNode{
		nodeType: ASTMultiSelectList,
		children: nodes,
	}, nil
}

func (p *parser) parseMultiSelectHash() (ASTNode, error) {
	nodes := []ASTNode{}
	for {
		keyToken := p.lookaheadToken(0)
		if err := p.match(tUnquotedIdentifier); err != nil {
			if err := p.match(tQuotedIdentifier); err != nil {
				return ASTNode{}, p.syntaxError("Expected tQuotedIdentifier or tUnquotedIdentifier")
			}
		}
		keyName := keyToken.value.(string)
		if err := p.match(tColon); err != nil {
			return ASTNode{}, err
		}
		value, err := p.parseExpression(0)
		if err != nil {
			return ASTNode{}, err
		}
		nodes = append(nodes, ASTNode{
			nodeType: ASTKeyValPair,
			value:    keyName,
			children: []ASTNode{value},
		})
		if p.current() == tRbrace {
			break
		}
		if err := p.match(tComma); err != nil {
			return ASTNode{}, err
		}
	}
	p.advance()
	return ASTNode{
		nodeType: ASTMultiSelectHash,
		children: nodes,
	}, nil
}

func (p *parser) projectIfSlice(left ASTNode, right ASTNode) (ASTNode, error) {
	if right.nodeType == ASTSlice {
		rightNode, err := p.parseProjectionRHS(bindingPowers[tStar])
		if err != nil {
			return ASTNode{}, err
		}
		return ASTNode{
			nodeType: ASTProjection,
			children: []ASTNode{ASTNode{nodeType: ASTIndexExpression, children: []ASTNode{left, right}}, rightNode},
		}, nil
	}
	return ASTNode{
		nodeType: ASTIndexExpression,
		children: []ASTNode{left, right},
	}, nil
}

func (p *parser) parseProjectionRHS(bindingPower int) (ASTNode, error) {
	if bindingPowers[p.current()] < 10 {
		return ASTNode{nodeType: ASTIdentity}, nil
	}
	switch p.current() {
	case tDot:
		p.advance()
		return p.parseSubexpression(bindingPower)
	case tLbracket:
		return p.parseExpression(bindingPower)
	case tFilter:
		return p.parseExpression(bindingPower)
	default:
		return ASTNode{}, p.syntaxError(fmt.Sprintf("Error in parse projection RHS, unexpected token: %s", p.current()))
	}
}

func (p *parser) lookahead(offset int) tokType {
	return p.lookaheadToken(offset).tokType
}

func (p *parser) current() tokType {
	return p.lookahead(0)
}

func (p *parser) lookaheadToken(offset int) token {
	remainingTokens := len(p.tokens) - p.index
	if remainingTokens <= offset {
		return token{tokType: tEOF}
	}
	return p.tokens[p.index+offset]
}

func (p *parser) advance() {
	p.index++
}

func (p *parser) syntaxError(msg string) SyntaxError {
	return SyntaxError{
		msg:        msg,
		Expression: p.expression,
		Offset:     p.lookaheadToken(0).position,
	}
}
func (p *parser) syntaxErrorToken(msg string, t token) SyntaxError {
	return SyntaxError{
		msg:        msg,
		Expression: p.expression,
		Offset:     t.position,
	}
}

// Bind power based on precedence
var bindingPowers = map[tokType]int{
	tEOF:              0,
	tUnquotedIdentifier: 0,
	tQuotedIdentifier:   0,
	tRbracket:        0,
	tRparen:          0,
	tComma:           0,
	tRbrace:          0,
	tNumber:          0,
	tCurrent:         0,
	tExpref:          0,
	tColon:           0,
	tPipe:            1,
	tOr:              2,
	tAnd:             3,
	tEQ:              5,
	tGT:              5,
	tLT:              5,
	tGTE:             5,
	tLTE:             5,
	tNE:              5,
	tFlatten:         9,
	tStar:            20,
	tFilter:          21,
	tDot:             40,
	tNot:             45,
	tLbrace:          50,
	tLbracket:        55,
	tLparen:          60,
}
