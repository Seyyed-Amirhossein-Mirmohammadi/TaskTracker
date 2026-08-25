package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"task-cli/internal/task"
	"testing"
)

func setupTestCLI(t *testing.T) (string, func()) {
	t.Helper()
	tmpDir, err := os.MkdirTemp("", "cli_test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	absDir, err := filepath.Abs(tmpDir)
	if err != nil {
		t.Fatalf("failed to get absolute path: %v", err)
	}
	filePath := filepath.Join(absDir, "tasks.json")

	s, err := task.NewTaskStore(filePath)
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}
	task.SetTestStore(s)

	cleanup := func() {
		os.RemoveAll(absDir)
	}
	return absDir, cleanup
}

func captureOutput(f func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	f()
	w.Close()
	os.Stdout = old
	var buf bytes.Buffer
	buf.ReadFrom(r)
	return buf.String()
}

func contains(str, substr string) bool {
	return bytes.Contains([]byte(str), []byte(substr))
}

func TestRun_AddCommand(t *testing.T) {
	_, cleanup := setupTestCLI(t)
	defer cleanup()

	output := captureOutput(func() {
		err := Run([]string{"add", "test task"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	if !contains(output, "Task added successfully (ID: 1)") {
		t.Errorf("expected success message, got %s", output)
	}
}

func TestRun_UpdateCommand(t *testing.T) {
	_, cleanup := setupTestCLI(t)
	defer cleanup()

	Run([]string{"add", "original"})

	output := captureOutput(func() {
		err := Run([]string{"update", "1", "updated"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	if !contains(output, "Task 1 updated successfully") {
		t.Errorf("expected update message, got %s", output)
	}
}

func TestRun_DeleteCommand(t *testing.T) {
	_, cleanup := setupTestCLI(t)
	defer cleanup()

	Run([]string{"add", "to delete"})

	output := captureOutput(func() {
		err := Run([]string{"delete", "1"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	if !contains(output, "Task 1 deleted successfully") {
		t.Errorf("expected delete message, got %s", output)
	}
}

func TestRun_ListCommand(t *testing.T) {
	_, cleanup := setupTestCLI(t)
	defer cleanup()

	Run([]string{"add", "task1"})
	Run([]string{"add", "task2"})

	output := captureOutput(func() {
		err := Run([]string{"list"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	if !contains(output, "task1") || !contains(output, "task2") {
		t.Errorf("expected both tasks in output, got %s", output)
	}
}

func TestRun_ListFiltered(t *testing.T) {
	_, cleanup := setupTestCLI(t)
	defer cleanup()

	Run([]string{"add", "todo task"})
	Run([]string{"add", "inprogress"})
	Run([]string{"mark-in-progress", "2"})

	output := captureOutput(func() {
		err := Run([]string{"list", "in-progress"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	if !contains(output, "inprogress") {
		t.Errorf("expected in-progress task, got %s", output)
	}
	if contains(output, "todo task") {
		t.Error("should not see todo task")
	}
}

func TestRun_MarkInProgress(t *testing.T) {
	_, cleanup := setupTestCLI(t)
	defer cleanup()

	Run([]string{"add", "task"})

	output := captureOutput(func() {
		err := Run([]string{"mark-in-progress", "1"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	if !contains(output, "Task 1 marked as in-progress") {
		t.Errorf("expected mark message, got %s", output)
	}
}

func TestRun_MarkDone(t *testing.T) {
	_, cleanup := setupTestCLI(t)
	defer cleanup()

	Run([]string{"add", "task"})

	output := captureOutput(func() {
		err := Run([]string{"mark-done", "1"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	if !contains(output, "Task 1 marked as done") {
		t.Errorf("expected mark done message, got %s", output)
	}
}

func TestRun_UnknownCommand(t *testing.T) {
	_, cleanup := setupTestCLI(t)
	defer cleanup()

	err := Run([]string{"invalid"})
	if err == nil {
		t.Error("expected error for unknown command")
	}
}

func TestRun_HelpCommand(t *testing.T) {
	_, cleanup := setupTestCLI(t)
	defer cleanup()

	output := captureOutput(func() {
		err := Run([]string{"help"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	if !contains(output, "Usage:") {
		t.Errorf("expected help output, got %s", output)
	}
}

func TestRun_NoArgs(t *testing.T) {
	_, cleanup := setupTestCLI(t)
	defer cleanup()

	output := captureOutput(func() {
		err := Run([]string{})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	if !contains(output, "Usage:") {
		t.Errorf("expected help output, got %s", output)
	}
}
