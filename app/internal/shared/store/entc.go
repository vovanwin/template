//go:build ignore

package main

import (
	"log"

	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
)

func main() {
	opts := []entc.Option{
		entc.TemplateDir("./templates"),
	}
	config := &gen.Config{
		Package: "github.com/vovanwin/template/internal/shared/store/gen",
		Features: []gen.Feature{
			gen.FeatureVersionedMigration,
			gen.FeatureUpsert,
			gen.FeatureIntercept,
			gen.FeatureNamedEdges,
		},
		Target: "gen",
		Schema: "schema",
	}
	if err := entc.Generate("./schema", config, opts...); err != nil {
		log.Fatalf("running ent codegen: %v", err)
	}
}
