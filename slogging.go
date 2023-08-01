package slogging

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"time"
)

type source int

const (
	SourceNone source = iota
	SourceFile
	SourcePath
)

type configuration struct {
	source   source
	stampFmt string
	level    slog.Level
}

type option func(*configuration)

var Options options

type options struct{}

func (options) Source(s source) option    { return func(c *configuration) { c.source = s } }
func (options) StampFmt(s string) option  { return func(c *configuration) { c.stampFmt = s } }
func (options) Level(l slog.Level) option { return func(c *configuration) { c.level = l } }

func HandlerOptions(options ...option) *slog.HandlerOptions {
	var config configuration
	for _, option := range options {
		option(&config)
	}
	return &slog.HandlerOptions{
		Level:     config.level,
		AddSource: config.source != SourceNone,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			switch a.Key {
			case "source":
				if config.source == SourceFile {
					source := a.Value.Any().(*slog.Source)
					return slog.String("source", fmt.Sprintf("%s:%d", filepath.Base(source.File), source.Line))
				}
			case "time":
				if len(config.stampFmt) > 0 {
					return slog.String("time", a.Value.Time().Format(config.stampFmt))
				}
			}
			return a
		},
	}
}

func SetScriptingLogger(outs ...io.Writer) {
	slog.SetDefault(slog.New(slog.NewTextHandler(io.MultiWriter(append(outs, os.Stderr)...), HandlerOptions(
		Options.StampFmt(time.TimeOnly),
		Options.Source(SourceFile),
	))))
}
