package server

import (
	"context"
	"encoding/json"
	"fmt"
	"io/fs"
	"log/slog"
	"net"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi/v5"
)

type specEntry struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

// discoverEmbedFiles сканирует fs.FS и возвращает файлы с данным суффиксом.
func discoverEmbedFiles(fsys fs.FS, suffix string) []string {
	var files []string
	_ = fs.WalkDir(fsys, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if !d.IsDir() && strings.HasSuffix(d.Name(), suffix) {
			files = append(files, path)
		}
		return nil
	})
	return files
}

func (s *Server) initSwagger(log *slog.Logger) {
	if s.cfg.SwaggerFS == nil {
		log.Warn("SwaggerFS не задан, Swagger UI отключён")
		return
	}

	r := chi.NewRouter()

	swaggerFiles := discoverEmbedFiles(s.cfg.SwaggerFS, ".swagger.json")

	var protoFiles []string
	if s.cfg.ProtoFS != nil {
		protoFiles = discoverEmbedFiles(s.cfg.ProtoFS, ".proto")
	}

	log.Debug("swagger specs found", slog.Any("specs", swaggerFiles))
	log.Debug("proto files found", slog.Any("protos", protoFiles))

	// Статика: swagger JSON файлы из embed.FS
	r.Handle("/spec/*", http.StripPrefix("/spec/", http.FileServerFS(s.cfg.SwaggerFS)))

	// Статика: proto файлы как text/plain из embed.FS
	if s.cfg.ProtoFS != nil {
		r.Get("/proto/*", func(w http.ResponseWriter, req *http.Request) {
			relPath := chi.URLParam(req, "*")

			data, err := fs.ReadFile(s.cfg.ProtoFS, relPath)
			if err != nil {
				http.NotFound(w, req)
				return
			}
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			_, _ = w.Write(data)
		})
	}

	// API: метаданные
	r.Get("/api/specs", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		specs := make([]specEntry, 0, len(swaggerFiles))
		for _, f := range swaggerFiles {
			name := strings.TrimSuffix(filepath.Base(f), ".swagger.json")
			specs = append(specs, specEntry{Name: name, Path: f})
		}

		_ = json.NewEncoder(w).Encode(map[string]any{
			"specs":  specs,
			"protos": protoFiles,
		})
	})

	// Главная страница
	r.Get("/", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write([]byte(buildSwaggerHTML(swaggerFiles, protoFiles)))
	})

	addr := net.JoinHostPort(s.cfg.Host, s.cfg.SwaggerPort)

	s.swaggerSrv = &http.Server{
		Addr:    addr,
		Handler: r,
	}

	go func() {
		log.Info("Swagger UI запущен", slog.String("addr", addr), slog.Int("specs", len(swaggerFiles)))
		if err := s.swaggerSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("Swagger сервер остановлен с ошибкой", slog.String("error", err.Error()))
		}
	}()
}

func buildSwaggerHTML(specs []string, protos []string) string {
	if len(specs) == 0 {
		return `<!doctype html><html><body><h2>No swagger specs found</h2></body></html>`
	}

	// Собираем навигационные ссылки для API спеков
	var specLinks strings.Builder
	for _, spec := range specs {
		name := strings.TrimSuffix(filepath.Base(spec), ".swagger.json")
		specLinks.WriteString(fmt.Sprintf(
			`        <a href="?spec=%s" class="nav-link" data-spec="%s" data-url="/spec/%s">%s</a>`+"\n",
			name, name, spec, name,
		))
	}

	// Собираем JSON массив proto файлов для дерева на клиенте
	protoPaths := make([]string, 0, len(protos))
	for _, p := range protos {
		protoPaths = append(protoPaths, fmt.Sprintf(`"%s"`, p))
	}
	protoJSON := "[" + strings.Join(protoPaths, ",") + "]"

	return fmt.Sprintf(`<!doctype html>
<html>
<head>
  <meta charset="utf-8">
  <title>API Documentation</title>
  <script type="module" src="https://unpkg.com/rapidoc/dist/rapidoc-min.js"></script>
  <style>
    * { margin: 0; padding: 0; box-sizing: border-box; }
    html, body { height: 100%%; font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif; }
    body { display: flex; flex-direction: column; }

    .nav {
      display: flex;
      align-items: center;
      gap: 4px;
      padding: 8px 16px;
      background: #1e293b;
      color: #fff;
      flex-shrink: 0;
      flex-wrap: wrap;
    }
    .nav-group { display: flex; align-items: center; gap: 4px; }
    .nav-label { font-size: 11px; text-transform: uppercase; color: #94a3b8; margin-right: 4px; letter-spacing: 0.5px; }
    .nav-sep { width: 1px; height: 24px; background: #475569; margin: 0 12px; }
    .nav-link {
      padding: 4px 12px;
      border-radius: 4px;
      color: #cbd5e1;
      text-decoration: none;
      font-size: 13px;
      cursor: pointer;
      transition: background 0.15s, color 0.15s;
    }
    .nav-link:hover { background: #334155; color: #fff; }
    .nav-link.active { background: #3b82f6; color: #fff; }

    #content { flex: 1; overflow: hidden; position: relative; }
    rapi-doc { height: 100%%; }

    /* Proto browser: sidebar + code */
    #proto-browser {
      display: none;
      height: 100%%;
      flex-direction: row;
    }

    #proto-sidebar {
      width: 260px;
      min-width: 200px;
      background: #0f172a;
      border-right: 1px solid #1e293b;
      overflow-y: auto;
      padding: 12px 0;
      flex-shrink: 0;
    }
    .tree-folder {
      user-select: none;
    }
    .tree-folder-label {
      display: flex;
      align-items: center;
      gap: 6px;
      padding: 6px 12px 6px calc(12px + var(--depth, 0) * 16px);
      color: #94a3b8;
      font-size: 13px;
      font-weight: 600;
      cursor: pointer;
      transition: background 0.1s;
    }
    .tree-folder-label:hover { background: #1e293b; }
    .tree-folder-label .arrow {
      display: inline-block;
      width: 16px;
      text-align: center;
      font-size: 10px;
      transition: transform 0.15s;
    }
    .tree-folder.collapsed > .tree-children { display: none; }
    .tree-folder.collapsed > .tree-folder-label .arrow { transform: rotate(-90deg); }
    .tree-children { }
    .tree-file {
      display: flex;
      align-items: center;
      gap: 6px;
      padding: 5px 12px 5px calc(12px + var(--depth, 0) * 16px);
      color: #cbd5e1;
      font-size: 13px;
      cursor: pointer;
      transition: background 0.1s, color 0.1s;
    }
    .tree-file:hover { background: #1e293b; color: #fff; }
    .tree-file.active { background: #1e293b; color: #60a5fa; }
    .tree-file .icon { opacity: 0.5; font-size: 12px; }

    #proto-content {
      flex: 1;
      overflow: auto;
      background: #0f172a;
      padding: 20px 24px;
    }
    #proto-content .proto-path {
      font-size: 12px;
      color: #64748b;
      margin-bottom: 12px;
      font-family: monospace;
    }
    #proto-content pre {
      font-family: "JetBrains Mono", "Fira Code", monospace;
      font-size: 14px;
      line-height: 1.6;
      color: #e2e8f0;
      white-space: pre;
      tab-size: 2;
    }
    #proto-content .empty {
      color: #475569;
      font-size: 14px;
      margin-top: 40px;
      text-align: center;
    }
  </style>
</head>
<body>
  <nav class="nav">
    <div class="nav-group">
      <span class="nav-label">API</span>
%s    </div>
    <div class="nav-sep"></div>
    <div class="nav-group">
      <span class="nav-label">Proto</span>
      <a class="nav-link" id="proto-nav-btn">Proto Browser</a>
    </div>
  </nav>

  <div id="content">
    <rapi-doc id="api-doc"
      spec-url="/spec/%s"
      theme="dark"
      render-style="read"
      show-header="false"
      allow-try="true"
    ></rapi-doc>
    <div id="proto-browser">
      <div id="proto-sidebar"></div>
      <div id="proto-content">
        <div class="empty">Select a .proto file from the tree</div>
      </div>
    </div>
  </div>

  <script>
    const apiDoc = document.getElementById('api-doc');
    const protoBrowser = document.getElementById('proto-browser');
    const protoSidebar = document.getElementById('proto-sidebar');
    const protoContent = document.getElementById('proto-content');
    const protoNavBtn = document.getElementById('proto-nav-btn');
    const navLinks = document.querySelectorAll('.nav-link[data-spec]');

    const protoFiles = %s;

    // --- File tree builder ---
    function buildTree(paths) {
      const root = {};
      paths.forEach(p => {
        const parts = p.split('/');
        let node = root;
        parts.forEach((part, i) => {
          if (i === parts.length - 1) {
            if (!node._files) node._files = [];
            node._files.push({ name: part, path: p });
          } else {
            if (!node[part]) node[part] = {};
            node = node[part];
          }
        });
      });
      return root;
    }

    function renderTree(node, depth) {
      let html = '';
      const dirs = Object.keys(node).filter(k => k !== '_files').sort();
      dirs.forEach(dir => {
        html += '<div class="tree-folder" style="--depth:' + depth + '">';
        html += '<div class="tree-folder-label" style="--depth:' + depth + '"><span class="arrow">&#9660;</span>' + dir + '</div>';
        html += '<div class="tree-children">' + renderTree(node[dir], depth + 1) + '</div>';
        html += '</div>';
      });
      if (node._files) {
        node._files.sort((a, b) => a.name.localeCompare(b.name));
        node._files.forEach(f => {
          html += '<div class="tree-file" style="--depth:' + (depth) + '" data-path="' + f.path + '" data-url="/proto/' + f.path + '">';
          html += '<span class="icon">&#9679;</span>' + f.name;
          html += '</div>';
        });
      }
      return html;
    }

    const tree = buildTree(protoFiles);
    protoSidebar.innerHTML = renderTree(tree, 0);

    // --- Folder toggle ---
    protoSidebar.addEventListener('click', e => {
      const label = e.target.closest('.tree-folder-label');
      if (label) {
        label.parentElement.classList.toggle('collapsed');
        return;
      }
      const file = e.target.closest('.tree-file');
      if (file) {
        loadProto(file.dataset.path, file.dataset.url);
      }
    });

    function loadProto(path, url) {
      document.querySelectorAll('.tree-file').forEach(f => f.classList.remove('active'));
      const el = protoSidebar.querySelector('[data-path="' + path + '"]');
      if (el) el.classList.add('active');

      fetch(url)
        .then(r => r.text())
        .then(text => {
          protoContent.innerHTML =
            '<div class="proto-path">' + path + '</div>' +
            '<pre>' + escapeHtml(text) + '</pre>';
        });

      history.pushState(null, '', '?proto=' + encodeURIComponent(path));
      setActiveNav(null);
      protoNavBtn.classList.add('active');
    }

    function escapeHtml(t) {
      return t.replace(/&/g,'&amp;').replace(/</g,'&lt;').replace(/>/g,'&gt;');
    }

    // --- Navigation ---
    function showSpec(slug) {
      const link = document.querySelector('[data-spec="' + slug + '"]');
      if (!link) return;
      apiDoc.setAttribute('spec-url', link.dataset.url);
      apiDoc.style.display = '';
      protoBrowser.style.display = 'none';
      setActiveNav(link);
    }

    function showProtoBrowser(filePath) {
      apiDoc.style.display = 'none';
      protoBrowser.style.display = 'flex';
      setActiveNav(null);
      protoNavBtn.classList.add('active');
      if (filePath) {
        const el = protoSidebar.querySelector('[data-path="' + filePath + '"]');
        if (el) loadProto(filePath, el.dataset.url);
      }
    }

    function setActiveNav(activeLink) {
      document.querySelectorAll('.nav-link').forEach(l => l.classList.remove('active'));
      if (activeLink) activeLink.classList.add('active');
    }

    // Nav click handlers
    navLinks.forEach(link => {
      link.addEventListener('click', e => {
        e.preventDefault();
        history.pushState(null, '', '?spec=' + link.dataset.spec);
        route();
      });
    });

    protoNavBtn.addEventListener('click', e => {
      e.preventDefault();
      history.pushState(null, '', '?proto=');
      route();
    });

    // --- Routing ---
    function route() {
      const params = new URLSearchParams(location.search);
      const spec = params.get('spec');
      const proto = params.get('proto');
      if (proto !== null) {
        showProtoBrowser(proto || '');
      } else if (spec) {
        showSpec(spec);
      } else {
        const first = document.querySelector('[data-spec]');
        if (first) showSpec(first.dataset.spec);
      }
    }

    window.addEventListener('popstate', route);
    route();
  </script>
</body>
</html>`, specLinks.String(), specs[0], protoJSON)
}

func (s *Server) stopSwagger(ctx context.Context, log *slog.Logger) error {
	if s.swaggerSrv != nil {
		log.Info("Swagger сервер завершает работу...")
		return s.swaggerSrv.Shutdown(ctx)
	}
	return nil
}
