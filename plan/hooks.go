package plan

func HooksFromDecls(event string, run string) PlannedHook {
	return PlannedHook{
		Event: event,
		Run:   run,
	}
}
