package plan

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
	Files        []FileOutput
	PlannedHooks []PlannedHook
}
