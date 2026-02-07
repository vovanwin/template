package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"unicode"
)

type rpcMethod struct {
	Name     string
	Request  string
	Response string
}

type service struct {
	Name       string
	Methods    []rpcMethod
	GoPackage  string // full import path, e.g. "github.com/vovanwin/template/pkg/template"
	PbAlias    string // package alias, e.g. "template"
	StructName string // e.g. "TemplateGRPCServer"
	DirName    string // e.g. "template"
}

func main() {
	root, err := os.Getwd()
	if err != nil {
		fatal("getwd: %v", err)
	}

	goModule := parseGoModule(filepath.Join(root, "go.mod"))
	if goModule == "" {
		fatal("could not parse module path from go.mod")
	}

	protoFiles, err := filepath.Glob(filepath.Join(root, "api", "*", "*.proto"))
	if err != nil {
		fatal("glob: %v", err)
	}

	if len(protoFiles) == 0 {
		fmt.Println("no .proto files found in api/")
		return
	}

	for _, pf := range protoFiles {
		services, err := parseProto(pf)
		if err != nil {
			fatal("parse %s: %v", pf, err)
		}

		for _, svc := range services {
			controllerDir := filepath.Join(root, "internal", "controller", svc.DirName)
			if err := os.MkdirAll(controllerDir, 0o755); err != nil {
				fatal("mkdir %s: %v", controllerDir, err)
			}

			// generate controller.go if missing
			controllerFile := filepath.Join(controllerDir, "controller.go")
			if _, err := os.Stat(controllerFile); os.IsNotExist(err) {
				content := genController(svc)
				if err := os.WriteFile(controllerFile, []byte(content), 0o644); err != nil {
					fatal("write %s: %v", controllerFile, err)
				}
				fmt.Printf("created %s\n", controllerFile)
			}

			// generate module.go if missing
			moduleFile := filepath.Join(controllerDir, "module.go")
			if _, err := os.Stat(moduleFile); os.IsNotExist(err) {
				content := genModule(svc, goModule)
				if err := os.WriteFile(moduleFile, []byte(content), 0o644); err != nil {
					fatal("write %s: %v", moduleFile, err)
				}
				fmt.Printf("created %s\n", moduleFile)
			}

			// generate method stubs for missing methods
			for _, m := range svc.Methods {
				fileName := camelToSnake(m.Name) + ".go"
				filePath := filepath.Join(controllerDir, fileName)
				if _, err := os.Stat(filePath); os.IsNotExist(err) {
					content := genMethod(svc, m)
					if err := os.WriteFile(filePath, []byte(content), 0o644); err != nil {
						fatal("write %s: %v", filePath, err)
					}
					fmt.Printf("created %s\n", filePath)
				}
			}
		}
	}

	fmt.Println("done")
}

var (
	reGoPackage = regexp.MustCompile(`option\s+go_package\s*=\s*"([^"]+)"`)
	reService   = regexp.MustCompile(`service\s+(\w+)\s*\{`)
	reRPC       = regexp.MustCompile(`rpc\s+(\w+)\s*\(\s*(\w+)\s*\)\s*returns\s*\(\s*(\w+)\s*\)`)
	reModule    = regexp.MustCompile(`(?m)^module\s+(\S+)`)
)

func parseGoModule(goModPath string) string {
	data, err := os.ReadFile(goModPath)
	if err != nil {
		return ""
	}
	m := reModule.FindStringSubmatch(string(data))
	if m == nil {
		return ""
	}
	return m[1]
}

func parseProto(path string) ([]service, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	content := string(data)

	// parse go_package
	goPackage, pbAlias := parseGoPackage(content)
	if goPackage == "" {
		return nil, fmt.Errorf("go_package option not found in %s", path)
	}

	// parse services and their RPCs
	var services []service

	svcMatches := reService.FindAllStringIndex(content, -1)
	for i, loc := range svcMatches {
		svcName := reService.FindStringSubmatch(content[loc[0]:loc[1]])[1]

		// determine the block for this service (up to next service or EOF)
		start := loc[1]
		end := len(content)
		if i+1 < len(svcMatches) {
			end = svcMatches[i+1][0]
		}
		block := content[start:end]

		var methods []rpcMethod
		for _, m := range reRPC.FindAllStringSubmatch(block, -1) {
			methods = append(methods, rpcMethod{
				Name:     m[1],
				Request:  m[2],
				Response: m[3],
			})
		}

		dirName := strings.TrimSuffix(svcName, "Service")
		dirName = strings.ToLower(dirName)

		structName := strings.TrimSuffix(svcName, "Service") + "GRPCServer"

		services = append(services, service{
			Name:       svcName,
			Methods:    methods,
			GoPackage:  goPackage,
			PbAlias:    pbAlias,
			StructName: structName,
			DirName:    dirName,
		})
	}

	return services, nil
}

func parseGoPackage(content string) (goPackage, alias string) {
	m := reGoPackage.FindStringSubmatch(content)
	if m == nil {
		return "", ""
	}
	raw := m[1]
	if idx := strings.LastIndex(raw, ";"); idx != -1 {
		return raw[:idx], raw[idx+1:]
	}
	// no alias — use last path segment
	parts := strings.Split(raw, "/")
	return raw, parts[len(parts)-1]
}

func genController(svc service) string {
	return fmt.Sprintf(`package %s

import (
	"log/slog"

	%spb "%s"
	"go.uber.org/fx"
)

// Deps содержит зависимости для %s.
type Deps struct {
	fx.In

	Log *slog.Logger
}

// %s реализует gRPC сервис %s.
type %s struct {
	%spb.Unimplemented%sServer
	log *slog.Logger
}

// New%s создаёт новый %s.
func New%s(deps Deps) *%s {
	return &%s{log: deps.Log}
}
`, svc.DirName,
		svc.PbAlias, svc.GoPackage,
		svc.StructName,
		svc.StructName, svc.Name,
		svc.StructName,
		svc.PbAlias, svc.Name,
		svc.StructName, svc.StructName,
		svc.StructName, svc.StructName,
		svc.StructName)
}

func genModule(svc service, goModule string) string {
	serverImport := goModule + "/internal/pkg/server"

	return fmt.Sprintf(`package %s

import (
	"context"

	"%s"
	%spb "%s"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/fx"
	"google.golang.org/grpc"
)

// Module возвращает fx.Option для подключения %s.
func Module() fx.Option {
	return fx.Options(
		fx.Provide(New%s),
		fx.Provide(
			fx.Annotate(
				func(srv *%s) server.GRPCRegistrator {
					return func(s *grpc.Server) {
						%spb.Register%sServer(s, srv)
					}
				},
				fx.ResultTags(`+"`"+`group:"grpc_registrators"`+"`"+`),
			),
		),
		fx.Provide(
			fx.Annotate(
				func(srv *%s) server.GatewayRegistrator {
					return func(ctx context.Context, mux *runtime.ServeMux, _ *grpc.Server) error {
						return %spb.Register%sHandlerServer(ctx, mux, srv)
					}
				},
				fx.ResultTags(`+"`"+`group:"gateway_registrators"`+"`"+`),
			),
		),
	)
}
`, svc.DirName,
		serverImport,
		svc.PbAlias, svc.GoPackage,
		svc.Name,
		svc.StructName,
		svc.StructName,
		svc.PbAlias, svc.Name,
		svc.StructName,
		svc.PbAlias, svc.Name)
}

func genMethod(svc service, m rpcMethod) string {
	return fmt.Sprintf(`package %s

import (
	"context"

	%spb "%s"
)

func (s *%s) %s(_ context.Context, req *%spb.%s) (*%spb.%s, error) {
	// TODO: implement
	panic("not implemented")
}
`, svc.DirName,
		svc.PbAlias, svc.GoPackage,
		svc.StructName, m.Name, svc.PbAlias, m.Request, svc.PbAlias, m.Response)
}

func camelToSnake(s string) string {
	var result strings.Builder
	for i, r := range s {
		if unicode.IsUpper(r) {
			if i > 0 {
				result.WriteByte('_')
			}
			result.WriteRune(unicode.ToLower(r))
		} else {
			result.WriteRune(r)
		}
	}
	return result.String()
}

func fatal(format string, args ...any) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
