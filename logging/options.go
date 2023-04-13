package logging

import (
	"io"
	"os"

	"github.com/go-kit/kit/log"
	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	StdoutFile = "stdout"
)

// Options stores the configuration of a Logger.  Lumberjack is used for rolling files.
type Options struct {
	// File is the system file path for the log file.  If set to "stdout", this will log to os.Stdout.
	// Otherwise, a lumberjack.Logger is created
	File string `json:"file"`

	// MaxSize is the lumberjack MaxSize
	MaxSize int `json:"maxsize"`

	// MaxAge is the lumberjack MaxAge
	MaxAge int `json:"maxage"`

	// MaxBackups is the lumberjack MaxBackups
	MaxBackups int `json:"maxbackups"`

	// JSON is a flag indicating whether JSON logging output is used.  The default is false,
	// meaning that logfmt output is used.
	JSON bool `json:"json"`

	// Level is the error level to output: ERROR, INFO, WARN, or DEBUG.  Any unrecognized string,
	// including the empty string, is equivalent to passing ERROR.
	Level string `json:"level"`
}

func (o *Options) output() io.Writer {
	if o != nil && len(o.File) > 0 && o.File != StdoutFile {
		return &lumberjack.Logger{
			Filename:   o.File,
			MaxSize:    o.MaxSize,
			MaxAge:     o.MaxAge,
			MaxBackups: o.MaxBackups,
		}
	}

	return log.NewSyncWriter(os.Stdout)
}

func (o *Options) loggerFactory() func(io.Writer) log.Logger {
	if o != nil && o.JSON {
		return log.NewJSONLogger
	}

	return log.NewLogfmtLogger
}

func (o *Options) level() string {
	if o != nil {
		return o.Level
	}

	return ""
}
