package flagsui

import (
	"encoding/json"
	"net/http"

	"github.com/vovanwin/template/config"
)

// Handler возвращает HTTP handler для просмотра текущих значений feature flags.
// Монтируется на debug-порт: server.WithDebugHandler("/flags", flagsui.Handler(flags))
func Handler(flags *config.Flags) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /flags", func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Accept") == "application/json" {
			handleJSON(w, flags)
			return
		}
		handleHTML(w, flags)
	})

	mux.HandleFunc("GET /flags/api", func(w http.ResponseWriter, r *http.Request) {
		handleJSON(w, flags)
	})

	return mux
}

type flagValue struct {
	Name    string `json:"name"`
	Value   any    `json:"value"`
	Default any    `json:"default"`
	Type    string `json:"type"`
}

func collectFlags(flags *config.Flags) []flagValue {
	defaults := config.DefaultFlagValues()
	var result []flagValue

	for key, def := range defaults {
		fv := flagValue{
			Name:    key,
			Default: def,
		}
		switch def.(type) {
		case bool:
			fv.Type = "bool"
			fv.Value = flags.Store().GetBool(key, def.(bool))
		case int:
			fv.Type = "int"
			fv.Value = flags.Store().GetInt(key, def.(int))
		case float64:
			fv.Type = "float"
			fv.Value = flags.Store().GetFloat(key, def.(float64))
		case string:
			fv.Type = "string"
			fv.Value = flags.Store().GetString(key, def.(string))
		}
		result = append(result, fv)
	}

	return result
}

func handleJSON(w http.ResponseWriter, flags *config.Flags) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(collectFlags(flags))
}

func handleHTML(w http.ResponseWriter, flags *config.Flags) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	data := collectFlags(flags)

	w.Write([]byte(`<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<title>Feature Flags</title>
<style>
  body { font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif; margin: 40px; background: #f5f5f5; }
  h1 { color: #333; }
  table { border-collapse: collapse; width: 100%; max-width: 800px; background: white; border-radius: 8px; overflow: hidden; box-shadow: 0 1px 3px rgba(0,0,0,0.1); }
  th, td { padding: 12px 16px; text-align: left; border-bottom: 1px solid #eee; }
  th { background: #f8f9fa; font-weight: 600; color: #555; }
  .type { color: #888; font-size: 0.85em; }
  .value { font-family: monospace; font-weight: 600; }
  .bool-true { color: #22863a; }
  .bool-false { color: #cb2431; }
  .default { color: #888; font-family: monospace; font-size: 0.85em; }
  .info { color: #666; margin-bottom: 20px; }
  a { color: #0366d6; }
</style>
</head>
<body>
<h1>Feature Flags</h1>
<p class="info">Текущие значения флагов. <a href="/flags/api">JSON API</a></p>
<table>
<tr><th>Flag</th><th>Type</th><th>Value</th><th>Default</th></tr>`))

	for _, f := range data {
		valClass := "value"
		if f.Type == "bool" {
			if f.Value == true {
				valClass += " bool-true"
			} else {
				valClass += " bool-false"
			}
		}

		valJSON, _ := json.Marshal(f.Value)
		defJSON, _ := json.Marshal(f.Default)

		w.Write([]byte(`<tr><td>` + f.Name + `</td><td class="type">` + f.Type + `</td><td class="` + valClass + `">` + string(valJSON) + `</td><td class="default">` + string(defJSON) + `</td></tr>`))
	}

	w.Write([]byte(`</table></body></html>`))
}
