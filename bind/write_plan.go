package bind

import (
	"os"
	"path/filepath"
)

type plannedWriteAction struct {
	Path    string
	Action  WriteActionKind
	Content string
}

func planWriteActions(out Output, dir string) ([]plannedWriteAction, error) {
	active := make(map[string]struct{}, len(out.Files))
	actions := make([]plannedWriteAction, 0, len(out.Files))

	for _, file := range out.Files {
		target := filepath.Join(dir, filepath.Base(file.Path))
		active[filepath.Base(file.Path)] = struct{}{}

		action := WriteActionCreate
		if _, err := os.Stat(target); err == nil {
			action = WriteActionOverwrite
		} else if !os.IsNotExist(err) {
			return nil, err
		}

		actions = append(actions, plannedWriteAction{
			Path:    target,
			Action:  action,
			Content: file.Content,
		})
	}

	if entries, err := os.ReadDir(dir); err == nil {
		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}
			name := entry.Name()
			if _, ok := active[name]; ok {
				continue
			}

			target := filepath.Join(dir, name)
			content, err := os.ReadFile(target)
			if err != nil {
				return nil, err
			}
			if !isOwnedBindingFile(out.BindingName, out.PackageName, target, string(content)) {
				continue
			}

			actions = append(actions, plannedWriteAction{
				Path:   target,
				Action: WriteActionRemove,
			})
		}
	} else if !os.IsNotExist(err) {
		return nil, err
	}

	return actions, nil
}
