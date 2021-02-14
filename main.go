package main

import (
	"fmt"
	"os"

	"github.com/alecthomas/kong"
	_ "github.com/joho/godotenv/autoload"

	"github.com/sjansen/magnet/internal/build"
	"github.com/sjansen/magnet/internal/config"
	"github.com/sjansen/magnet/internal/server"
)

type context struct {
	cfg *config.Config
}

var cli struct {
	Runserver runserverCmd `cmd:"cmd"`
}

func main() {
	if os.Getenv("AWS_LAMBDA_RUNTIME_API") != "" {
		startLambda()
	} else {
		startCLI()
	}
}

func startCLI() {
	ctx := kong.Parse(&cli,
		kong.UsageOnError(),
	)

	cfg, err := config.New()
	ctx.FatalIfErrorf(err)

	err = ctx.Run(&context{
		cfg: cfg,
	})
	ctx.FatalIfErrorf(err)
}

func startLambda() {
	fmt.Println("GitSHA:", build.GitSHA)
	fmt.Println("Timestamp:", build.Timestamp)

	cfg, err := config.New()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	s, err := server.New(cfg)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	s.HandleLambda()
}

type runserverCmd struct{}

func (cmd *runserverCmd) Run(ctx *context) error {
	s, err := server.New(ctx.cfg)
	if err != nil {
		return err
	}
	return s.ListenAndServe()
}
