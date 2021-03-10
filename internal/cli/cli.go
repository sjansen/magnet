package cli

import (
	"fmt"

	"github.com/alecthomas/kong"

	"github.com/sjansen/magnet/internal/config"
)

type context struct {
	cfg *config.WebUI
}

var cli struct {
	Runserver runserverCmd `cmd:"cmd"`
}

func ParseAndRun() {
	ctx := kong.Parse(&cli,
		kong.UsageOnError(),
	)

	fmt.Println("Loading config...")
	cfg, err := config.LoadWebUIConfig()
	ctx.FatalIfErrorf(err)

	err = ctx.Run(&context{
		cfg: cfg,
	})
	ctx.FatalIfErrorf(err)
}
