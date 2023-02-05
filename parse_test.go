package hashbrown

import (
	"fmt"
	"log"
	"testing"
)

/* Really simple parser test */

const Package TokenType = "package"
const Build TokenType = "go:build"

func TestParser(t *testing.T) {
	var fileError = Init("unix.go", func(_ *rune) {
		PrintCurrentChar()
		switch {
		case IsToken(EOL) || IsToken(Tab) || IsToken(Space):
			AdvanceChar()
		case IsToken(ForwardSlash):
			AdvanceChar()
			if TokenAhead(Build) {
				AdvanceChar()
				var platform = Collect('\n')
				Tokens = append(Tokens, Token{
					Kind:       Build,
					Identifier: "build",
					ValueType:  String,
					Value:      platform,
				})
			} else {
				Exit("Expected build token")
			}
		case TokenAhead(Package):
			AdvanceChar()
			var identifier = Collect('\n')
			Tokens = append(Tokens, Token{
				Kind:       Constant,
				Identifier: "package",
				ValueType:  String,
				Value:      identifier,
			})
		case TokenAhead(Constant):
			AdvanceChar()
			var identifier = Collect(' ')
			AdvanceChar()
			var value string
			if IsToken(Equality) {
				AdvanceChars(2)
				value = CollectString()
				Tokens = append(Tokens, Token{
					Kind:       Constant,
					Identifier: identifier,
					ValueType:  String,
					Value:      value,
				})
			} else {
				Exit("Expected equality operator")
			}
		default:
			Exit(fmt.Sprintf("Illegal character '%s'", string(CurrentChar)))
		}
	})
	if fileError != nil {
		log.Fatalf("File error: %s", fileError)
	}
	fmt.Println("tokens", Tokens)
}
