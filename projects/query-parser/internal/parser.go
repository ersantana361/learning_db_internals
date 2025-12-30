package internal

import (
	"fmt"
)

// ASTNodeType represents the type of an AST node
type ASTNodeType string

const (
	NodeSelect     ASTNodeType = "SELECT"
	NodeFrom       ASTNodeType = "FROM"
	NodeWhere      ASTNodeType = "WHERE"
	NodeColumns    ASTNodeType = "COLUMNS"
	NodeColumn     ASTNodeType = "COLUMN"
	NodeTable      ASTNodeType = "TABLE"
	NodeCondition  ASTNodeType = "CONDITION"
	NodeBinaryExpr ASTNodeType = "BINARY_EXPR"
	NodeLiteral    ASTNodeType = "LITERAL"
	NodeIdentifier ASTNodeType = "IDENTIFIER"
	NodeFunction   ASTNodeType = "FUNCTION"
	NodeOrderBy    ASTNodeType = "ORDER_BY"
	NodeLimit      ASTNodeType = "LIMIT"
	NodeJoin       ASTNodeType = "JOIN"
	NodeStatement  ASTNodeType = "STATEMENT"
)

// ASTNode represents a node in the abstract syntax tree
type ASTNode struct {
	ID       string                 `json:"id"`
	Type     ASTNodeType            `json:"type"`
	Value    string                 `json:"value,omitempty"`
	Children []string               `json:"children"`
	Parent   string                 `json:"parent,omitempty"`
	Meta     map[string]interface{} `json:"meta,omitempty"`
}

// Parser parses SQL tokens into an AST
type Parser struct {
	tokens   []Token
	pos      int
	nodes    map[string]*ASTNode
	nodeSeq  int
	rootID   string
	errors   []string
}

// NewParser creates a new parser for the given tokens
func NewParser(tokens []Token) *Parser {
	return &Parser{
		tokens:  tokens,
		pos:     0,
		nodes:   make(map[string]*ASTNode),
		nodeSeq: 0,
		errors:  []string{},
	}
}

// Parse parses all tokens into an AST
func (p *Parser) Parse() (*ASTNode, error) {
	if len(p.tokens) == 0 {
		return nil, fmt.Errorf("no tokens to parse")
	}

	root := p.parseStatement()
	if root == nil {
		if len(p.errors) > 0 {
			return nil, fmt.Errorf("parse errors: %v", p.errors)
		}
		return nil, fmt.Errorf("failed to parse statement")
	}

	p.rootID = root.ID
	return root, nil
}

// GetNodes returns all AST nodes
func (p *Parser) GetNodes() map[string]*ASTNode {
	return p.nodes
}

// GetRootID returns the root node ID
func (p *Parser) GetRootID() string {
	return p.rootID
}

// GetErrors returns parse errors
func (p *Parser) GetErrors() []string {
	return p.errors
}

func (p *Parser) createNode(nodeType ASTNodeType, value string) *ASTNode {
	p.nodeSeq++
	node := &ASTNode{
		ID:       fmt.Sprintf("ast-%d", p.nodeSeq),
		Type:     nodeType,
		Value:    value,
		Children: []string{},
		Meta:     make(map[string]interface{}),
	}
	p.nodes[node.ID] = node
	return node
}

func (p *Parser) addChild(parent, child *ASTNode) {
	parent.Children = append(parent.Children, child.ID)
	child.Parent = parent.ID
}

func (p *Parser) current() Token {
	if p.pos >= len(p.tokens) {
		return Token{Type: TokenEOF}
	}
	return p.tokens[p.pos]
}

func (p *Parser) peek() Token {
	if p.pos+1 >= len(p.tokens) {
		return Token{Type: TokenEOF}
	}
	return p.tokens[p.pos+1]
}

func (p *Parser) advance() Token {
	token := p.current()
	p.pos++
	return token
}

func (p *Parser) expect(tokenType TokenType) (Token, bool) {
	token := p.current()
	if token.Type != tokenType {
		p.errors = append(p.errors, fmt.Sprintf("expected %s but got %s", tokenType, token.Type))
		return token, false
	}
	return p.advance(), true
}

func (p *Parser) parseStatement() *ASTNode {
	token := p.current()

	switch token.Type {
	case TokenSelect:
		return p.parseSelect()
	default:
		p.errors = append(p.errors, fmt.Sprintf("unexpected token: %s", token.Type))
		return nil
	}
}

func (p *Parser) parseSelect() *ASTNode {
	selectNode := p.createNode(NodeStatement, "SELECT")
	selectNode.Meta["type"] = "SELECT"

	p.advance() // consume SELECT

	// Parse columns
	columns := p.parseColumns()
	if columns != nil {
		p.addChild(selectNode, columns)
	}

	// Parse FROM clause
	if p.current().Type == TokenFrom {
		fromNode := p.parseFrom()
		if fromNode != nil {
			p.addChild(selectNode, fromNode)
		}
	}

	// Parse WHERE clause
	if p.current().Type == TokenWhere {
		whereNode := p.parseWhere()
		if whereNode != nil {
			p.addChild(selectNode, whereNode)
		}
	}

	// Parse ORDER BY clause
	if p.current().Type == TokenOrderBy {
		orderNode := p.parseOrderBy()
		if orderNode != nil {
			p.addChild(selectNode, orderNode)
		}
	}

	// Parse LIMIT clause
	if p.current().Type == TokenLimit {
		limitNode := p.parseLimit()
		if limitNode != nil {
			p.addChild(selectNode, limitNode)
		}
	}

	return selectNode
}

func (p *Parser) parseColumns() *ASTNode {
	columnsNode := p.createNode(NodeColumns, "")

	for {
		if p.current().Type == TokenStar {
			col := p.createNode(NodeColumn, "*")
			p.addChild(columnsNode, col)
			p.advance()
		} else if p.current().Type == TokenIdentifier {
			col := p.parseColumnExpression()
			if col != nil {
				p.addChild(columnsNode, col)
			}
		} else {
			break
		}

		if p.current().Type != TokenComma {
			break
		}
		p.advance() // consume comma
	}

	return columnsNode
}

func (p *Parser) parseColumnExpression() *ASTNode {
	token := p.advance()
	col := p.createNode(NodeColumn, token.Value)

	// Check for alias
	if p.current().Type == TokenAs {
		p.advance()
		if alias, ok := p.expect(TokenIdentifier); ok {
			col.Meta["alias"] = alias.Value
		}
	}

	return col
}

func (p *Parser) parseFrom() *ASTNode {
	fromNode := p.createNode(NodeFrom, "")
	p.advance() // consume FROM

	// Parse table name
	if token, ok := p.expect(TokenIdentifier); ok {
		table := p.createNode(NodeTable, token.Value)
		p.addChild(fromNode, table)

		// Check for alias
		if p.current().Type == TokenAs {
			p.advance()
			if alias, ok := p.expect(TokenIdentifier); ok {
				table.Meta["alias"] = alias.Value
			}
		} else if p.current().Type == TokenIdentifier {
			// Alias without AS keyword
			alias := p.advance()
			table.Meta["alias"] = alias.Value
		}
	}

	// Parse JOINs
	for p.current().Type == TokenJoin {
		joinNode := p.parseJoin()
		if joinNode != nil {
			p.addChild(fromNode, joinNode)
		}
	}

	return fromNode
}

func (p *Parser) parseJoin() *ASTNode {
	joinNode := p.createNode(NodeJoin, "")
	p.advance() // consume JOIN

	// Parse table
	if token, ok := p.expect(TokenIdentifier); ok {
		table := p.createNode(NodeTable, token.Value)
		p.addChild(joinNode, table)
	}

	// Parse ON condition
	if p.current().Type == TokenOn {
		p.advance()
		condition := p.parseExpression()
		if condition != nil {
			p.addChild(joinNode, condition)
		}
	}

	return joinNode
}

func (p *Parser) parseWhere() *ASTNode {
	whereNode := p.createNode(NodeWhere, "")
	p.advance() // consume WHERE

	condition := p.parseExpression()
	if condition != nil {
		p.addChild(whereNode, condition)
	}

	return whereNode
}

func (p *Parser) parseExpression() *ASTNode {
	return p.parseOrExpression()
}

func (p *Parser) parseOrExpression() *ASTNode {
	left := p.parseAndExpression()

	for p.current().Type == TokenOr {
		p.advance()
		right := p.parseAndExpression()

		orNode := p.createNode(NodeBinaryExpr, "OR")
		p.addChild(orNode, left)
		if right != nil {
			p.addChild(orNode, right)
		}
		left = orNode
	}

	return left
}

func (p *Parser) parseAndExpression() *ASTNode {
	left := p.parseComparisonExpression()

	for p.current().Type == TokenAnd {
		p.advance()
		right := p.parseComparisonExpression()

		andNode := p.createNode(NodeBinaryExpr, "AND")
		p.addChild(andNode, left)
		if right != nil {
			p.addChild(andNode, right)
		}
		left = andNode
	}

	return left
}

func (p *Parser) parseComparisonExpression() *ASTNode {
	left := p.parsePrimaryExpression()

	if p.current().Type == TokenOperator {
		op := p.advance()
		right := p.parsePrimaryExpression()

		compNode := p.createNode(NodeBinaryExpr, op.Value)
		p.addChild(compNode, left)
		if right != nil {
			p.addChild(compNode, right)
		}
		return compNode
	}

	return left
}

func (p *Parser) parsePrimaryExpression() *ASTNode {
	token := p.current()

	switch token.Type {
	case TokenIdentifier:
		p.advance()
		return p.createNode(NodeIdentifier, token.Value)
	case TokenNumber:
		p.advance()
		node := p.createNode(NodeLiteral, token.Value)
		node.Meta["literalType"] = "number"
		return node
	case TokenString:
		p.advance()
		node := p.createNode(NodeLiteral, token.Value)
		node.Meta["literalType"] = "string"
		return node
	case TokenLParen:
		p.advance() // consume (
		expr := p.parseExpression()
		p.expect(TokenRParen)
		return expr
	default:
		p.errors = append(p.errors, fmt.Sprintf("unexpected token in expression: %s", token.Type))
		return nil
	}
}

func (p *Parser) parseOrderBy() *ASTNode {
	orderNode := p.createNode(NodeOrderBy, "")
	p.advance() // consume ORDER

	// Consume BY if present
	if p.current().Type == TokenIdentifier && p.current().Value == "BY" {
		p.advance()
	}

	// Parse column
	if p.current().Type == TokenIdentifier {
		col := p.createNode(NodeColumn, p.advance().Value)
		p.addChild(orderNode, col)

		// Check for ASC/DESC
		if p.current().Type == TokenIdentifier {
			val := p.current().Value
			if val == "ASC" || val == "DESC" {
				col.Meta["direction"] = val
				p.advance()
			}
		}
	}

	return orderNode
}

func (p *Parser) parseLimit() *ASTNode {
	limitNode := p.createNode(NodeLimit, "")
	p.advance() // consume LIMIT

	if p.current().Type == TokenNumber {
		limitNode.Value = p.advance().Value
	}

	// Parse OFFSET
	if p.current().Type == TokenOffset {
		p.advance()
		if p.current().Type == TokenNumber {
			limitNode.Meta["offset"] = p.advance().Value
		}
	}

	return limitNode
}

// CurrentPosition returns the current token position
func (p *Parser) CurrentPosition() int {
	return p.pos
}

// CurrentToken returns the current token
func (p *Parser) CurrentToken() *Token {
	if p.pos >= len(p.tokens) {
		return nil
	}
	return &p.tokens[p.pos]
}
