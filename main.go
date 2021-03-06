package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/sjansen/magnet/internal/cli"
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

func main() {
	if os.Getenv("AWS_LAMBDA_RUNTIME_API") != "" {
		fmt.Fprintln(os.Stderr, "This executable should not be run on AWS Lambda.")
		os.Exit(1)
	}

	cli.ParseAndRun()
}
