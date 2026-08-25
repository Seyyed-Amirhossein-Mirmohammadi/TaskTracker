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

	data, err := json.MarshalIndent(s.Data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	if err := os.WriteFile(s.filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

func generateTaskId() int {
	store.mu.Lock()
	defer store.mu.Unlock()

	store.Data.LastID++
	return store.Data.LastID
}

func loadTaskById(id int) (Task, error) {
	store.mu.Lock()
	defer store.mu.Unlock()

	for _, task := range store.Data.Tasks {
		if task.Id == id {
			return task, nil
		}
	}
	return Task{}, fmt.Errorf("task with ID %d not found", id)
}

func saveNewTaskToJson(task *Task) {
	store.mu.Lock()
	defer store.mu.Unlock()

	store.Data.Tasks = append(store.Data.Tasks, *task)
	if err := store.Save(); err != nil {
		fmt.Printf("Error saving task: %v\n", err)
	}
}

func updateTaskInJson(task *Task) {
	store.mu.Lock()
	defer store.mu.Unlock()

	for i, t := range store.Data.Tasks {
		if t.Id == task.Id {
			store.Data.Tasks[i] = *task
			if err := store.Save(); err != nil {
				fmt.Printf("Error updating task: %v\n", err)
			}
			return
		}
	}
	fmt.Printf("Task with ID %d not found for update\n", task.Id)
}

func deleteTaskFromJson(task *Task) {
	store.mu.Lock()
	defer store.mu.Unlock()

	for i, t := range store.Data.Tasks {
		if t.Id == task.Id {
			store.Data.Tasks = append(store.Data.Tasks[:i], store.Data.Tasks[i+1:]...)
			if err := store.Save(); err != nil {
				fmt.Printf("Error deleting task: %v\n", err)
			}
			return
		}
	}
	fmt.Printf("Task with ID %d not found for deletion\n", task.Id)
}

func listAll() {
	store.mu.Lock()
	defer store.mu.Unlock()

	if len(store.Data.Tasks) == 0 {
		fmt.Println("No tasks found.")
		return
	}

	fmt.Printf("All Tasks (%d total):\n\n", len(store.Data.Tasks))
	for _, task := range store.Data.Tasks {
		displayTask(task)
	}
}

func list(status Status) {
	store.mu.Lock()
	defer store.mu.Unlock()

	var filtered []Task
	for _, task := range store.Data.Tasks {
		if task.Status == status {
			filtered = append(filtered, task)
		}
	}

	if len(filtered) == 0 {
		fmt.Printf("No %s tasks found.\n", status)
		return
	}

	fmt.Printf("%s Tasks (%d total):\n\n", status, len(filtered))
	for _, task := range filtered {
		displayTask(task)
	}
}

func displayTask(task Task) {
	statusStr := ""
	switch task.Status {
	case Todo:
		statusStr = "Todo"
	case InProgress:
		statusStr = "In Progress"
	case Done:
		statusStr = "Done"
	}

	fmt.Printf("[%d] %s\n", task.Id, task.Description)
	fmt.Printf("  Status: %s | Created: %s\n",
		statusStr,
		task.CreatedAt.Format("2006-01-02 15:04"))

	if task.UpdatedAt != nil {
		fmt.Printf("  Updated: %s\n", task.UpdatedAt.Format("2006-01-02 15:04"))
	} else {
		fmt.Printf("  Updated: Never\n")
	}
	fmt.Println()
}

func (s Status) String() string {
	switch s {
	case Todo:
		return "Todo"
	case InProgress:
		return "In Progress"
	case Done:
		return "Done"
	default:
		return "Unknown"
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
	case "Todo":
		*s = Todo
	case "In Progress":
		*s = InProgress
	case "Done":
		*s = Done
	default:
		return fmt.Errorf("invalid status: %s", str)
	}
	return nil
}
