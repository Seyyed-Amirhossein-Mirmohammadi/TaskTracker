package task

import (
	"fmt"
	"time"
)

func AddTask(description string) {
	task := Task{Id: generateTaskId(), Description: description, Status: Todo, CreatedAt: time.Now(), UpdatedAt: nil}
	saveNewTaskToJson(&task)
}

func UpdateTaskTitle(id int, description string) {
	task, err := loadTaskById(id)
	if err != nil {
		fmt.Println("Not a valid Id. Task not found")
		return
	}
	task.Description = description
	updateTaskInJson(&task)
}

func DeleteTask(id int) {
	task, err := loadTaskById(id)
	if err != nil {
		fmt.Println("Not a valid Id. Task not found")
		return
	}
	deleteTaskFromJson(&task)
}

func ChangeState(id int, status Status) {
	task, err := loadTaskById(id)
	if err != nil {
		fmt.Println("Not a valid Id. Task not found")
		return
	}
	task.Status = status
	updateTaskInJson(&task)
}

func ListAll() {
	listAll()
}

func List(status Status) {
	list(status)
}
