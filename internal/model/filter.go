package model

type FilterType string

const (
	FilterEnum      FilterType = "enum"
	FilterDateRange FilterType = "date_range"
	FilterString    FilterType = "string"
	FilterNumber    FilterType = "number"
)

type FilterOption struct {
	Value string
	Label string
}

type FilterDef struct {
	Type        FilterType
	Key         string         // ключ в query params и имя колонки
	Label       string         // отображаемое название
	Options     []FilterOption // только для FilterEnum
	Placeholder string         // для string/number
}

type ActiveFilter struct {
	Key      string
	Type     FilterType
	Value    string // значение (для enum/string/number/date_from)
	ValueTo  string // для date_range (date_to)
	Operator string // для number: "gt", "lt", "eq"
}
