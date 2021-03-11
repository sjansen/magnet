package main

import (
	"fmt"
	"os"

	"github.com/sjansen/magnet/internal/build"
	"github.com/sjansen/magnet/internal/webui/server"
)

func main() {
	if os.Getenv("AWS_LAMBDA_RUNTIME_API") == "" {
		fmt.Fprintln(os.Stderr, "This executable should be run on AWS Lambda.")
		os.Exit(1)
	}

	fmt.Println("GitSHA:", build.GitSHA)
	fmt.Println("Timestamp:", build.Timestamp)

	s, err := server.New()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	s.StartLambdaHandler()
}
