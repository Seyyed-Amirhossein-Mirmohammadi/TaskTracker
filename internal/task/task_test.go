package task

import (
	"os"
	"path/filepath"
	"testing"
)

func setupTestTask(t *testing.T) (*TaskStore, string) {
	t.Helper()
	tmpDir, err := os.MkdirTemp("", "task_test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	absDir, err := filepath.Abs(tmpDir)
	if err != nil {
		t.Fatalf("failed to get absolute path: %v", err)
	}
	filePath := filepath.Join(absDir, "tasks.json")

	s, err := NewTaskStore(filePath)
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}
	store = s
	return s, absDir
}

func cleanupTestTask(tmpDir string) {
	os.RemoveAll(tmpDir)
}

func TestAddTask(t *testing.T) {
	s, tmpDir := setupTestTask(t)
	defer cleanupTestTask(tmpDir)

	s.Data.Tasks = []Task{}
	s.Data.LastID = 0

	id, err := AddTask("Buy milk")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if id != 1 {
		t.Errorf("expected ID 1, got %d", id)
	}
	if len(s.Data.Tasks) != 1 {
		t.Errorf("expected 1 task, got %d", len(s.Data.Tasks))
	}
	if s.Data.Tasks[0].Description != "Buy milk" {
		t.Errorf("expected 'Buy milk', got '%s'", s.Data.Tasks[0].Description)
	}
	if s.Data.Tasks[0].Status != Todo {
		t.Errorf("expected status Todo, got %v", s.Data.Tasks[0].Status)
	}
	if s.Data.Tasks[0].UpdatedAt != nil {
		t.Error("expected UpdatedAt to be nil")
	}
}

func TestAddTask_EmptyDescription(t *testing.T) {
	s, tmpDir := setupTestTask(t)
	defer cleanupTestTask(tmpDir)

	s.Data.Tasks = []Task{}
	s.Data.LastID = 0

	_, err := AddTask("")
	if err == nil {
		t.Error("expected error for empty description")
	}
}

func TestUpdateTaskTitle(t *testing.T) {
	s, tmpDir := setupTestTask(t)
	defer cleanupTestTask(tmpDir)

	s.Data.Tasks = []Task{}
	s.Data.LastID = 0
	AddTask("old")

	err := UpdateTaskTitle(1, "new")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.Data.Tasks[0].Description != "new" {
		t.Errorf("expected description 'new', got '%s'", s.Data.Tasks[0].Description)
	}
	if s.Data.Tasks[0].UpdatedAt == nil {
		t.Error("expected UpdatedAt to be set")
	}
}

func TestUpdateTaskTitle_NotFound(t *testing.T) {
	_, tmpDir := setupTestTask(t)
	defer cleanupTestTask(tmpDir)

	err := UpdateTaskTitle(99, "ignored")
	if err == nil {
		t.Error("expected error for non-existent ID")
	}
}

func TestUpdateTaskTitle_EmptyDescription(t *testing.T) {
	s, tmpDir := setupTestTask(t)
	defer cleanupTestTask(tmpDir)

	s.Data.Tasks = []Task{}
	s.Data.LastID = 0
	AddTask("test")

	err := UpdateTaskTitle(1, "")
	if err == nil {
		t.Error("expected error for empty description")
	}
}

func TestDeleteTask(t *testing.T) {
	s, tmpDir := setupTestTask(t)
	defer cleanupTestTask(tmpDir)

	s.Data.Tasks = []Task{}
	s.Data.LastID = 0
	AddTask("delete me")

	err := DeleteTask(1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(s.Data.Tasks) != 0 {
		t.Errorf("expected 0 tasks, got %d", len(s.Data.Tasks))
	}
}

func TestDeleteTask_NotFound(t *testing.T) {
	_, tmpDir := setupTestTask(t)
	defer cleanupTestTask(tmpDir)

	err := DeleteTask(99)
	if err == nil {
		t.Error("expected error for non-existent ID")
	}
}

func TestChangeState(t *testing.T) {
	s, tmpDir := setupTestTask(t)
	defer cleanupTestTask(tmpDir)

	s.Data.Tasks = []Task{}
	s.Data.LastID = 0
	AddTask("task")

	err := ChangeState(1, InProgress)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.Data.Tasks[0].Status != InProgress {
		t.Errorf("expected status InProgress, got %v", s.Data.Tasks[0].Status)
	}
	if s.Data.Tasks[0].UpdatedAt == nil {
		t.Error("expected UpdatedAt to be set")
	}
}

func TestChangeState_NotFound(t *testing.T) {
	_, tmpDir := setupTestTask(t)
	defer cleanupTestTask(tmpDir)

	err := ChangeState(99, Done)
	if err == nil {
		t.Error("expected error for non-existent ID")
	}
}

func TestListAll_Task(t *testing.T) {
	s, tmpDir := setupTestTask(t)
	defer cleanupTestTask(tmpDir)

	s.Data.Tasks = []Task{}
	s.Data.LastID = 0
	AddTask("a")
	AddTask("b")

	tasks, err := ListAll()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(tasks) != 2 {
		t.Errorf("expected 2 tasks, got %d", len(tasks))
	}
}

func TestList_Task(t *testing.T) {
	s, tmpDir := setupTestTask(t)
	defer cleanupTestTask(tmpDir)

	s.Data.Tasks = []Task{}
	s.Data.LastID = 0
	AddTask("a")
	AddTask("b")
	ChangeState(2, Done)

	tasks, err := List(Done)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(tasks) != 1 {
		t.Errorf("expected 1 done task, got %d", len(tasks))
	}
	if tasks[0].ID != 2 {
		t.Errorf("expected ID 2, got %d", tasks[0].ID)
	}
}
