package main

/*
 simple command line layer for making tab complete interactive CLIs

 cd [path]
 pwd

 show ?
 show object [path]

 delete [path]
 makedir [path]
 upload [localpath] [path]

 object set mtime [mtime]



*/

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/sigmonsays/tabsy/tab"
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
		Exec: func(ctx tab.CommandContext) error {
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
		Exec: func(ctx tab.CommandContext) error {
			fmt.Printf("profiles..\n")
			return nil
		},
	}
	show.Add(profiles)

	// show cpu
	cpu := &tab.Command{
		Name:        "cpu",
		Description: "cpu utilization",
		Exec: func(ctx tab.CommandContext) error {
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
		Exec: func(ctx tab.CommandContext) error {
			path := ctx.Arg(0)
			prompt.Path = path
			fmt.Printf("cd %s\n", prompt.Path)

			return nil
		},
	}
	cli.Add(cd)

	echo := &tab.Command{
		Name:        "echo",
		Description: "echo",
		Exec: func(ctx tab.CommandContext) error {
			fmt.Printf("%s\n", ctx.Args())
			return nil
		},
	}
	cli.Add(echo)

	// ls
	ls := &tab.Command{
		Name:        "ls",
		Description: "list files",
		Exec: func(ctx tab.CommandContext) error {
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
		Exec: func(ctx tab.CommandContext) error {
			fmt.Printf("%s\n", prompt.Path)
			return nil
		},
	}
	cli.Add(pwd)

	// exit
	exit := &tab.Command{
		Name: "exit",
		Exec: func(ctx tab.CommandContext) error {
			return io.EOF
		},
	}
	cli.Add(exit)
	cli.Add(exit.Alias("quit"))

	// foo1 and foo2
	foo1 := &tab.Command{
		Name:        "foo1",
		Description: "print foo",
		Exec: func(ctx tab.CommandContext) error {
			fmt.Printf("foo1\n")
			return nil
		},
	}
	foo2 := foo1.Alias("foo2")
	foo2.Exec = func(ctx tab.CommandContext) error {
		fmt.Printf("foo2\n")
		return nil
	}
	cli.Add(foo1)
	cli.Add(foo2)

	cli.InitTerm()
	tab.Loop(cli, nil)
	cli.ReleaseTerm()
	fmt.Println()
}
