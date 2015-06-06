package tab

import (
	"fmt"
	"os"
	"strings"
)

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

	dlog *os.File
}

func (c *RootCommand) OpenDebugLog(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	c.dlog = f
	return nil
}

func (c *RootCommand) dbg(s string, args ...interface{}) {
	if c.dlog == nil {
		return
	}
	line := fmt.Sprintf(s, args...)
	fmt.Fprintf(c.dlog, line+"\n")
}

// find which command will be called to execute given a full text
// returns root command if no match
func (c *RootCommand) FindCommand(line string) (*Command, error) {
	if line == "" {
		return c.Command, ErrEmptyCommand
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
	// trailing_space := strings.HasSuffix(line, " ")

	var err error
	var cmd *Command
	var cmds []*Command
	cmd = c.Command
	fields := strings.Fields(line)
	for i := 0; i < len(fields); i++ {

		field := fields[i]
		cmds, err = cmd.Find(field)
		if err != nil {
			c.dbg("cmd=%s find %s: %s", cmd.Name, field, err)
			break
		}

		if len(cmds) == 0 {
			err = ErrNoSuchCommand
		} else if len(cmds) == 1 {
			cmd = cmds[0]
		} else {
			err = ErrAmbiguousCommand
		}
	}
	if err != nil {
		return err
	}

	if cmd.Exec == nil {
		return ErrCommandNotExec.Errorf("%s", cmd.Name)
	}
	c.Ctx.args = fields[1:]
	err = cmd.Exec(c.Ctx)
	if err != nil {
		c.dbg("%s: %s", cmd.Name, err)
	}

	return err
}
