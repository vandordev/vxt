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

type Node interface {
	isNode()
}

type TextNode struct {
	Text string
}

type ExprNode struct {
	Expr string
}

type IfNode struct {
	Cond string
	Then []Node
	Else []Node
}

type EachNode struct {
	Collection string
	Item       string
	Body       []Node
}

func (TextNode) isNode() {}
func (ExprNode) isNode() {}
func (IfNode) isNode()   {}
func (EachNode) isNode() {}

func ParseTemplate(src source.Source) ([]Node, error) {
	tokens, err := lexTemplate(src.Text)
	if err != nil {
		return nil, err
	}

	pos := 0
	nodes, stop, err := parseNodes(tokens, &pos, nil)
	if err != nil {
		return nil, err
	}
	if stop != "" {
		return nil, fmt.Errorf("unexpected control terminator %q", stop)
	}
	return nodes, nil
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

func parseNodes(tokens []token, pos *int, stop map[string]bool) ([]Node, string, error) {
	nodes := make([]Node, 0)

	for *pos < len(tokens) {
		tok := tokens[*pos]
		*pos++

		switch tok.kind {
		case tokenText:
			nodes = append(nodes, TextNode{Text: tok.value})
		case tokenExpr:
			expr := strings.TrimSpace(tok.value)
			if stop != nil && stop[expr] {
				return nodes, expr, nil
			}

			switch {
			case strings.HasPrefix(expr, "if "):
				cond := strings.TrimSpace(strings.TrimPrefix(expr, "if "))
				thenNodes, found, err := parseNodes(tokens, pos, map[string]bool{
					"else":   true,
					"end if": true,
				})
				if err != nil {
					return nil, "", err
				}
				if found == "" {
					return nil, "", fmt.Errorf("missing end if")
				}

				node := IfNode{Cond: cond, Then: thenNodes}
				if found == "else" {
					elseNodes, endTok, err := parseNodes(tokens, pos, map[string]bool{
						"end if": true,
					})
					if err != nil {
						return nil, "", err
					}
					if endTok != "end if" {
						return nil, "", fmt.Errorf("missing end if")
					}
					node.Else = elseNodes
				}
				nodes = append(nodes, node)
			case strings.HasPrefix(expr, "each "):
				body := strings.TrimSpace(strings.TrimPrefix(expr, "each "))
				before, after, ok := strings.Cut(body, " as ")
				if !ok {
					return nil, "", fmt.Errorf("invalid each syntax")
				}
				eachNodes, found, err := parseNodes(tokens, pos, map[string]bool{
					"end each": true,
				})
				if err != nil {
					return nil, "", err
				}
				if found != "end each" {
					return nil, "", fmt.Errorf("missing end each")
				}
				nodes = append(nodes, EachNode{
					Collection: strings.TrimSpace(before),
					Item:       strings.TrimSpace(after),
					Body:       eachNodes,
				})
			case expr == "else", expr == "end if", expr == "end each":
				return nil, "", fmt.Errorf("unexpected control terminator %q", expr)
			default:
				nodes = append(nodes, ExprNode{Expr: expr})
			}
		}
	}

	return nodes, "", nil
}
