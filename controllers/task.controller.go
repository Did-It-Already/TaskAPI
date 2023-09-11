package controllers

import (
	"context"
	"dia_tasks_ms/dtos"
	"dia_tasks_ms/models"
	"dia_tasks_ms/services"
	"dia_tasks_ms/taskMiddlewares"
	"dia_tasks_ms/utils"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2/bson"
)

type TaskController struct {
	DB  *mongo.Database
	ctx context.Context
}

type Message struct {
	Msg string
}

func NewTaskController(DB *mongo.Database, ctx context.Context) *TaskController {
	return &TaskController{DB, ctx}
}

func (taskController TaskController) GetAllTasks(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	user_id := params.ByName("user_id")
	taskCollection := taskController.DB.Collection("tasks")
	cursor, err := taskCollection.Find(taskController.ctx, bson.M{"user_id": user_id})
	if err != nil {
		utils.FormatMessage("Usuario no encontrado", http.StatusNotFound, res)
		return
	}
	defer cursor.Close(taskController.ctx)

	var tasks []bson.M
	if err := cursor.All(taskController.ctx, &tasks); err != nil {
		utils.FormatMessage("Error al obtener los tasks del usuario", http.StatusInternalServerError, res)
		return
	}
	tasksJSON, err := json.Marshal(tasks)
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	fmt.Fprintf(res, "%s\n", tasksJSON)
}

func (taskController TaskController) GetTask(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	res.Header().Set("Content-Type", "application/json")
	id := params.ByName("id")
	user_id := params.ByName("user_id")
	// Verify ObjectId
	taskObjectId := services.VerifyObjectId(id, res)

	if taskObjectId != false {
		task := models.Task{}
		filter := bson.M{"_id": taskObjectId, "user_id": user_id}

		// Find Task
		err := taskController.DB.Collection("tasks").FindOne(taskController.ctx, filter).Decode(&task)
		if err != nil {
			utils.FormatMessage("Task no encontrado", http.StatusNotFound, res)
			return
		}

		taskJSON, _ := json.Marshal(task)
		res.WriteHeader(http.StatusOK)
		fmt.Fprintf(res, "%s\n", taskJSON)
	}

}

func (taskController TaskController) CreateTask(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	res.Header().Set("Content-Type", "application/json")
	// Create Model
	newTask := models.Task{}
	json.NewDecoder(req.Body).Decode(&newTask)
	isTaskComplete := taskMiddlewares.ValidateEntry(newTask, res)
	if !isTaskComplete {
		utils.FormatMessage("Recuerda que los campos Name, Description, Date y User_id son obligatorios y que Date debe tener el formato AAAA-MM-DD", http.StatusBadRequest, res)
	} else {
		newTask.Is_done = false

		// Add to DB
		result, err := taskController.DB.Collection("tasks").InsertOne(taskController.ctx, newTask)
		if err != nil {
			utils.FormatMessage("Error al insertar en la base de datos", http.StatusInternalServerError, res)
		}
		taskJSON, _ := json.Marshal(result)

		res.WriteHeader(http.StatusOK)
		fmt.Fprintf(res, "%s\n", taskJSON)
	}
}

func (taskController TaskController) DeleteTask(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	res.Header().Set("Content-Type", "application/json")
	id := params.ByName("id")
	user_id := params.ByName("user_id")

	if !bson.IsObjectIdHex(id) {
		utils.FormatMessage("Formato erróneo ObjectId", http.StatusNotFound, res)
	} else {
		taskObjectId, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			utils.FormatMessage("Error al obtener el ObjectId", http.StatusBadRequest, res)
		}
		filter := bson.M{"_id": taskObjectId, "user_id": user_id}

		deleteResult, err := taskController.DB.Collection("tasks").DeleteOne(taskController.ctx, filter)
		if err != nil || deleteResult.DeletedCount == 0 {
			utils.FormatMessage("Error al eliminar registro", http.StatusBadRequest, res)

		} else {
			utils.FormatMessage("Deleted tasks: "+fmt.Sprint(deleteResult.DeletedCount), http.StatusOK, res)
		}
	}
}

func (taskController TaskController) UpdateTask(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	res.Header().Set("Content-Type", "application/json")
	id := params.ByName("id")
	user_id := params.ByName("user_id")
	// Verify ObjectId
	taskObjectId := services.VerifyObjectId(id, res)

	if taskObjectId != false {
		task := models.Task{}
		filter := bson.M{"_id": taskObjectId, "user_id": user_id}

		// Find Task
		err := taskController.DB.Collection("tasks").FindOne(taskController.ctx, filter).Decode(&task)
		if err != nil {
			utils.FormatMessage("Task no encontrado", http.StatusNotFound, res)
			return
		}

		newTask := dtos.TaskEntryDTO{}
		json.NewDecoder(req.Body).Decode(&newTask)
		taskCompleted := taskMiddlewares.ValidateUpdateEntry(task, newTask, res)
		if taskCompleted.Name == "" {
			utils.FormatMessage("Debe modificar al menos uno de estos campos: Name, Description, Date. Y Date debe tener el formato AAAA-MM-DD", http.StatusBadRequest, res)
			return
		} else {
			update := bson.M{
				"$set": bson.M{
					// "user_id": taskCompleted.User_id,
					"name":        taskCompleted.Name,
					"description": taskCompleted.Description,
					"date":        taskCompleted.Date,
				},
			}
			result, err := taskController.DB.Collection("tasks").UpdateOne(taskController.ctx, filter, update)
			if err != nil {
				utils.FormatMessage("Error al actualizar en la base de datos", http.StatusInternalServerError, res)
				return
			}
			taskJSON, _ := json.Marshal(result)
			res.WriteHeader(http.StatusOK)
			fmt.Fprintf(res, "%s\n", taskJSON)
		}
	}
}

func (taskController TaskController) UpdateTaskIsDone(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	res.Header().Set("Content-Type", "application/json")
	id := params.ByName("id")
	user_id := params.ByName("user_id")
	// Verify ObjectId
	taskObjectId := services.VerifyObjectId(id, res)

	if taskObjectId != false {
		task := models.Task{}
		filter := bson.M{"_id": taskObjectId, "user_id": user_id}

		// Find Task
		err := taskController.DB.Collection("tasks").FindOne(taskController.ctx, filter).Decode(&task)
		if err != nil {
			utils.FormatMessage("Task no encontrado", http.StatusNotFound, res)
			return
		}

		update := bson.M{
			"$set": bson.M{
				"is_done": !task.Is_done,
			},
		}
		result, err := taskController.DB.Collection("tasks").UpdateOne(taskController.ctx, filter, update)
		if err != nil {
			utils.FormatMessage("Error al actualizar en la base de datos", http.StatusInternalServerError, res)
			return
		}
		taskJSON, _ := json.Marshal(result)
		res.WriteHeader(http.StatusOK)
		fmt.Fprintf(res, "%s\n", taskJSON)
	}
}
