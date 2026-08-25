package task

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

type TaskStore struct {
	mu       sync.Mutex
	filePath string
	Data     StoreData
}

type StoreData struct {
	Tasks  []Task `json:"tasks"`
	LastID int    `json:"last_id"`
}

var store *TaskStore

func InitStore() error {
	var err error
	store, err = NewTaskStore("data/tasks.json")
	return err
}

func NewTaskStore(filePath string) (*TaskStore, error) {
	store := &TaskStore{
		filePath: filePath,
		Data: StoreData{
			Tasks:  []Task{},
			LastID: 0,
		},
	}

	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create data directory: %w", err)
	}

	if err := store.Load(); err != nil {
		if os.IsNotExist(err) {
			if err := store.Save(); err != nil {
				return nil, fmt.Errorf("failed to create new task file: %w", err)
			}
		} else {
			return nil, fmt.Errorf("failed to load tasks: %w", err)
		}
	}

	return store, nil
}

func (s *TaskStore) Load() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.loadLocked()
}

func (s *TaskStore) loadLocked() error {
	data, err := os.ReadFile(s.filePath)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(data, &s.Data); err != nil {
		return fmt.Errorf("failed to parse JSON: %w", err)
	}
	return nil
}

func (s *TaskStore) Save() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.saveLocked()
}

func (s *TaskStore) saveLocked() error {
	data, err := json.MarshalIndent(s.Data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}
	if err := os.WriteFile(s.filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}
	return nil
}

func generateTaskID() int {
	store.mu.Lock()
	defer store.mu.Unlock()
	store.Data.LastID++
	return store.Data.LastID
}

func loadTaskByID(id int) (Task, error) {
	store.mu.Lock()
	defer store.mu.Unlock()
	for _, task := range store.Data.Tasks {
		if task.ID == id {
			return task, nil
		}
	}
	return Task{}, fmt.Errorf("task with ID %d not found", id)
}

func saveNewTaskToJSON(task *Task) error {
	store.mu.Lock()
	defer store.mu.Unlock()

	store.Data.Tasks = append(store.Data.Tasks, *task)
	if err := store.saveLocked(); err != nil {
		return fmt.Errorf("failed to save task to file: %w", err)
	}
	return nil
}

func updateTaskInJSON(task *Task) error {
	store.mu.Lock()
	defer store.mu.Unlock()

	for i, t := range store.Data.Tasks {
		if t.ID == task.ID {
			store.Data.Tasks[i] = *task
			if err := store.saveLocked(); err != nil {
				return fmt.Errorf("failed to update task in file: %w", err)
			}
			return nil
		}
	}
	return fmt.Errorf("task with ID %d not found for update", task.ID)
}

func deleteTaskFromJSON(task *Task) error {
	store.mu.Lock()
	defer store.mu.Unlock()

	for i, t := range store.Data.Tasks {
		if t.ID == task.ID {
			store.Data.Tasks = append(store.Data.Tasks[:i], store.Data.Tasks[i+1:]...)
			if err := store.saveLocked(); err != nil {
				return fmt.Errorf("failed to delete task from file: %w", err)
			}
			return nil
		}
	}
	return fmt.Errorf("task with ID %d not found for deletion", task.ID)
}

func listAll() ([]Task, error) {
	store.mu.Lock()
	defer store.mu.Unlock()

	if len(store.Data.Tasks) == 0 {
		return []Task{}, nil
	}

	tasks := make([]Task, len(store.Data.Tasks))
	copy(tasks, store.Data.Tasks)
	return tasks, nil
}

func list(status Status) ([]Task, error) {
	store.mu.Lock()
	defer store.mu.Unlock()

	var filtered []Task
	for _, task := range store.Data.Tasks {
		if task.Status == status {
			filtered = append(filtered, task)
		}
	}

	tasks := make([]Task, len(filtered))
	copy(tasks, filtered)
	return tasks, nil
}

func (s Status) String() string {
	switch s {
	case Todo:
		return "todo"
	case InProgress:
		return "in-progress"
	case Done:
		return "done"
	default:
		return "unknown"
	}
}

func (s Status) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}

func (s *Status) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}
	switch str {
	case "todo":
		*s = Todo
	case "in-progress":
		*s = InProgress
	case "done":
		*s = Done
	default:
		return fmt.Errorf("invalid status: %s", str)
	}
	return nil
}
