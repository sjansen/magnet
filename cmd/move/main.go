package main

import (
	"fmt"
	"os"

	"github.com/sjansen/magnet/internal/build"
	"github.com/sjansen/magnet/internal/move"
)

func main() {
	if os.Getenv("AWS_LAMBDA_RUNTIME_API") == "" {
		fmt.Fprintln(os.Stderr, "This executable should be run on AWS Lambda.")
		os.Exit(1)
	}

	fmt.Println("GitSHA:", build.GitSHA)
	fmt.Println("Timestamp:", build.Timestamp)

	m, err := move.New()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	m.StartLambdaHandler()
}
