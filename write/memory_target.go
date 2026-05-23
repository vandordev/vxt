package write

type MemoryTarget struct {
	files map[string][]byte
	dirs  map[string]struct{}
}

func NewMemoryTarget() *MemoryTarget {
	return &MemoryTarget{
		files: map[string][]byte{},
		dirs:  map[string]struct{}{},
	}
}

func (m *MemoryTarget) MkdirAll(path string) error {
	m.dirs[path] = struct{}{}
	return nil
}

func (m *MemoryTarget) WriteFile(path string, content []byte, mode string) (bool, error) {
	switch mode {
	case "", "overwrite":
	case "skip-if-exists":
		if _, exists := m.files[path]; exists {
			return false, nil
		}
	case "create":
		if _, exists := m.files[path]; exists {
			return false, ErrFileExists{Path: path}
		}
	default:
		return false, ErrUnsupportedWriteMode{Mode: mode}
	}

	copied := make([]byte, len(content))
	copy(copied, content)
	m.files[path] = copied
	return true, nil
}

func (m *MemoryTarget) Files() map[string][]byte {
	return m.files
}
