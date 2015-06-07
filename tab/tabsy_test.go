package tab

import (
	"testing"
)

var foo_called bool
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

	//
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
