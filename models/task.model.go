package models

import "gopkg.in/mgo.v2/bson"

type Task struct {
	_id         bson.ObjectId
	Name        string
	Description string
	Date        string
	Is_done     bool
	User_id     string
}
