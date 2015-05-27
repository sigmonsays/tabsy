package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/sigmonsays/test/tab"
)

type prompt struct {
	Path string
}

func (p *prompt) String() string {
	return fmt.Sprintf("%s > ", p.Path)
}

func main() {
	cli := tab.NewCommandSet("test1")
	cli.Description = "do something useful.."

	help := &tab.Command{
		Name:        "help",
		Description: "show help for commands",
		Exec: func(ctx *tab.Context) error {
			return nil
		},
	}
	cli.Add(help)

	// show
	show := &tab.Command{
		Name:        "show",
		Description: "top level show command, use ? for help",
	}

	// show profiles
	profiles := &tab.Command{
		Name:        "profiles",
		Description: "show connection profiles",
		Exec: func(ctx *tab.Context) error {
			fmt.Printf("profiles..\n")
			return nil
		},
	}
	show.Add(profiles)

	// show cpu
	cpu := &tab.Command{
		Name:        "cpu",
		Description: "cpu utilization",
		Exec: func(ctx *tab.Context) error {
			buf, err := ioutil.ReadFile("/proc/loadavg")
			if err != nil {
				return err
			}
			fmt.Printf("%s\n", buf)
			return nil
		},
	}
	show.Add(cpu)

	cli.Add(show)

	prompt := &prompt{"/"}

	// cd
	cd := &tab.Command{
		Name:        "cd",
		Description: "change current directory",
		Exec: func(ctx *tab.Context) error {
			path := ctx.Arg(0)
			prompt.Path = path
			ctx.SetPrompt(prompt)
			fmt.Printf("cd %s\n", prompt.Path)

			return nil
		},
	}
	cli.Add(cd)

	// ls
	ls := &tab.Command{
		Name:        "ls",
		Description: "list files",
		Exec: func(ctx *tab.Context) error {
			args := ctx.Args()
			cmd := exec.Command("ls", args...)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			err := cmd.Run()
			if err != nil {
				fmt.Printf("ERROR: %s\n", err)
			}
			return err
		},
	}
	cli.Add(ls)

	// pwd
	pwd := &tab.Command{
		Name:        "pwd",
		Description: "print working directory",
		Exec: func(ctx *tab.Context) error {
			fmt.Printf("%s\n", prompt.Path)
			return nil
		},
	}
	cli.Add(pwd)

	// exit
	exit := &tab.Command{
		Name: "exit",
		Exec: func(ctx *tab.Context) error {
			return io.EOF
		},
	}
	cli.Add(exit)

	quit := make(chan bool, 0)

	cli.InitTerm()
	cli.Ctx.SetPrompt(prompt)
	tab.Loop(cli, quit)
	cli.ReleaseTerm()
	fmt.Println()
}
