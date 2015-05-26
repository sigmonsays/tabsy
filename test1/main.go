package main

import (
	"fmt"
	"io"
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

	// show
	show := &Command{
		Name:        "show",
		Description: "top level show command, use ? for help",
	}

	// show profiles
	profiles := &Command{
		Name:        "profiles",
		Description: "show connection profiles",
		Exec: func(ctx *Context) error {
			fmt.Printf("profiles..\n")
			return nil
		},
	}
	show.Add(profiles)

	// show cpu
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

	// cd
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

	// pwd
	pwd := &Command{
		Name:        "pwd",
		Description: "print working directory",
		Exec: func(ctx *Context) error {
			fmt.Printf("%s\n", prompt.Path)
			return nil
		},
	}
	cli.Add(pwd)

	// exit
	exit := &Command{
		Name: "exit",
		Exec: func(ctx *Context) error {
			return io.EOF
		},
	}
	cli.Add(exit)

	quit := make(chan bool, 0)

	cli.InitTerm()
	cli.Ctx.SetPrompt(prompt)
	Loop(cli, quit)
	cli.ReleaseTerm()
	fmt.Println()
}
