package main

import (
	"context"
	"fmt"
	"os"

	"github.com/alecthomas/kingpin/v2"
	"github.com/hashicorp/go-hclog"
	"github.com/pkg/errors"
)

var (
	Build   string
	Version string
)

const (
	level = "INFO"
	name  = "agent"
)

var (
	app          = kingpin.New(name, "insight agent")
	durationTime = app.Flag("duration-time", "Duration time ((h:hour, m:minute, s:second)").Required().String()
	logLevel     = app.Flag("log-level", "Log level (DEBUG|INFO|WARN|ERROR)").Default(level).String()
)

func main() {
	ctx := context.Background()

	if err := Run(ctx); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	os.Exit(0)
}

func Run(ctx context.Context) error {
	kingpin.MustParse(app.Parse(os.Args[1:]))

	logger, err := initLogger(ctx, *logLevel)
	if err != nil {
		return errors.Wrap(err, "failed to init logger")
	}

	if err := runAgent(ctx, logger, *durationTime); err != nil {
		return errors.Wrap(err, "failed to run agent")
	}

	return nil
}

func initLogger(_ context.Context, level string) (hclog.Logger, error) {
	return hclog.New(&hclog.LoggerOptions{
		Name:  name,
		Level: hclog.LevelFromString(level),
	}), nil
}

func runAgent(_ context.Context, logger hclog.Logger, duration string) error {
	logger.Debug("agent: runAgent")

	// TBD: FIXME

	return nil
}
