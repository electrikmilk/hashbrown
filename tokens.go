package hashbrown

type TokenType string

type Token struct {
	kind       TokenType
	identifier string
	valueType  TokenType
	value      any
}

var Tokens []Token

/* Types */

const (
	String      TokenType = "\""
	Integer     TokenType = "0123456789.-"
	Dictionary  TokenType = "{"
	Array       TokenType = "["
	Bool        TokenType = "boolean"
	Date        TokenType = "date"
	True        TokenType = "true"
	False       TokenType = "false"
	Comment     TokenType = "comment"
	Expression  TokenType = "expression"
	Variable    TokenType = "var"
	Constant    TokenType = "const"
	Conditional TokenType = "conditional"
	Iterator    TokenType = "iterator"
)

/* Operators */

const (
	At             TokenType = "@"
	DollarSign     TokenType = "$"
	Asterisk       TokenType = "*"
	Hash           TokenType = "#"
	Colon          TokenType = ":"
	Semicolon      TokenType = ";"
	Period         TokenType = "."
	AddTo          TokenType = "+="
	Is             TokenType = "=="
	Not            TokenType = "!="
	GreaterThan    TokenType = ">"
	GreaterOrEqual TokenType = ">="
	LessThan       TokenType = "<"
	LessOrEqual    TokenType = "<="
	Exclamation    TokenType = "!"
	Between        TokenType = "<>"
	ForwardSlash   TokenType = "/"
	BackSlash      TokenType = "\\"
	Plus           TokenType = "+"
	Modulus        TokenType = "%"
	LeftParen      TokenType = "("
	RightParen     TokenType = ")"
	LeftBrace      TokenType = "{"
	RightBrace     TokenType = "}"
	LeftBracket    TokenType = "["
	RightBracket   TokenType = "]"
	Dash           TokenType = "-"
	DoubleQuote    TokenType = "\""
	SingleQuote    TokenType = "'"
	Equality       TokenType = "="
	Space          TokenType = " "
	Tab            TokenType = "\t"
)
