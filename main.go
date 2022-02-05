package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"
)

type ctxKey uint8

const lookupKey ctxKey = 0

func main() {
	var (
		list   bool
		lookup bool
	)

	flag.BoolVar(&list, "list", false, "")
	flag.BoolVar(&list, "l", false, "")
	flag.BoolVar(&lookup, "lookup", false, "")
	flag.Parse()

	if list && lookup {
		fmt.Println("cannot use --list with --force")
		os.Exit(1)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	if lookup {
		ctx = context.WithValue(ctx, lookupKey, lookup)
	}

	var run func(context.Context, io.Writer, io.Writer, string) error = chgoRun

	if list {
		run = listRun
	}

	if err := run(ctx, os.Stdout, os.Stderr, flag.Arg(0)); err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
}
