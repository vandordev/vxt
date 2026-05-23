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

func (m *MemoryTarget) WriteFile(path string, content []byte) error {
	copied := make([]byte, len(content))
	copy(copied, content)
	m.files[path] = copied
	return nil
}

func (m *MemoryTarget) Files() map[string][]byte {
	return m.files
}
