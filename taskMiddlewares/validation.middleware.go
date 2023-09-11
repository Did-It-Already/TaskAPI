package taskMiddlewares

import (
	"dia_tasks_ms/dtos"
	"dia_tasks_ms/models"
	"net/http"
	"regexp"
)

func ValidateEntry(newTask models.Task, res http.ResponseWriter) bool {
	dateRegex := regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)
	// Verifica si los campos requeridos existen y contienen datos
	if newTask.Name == "" || newTask.Description == "" || newTask.User_id == "" {
		return false
	}
	if !dateRegex.MatchString(newTask.Date) {
		return false
	}
	return true

}

func ValidateUpdateEntry(prevTask models.Task, newTask dtos.TaskEntryDTO, res http.ResponseWriter) models.Task {
	dateRegex := regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)
	// dto := dtos.TaskEntryDTO{}
	task := models.Task{
		User_id: prevTask.User_id,
		Is_done: prevTask.Is_done,
	}
	if newTask.Name == "" && newTask.Description == "" && newTask.Date == "" || !dateRegex.MatchString(newTask.Date) {
		return task
	}
	if newTask.Name == "" {
		task.Name = prevTask.Name
	} else {
		task.Name = newTask.Name
	}

	if newTask.Description == "" {
		task.Description = prevTask.Description
	} else {
		task.Description = newTask.Description
	}

	if newTask.Date == "" {
		task.Date = prevTask.Date
	} else {
		task.Date = newTask.Date
	}
	return task
}
