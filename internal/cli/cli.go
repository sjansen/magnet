package cli

import (
	"github.com/alecthomas/kong"
)

var cli struct {
	Runserver runserverCmd `cmd:"cmd"`
}

// ParseAndRun parses command line arguments then runs the matching command.
func ParseAndRun() {
	ctx := kong.Parse(&cli,
		kong.UsageOnError(),
	)

	err := ctx.Run()
	ctx.FatalIfErrorf(err)
}
