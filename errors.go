package hashbrown

import (
	"fmt"
	"os"
	"strings"
)

type sgr int

const (
	bold      sgr = 1
	dim       sgr = 2
	underline sgr = 4
	red       sgr = 31
	green     sgr = 32
	yellow    sgr = 33
)

func ansi(message string, sgr sgr) string {
	return fmt.Sprintf("%s%dm%s", "\033[", sgr, message) + "\033[0m"
}

// Warning displays `message` along with the filepath, and the line and column number that the warning is being given for.
func Warning(message string) {
	fmt.Print(ansi(fmt.Sprintf("\nWarning: %s %s:%d:%d\n", message, filepath, LineIndex+1, LineCharIndex+1), yellow))
}

// Error is an alias for Exit, if `err` is not nil, Exit is called using `err` converted into a string.
func Error(err error) {
	if err != nil {
		Exit(fmt.Sprintf("%s", err))
	}
}

// Exit prints a parser error message and then exits.
// This will print out the current line and column number when Exit was called
// the lines applicable to the error message in the file contents given in `Init(path,contents)`
// and highlight the characters in current line at the current column number that has been determined as the error location.
func Exit(message string) {
	var lines = strings.Split(contents, "\n")
	if CurrentChar == '\n' || PrevChar(1) == '\n' {
		LineIndex--
	}
	if LineIndex != -1 {
		fmt.Print("\033[31m")
		fmt.Println("\n" + ansi(message, bold))
		fmt.Printf("\n\033[2m----- \033[0m%s:%d:%d\n", filepath, LineIndex+1, LineCharIndex+1)
		if len(lines) > (LineIndex-1) && LineIndex != 0 {
			fmt.Printf("\033[2m%d | %s\033[0m\n", LineIndex, lines[LineIndex-1])
		}
		if len(lines) > LineIndex {
			fmt.Printf("\033[31m\033[1m%d | ", LineIndex+1)
			for c, chr := range strings.Split(lines[LineIndex], "") {
				if c == Index {
					fmt.Print(ansi(chr, underline))
				} else {
					fmt.Print(chr)
				}
			}
			fmt.Print("\033[0m\n")
		}
		var spaces string
		for i := 0; i < (LineCharIndex + 4); i++ {
			spaces += " "
		}
		fmt.Println("\033[31m" + spaces + "^\033[0m")
		if len(lines) > (LineIndex + 1) {
			fmt.Printf("\033[2m%d | %s\n-----\033[0m\n\n", LineIndex+2, lines[LineIndex+1])
		}
	} else {
		fmt.Printf("Error: %s %s:%d:%d\n", message, filepath, LineIndex+1, LineCharIndex+1)
	}
	os.Exit(1)
}
