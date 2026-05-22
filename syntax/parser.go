package syntax

import (
	"errors"
	"fmt"
	"strings"

	"github.com/alfariiizi/vxt/model"
	"github.com/alfariiizi/vxt/source"
)

var errUnexpectedEOF = errors.New("unterminated template expression")

type Part struct {
	Text string
	Expr string
}

func ParseTemplate(src source.Source) ([]Part, error) {
	tokens, err := lexTemplate(src.Text)
	if err != nil {
		return nil, err
	}

	parts := make([]Part, 0, len(tokens))
	for _, tok := range tokens {
		switch tok.kind {
		case tokenText:
			parts = append(parts, Part{Text: tok.value})
		case tokenExpr:
			parts = append(parts, Part{Expr: tok.value})
		}
	}

	return parts, nil
}

func ParseDocument(src source.Source) (*model.CompiledDocument, error) {
	lines := strings.Split(src.Text, "\n")
	doc := &model.CompiledDocument{Source: src}

	for i := 0; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])
		if line == "" {
			continue
		}

		switch {
		case strings.HasPrefix(line, "@template "):
			doc.Template = strings.TrimSpace(strings.TrimPrefix(line, "@template "))
		case strings.HasPrefix(line, "@input "):
			parts := strings.Fields(strings.TrimPrefix(line, "@input "))
			if len(parts) != 2 {
				return nil, fmt.Errorf("invalid input declaration")
			}
			doc.Inputs = append(doc.Inputs, model.InputDecl{
				Name:     parts[0],
				TypeName: parts[1],
			})
		case strings.HasPrefix(line, "@hook "):
			payload := strings.TrimSpace(strings.TrimPrefix(line, "@hook "))
			event, run, ok := strings.Cut(payload, " ")
			if !ok {
				return nil, fmt.Errorf("invalid hook declaration")
			}
			doc.Hooks = append(doc.Hooks, model.HookDecl{
				Event: event,
				Run:   strings.Trim(run, `"`),
			})
		case strings.HasPrefix(line, "@file "):
			path := strings.Trim(strings.TrimSpace(strings.TrimPrefix(line, "@file ")), `"`)
			var bodyLines []string
			foundEnd := false
			for j := i + 1; j < len(lines); j++ {
				if strings.TrimSpace(lines[j]) == "@endfile" {
					i = j
					foundEnd = true
					break
				}
				bodyLines = append(bodyLines, lines[j])
			}
			if !foundEnd {
				return nil, errUnexpectedEOF
			}
			doc.Files = append(doc.Files, model.FileBlock{
				Path: path,
				Body: strings.Join(bodyLines, "\n"),
				Mode: "create",
			})
		}
	}

	if doc.Template == "" {
		return nil, fmt.Errorf("missing @template")
	}

	return doc, nil
}
