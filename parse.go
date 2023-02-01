package hashbrown

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// Index is the current position in the character map `Chars`.
var Index int

// LineIndex keeps track of the current line position in the file.
var LineIndex int

// LineCharIndex keeps track of the current column position in the current line we are parsing in the file.
var LineCharIndex int

// CurrentChar is the current character we are looking at as an int32.
// This is set based on `Chars[Index]` when the parser is initialized using `Init()` and when `AdvanceChar()` is called.
var CurrentChar rune

// Chars is a string slice of all the characters in the contents given in `Init()`.
var Chars []string

var contents string
var filepath string

// Init initializes the parser and starts the main parsing loop.
// Provide a `path` to the file that will be parsed and a callback to pass the current character for each loop.
// Returns an error if the file at `path` does not exist or cannot be read.
func Init(path string, loop func(char *rune)) (fileError error) {
	if _, statErr := os.Stat(path); os.IsNotExist(statErr) {
		return statErr
	}
	var bytes, readErr = os.ReadFile(path)
	if readErr != nil {
		return readErr
	}
	filepath = path
	contents = string(bytes)
	Chars = strings.Split(contents, "")
	Index = -1
	AdvanceChar()
	for CurrentChar != -1 {
		loop(&CurrentChar)
	}
	return nil
}

// Collect advances and collects `CurrentChar` until we reach `until`, then returns `collected`.
func Collect(until rune) (collected string) {
	for CurrentChar != -1 && CurrentChar != until {
		collected += string(CurrentChar)
		AdvanceChar()
	}
	return
}

// CollectInteger advances and collects `CurrentChar` until it stops matching the `Integer` type, then returns `collected`.
func CollectInteger() (collected string) {
	for strings.Contains(
		string(Integer),
		string(CurrentChar),
	) {
		collected += string(CurrentChar)
		AdvanceChar()
	}
	return
}

// CollectString advances and collects `CurrentChar` until we reach an unescaped double quote, then returns `collected`.
func CollectString() (collected string) {
	for CurrentChar != -1 {
		if IsToken(DoubleQuote) && PrevChar(1) != '\\' {
			break
		}
		if IsToken(BackSlash) && NextChar(1) == '"' {
			AdvanceChar()
			continue
		}
		collected += string(CurrentChar)
		AdvanceChar()
	}
	AdvanceChar()
	collected = strings.Trim(collected, " ")
	return
}

// CollectArray collects an array assuming the standard JSON array syntax.
// It then inserts the collected array syntax into a JSON object string and unmarshal it.
// Finally, it returns the resulting interface.
func CollectArray() (array interface{}) {
	var rawJson = "{\"array\":["
	for !IsToken(RightBracket) && CurrentChar != -1 {
		rawJson += string(CurrentChar)
		AdvanceChar()
	}
	rawJson += "]}"
	if err := json.Unmarshal([]byte(rawJson), &array); err != nil {
		LineIndex -= 2
		Error(err)
	}
	array = array.(map[string]interface{})["array"]
	AdvanceChar()
	return
}

// CollectDictionary collects a dictionary assuming the standard JSON dictionary syntax.
// It will then unmarshal the collected dictionary syntax and returns the resulting interface.
func CollectDictionary() (dictionary interface{}) {
	var rawJson = "{"
	var insideInnerObject = false
	for {
		rawJson += string(CurrentChar)
		if CurrentChar == '{' {
			insideInnerObject = true
		} else if CurrentChar == '}' {
			if !insideInnerObject {
				break
			}
			insideInnerObject = false
		}
		AdvanceChar()
	}
	if err := json.Unmarshal([]byte(rawJson), &dictionary); err != nil {
		Error(err)
	}
	AdvanceChar()
	return
}

// LookAhead will pseudo advance until `until` is reached and returns the collected characters.
// This is a way to see what characters lie ahead without advancing.
func LookAhead(until rune) (ahead string) {
	var nextIdx = Index
	var nextChar rune
	for nextChar != until {
		if len(Chars) > nextIdx {
			nextChar = []rune(Chars[nextIdx])[0]
			ahead += Chars[nextIdx]
			nextIdx++
		} else {
			break
		}
	}
	ahead = strings.Trim(strings.ToLower(ahead), " \t\n")
	return
}

// AdvanceChar advances `CurrentChar` to the next character that lies at the next index in `Chars` determined by `Index`.
// If the value of `Index` exceeds the length of `Chars`, it is given a negative value to indicate that the end of the file has been reached.
func AdvanceChar() {
	Index++
	if len(Chars) > Index {
		CurrentChar = []rune(Chars[Index])[0]
		if CurrentChar == '\n' {
			LineCharIndex = 0
			LineIndex++
		} else {
			LineCharIndex++
		}
	} else {
		CurrentChar = -1
	}
}

// AdvanceChars is an alias for AdvanceChar to advance multiple times in one function call.
func AdvanceChars(times int) {
	for i := 0; i < times; i++ {
		AdvanceChar()
	}
}

// IsToken compares `CurrentChar` to `t`.
// If `CurrentChar` is equal to `t`, we advance the length of `t`.
func IsToken(t TokenType) bool {
	if strings.ToLower(string(CurrentChar)) != string(t) {
		return false
	}
	var tokenLength = len(string(t))
	AdvanceChars(tokenLength)
	return true
}

// TokenAhead iterates through the characters of `t` and pseudo advances to check if the upcoming characters match the value of `t`.
// If the value of `t` matches the upcoming characters, we advance the length of `t`.
func TokenAhead(t TokenType) (isAhead bool) {
	var tokenChars = strings.Split(string(t), "")
	isAhead = true
	for i, tchar := range tokenChars {
		if tchar == " " || tchar == "\t" || tchar == "\n" {
			continue
		}
		if i == 0 {
			if strings.ToLower(string(CurrentChar)) != tchar {
				isAhead = false
				break
			}
		} else if NextChar(i) != []rune(tchar)[0] {
			isAhead = false
			break
		}
	}
	if isAhead {
		AdvanceChars(len(tokenChars))
	}
	return
}

// TokensAhead is an alias of TokensAhead that checks for multiple tokens ahead in one function call.
func TokensAhead(t ...TokenType) bool {
	for _, aheadToken := range t {
		if TokenAhead(aheadToken) {
			return true
		}
	}
	return false
}

// NextChar returns the character ahead of the current character.
func NextChar(mov int) rune {
	return seek(&mov, false)
}

// PrevChar returns the character previous to the current character.
func PrevChar(mov int) rune {
	return seek(&mov, true)
}

func seek(mov *int, reverse bool) (requestedChar rune) {
	var nextChar = Index
	if reverse {
		nextChar -= *mov
	} else {
		nextChar += *mov
	}
	requestedChar = GetChar(nextChar)
	for requestedChar == ' ' || requestedChar == '\t' || requestedChar == '\n' {
		if reverse {
			nextChar -= 1
		} else {
			nextChar += 1
		}
		requestedChar = GetChar(nextChar)
	}
	return
}

// GetChar is a safe way to retrieve the character at the index `i` in `Chars`.
// If the value of `i` is larger than the length of `Chars`, -1 will be returned.
// If `i` is -1, the first character will be returned.
func GetChar(i int) rune {
	if i == -1 {
		return []rune(Chars[0])[0]
	}
	if len(Chars) > i {
		return []rune(Chars[i])[0]
	}
	return -1
}

// PrintCurrentChar is a debug function prints out the character at the parsers current index, the unicode number of that character and the line and column number of where it is located.
// If the current character is a tab, new line, carriage return, or space, it will be substituted with a descriptor.
// Output Format:
// [line:column] (character) (code)
func PrintCurrentChar() {
	var char string
	switch CurrentChar {
	case '\t':
		char = "TAB"
	case '\n':
		char = "NEW LINE (\\n)"
	case '\r':
		char = "CARRIAGE RETURN (\\r)"
	case ' ':
		char = "SPACE"
	default:
		char = string(CurrentChar)
	}
	fmt.Printf("%s %s %s\n", ansi(fmt.Sprintf("[%d:%d]", LineIndex+1, LineCharIndex+1), bold), char, ansi(fmt.Sprintf("%d", CurrentChar), dim))
}
