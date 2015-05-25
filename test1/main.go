package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"golang.org/x/crypto/ssh/terminal"
)

var regularState *terminal.State

func main() {
	oldState, err := terminal.MakeRaw(0)
	if err != nil {
		panic(err)
	}
	regularState = oldState

	prompt := "> "
	fmt.Printf("starting..\n")
	t := terminal.NewTerminal(os.Stdin, prompt)
	t.AutoCompleteCallback = handleKey

	term(t, oldState)

	fmt.Printf("exiting..\n")
}

func term(t *terminal.Terminal, oldState *terminal.State) {
	defer terminal.Restore(0, oldState)

	for {
		text, err := t.ReadLine()
		if err != nil {
			if err == io.EOF {
				// Quit without error on Ctrl^D
				fmt.Println()
				break
			}
			panic(err)
		}

		text = strings.Replace(text, " ", "", -1)
		if text == "exit" || text == "quit" {
			break
		}

		fmt.Printf("input: %s\n", text)

	}
}
func handleKey(line string, pos int, key rune) (newLine string, newPos int, ok bool) {
	if key == 0x03 {
		fmt.Println()
		terminal.Restore(0, regularState)
		os.Exit(0)
	}
	return "", 0, false
}
