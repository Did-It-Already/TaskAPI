package main

import (
	"context"
	"dia_tasks_ms/controllers"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/julienschmidt/httprouter"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	router := httprouter.New()
	taskController := controllers.NewTaskController(getDBSession())
	router.GET("/task/:user_id/:id", taskController.GetTask)
	router.GET("/tasks/:user_id", taskController.GetAllTasks)
	// router.POST("/task", taskMiddlewares.ValidationMiddleware(taskController.CreateTask))
	router.POST("/task", taskController.CreateTask)
	router.DELETE("/task/:user_id/:id", taskController.DeleteTask)
	router.PUT("/task/:user_id/:id", taskController.UpdateTask)
	router.PATCH("/task/:user_id/:id", taskController.UpdateTaskIsDone)
	http.ListenAndServe("localhost:9000", router)
}

func getDBSession() (*mongo.Database, context.Context) {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error cargando el archivo .env: %v", err)
	}

	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	DBConnectionURI := "mongodb+srv://" + dbUser + ":" + dbPassword + "@diatasksdb.xkfycmi.mongodb.net/?retryWrites=true&w=majority"
	fmt.Println("Estableciendo conexión con MongoDB Atlas...")
	clientOptions := options.Client().ApplyURI(DBConnectionURI)
	ctx := context.Background()

	// Conectarse al servidor de MongoDB Atlas
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		fmt.Println(err)
	} else {
		err = client.Ping(ctx, nil)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("Conexión exitosa")
		}

	}

	return client.Database(dbName), ctx
}
