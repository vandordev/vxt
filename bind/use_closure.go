package bind

import (
	"fmt"
	"strings"

	"github.com/vandordev/vxt/runtime"
	"github.com/vandordev/vxt/source"
)

func collectLocalUseSources(doc *runtime.CompiledDocument, uses map[string]source.Source) (map[string]source.Source, error) {
	if doc == nil || len(doc.Uses) == 0 {
		return nil, nil
	}

	out := map[string]source.Source{}
	for _, use := range doc.Uses {
		if !isLocalUsePath(use.Path) {
			continue
		}
		src, ok := uses[use.Path]
		if !ok {
			return nil, fmt.Errorf("bind: missing local use source for %q", use.Path)
		}
		out[use.Path] = src
	}
	if len(out) == 0 {
		return nil, nil
	}
	return out, nil
}

func isLocalUsePath(path string) bool {
	return strings.HasPrefix(path, "./") || strings.HasPrefix(path, "../")
}
