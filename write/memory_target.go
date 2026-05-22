package write

type MemoryTarget struct {
	files map[string][]byte
}

func NewMemoryTarget() *MemoryTarget {
	return &MemoryTarget{
		files: map[string][]byte{},
	}
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
