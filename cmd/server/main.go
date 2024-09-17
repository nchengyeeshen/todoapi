package main

import (
	"context"
	"errors"
	"flag"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/charmbracelet/log"
)

func main() {
	os.Exit(run(os.Args, os.Stderr))
}

func run(args []string, stderr io.Writer) int {
	charmlog := log.NewWithOptions(
		stderr,
		log.Options{
			ReportTimestamp: true,
			TimeFormat:      time.RFC3339,
		},
	)
	logger := slog.New(ContextHandler{charmlog})

	env := Environment{
		Env: "development",
	}
	fs := flag.NewFlagSet(args[0], flag.ContinueOnError)
	fs.BoolVar(&env.Debug, "debug", false, "Debug mode")
	fs.StringVar(&env.ListenAddr, "listen-addr", ":8000", "HTTP listen address")
	fs.Func(
		"env",
		`Environment. Must be one of: production, development (default "development")`,
		func(s string) error {
			v := strings.ToLower(s)
			switch v {
			case "production", "development":
				env.Env = v
				return nil
			default:
				return errors.New("must be one of: production, development")
			}
		},
	)
	if err := fs.Parse(args[1:]); err != nil {
		if !errors.Is(err, flag.ErrHelp) {
			logger.Error("Parse command line flags", "err", err)
		}
		return 1
	}

	if env.Env == "production" {
		charmlog.SetFormatter(log.LogfmtFormatter)
	}
	if env.Debug {
		charmlog.SetLevel(log.DebugLevel)
	}

	todoRepo := NewInMemoryTodoRepository()

	app := NewApplication(logger, todoRepo)

	httpSrv := &http.Server{
		Addr:    env.ListenAddr,
		Handler: app.routes(),
	}

	httpSrvShutdownErrChan := make(chan error, 1)
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

		signal := <-sigChan
		logger.Info("Caught shutdown signal", "signal", signal.String())

		gracefulTimeout := 5 * time.Second
		ctx, cancel := context.WithTimeout(context.Background(), gracefulTimeout)
		defer cancel()

		logger.Info("HTTP server shutdown", "gracefulTimeout", gracefulTimeout.String())
		httpSrvShutdownErrChan <- httpSrv.Shutdown(ctx)
	}()

	logger.Info("HTTP server listen", "addr", httpSrv.Addr)
	err := httpSrv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		logger.Error("HTTP listen and serve", "err", err)
		return 1
	}

	return 0
}

type Environment struct {
	Env        string
	Debug      bool
	ListenAddr string
}
