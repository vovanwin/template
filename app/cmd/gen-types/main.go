package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var errEmpty = errors.New("empty")

func main() {
	if len(os.Args) != 4 {
		log.Fatalf("invalid args count: %d", len(os.Args)-1)
	}

	pkg, types, out := os.Args[1], strings.Split(os.Args[2], ","), os.Args[3]
	if err := run(pkg, types, out); err != nil {
		log.Fatal(err)
	}

	p, _ := os.Getwd()
	fmt.Printf("%v generated\n", filepath.Join(p, out)) //nolint:forbidigo // It has to report about results
}

func run(pkg string, types []string, outFile string) error {
	if len(pkg) == 0 {
		return fmt.Errorf("package: %v", errEmpty)
	}
	if len(types) == 0 {
		return fmt.Errorf("types: %v", errEmpty)
	}
	if len(outFile) == 0 {
		return fmt.Errorf("output filename: %v", errEmpty)
	}
	w, err := os.Create(outFile)
	if err != nil {
		return fmt.Errorf("open file %q: %v", outFile, err)
	}
	err = headerTpl.Execute(w, map[string]interface{}{
		"package": pkg,
		"types":   strings.Join(types, " | "),
	})
	if err != nil {
		return fmt.Errorf("exec header tpl: %v", err)
	}
	for _, t := range types {
		err = bodyTpl.Execute(w, t)
		if err != nil {
			return fmt.Errorf("exec body tpl with type %q: %v", t, err)
		}
	}
	return w.Close()
}
