package repository

import (
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/vovanwin/template/internal/model"
)

// FilterWhitelist маппит ключ фильтра на имя колонки в БД.
type FilterWhitelist map[string]string

// ApplyFilters добавляет условия фильтрации к squirrel SelectBuilder.
func ApplyFilters(b squirrel.SelectBuilder, filters []model.ActiveFilter, wl FilterWhitelist) squirrel.SelectBuilder {
	for _, f := range filters {
		col, ok := wl[f.Key]
		if !ok {
			continue
		}

		switch f.Type {
		case model.FilterEnum:
			b = b.Where(squirrel.Eq{col: f.Value})

		case model.FilterString:
			escaped := strings.ReplaceAll(f.Value, "%", "\\%")
			escaped = strings.ReplaceAll(escaped, "_", "\\_")
			b = b.Where(squirrel.ILike{col: "%" + escaped + "%"})

		case model.FilterDateRange:
			if f.Value != "" {
				b = b.Where(squirrel.GtOrEq{col: f.Value})
			}
			if f.ValueTo != "" {
				b = b.Where(squirrel.LtOrEq{col: f.ValueTo + " 23:59:59"})
			}

		case model.FilterNumber:
			switch f.Operator {
			case "gt":
				b = b.Where(squirrel.Gt{col: f.Value})
			case "lt":
				b = b.Where(squirrel.Lt{col: f.Value})
			default:
				b = b.Where(squirrel.Eq{col: f.Value})
			}
		}
	}
	return b
}
