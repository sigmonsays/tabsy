package tab

import (
	"fmt"
	"os"

	"github.com/chzyer/readline"
)

type Context struct {
	args []string
	dlog *os.File
	rl   *readline.Instance

	Prompt Prompt
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

	c.rl.SetPrompt(c.Prompt.String())
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
