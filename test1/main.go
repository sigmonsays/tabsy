package main

import (
	"fmt"
)

type prompt struct {
	Path string
}

func (p *prompt) String() string {
	return fmt.Sprintf("%s > ", p.Path)
}

func main() {
	cli := NewCommandSet("test1")
	cli.Description = "do something useful.."

	help := &Command{
		Name:        "help",
		Description: "show help for commands",
	}
	cli.Add(help)

	show := &Command{
		Name:        "show",
		Description: "top level show command, use ? for help",
	}

	profiles := &Command{
		Name:        "profiles",
		Description: "show connection profiles",
		Exec: func(ctx *Context) error {
			fmt.Printf("profiles..\n")
			return nil
		},
	}
	show.Add(profiles)

	cpu := &Command{
		Name:        "cpu",
		Description: "cpu utilization",
		Exec: func(ctx *Context) error {
			fmt.Printf("CPUs..\n")
			return nil
		},
	}
	show.Add(cpu)

	cli.Add(show)

	prompt := &prompt{"/"}

	cd := &Command{
		Name:        "cd",
		Description: "change current directory",
		Exec: func(ctx *Context) error {
			path := ctx.Arg(0)
			prompt.Path = path
			ctx.SetPrompt(prompt)
			fmt.Printf("cd %s\n", prompt.Path)

			return nil
		},
	}
	cli.Add(cd)

	quit := make(chan bool, 0)

	cli.InitTerm()
	cli.Ctx.SetPrompt(prompt)
	Loop(cli, quit)
	cli.ReleaseTerm()
}
