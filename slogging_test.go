package slogging_test

import (
	"log/slog"
	"strings"
	"testing"
	"time"

	"github.com/mdwhatcott/slogging"
)

func Test(t *testing.T) {
	options := slogging.HandlerOptions(
		slogging.Options.Level(slog.LevelWarn),
		slogging.Options.Source(slogging.SourceFile),
		slogging.Options.StampFmt(time.TimeOnly),
	)
	slog.SetDefault(slog.New(slog.NewTextHandler(&tWriter{T: t}, options)))
	slog.Warn("hello", slog.String("recipient", "world"))
}

type tWriter struct{ *testing.T }

func (this *tWriter) Write(p []byte) (int, error) {
	this.Log(strings.TrimSpace(string(p)))
	return len(p), nil
}
