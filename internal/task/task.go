package task

import (
	"fmt"
	"time"
)

func InitTaskStore() error {
	return InitStore()
}

func AddTask(description string) (int, error) {
	if description == "" {
		return 0, fmt.Errorf("task description cannot be empty")
	}

	task := Task{
		ID:          generateTaskID(),
		Description: description,
		Status:      Todo,
		CreatedAt:   time.Now(),
		UpdatedAt:   nil,
	}

	if err := saveNewTaskToJSON(&task); err != nil {
		return 0, fmt.Errorf("failed to save task: %w", err)
	}

	return task.ID, nil
}

func UpdateTaskTitle(id int, description string) error {
	if description == "" {
		return fmt.Errorf("task description cannot be empty")
	}

	task, err := loadTaskByID(id)
	if err != nil {
		return fmt.Errorf("task with ID %d not found", id)
	}

	task.Description = description
	now := time.Now()
	task.UpdatedAt = &now

	if err := updateTaskInJSON(&task); err != nil {
		return fmt.Errorf("failed to update task %d: %w", id, err)
	}

	return nil
}

func DeleteTask(id int) error {
	task, err := loadTaskByID(id)
	if err != nil {
		return fmt.Errorf("task with ID %d not found", id)
	}

	if err := deleteTaskFromJSON(&task); err != nil {
		return fmt.Errorf("failed to delete task %d: %w", id, err)
	}

	return nil
}

func ChangeState(id int, status Status) error {
	task, err := loadTaskByID(id)
	if err != nil {
		return fmt.Errorf("task with ID %d not found", id)
	}

	task.Status = status
	now := time.Now()
	task.UpdatedAt = &now

	if err := updateTaskInJSON(&task); err != nil {
		return fmt.Errorf("failed to update task %d status: %w", id, err)
	}

	return nil
}

func ListAll() ([]Task, error) {
	return listAll()
}

func List(status Status) ([]Task, error) {
	return list(status)
}
