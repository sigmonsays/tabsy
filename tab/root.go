package tab

import (
	"log"
	"strings"

	"github.com/chzyer/readline"
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

type ErrorHandler func(err error)

type RootCommand struct {
	*Command
	Ctx          *Context
	ErrorHandler ErrorHandler
}

func buildCompletion(root *RootCommand, cmdline string, pclist []*readline.PrefixCompleter, pc *readline.PrefixCompleter, c *Command) []*readline.PrefixCompleter {
	if cmdline != "" {
		cmdline += " "
	}
	cmdline += c.Name

	root.dbg("complete cmdline [%s] pclist %d: cmd %s", cmdline, len(pclist), c.Name)

	if pc == nil {
		p := &readline.PrefixCompleter{
			Name:     []rune(strings.Trim(cmdline, " ") + " "),
			Children: make([]*readline.PrefixCompleter, 0),
		}
		pclist = append(pclist, p)
		pc = p
	}

	for _, s := range c.SubCmd {

		p2 := &readline.PrefixCompleter{
			Name:     []rune(strings.Trim(s.Name, " ") + " "),
			Children: make([]*readline.PrefixCompleter, 0),
		}

		pc.Children = append(pc.Children, p2)

		root.dbg("scmd %s cmdline %s", s.Name, cmdline)
		pclist = buildCompletion(root, cmdline, pclist, p2, s)

	}

	return pclist
}

func (c *RootCommand) InitTerm() error {

	pc := make([]*readline.PrefixCompleter, 0)

	for _, s := range c.SubCmd {
		pc = buildCompletion(c, "", pc, nil, s)
	}

	completer := readline.NewPrefixCompleter(pc...)

	// prompt := "\033[31m»\033[0m "
	prompt := "(init)> "

	rlconf := &readline.Config{
		Prompt:       prompt,
		HistoryFile:  "/tmp/readline.tmp",
		AutoComplete: completer,
	}
	c.dbg("readline init")

	rl, err := readline.NewEx(rlconf)
	if err != nil {
		return err
	}

	log.SetOutput(rl.Stderr())

	c.Ctx.rl = rl
	return nil
}

func (c *RootCommand) ReleaseTerm() error {
	c.dbg("release term")
	c.Ctx.rl.Close()
	return nil
}

func (r *RootCommand) Context() *Context {
	return r.Ctx
}

func (c *RootCommand) dbg(s string, args ...interface{}) {
	if c.Ctx.dlog == nil {
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
		}
	}

	/*

		if err != nil || cmd == nil {
			return c.withError(err)
		}
	*/

	if cmd.IsRoot {
		c.dbg("no such top level command: %s", fields[0])
		return c.withError(ErrNoSuchCommand.Errorf("%s", line))
	}

	if cmd.Exec == nil {
		c.dbg("command not executable: %s", cmd.Name)
		return c.withError(ErrCommandNotExec.Errorf("%s", cmd.Name))
	}
	c.Ctx.args = fields[1:]
	c.dbg("executing %s (args %s)", cmd.Name, c.Ctx.args)
	err = cmd.Exec(c.Ctx)
	if err != nil {
		c.dbg("%s: %s", cmd.Name, err)
		return c.withError(err)
	}

	return err
}
