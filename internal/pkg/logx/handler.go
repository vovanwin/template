package logx

import (
	"context"
	"log/slog"
	"strings"
)

// ComponentHandler — обработчик, позволяющий переопределять уровень логирования для отдельных компонентов.
type ComponentHandler struct {
	base      slog.Handler
	overrides map[string]slog.Level
	component string
}

// NewComponentHandler создает новый ComponentHandler.
func NewComponentHandler(base slog.Handler, overrides map[string]string) *ComponentHandler {
	lvlOverrides := make(map[string]slog.Level)
	for k, v := range overrides {
		var level slog.Level
		if err := level.UnmarshalText([]byte(strings.ToUpper(v))); err == nil {
			lvlOverrides[k] = level
		}
	}
	return &ComponentHandler{
		base:      base,
		overrides: lvlOverrides,
	}
}

func (h *ComponentHandler) Enabled(ctx context.Context, level slog.Level) bool {
	if h.component != "" {
		// Сначала ищем точное совпадение
		if overrideLevel, ok := h.overrides[h.component]; ok {
			return level >= overrideLevel
		}

		// Опционально: поиск по префиксу (например, "telegram" для "telegram:reminder")
		if idx := strings.LastIndex(h.component, ":"); idx != -1 {
			prefix := h.component[:idx]
			if overrideLevel, ok := h.overrides[prefix]; ok {
				return level >= overrideLevel
			}
		}
	}
	return h.base.Enabled(ctx, level)
}

func (h *ComponentHandler) Handle(ctx context.Context, r slog.Record) error {
	return h.base.Handle(ctx, r)
}

func (h *ComponentHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	newHandler := &ComponentHandler{
		base:      h.base.WithAttrs(attrs),
		overrides: h.overrides,
		component: h.component,
	}
	for _, attr := range attrs {
		if attr.Key == "component" {
			newHandler.component = attr.Value.String()
		}
	}
	return newHandler
}

func (h *ComponentHandler) WithGroup(name string) slog.Handler {
	return &ComponentHandler{
		base:      h.base.WithGroup(name),
		overrides: h.overrides,
		component: h.component,
	}
}

// WithComponent возвращает логгер с установленным атрибутом "component".
func WithComponent(log *slog.Logger, name string) *slog.Logger {
	return log.With("component", name)
}
