package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"
)

func main() {
	var list bool

	flag.BoolVar(&list, "list", false, "")
	flag.BoolVar(&list, "l", false, "")
	flag.Parse()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	var run func(context.Context, io.Writer, io.Writer, string) error = chgoRun

	if list {
		run = listRun
	}

	if err := run(ctx, os.Stdout, os.Stderr, flag.Arg(0)); err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
}
