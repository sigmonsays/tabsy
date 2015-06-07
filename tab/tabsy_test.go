package tab

import (
	"fmt"
	"strings"
	"testing"
)

var foo_called bool
var foo_args string
var bar_called bool

func TestDispatch(t *testing.T) {

	cli := build_cli()

	// top level commands work
	cli.Dispatch("foo")
	if foo_called != true {
		t.Errorf("foo was not called")
	}

	// sub commands work
	cli.Dispatch("foo bar")
	if bar_called != true {
		t.Errorf("bar was not called")
	}

	// arguments are properly parsed
	cli.Dispatch("foo baz boo")
	expected := "baz boo"
	if foo_args != expected {
		fmt.Printf("foo blah args -- %q\n", foo_args)
		t.Errorf("bad arguments got %s, expected %s", foo_args, expected)
	}

}

func build_cli() *RootCommand {
	cli := NewCommandSet("test1")
	cli.Description = "do something useful.."

	cli.OpenDebugLog("/tmp/term-test.log")

	foo := &Command{
		Name:        "foo",
		Description: "foo",
		Exec: func(ctx *Context) error {
			foo_called = true
			foo_args = strings.Join(ctx.Args(), " ")
			return nil
		},
	}
	cli.Add(foo)

	// bar
	bar := &Command{
		Name:        "bar",
		Description: "foo bar command",
		Exec: func(ctx *Context) error {
			bar_called = true
			return nil
		},
	}
	foo.Add(bar)

	return cli
}
