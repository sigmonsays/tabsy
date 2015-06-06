package tab

import (
	"fmt"
	"strings"
)

var ErrNoSuchCommand = fmt.Errorf("no such command")
var ErrAmbiguousCommand = fmt.Errorf("ambiguous")
var ErrCommandNotExec = fmt.Errorf("command not executable")
var ErrEmptyCommand = fmt.Errorf("empty command")

type Prompt interface {
	String() string
}

type ExecuteCommand func(*Context) error

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

// return all matching sub commands
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
