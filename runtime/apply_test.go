package runtime_test

import (
	"errors"
	"testing"

	"github.com/vandordev/vxt/runtime"
	"github.com/vandordev/vxt/write"
)

type recordingExecutor struct {
	calls []runtime.PlannedHook
	ctxs  []runtime.HookContext
}

func (r *recordingExecutor) Execute(ctx runtime.HookContext, hook runtime.PlannedHook) error {
	r.calls = append(r.calls, hook)
	r.ctxs = append(r.ctxs, ctx)
	return nil
}

type selectiveFailExecutor struct {
	failRun string
	calls   []runtime.PlannedHook
}

func (s *selectiveFailExecutor) Execute(_ runtime.HookContext, hook runtime.PlannedHook) error {
	s.calls = append(s.calls, hook)
	if hook.Run == s.failRun {
		return errors.New("hook failed")
	}
	return nil
}

func TestApplyPlanRunsAfterWriteHooks(t *testing.T) {
	plan := runtime.Plan{
		PlannedHooks: []runtime.PlannedHook{
			{Event: "after:write", Run: "echo one"},
		},
	}
	target := write.NewMemoryTarget()
	exec := &recordingExecutor{}

	result := runtime.ApplyPlan(plan, target, exec)

	if result.WriteResult.Err != nil {
		t.Fatalf("unexpected write err: %v", result.WriteResult.Err)
	}
	if len(exec.calls) != 1 {
		t.Fatalf("got %d hook calls", len(exec.calls))
	}
}

func TestApplyPlanWritesBeforeRunningHooksAndPreservesOrder(t *testing.T) {
	plan := runtime.Plan{
		Files: []runtime.FileOutput{{
			Path:    "hello.txt",
			Content: "hello",
			Mode:    "create",
		}},
		PlannedHooks: []runtime.PlannedHook{
			{Event: "after:write", Run: "echo one"},
			{Event: "before:write", Run: "echo skipped"},
			{Event: "after:write", Run: "echo two"},
		},
	}
	target := write.NewMemoryTarget()
	exec := &recordingExecutor{}

	result := runtime.ApplyPlan(plan, target, exec)

	if result.WriteResult.Err != nil {
		t.Fatalf("unexpected write err: %v", result.WriteResult.Err)
	}
	if result.WriteResult.Report.FilesWritten != 1 {
		t.Fatalf("got files written %d", result.WriteResult.Report.FilesWritten)
	}
	if len(exec.calls) != 2 {
		t.Fatalf("got %d hook calls", len(exec.calls))
	}
	if exec.calls[0].Run != "echo one" || exec.calls[1].Run != "echo two" {
		t.Fatalf("unexpected hook order: %#v", exec.calls)
	}
	if len(exec.ctxs) != 2 {
		t.Fatalf("got %d contexts", len(exec.ctxs))
	}
	if exec.ctxs[0].Event != "after:write" || exec.ctxs[1].Event != "after:write" {
		t.Fatalf("unexpected context events: %#v", exec.ctxs)
	}
	if exec.ctxs[0].WriteReport.FilesWritten != 1 {
		t.Fatalf("unexpected write report in hook context: %#v", exec.ctxs[0].WriteReport)
	}
}

func TestApplyPlanSkipsHooksWhenWriteFails(t *testing.T) {
	plan := runtime.Plan{
		Files: []runtime.FileOutput{{
			Path:    "hello.txt",
			Content: "hello",
			Mode:    "unsupported",
		}},
		PlannedHooks: []runtime.PlannedHook{
			{Event: "after:write", Run: "echo one"},
		},
	}
	target := write.NewMemoryTarget()
	exec := &recordingExecutor{}

	result := runtime.ApplyPlan(plan, target, exec)

	if result.WriteResult.Err == nil {
		t.Fatal("expected write error")
	}
	if len(exec.calls) != 0 {
		t.Fatalf("expected no hook calls, got %d", len(exec.calls))
	}
	if len(result.HookErrors) != 0 {
		t.Fatalf("expected no hook errors, got %d", len(result.HookErrors))
	}
}

func TestApplyPlanCollectsHookErrorsAndContinues(t *testing.T) {
	plan := runtime.Plan{
		PlannedHooks: []runtime.PlannedHook{
			{Event: "after:write", Run: "echo one"},
			{Event: "after:write", Run: "echo two"},
		},
	}
	target := write.NewMemoryTarget()
	exec := &selectiveFailExecutor{failRun: "echo one"}

	result := runtime.ApplyPlan(plan, target, exec)

	if result.WriteResult.Err != nil {
		t.Fatalf("unexpected write err: %v", result.WriteResult.Err)
	}
	if len(exec.calls) != 2 {
		t.Fatalf("expected both hooks to run, got %d", len(exec.calls))
	}
	if len(result.HookErrors) != 1 {
		t.Fatalf("expected 1 hook error, got %d", len(result.HookErrors))
	}
}
