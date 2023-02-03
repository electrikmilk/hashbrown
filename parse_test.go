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
					kind:       Build,
					identifier: "build",
					valueType:  String,
					value:      platform,
				})
			} else {
				Exit("Expected build token")
			}
		case TokenAhead(Package):
			AdvanceChar()
			var identifier = Collect('\n')
			Tokens = append(Tokens, Token{
				kind:       Constant,
				identifier: "package",
				valueType:  String,
				value:      identifier,
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
					kind:       Constant,
					identifier: identifier,
					valueType:  String,
					value:      value,
				})
			} else {
				Exit("Expected equality operator")
			}
		}
	})
	if fileError != nil {
		log.Fatalf("File error: %s", fileError)
	}
	fmt.Println("tokens", Tokens)
}
