package cli

import (
	"context"
	"fmt"

	"github.com/sjansen/magnet/internal/webui/server"
)

type runserverCmd struct{}

func (cmd *runserverCmd) Run() error {
	fmt.Println("Starting server...")
	ctx := context.Background()
	s, err := server.New(ctx)
	if err != nil {
		return err
	}
	return s.ListenAndServe()
}
