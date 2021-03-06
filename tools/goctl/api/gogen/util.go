package gogen

import (
	"fmt"
	goformat "go/format"
	"io"
	"path/filepath"
	"strings"

	"zero/core/collection"
	"zero/tools/goctl/api/spec"
	"zero/tools/goctl/api/util"
	"zero/tools/goctl/vars"
)

func getParentPackage(dir string) (string, error) {
	absDir, err := filepath.Abs(dir)
	if err != nil {
		return "", err
	}
	pos := strings.Index(absDir, vars.ProjectName)
	if pos < 0 {
		return "", fmt.Errorf("%s not in project directory", dir)
	}

	return absDir[pos:], nil
}

func writeIndent(writer io.Writer, indent int) {
	for i := 0; i < indent; i++ {
		fmt.Fprint(writer, "\t")
	}
}

func writeProperty(writer io.Writer, name, tp, tag, comment string, indent int) error {
	writeIndent(writer, indent)
	var err error
	if len(comment) > 0 {
		comment = strings.TrimPrefix(comment, "//")
		comment = "//" + comment
		_, err = fmt.Fprintf(writer, "%s %s %s %s\n", strings.Title(name), tp, tag, comment)
	} else {
		_, err = fmt.Fprintf(writer, "%s %s %s\n", strings.Title(name), tp, tag)
	}
	return err
}

func getAuths(api *spec.ApiSpec) []string {
	var authNames = collection.NewSet()
	for _, g := range api.Service.Groups {
		if value, ok := util.GetAnnotationValue(g.Annotations, "server", "jwt"); ok {
			authNames.Add(value)
		}
		if value, ok := util.GetAnnotationValue(g.Annotations, "server", "signature"); ok {
			authNames.Add(value)
		}
	}
	return authNames.KeysStr()
}

func formatCode(code string) string {
	ret, err := goformat.Source([]byte(code))
	if err != nil {
		return code
	}
	return string(ret)
}
