package ui

import (
	"net/http"
	"regexp"

	"github.com/vovanwin/template/internal/model"
)

var dateRegexp = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)

// ParseFilters извлекает активные фильтры из GET-параметров запроса.
func ParseFilters(r *http.Request, filters []model.FilterDef) []model.ActiveFilter {
	q := r.URL.Query()
	var result []model.ActiveFilter

	for _, fd := range filters {
		switch fd.Type {
		case model.FilterEnum:
			val := q.Get("f_" + fd.Key)
			if val == "" {
				continue
			}
			// Валидируем что значение из допустимых опций
			valid := false
			for _, opt := range fd.Options {
				if opt.Value == val {
					valid = true
					break
				}
			}
			if !valid {
				continue
			}
			result = append(result, model.ActiveFilter{
				Key:   fd.Key,
				Type:  model.FilterEnum,
				Value: val,
			})

		case model.FilterString:
			val := q.Get("f_" + fd.Key)
			if val == "" {
				continue
			}
			result = append(result, model.ActiveFilter{
				Key:   fd.Key,
				Type:  model.FilterString,
				Value: val,
			})

		case model.FilterDateRange:
			from := q.Get("f_" + fd.Key + "_from")
			to := q.Get("f_" + fd.Key + "_to")
			if from == "" && to == "" {
				continue
			}
			// Валидируем формат дат
			if from != "" && !dateRegexp.MatchString(from) {
				from = ""
			}
			if to != "" && !dateRegexp.MatchString(to) {
				to = ""
			}
			if from == "" && to == "" {
				continue
			}
			result = append(result, model.ActiveFilter{
				Key:     fd.Key,
				Type:    model.FilterDateRange,
				Value:   from,
				ValueTo: to,
			})

		case model.FilterNumber:
			val := q.Get("f_" + fd.Key)
			if val == "" {
				continue
			}
			op := q.Get("f_" + fd.Key + "_op")
			if op != "gt" && op != "lt" && op != "eq" {
				op = "eq"
			}
			result = append(result, model.ActiveFilter{
				Key:      fd.Key,
				Type:     model.FilterNumber,
				Value:    val,
				Operator: op,
			})
		}
	}

	return result
}
