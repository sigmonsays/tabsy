package tab

import (
	"fmt"
	"io"
	"os"
	"strings"
	"text/tabwriter"
)

func ShowHelp(root *RootCommand, cmd *Command, find_error error) {
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

func ExecuteLine(c *RootCommand, line []string) error {

	text := strings.Join(line, " ")

	return c.Dispatch(text)
}

// send on quit channel to terminate loop
func Loop(c *RootCommand, quit func()) error {

	c.dbg("loop start")

Dance:
	for {
		c.dbg("readline wait")

		text, err := c.Ctx.Readline()
		if err == io.EOF {
			// Quit without error on Ctrl^D
			c.dbg("error %s", err)
		}
		if err != nil {
			c.dbg("error %s", err)
			if quit != nil {
				quit()
			}
		}
		c.dbg("text %q", text)
		err = c.Dispatch(text)
		if err == io.EOF {
			break Dance
		}
		if err != nil {
			c.dbg("dispatch error %s\n", err)
		}
	}
	c.dbg("loop exit")

	return nil
}
