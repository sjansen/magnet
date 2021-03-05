package main

import (
	"fmt"
	"os"

	"github.com/alecthomas/kong"
	"github.com/joho/godotenv"

	"github.com/sjansen/magnet/internal/build"
	"github.com/sjansen/magnet/internal/config"
	"github.com/sjansen/magnet/internal/server"
)

func init() {
	env := os.Getenv("MAGNET_ENV")
	if "" == env {
		env = "development"
	}

	_ = godotenv.Load(".env." + env + ".local")
	if env != "test" {
		_ = godotenv.Load(".env.local")
	}
	_ = godotenv.Load(".env." + env)
	_ = godotenv.Load()
}

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

	fmt.Println("Loading config...")
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
	fmt.Println("Starting server...")
	s, err := server.New(ctx.cfg)
	if err != nil {
		return err
	}
	return s.ListenAndServe()
}
