package tab

import (
	"fmt"
	"io"
	"os"
	"text/tabwriter"

	"golang.org/x/crypto/ssh/terminal"
)

const (
	KeyTab          = 0x9
	KeyControlC     = 0x03
	KeyQuestionMark = 0x1f
)

func (c *RootCommand) InitTerm() error {
	oldState, err := terminal.MakeRaw(0)
	if err != nil {
		panic(err)
	}
	c.Ctx.RegularState = oldState

	prompt := "> "
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

			if k.key == KeyControlC { // ^C
				k.response <- &KeyResponse{"", 0, true}

			} else if k.key == KeyTab { // TAB
				res, _ := tabComplete(c, k)
				if res == nil {
					k.response <- &KeyResponse{"", 0, false}
				} else {
					c.dbg("complete line=%q ok=%v newpos=%d newline=%s", k.line, res.ok, res.newpos, res.newline)
					k.response <- res
				}

			} else if k.key == KeyQuestionMark { // ^?
				cmd, err := c.FindCommand(k.line)
				if cmd == nil {
					cmd = c.Command
				}

				showHelp(cmd, err)

				newline := k.line
				k.response <- &KeyResponse{newline, len(newline), true}
				fmt.Printf("\n%s%s", c.Ctx.Prompt, k.line)

			} else {
				switch true {
				case k.key >= 'a' && k.key <= 'z':
				case k.key >= 'A' && k.key <= 'Z':
				case k.key >= '0' && k.key <= '9':
				default:
					c.dbg("key %x", k.key)
				}
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
				fmt.Printf("%s\n", err)
			}
			cli_ready <- true
		}
	}

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
