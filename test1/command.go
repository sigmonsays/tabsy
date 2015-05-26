package main

import (
	"fmt"
	"strings"

	"golang.org/x/crypto/ssh/terminal"
)

/*
 simple command line layer for making tab complete interactive CLIs

 cd [path]
 pwd

 show ?
 show object [path]

 delete [path]
 makedir [path]
 upload [localpath] [path]

 object set mtime [mtime]



*/

type Prompt interface {
	String() string
}

type ExecuteCommand func(*Context) error

func NewCommandSet(name string) *RootCommand {
	root := &RootCommand{
		Ctx: &Context{},
		Command: &Command{
			Name: name,
		},
	}
	return root
}

type RootCommand struct {
	*Command
	Ctx *Context
}

// search one level deep for a command
func (c *Command) Find(name string) (*Command, error) {
	matches := make([]*Command, 0)
	for _, scmd := range c.SubCmd {
		if scmd.Name == name {
			return scmd, nil
		}
		if strings.HasPrefix(scmd.Name, name) {
			matches = append(matches, scmd)
		}
	}

	if len(matches) == 1 {
		return matches[0], nil
	}

	return nil, fmt.Errorf("cmd=%s: no such command %q", c.Name, name)
}

// find which command will be called to execute
// returns root command if no match
func (c *RootCommand) FindCommand(line string) (*Command, error) {
	if line == "" {
		return nil, fmt.Errorf("empty command")
	}
	fields := strings.Fields(line)

	var ret *Command
	var scmd *Command
	scmd = c.Command
	ret = scmd
	for i := 0; i < len(fields); i++ {
		tcmd, err := scmd.Find(fields[i])
		if err != nil {
			break
		}
		if tcmd != nil {
			scmd = tcmd
		}
		ret = scmd
	}

	return ret, nil
}

func (c *RootCommand) Dispatch(line string) error {
	if line == "" {
		return nil
	}

	cmd, err := c.FindCommand(line)
	if err != nil {
		return err
	}

	if cmd == nil {
		return err
	}

	if cmd.Exec == nil {
		// fmt.Printf("%s", c.Ctx.Prompt)
		return err
	}

	fields := strings.Fields(line)
	c.Ctx.args = fields[1:]
	err = cmd.Exec(c.Ctx)
	if err != nil {
		fmt.Printf("%s: %s", cmd.Name, err)
	}

	return err
}

type Command struct {
	Name        string
	Description string // short one line description
	Exec        ExecuteCommand

	SubCmd []*Command
}

func (c *Command) Add(cmd *Command) error {
	if c.SubCmd == nil {
		c.SubCmd = make([]*Command, 0)
	}
	c.SubCmd = append(c.SubCmd, cmd)
	return nil
}

type Context struct {
	args []string

	Prompt       Prompt
	Term         *terminal.Terminal
	RegularState *terminal.State
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
