package logger

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/6543/logfile-open"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v3"
	"golang.org/x/term"
)

var GlobalLoggerFlags = []cli.Flag{
	&cli.StringFlag{
		Sources: cli.EnvVars("LOG_LEVEL"),
		Name:    "log-level",
		Usage:   "日志level(默认info)",
		Value:   "info",
	},
	&cli.StringFlag{
		Sources: cli.EnvVars("LOG_FILE"),
		Name:    "log-file",
		Usage:   "日志输出(stdout或者stderr，默认stderr)",
		Value:   "stderr",
	},
	&cli.BoolFlag{
		Sources: cli.EnvVars("DEBUG_PRETTY"),
		Name:    "pretty",
		Usage:   "美化日志打印",
		Value:   isInteractiveTerminal(), // make pretty on interactive terminal by default
	},
	&cli.BoolFlag{
		Sources: cli.EnvVars("DEBUG_NOCOLOR"),
		Name:    "nocolor",
		Usage:   "禁用彩色调试输出(只有设置了pretty时才有效)",
		Value:   !isInteractiveTerminal(), // do color on interactive terminal by default
	},
}

func SetupGlobalLogger(ctx context.Context, c *cli.Command, outputLvl bool) error {
	logLevel := c.String("log-level")
	pretty := c.Bool("pretty")
	logFile := c.String("log-file")
	noColor := c.Bool("nocolor")

	var file io.ReadWriteCloser
	switch logFile {
	case "", "stderr": // default case
		file = os.Stderr
	case "stdout":
		file = os.Stdout
	default: // a file was set
		openFile, err := logfile.OpenFileWithContext(ctx, logFile, 0o660)
		if err != nil {
			return fmt.Errorf("could not open log file '%s': %w", logFile, err)
		}
		file = openFile
		noColor = true
	}

	log.Logger = zerolog.New(file).With().Timestamp().Logger()

	if pretty {
		log.Logger = log.Output(
			zerolog.ConsoleWriter{
				Out:     file,
				NoColor: noColor,
			},
		)
	}

	lvl, err := zerolog.ParseLevel(logLevel)
	if err != nil {
		return fmt.Errorf("unknown logging level: %s", logLevel)
	}
	zerolog.SetGlobalLevel(lvl)

	// if debug or trace also log the caller
	if zerolog.GlobalLevel() <= zerolog.DebugLevel {
		log.Logger = log.With().Caller().Logger()
	}

	if outputLvl {
		log.Info().Msgf("log level: %s", zerolog.GlobalLevel().String())
	}

	return nil
}

func isInteractiveTerminal() bool {
	return term.IsTerminal(int(os.Stdout.Fd()))
}
