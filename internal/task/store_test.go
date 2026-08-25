package task

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func setupTestStore(t *testing.T) (*TaskStore, string) {
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
	return s, absDir
}

func cleanupTestStore(tmpDir string) {
	os.RemoveAll(tmpDir)
}

func TestNewTaskStore_CreatesFileAndDirectory(t *testing.T) {
	s, tmpDir := setupTestStore(t)
	defer cleanupTestStore(tmpDir)

	if _, err := os.Stat(s.filePath); os.IsNotExist(err) {
		t.Error("expected file to be created")
	}
	if len(s.Data.Tasks) != 0 {
		t.Errorf("expected empty tasks, got %d", len(s.Data.Tasks))
	}
	if s.Data.LastID != 0 {
		t.Errorf("expected last_id 0, got %d", s.Data.LastID)
	}
}

func TestStore_LoadAndSave(t *testing.T) {
	s, tmpDir := setupTestStore(t)
	defer cleanupTestStore(tmpDir)

	s.Data.Tasks = []Task{
		{ID: 1, Description: "Test", Status: Todo, CreatedAt: time.Now()},
	}
	s.Data.LastID = 1
	err := s.Save()
	if err != nil {
		t.Fatalf("save failed: %v", err)
	}

	newStore, err := NewTaskStore(s.filePath)
	if err != nil {
		t.Fatalf("new store failed: %v", err)
	}
	if len(newStore.Data.Tasks) != 1 {
		t.Errorf("expected 1 task, got %d", len(newStore.Data.Tasks))
	}
	if newStore.Data.Tasks[0].Description != "Test" {
		t.Errorf("expected 'Test', got '%s'", newStore.Data.Tasks[0].Description)
	}
	if newStore.Data.LastID != 1 {
		t.Errorf("expected last_id 1, got %d", newStore.Data.LastID)
	}
}

func TestGenerateTaskID(t *testing.T) {
	s, tmpDir := setupTestStore(t)
	defer cleanupTestStore(tmpDir)
	s.Data.LastID = 5
	store = s
	id := generateTaskID()
	if id != 6 {
		t.Errorf("expected 6, got %d", id)
	}
	if s.Data.LastID != 6 {
		t.Errorf("expected last_id 6, got %d", s.Data.LastID)
	}
}

func TestLoadTaskByID_Found(t *testing.T) {
	s, tmpDir := setupTestStore(t)
	defer cleanupTestStore(tmpDir)

	s.Data.Tasks = []Task{
		{ID: 1, Description: "task1"},
		{ID: 2, Description: "task2"},
	}
	store = s

	task, err := loadTaskByID(2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if task.Description != "task2" {
		t.Errorf("expected task2, got %s", task.Description)
	}
}

func TestLoadTaskByID_NotFound(t *testing.T) {
	s, tmpDir := setupTestStore(t)
	defer cleanupTestStore(tmpDir)

	store = s

	_, err := loadTaskByID(99)
	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestSaveNewTaskToJSON(t *testing.T) {
	s, tmpDir := setupTestStore(t)
	defer cleanupTestStore(tmpDir)

	store = s

	task := &Task{ID: 1, Description: "new task", CreatedAt: time.Now()}
	err := saveNewTaskToJSON(task)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(s.Data.Tasks) != 1 {
		t.Errorf("expected 1 task, got %d", len(s.Data.Tasks))
	}
	if s.Data.Tasks[0].Description != "new task" {
		t.Errorf("expected 'new task', got '%s'", s.Data.Tasks[0].Description)
	}
}

func TestUpdateTaskInJSON(t *testing.T) {
	s, tmpDir := setupTestStore(t)
	defer cleanupTestStore(tmpDir)

	s.Data.Tasks = []Task{{ID: 1, Description: "old"}}
	store = s

	task := &Task{ID: 1, Description: "updated"}
	err := updateTaskInJSON(task)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.Data.Tasks[0].Description != "updated" {
		t.Errorf("expected 'updated', got '%s'", s.Data.Tasks[0].Description)
	}
}

func TestUpdateTaskInJSON_NotFound(t *testing.T) {
	s, tmpDir := setupTestStore(t)
	defer cleanupTestStore(tmpDir)

	store = s

	task := &Task{ID: 99, Description: "missing"}
	err := updateTaskInJSON(task)
	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestDeleteTaskFromJSON(t *testing.T) {
	s, tmpDir := setupTestStore(t)
	defer cleanupTestStore(tmpDir)

	s.Data.Tasks = []Task{{ID: 1, Description: "delete me"}, {ID: 2, Description: "keep"}}
	store = s

	task := &Task{ID: 1}
	err := deleteTaskFromJSON(task)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(s.Data.Tasks) != 1 {
		t.Errorf("expected 1 task, got %d", len(s.Data.Tasks))
	}
	if s.Data.Tasks[0].ID != 2 {
		t.Errorf("expected remaining task ID 2, got %d", s.Data.Tasks[0].ID)
	}
}

func TestDeleteTaskFromJSON_NotFound(t *testing.T) {
	s, tmpDir := setupTestStore(t)
	defer cleanupTestStore(tmpDir)

	store = s

	task := &Task{ID: 99}
	err := deleteTaskFromJSON(task)
	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestListAll(t *testing.T) {
	s, tmpDir := setupTestStore(t)
	defer cleanupTestStore(tmpDir)

	s.Data.Tasks = []Task{
		{ID: 1, Description: "a"},
		{ID: 2, Description: "b"},
	}
	store = s

	tasks, err := listAll()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(tasks) != 2 {
		t.Errorf("expected 2 tasks, got %d", len(tasks))
	}
}

func TestList_ByStatus(t *testing.T) {
	s, tmpDir := setupTestStore(t)
	defer cleanupTestStore(tmpDir)

	s.Data.Tasks = []Task{
		{ID: 1, Status: Todo},
		{ID: 2, Status: InProgress},
		{ID: 3, Status: Done},
	}
	store = s

	tasks, err := list(InProgress)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(tasks) != 1 {
		t.Errorf("expected 1 in-progress task, got %d", len(tasks))
	}
	if tasks[0].ID != 2 {
		t.Errorf("expected ID 2, got %d", tasks[0].ID)
	}
}

func TestList_Empty(t *testing.T) {
	s, tmpDir := setupTestStore(t)
	defer cleanupTestStore(tmpDir)

	store = s

	tasks, err := listAll()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(tasks) != 0 {
		t.Errorf("expected empty slice, got %d", len(tasks))
	}
}
