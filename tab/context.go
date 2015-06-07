package tab

import (
	"golang.org/x/crypto/ssh/terminal"
)

type Context struct {
	args []string

	Prompt       Prompt
	Term         *terminal.Terminal
	RegularState *terminal.State
}

func (c *Context) Args() []string {
	return c.args
}

func (c *Context) HasArg(n int) bool {
	return n < len(c.args)-1
}
func (c *Context) Arg(n int) (arg string) {
	if c.HasArg(n) {
		arg = c.args[n]
	}
	return
}

func (c *Context) SetPrompt(prompt Prompt) {
	c.Prompt = prompt
	c.Term.SetPrompt(prompt.String())
}
