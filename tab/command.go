package tab

import (
	"fmt"
	"strings"
)

var ErrNoSuchCommand = fmt.Errorf("no such command")
var ErrAmbiguousCommand = fmt.Errorf("ambiguous")
var ErrCommandNotExec = fmt.Errorf("command not executable")

type Prompt interface {
	String() string
}

type ExecuteCommand func(*Context) error

func NewCommandSet(name string) *RootCommand {
	root := &RootCommand{
		Ctx: &Context{},
		Command: &Command{
			Name:   name,
			IsRoot: true,
		},
	}
	return root
}

type RootCommand struct {
	*Command
	Ctx *Context
}

// find which command will be called to execute
// returns root command if no match
func (c *RootCommand) FindCommand(line string) (*Command, error) {
	if line == "" {
		return c.Command, fmt.Errorf("empty command")
	}
	fields := strings.Fields(line)

	var ret *Command
	var scmd *Command
	scmd = c.Command
	ret = scmd
	for i := 0; i < len(fields); i++ {
		tcmd, err := scmd.FindOne(fields[i])
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

	var cmd *Command
	cmds, err := c.Find(line)
	if err != nil {
		return err
	}
	if len(cmds) == 0 {
		fmt.Printf("no such command\n")
		return ErrNoSuchCommand
	} else if len(cmds) == 1 {
		cmd = cmds[0]
	} else {
		return ErrAmbiguousCommand
	}

	if cmd.Exec == nil {
		fmt.Printf("command not executable\n")
		return ErrCommandNotExec
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
	IsRoot      bool

	Parent *Command
	SubCmd []*Command
}

func (c *Command) Add(cmd *Command) error {
	if c.SubCmd == nil {
		c.SubCmd = make([]*Command, 0)
	}
	cmd.Parent = c
	c.SubCmd = append(c.SubCmd, cmd)
	return nil
}

func (c *Command) Alias(name string) *Command {
	c2 := *c
	c2.Name = name
	return &c2
}

// return single matching command if possible
// search one level deep for a command
// perform unique prefix matching
func (c *Command) FindOne(name string) (*Command, error) {
	matches, err := c.Find(name)
	if err != nil {
		return nil, err
	}

	if len(matches) == 1 {
		return matches[0], nil
	}
	return nil, fmt.Errorf("cmd=%s: no such command %q", c.Name, name)
}

// return all matching commands
// search one level deep for a command
// perform unique prefix matching
func (c *Command) Find(name string) ([]*Command, error) {
	name = strings.TrimRight(name, " ")
	matches := make([]*Command, 0)
	for _, scmd := range c.SubCmd {
		if strings.HasPrefix(scmd.Name, name) {
			matches = append(matches, scmd)
		}
	}
	if len(matches) == 0 {
		return nil, ErrNoSuchCommand
	}

	return matches, nil
}

// end
