package tab

import (
	"fmt"
	"strings"
)

func NewCommandSet(name string) *RootCommand {
	root := &RootCommand{
		Command: &Command{
			Name:   name,
			IsRoot: true,
		},
	}
	return root
}

type ErrorHandler func(err error)

type RootCommand struct {
	*Command
	Ctx          CommandContext
	Readline     Readliner
	ErrorHandler ErrorHandler
}

func (c *RootCommand) WithReadline(r Readliner) *RootCommand {
	c.Readline = r
	return c
}

func (c *RootCommand) UseContext(ctx CommandContext) {
	c.Ctx = ctx
}

func (c *RootCommand) InitTerm() error {

	// setup a buncha defaults if they dont configure things....

	if c.Readline == nil {
		rl, err := NewReadline(c)
		if err != nil {
			return err
		}
		c.Readline = rl
	}

	if c.Ctx == nil {
		c.Ctx = NewContext(c.Readline)
	}

	return nil
}

func (c *RootCommand) ReleaseTerm() error {
	c.dbg("release term")
	c.Readline.Release()
	return nil
}

func (r *RootCommand) Context() CommandContext {
	return r.Ctx
}

func (c *RootCommand) dbg(s string, args ...interface{}) {
	if c.Ctx == nil {
		fmt.Printf("[dbg] "+s+"\n", args...)
		return
	}
	if c.Ctx.Dbg == nil {
		fmt.Printf("[dbg] "+s+"\n", args...)
		return
	}
	c.Ctx.Dbg(s, args...)
}

// find which command will be called to execute given a full text
// returns root command if no match (should it?)
func (c *RootCommand) FindCommand(line string) (*Command, error) {
	if line == "" {
		return c.Command, ErrEmptyCommand
	}
	fields := strings.Fields(line)

	var ret *Command
	var scmd *Command
	scmd = c.Command
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

func (c *RootCommand) OnError(fn ErrorHandler) {
	c.ErrorHandler = fn
}

// calls the error handler with an error
func (c *RootCommand) withError(err error) error {
	if err == nil {
		return err
	}

	if c.ErrorHandler != nil {
		c.ErrorHandler(err)
	}
	return err
}

func (c *RootCommand) Dispatch(line string) error {
	if line == "" {
		return nil
	}
	c.dbg("dispatch %q", line)
	// trailing_space := strings.HasSuffix(line, " ")

	var err error
	var cmd *Command
	var cmds []*Command
	cmd = c.Command
	fields := strings.Fields(line)
	c.dbg("starting search at root cmd=%s", c.Name)
	for i := 0; i < len(fields); i++ {

		field := fields[i]
		cmds, err = cmd.Find(field)
		if err != nil {
			c.dbg("cmd=%s find %s: %s", cmd.Name, field, err)
			break
		}

		c.dbg("find field=%s returned %d entries", field, len(cmds))

		if len(cmds) == 0 {
			break
		} else if len(cmds) == 1 {
			cmd = cmds[0]
		} else {
			choices := make([]string, 0)
			for _, c := range cmds {
				choices = append(choices, c.Name)
			}
			err = ErrAmbiguousCommand.Errorf("%s", choices)
			c.dbg("command is ambiguous")

			// see if we have an exact match since Find will return based on prefix
			for _, c2 := range cmds {
				if c2.Name == field {
					cmd = c2
					c.dbg("exact match for command %s", field)
					err = nil
					break
				}
			}
			break
		}
	}

	// before
	if cmd != nil && cmd.Before != nil {
		c.Ctx.SetArgs(fields)
		c.dbg("executing Before func %s (args %s)", cmd.Name, fields)
		err = cmd.Before(c.Ctx)
		if err != nil {
			c.dbg("cmd=%s Before returned error: %s", cmd.Name, err)
			return c.withError(err)
		}
	}

	// command can be ambiguous
	if err != nil || cmd == nil {
		return c.withError(err)
	}

	if cmd.IsRoot {
		c.dbg("no such top level command: %s", fields[0])
		return c.withError(ErrNoSuchCommand.Errorf("%s", line))
	}

	if cmd.Exec == nil {
		c.dbg("command not executable: %s", cmd.Name)
		return c.withError(ErrCommandNotExec.Errorf("%s", cmd.Name))
	}
	c.Ctx.SetArgs(fields[1:])
	c.dbg("executing %s (args %s)", cmd.Name, fields[1:])

	err = cmd.Exec(c.Ctx)
	if err != nil {
		c.dbg("%s: %s", cmd.Name, err)
		return c.withError(err)
	}

	return err
}
