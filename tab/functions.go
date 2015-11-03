package tab

import (
	"io"
)

func ExitFunction(ctx CommandContext) error {
	return io.EOF
}
