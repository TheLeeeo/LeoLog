package leolog

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"sync"
)

const (
	defaultTimeFormat = "2006-01-02 [15:04:05]"
)

func NewHandler(slogOpts *slog.HandlerOptions, leoOpts ...HandlerOption) *Handler {
	if slogOpts == nil {
		slogOpts = &slog.HandlerOptions{}
	}
	b := &bytes.Buffer{}
	h := &Handler{
		b: b,
		h: slog.NewJSONHandler(b, &slog.HandlerOptions{
			Level:       slogOpts.Level,
			AddSource:   slogOpts.AddSource,
			ReplaceAttr: suppressDefaults(slogOpts.ReplaceAttr),
		}),
		m: &sync.Mutex{},

		cfg: config{
			timeFormat: defaultTimeFormat,
			escapeHTML: false,
		},
	}

	for _, opt := range leoOpts {
		opt(h)
	}

	return h
}

type Handler struct {
	h slog.Handler
	b *bytes.Buffer
	m *sync.Mutex

	cfg config
}

type config struct {
	timeFormat string
	escapeHTML bool
}

func (h *Handler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.h.Enabled(ctx, level)
}

func (h *Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &Handler{h: h.h.WithAttrs(attrs), b: h.b, m: h.m}
}

func (h *Handler) WithGroup(name string) slog.Handler {
	return &Handler{h: h.h.WithGroup(name), b: h.b, m: h.m}
}

func (h *Handler) Handle(ctx context.Context, r slog.Record) error {
	level := r.Level.String() + ":"

	switch r.Level {
	case slog.LevelDebug:
		level = colorize(darkGray, level)
	case slog.LevelInfo:
		level = colorize(cyan, level)
	case slog.LevelWarn:
		level = colorize(lightYellow, level)
	case slog.LevelError:
		level = colorize(lightRed, level)
	}

	logMessage := fmt.Sprintf("%s %s %s",
		colorize(lightGray, r.Time.Format(h.cfg.timeFormat)),
		level,
		colorize(white, r.Message),
	)

	if r.NumAttrs() > 0 {
		attrs, err := h.computeAttrs(ctx, r)
		if err != nil {
			return err
		}

		buffer := &bytes.Buffer{}
		encoder := json.NewEncoder(buffer)
		encoder.SetIndent("", "  ")
		encoder.SetEscapeHTML(h.cfg.escapeHTML)

		if err = encoder.Encode(attrs); err != nil {
			return fmt.Errorf("error when marshaling attrs: %w", err)
		}

		recordString := buffer.String()
		// It might begin with attributes but end up empty because they were skipped by a custom handler
		if recordString != "{}" {
			logMessage += " " + colorize(darkGray, recordString)
		}
	}

	fmt.Println(logMessage)

	return nil
}

func (h *Handler) computeAttrs(
	ctx context.Context,
	r slog.Record,
) (map[string]any, error) {
	h.m.Lock()
	defer func() {
		h.b.Reset()
		h.m.Unlock()
	}()
	if err := h.h.Handle(ctx, r); err != nil {
		return nil, fmt.Errorf("error when calling inner handler's Handle: %w", err)
	}

	var attrs map[string]any
	err := json.Unmarshal(h.b.Bytes(), &attrs)
	if err != nil {
		return nil, fmt.Errorf("error when unmarshaling inner handler's Handle result: %w", err)
	}
	return attrs, nil
}

func suppressDefaults(
	next func([]string, slog.Attr) slog.Attr,
) func([]string, slog.Attr) slog.Attr {
	return func(groups []string, a slog.Attr) slog.Attr {
		if a.Key == slog.TimeKey ||
			a.Key == slog.LevelKey ||
			a.Key == slog.MessageKey {
			return slog.Attr{}
		}
		if next == nil {
			return a
		}
		return next(groups, a)
	}
}
