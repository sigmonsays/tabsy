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

func (c *Context) Arg(n int) string {
	if n > len(c.args) {
		return ""
	}
	return c.args[n]
}

func (c *Context) SetPrompt(prompt Prompt) {
	c.Prompt = prompt
	c.Term.SetPrompt(prompt.String())
}
