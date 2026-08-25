package task

import (
	"fmt"
	"time"
)

func InitTaskStore() error {
	return InitStore()
}

func AddTask(description string) {
	task := Task{
		Id:          generateTaskId(),
		Description: description,
		Status:      Todo,
		CreatedAt:   time.Now(),
		UpdatedAt:   nil,
	}
	saveNewTaskToJson(&task)
	fmt.Printf("Task added successfully. ID: %d\n", task.Id)
}

func UpdateTaskTitle(id int, description string) {
	task, err := loadTaskById(id)
	if err != nil {
		fmt.Println("Not a valid Id. Task not found")
		return
	}
	task.Description = description
	now := time.Now()
	task.UpdatedAt = &now
	updateTaskInJson(&task)
	fmt.Printf("Task %d updated successfully\n", id)
}

func DeleteTask(id int) {
	task, err := loadTaskById(id)
	if err != nil {
		fmt.Println("Not a valid Id. Task not found")
		return
	}
	deleteTaskFromJson(&task)
	fmt.Printf("Task %d deleted successfully\n", id)
}

func ChangeState(id int, status Status) {
	task, err := loadTaskById(id)
	if err != nil {
		fmt.Println("Not a valid Id. Task not found")
		return
	}
	task.Status = status
	now := time.Now()
	task.UpdatedAt = &now
	updateTaskInJson(&task)
	fmt.Printf("Task %d status updated to %s\n", id, status)
}

func ListAll() {
	listAll()
}

func List(status Status) {
	list(status)
}
