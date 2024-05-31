package parsing

import (
	"github.com/floholz/mochatex/internal/job"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"text/template/parse"
)

func ParseTexFile(path *string, errLog, infoLog *log.Logger) *template.Template {
	if filepath.Ext(*path) != ".tex" {
		errLog.Fatalf("%s must be a valid .tex file", *path)
	}
	_, err := os.Stat(*path)
	if err != nil {
		errLog.Fatalf("error while reading info for %s: %v", *path, err)
	}

	tmpl, err := template.
		New(filepath.Base(*path)).
		Delims(job.DefaultDelimiters.Left, job.DefaultDelimiters.Right).
		ParseFiles(*path)
	if err != nil {
		errLog.Fatalf("error while parsing template %s: %v", *path, err)
	}

	return tmpl
}

func MapTemplateFields(t *template.Template) map[string]int {
	return mapNodeFields(t.Tree.Root, make(map[string]int))
}

func mapNodeFields(node parse.Node, res map[string]int) map[string]int {
	if node.Type() == parse.NodeAction {
		field := strings.TrimPrefix(node.String(), "{{")
		field = strings.TrimSuffix(field, "}}")

		cnt, ok := res[field]
		if ok {
			res[field] = cnt + 1
		} else {
			res[field] = 1
		}
	}

	if ln, ok := node.(*parse.ListNode); ok {
		for _, n := range ln.Nodes {
			res = mapNodeFields(n, res)
		}
	}
	return res
}
