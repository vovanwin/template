package ui

import (
	"net/http"

	"github.com/a-h/templ"
)

// Render рендерит компонент. Если это HTMX запрос (включая boosted), рендерит только фрагмент.
// Если это обычный запрос, рендерит полную страницу с layout.
func (c *UIController) Render(w http.ResponseWriter, r *http.Request, fullPage templ.Component, fragment templ.Component) {
	if isHTMX(r) {
		templ.Handler(fragment).ServeHTTP(w, r)
		return
	}
	templ.Handler(fullPage).ServeHTTP(w, r)
}

func isHTMX(r *http.Request) bool {
	return r.Header.Get("HX-Request") == "true"
}

// Redirect - умный редирект, учитывающий HTMX
func (c *UIController) Redirect(w http.ResponseWriter, r *http.Request, url string) {
	if isHTMX(r) {
		// HX-Location позволяет HTMX сделать "бесшовный" переход с сохранением истории
		w.Header().Set("HX-Location", url)
		w.WriteHeader(http.StatusOK)
		return
	}
	http.Redirect(w, r, url, http.StatusFound)
}
