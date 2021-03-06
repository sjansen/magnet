package cli

import (
	"fmt"

	"github.com/sjansen/magnet/internal/webui/server"
)

type runserverCmd struct{}

func (cmd *runserverCmd) Run(ctx *context) error {
	fmt.Println("Starting server...")
	s, err := server.New(ctx.cfg)
	if err != nil {
		return err
	}
	return s.ListenAndServe()
}
