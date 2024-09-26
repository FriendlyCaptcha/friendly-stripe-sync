package main

import (
	"context"
	"io"
	"os"
	"os/signal"
	"syscall"

	"github.com/friendlycaptcha/friendly-stripe-sync/cmd"
)

func run(ctx context.Context, w io.Writer, args []string) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer cancel()

	err := cmd.Execute(ctx, w, args)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	ctx := context.Background()
	if err := run(ctx, os.Stdout, os.Args[1:]); err != nil {
		os.Exit(1)
	}
}
