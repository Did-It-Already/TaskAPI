package taskMiddlewares

import (
	"dia_tasks_ms/dtos"
	"dia_tasks_ms/models"
	"net/http"
)

func ValidateEntry(newTask models.Task, res http.ResponseWriter) bool {

	// Crea una instancia vacía de TaskEntryDTO
	dto := dtos.TaskEntryDTO{
		Name:        newTask.Name,
		Description: newTask.Description,
		User_id:     newTask.User_id,
	}

	// Verifica si los campos requeridos existen y contienen datos
	if dto.Name == "" || dto.Description == "" || dto.User_id == "" {
		return false
	}
	return true

}

func ValidateUpdateEntry(prevTask models.Task, newTask dtos.TaskEntryDTO, res http.ResponseWriter) models.Task {
	// dto := dtos.TaskEntryDTO{}
	task := models.Task{
		User_id: prevTask.User_id,
		Is_done: prevTask.Is_done,
	}
	if newTask.Name == "" && newTask.Description == "" {
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
	return task
}
