package tab

import (
	"fmt"
	"strings"
)

func tabComplete(c *RootCommand, k *KeyPress) (*KeyResponse, error) {
	res := &KeyResponse{"", 0, false}
	trailing_space := strings.HasSuffix(k.line, " ")

	prefix := ""
	fields := strings.Fields(k.line)
	if len(fields) > 0 {
		if trailing_space == false {
			prefix = fields[len(fields)-1]
		}
	}

	c2, err := c.FindCommand(k.line)
	if err != nil {
		c.dbg("error %s\n", err)
		return nil, err
	}
	if c2 == nil {
		c2 = c.Command
	}

	var cmd *Command
	var cmds []*Command

	if trailing_space {
		cmds = c2.SubCmd
	} else {
		cmds, _ = c2.Find(k.line)
	}

	// list all sub commands
	ls := []*Command{}
	for _, cmd := range cmds {
		if prefix == "" || strings.HasPrefix(cmd.Name, prefix) {
			ls = append(ls, cmd)
		}

	}

	c.dbg("tabcomplete cmd=%s prefix=%s cmd-matches=%d ls=%d line=%q",
		c2.Name, prefix, len(cmds), len(ls), k.line)

	if len(ls) == 0 && c2 != nil && c2.IsRoot == false {
		// we have a command thats the exact match, lets complete the word
		p := c2.Name[len(prefix):] + " "
		res.ok = true
		res.newline += k.line + p
		res.newpos += k.pos + len(p)

	} else if len(ls) == 1 {
		cmd = cmds[0]
		p := cmd.Name[len(prefix):] + " "
		res.ok = true
		res.newline += k.line + p
		res.newpos += k.pos + len(p)

	} else if len(ls) > 0 {
		fmt.Println()
		for _, c := range ls {
			fmt.Printf("%-20s\n", c.Name)
		}
		fmt.Printf("\n%s%s", c.Ctx.Prompt, k.line)
	}
	return res, nil
}
