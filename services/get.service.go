package services

import (
	"dia_tasks_ms/utils"
	"net/http"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/mgo.v2/bson"
)

func VerifyObjectId(id string, res http.ResponseWriter) interface{} {
	if !bson.IsObjectIdHex(id) {
		utils.FormatMessage("ObjectId no aceptado", http.StatusNotFound, res)
		return false
	} else {
		taskObjectId, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			res.WriteHeader(http.StatusBadRequest)
			return false
		}
		return taskObjectId
	}
}
