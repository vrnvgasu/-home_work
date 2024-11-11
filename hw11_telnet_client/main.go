package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/spf13/pflag" //nolint:depguard
	"golang.org/x/sync/errgroup"
)

const (
	errTimeout = "incorrect timeout. Must be like `1s` or `1m` or `1h`"
)

type TelnetParams struct {
	Timeout time.Duration
	Address string
}

func main() {
	params, err := parseArgs()
	if err != nil {
		log.Fatalf("parse args: %s", err)
	}

	ctx, cansel := context.WithTimeout(context.Background(), params.Timeout)
	defer cansel()
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	tClient := NewTelnetClient(params.Address, params.Timeout, os.Stdin, os.Stdout)
	defer tClient.Close()
	if err = tClient.Connect(); err != nil {
		panic(fmt.Errorf("connect: %w", err))
	}

	listen(ctx, stop, tClient)

	<-ctx.Done()
	tClient.Close()
}

func listen(ctx context.Context, cancel context.CancelFunc, tClient TelnetClient) {
	go func() {
		eg, _ := errgroup.WithContext(ctx)
		eg.Go(func() error {
			if err := tClient.Send(); err != nil {
				return err
			}
			cancel()
			return nil
		})
		eg.Go(func() error {
			if err := tClient.Receive(); err != nil {
				return err
			}
			cancel()
			return nil
		})
		if err := eg.Wait(); err != nil {
			log.Fatalf("connect failed: %s", err)
		}
	}()
}

func parseArgs() (p TelnetParams, err error) {
	var timeoutFlag string
	pflag.StringVar(&timeoutFlag, "timeout", "10s", "time to wait for")
	pflag.Parse()

	if len(timeoutFlag) < 2 {
		return p, fmt.Errorf("%s: %s", errTimeout, timeoutFlag)
	}

	t, err := strconv.Atoi(timeoutFlag[0 : len(timeoutFlag)-1])
	if err != nil {
		return p, fmt.Errorf("%s: %w", errTimeout, err)
	}

	switch timeoutFlag[len(timeoutFlag)-1:] {
	case "s":
		p.Timeout = time.Duration(t) * time.Second
	case "m":
		p.Timeout = time.Duration(t) * time.Minute
	case "h":
		p.Timeout = time.Duration(t) * time.Hour
	default:
		return p, fmt.Errorf("%s: %s", errTimeout, timeoutFlag)
	}

	args := os.Args
	if len(args) < 3 {
		return p, errors.New("incorrect arguments, need like:$ go-telnet --timeout=10s host port")
	}
	p.Address = args[len(args)-2] + ":" + args[len(args)-1]

	return p, nil
}
