package main

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

func main() {
	interfaces()
	implementations()
}

func interfaces() {
	tpls, impls := buildFileLists("templates/pkg")

	t := template.New("interfaces")
	t.Funcs(map[string]any{
		"types": func(out, build, cli string) ImplementationFields {
			return ImplementationFields{out, build, cli}
		},
	})

	runTemplates(tpls, impls, t)
}

type ImplementationFields struct {
	OutputType, BuilderType, CLIFunc string
}

func implementations() {
	tpls, impls := buildFileLists("templates/internal")

	t := template.New("implementations")
	t.Funcs(map[string]any{
		"map": func(pairs ...any) (map[string]any, error) {
			if len(pairs)%2 != 0 {
				return nil, errors.New("misaligned map")
			}

			m := make(map[string]any, len(pairs)/2)

			for i := 0; i < len(pairs); i += 2 {
				key, ok := pairs[i].(string)

				if !ok {
					return nil, fmt.Errorf("cannot use type %T as map key", pairs[i])
				}

				m[key] = pairs[i+1]
			}

			return m, nil
		},
	})

	runTemplates(tpls, impls, t)
}

func buildFileLists(root string) ([]string, []string) {
	tpls := make([]string, 0, 10)
	impls := make([]string, 0, 10)

	filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		try(err)

		if d.IsDir() {
			return nil
		}

		if strings.HasPrefix(d.Name(), "_") {
			tpls = append(tpls, path)
		} else {
			impls = append(impls, path)
		}

		return nil
	})

	return tpls, impls
}

func runTemplates(tpls, impls []string, t *template.Template) {
	_ = tryGet(t.ParseFiles(append(tpls, impls...)...))

	for _, impl := range impls {
		executeTemplate(t, impl[10:])
	}
}

func executeTemplate(tpl *template.Template, path string) {
	try(os.MkdirAll(filepath.Dir(path), 0755))

	file := tryGet(os.OpenFile(strings.TrimSuffix(path, ".tpl"), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644))
	defer file.Close()

	try(tpl.ExecuteTemplate(file, filepath.Base(path), nil))
}

func try(err error) {
	if err != nil {
		panic(err)
	}
}

func tryGet[T any](value T, err error) T {
	if err != nil {
		panic(err)
	}

	return value
}
