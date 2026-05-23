package plan

type DirOutput struct {
	Path string
}

type FileOutput struct {
	Path    string
	Content string
	Mode    string
}

type PlannedHook struct {
	Event string
	Run   string
}

type Plan struct {
	Dirs         []DirOutput
	Files        []FileOutput
	PlannedHooks []PlannedHook
}
