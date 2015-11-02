package tab

import (
	"fmt"
	"os"
)

type CommandContext interface {
	Dbg(s string, args ...interface{})

	SetPrompt(prompt Prompt)

	Args() []string
	HasArg(n int) bool
	SetArgs(args []string)

	Readline() (string, error)
}

func NewContext(rl Readliner) *Context {
	ctx := &Context{
		rl: rl,
	}
	return ctx
}

type Context struct {
	args []string
	dlog *os.File
	rl   Readliner

	Prompt Prompt
}

func (c *Context) Args() []string {
	return c.args
}

func (c *Context) HasArg(n int) bool {
	return n <= len(c.args)-1
}
func (c *Context) Arg(n int) (arg string) {
	if c.HasArg(n) {
		arg = c.args[n]
	}
	return
}
func (c *Context) SetArgs(args []string) {
	c.args = args
}

func (c *Context) Readline() (string, error) {
	if c.rl == nil {
		return "", fmt.Errorf("init error: readline is nul")
	}
	return c.rl.Readline()
}

func (c *Context) SetPrompt(prompt Prompt) {
	c.Prompt = prompt
	if c.rl != nil {
		c.rl.SetPrompt(c.Prompt)
	}
}

func (c *Context) CloseDebugLog() {
	if c.dlog == nil {
		return
	}
	c.dlog.Close()
	c.dlog = nil
}

func (c *Context) OpenDebugLog(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	c.dlog = f
	return nil
}

func (c *Context) Dbg(s string, args ...interface{}) {
	if c.dlog != nil {
		line := fmt.Sprintf(s, args...)
		fmt.Fprintf(c.dlog, line+"\n")
	}
}

func (c *Context) Printf(s string, args ...interface{}) (int, error) {
	c.Dbg(s, args...)
	return fmt.Printf(s, args...)
}

func (c *Context) Errorf(s string, args ...interface{}) (int, error) {
	c.Dbg("ERROR "+s, args...)
	return fmt.Printf(s, args...)
}
