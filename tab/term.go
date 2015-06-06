package tab

import (
	"fmt"
	"io"
	"os"
	"strings"
	"text/tabwriter"

	"golang.org/x/crypto/ssh/terminal"
)

func (c *RootCommand) InitTerm() error {
	oldState, err := terminal.MakeRaw(0)
	if err != nil {
		panic(err)
	}
	c.Ctx.RegularState = oldState

	prompt := "> "
	fmt.Printf("starting..\n")
	t := terminal.NewTerminal(os.Stdin, prompt)
	c.Ctx.Term = t
	return nil
}

func (c *RootCommand) ReleaseTerm() error {
	terminal.Restore(0, c.Ctx.RegularState)
	return nil
}
func showHelp(cmd *Command, find_error error) {
	w := new(tabwriter.Writer)

	fmt.Printf("\nhelp %q - %s\n\n", cmd.Name, cmd.Description)
	if cmd.SubCmd == nil {
		return
	}
	w.Init(os.Stdout, 0, 9, 0, '\t', 0)
	for _, scmd := range cmd.SubCmd {
		fmt.Fprintln(w, fmt.Sprintf("%s\t%s\t", scmd.Name, scmd.Description))
	}
	fmt.Fprintln(w)
	w.Flush()
}

// send on quit channel to terminate loop
func Loop(c *RootCommand, quit chan bool) error {
	cli_ready := make(chan bool, 1)
	cli_text := make(chan string, 0)
	cli_keypress := make(chan *KeyPress, 0)
	readline_err := make(chan error, 0)

	ctx := c.Ctx
	_ = ctx

	c.Ctx.Term.AutoCompleteCallback = func(line string, pos int, key rune) (newLine string, newPos int, ok bool) {
		response := make(chan *KeyResponse, 1)
		cli_keypress <- &KeyPress{
			response: response,
			line:     line,
			pos:      pos,
			key:      key,
		}

		// wait for reply
		res := <-response

		return res.newline, res.newpos, res.ok
	}
	cli_ready <- true

	go func() {
		for {
			<-cli_ready
			text, err := c.Ctx.Term.ReadLine()
			if err == io.EOF {
				// Quit without error on Ctrl^D
				fmt.Println()
				fmt.Printf("EOF\n")
				quit <- true
			}
			if err != nil {
				readline_err <- err
			}
			cli_text <- text
		}
	}()

	var err error

Dance:
	for {
		select {

		case k := <-cli_keypress:

			if k.key == 0x03 { // ^C
				k.response <- &KeyResponse{"", 0, true}

			} else if k.key == 0x9 { // TAB

				res := &KeyResponse{"", 0, false}

				prefix := ""
				fields := strings.Fields(k.line)
				if len(fields) > 0 {
					if strings.HasSuffix(k.line, " ") == false {
						prefix = fields[len(fields)-1]
					}
				}

				cmd, _ := c.FindCommand(k.line)

				// fmt.Printf("cmd=%s prefix=%s\n", cmd.Name, prefix)

				if len(prefix) > 0 && cmd.IsRoot == false && strings.HasPrefix(cmd.Name, prefix) && cmd.Name != prefix {
					res.ok = true
					p := cmd.Name[len(prefix):] + " "
					res.newline += k.line + p
					res.newpos += k.pos + len(p)
				}

				// fmt.Printf("cmd=%s fields=%s prefix=%s\n", cmd.Name, fields, prefix)

				// list all sub commands
				ls := []*Command{}
				for _, cmd := range cmd.SubCmd {
					if prefix == "" || strings.HasPrefix(cmd.Name, prefix) {
						ls = append(ls, cmd)
					}

				}

				if len(ls) > 0 {
					fmt.Println()
					for _, c := range ls {
						fmt.Printf("%-20s\n", c.Name)
					}
					fmt.Printf("\n%s%s", ctx.Prompt, k.line)
				}

				k.response <- res

			} else if k.key == 0x3f { // ?
				cmd, err := c.FindCommand(k.line)
				if cmd == nil {
					cmd = c.Command
				}

				showHelp(cmd, err)

				newline := k.line
				k.response <- &KeyResponse{newline, len(newline), true}
				fmt.Printf("\n%s%s", ctx.Prompt, k.line)
			} else {
				k.response <- &KeyResponse{"", 0, false}
			}

		case err := <-readline_err:
			return err
		case <-quit:
			break Dance
		case text := <-cli_text:

			err = c.Dispatch(text)
			if err == io.EOF {
				break Dance
			}
			if err != nil {
				fmt.Printf("ERROR: %s\n", err)
			}
			cli_ready <- true
		}
	}

	fmt.Printf("Quit..\n")
	terminal.Restore(0, c.Ctx.RegularState)
	return nil
}

type KeyPress struct {
	response chan *KeyResponse
	line     string
	pos      int
	key      rune
}
type KeyResponse struct {
	newline string
	newpos  int
	ok      bool
}
