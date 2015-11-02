package tab

import (
	"log"
	"strings"

	"github.com/chzyer/readline"
)

type Readliner interface {
	Readline() (string, error)

	SetPrompt(Prompt)

	Release()
}

func NewReadline(root *RootCommand) (*Readline, error) {
	pc := make([]*readline.PrefixCompleter, 0)

	for _, s := range root.SubCmd {
		pc = buildCompletion(root, "", pc, nil, s)
	}

	completer := readline.NewPrefixCompleter(pc...)

	// prompt := "\033[31m»\033[0m "
	prompt := "(init)> "

	rlconf := &readline.Config{
		Prompt:       prompt,
		HistoryFile:  "/tmp/readline.history",
		AutoComplete: completer,
	}

	rl, err := readline.NewEx(rlconf)
	if err != nil {
		return nil, err
	}

	log.SetOutput(rl.Stderr())

	r := &Readline{
		rl: rl,
	}

	return r, nil
}

type Readline struct {
	rl *readline.Instance
}

func (r *Readline) Readline() (string, error) {
	return r.rl.Readline()
}

func (r *Readline) SetPrompt(prompt Prompt) {
	r.rl.SetPrompt(prompt.String())
}

func (r *Readline) Release() {
	r.rl.Close()
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
