package utils

import (
	"dia_tasks_ms/dtos"
	"encoding/json"
	"fmt"
	"net/http"
)

func FormatMessage(message string, status int, res http.ResponseWriter) {
	msgObject := dtos.Message{Msg: message}
	msgJson, _ := json.Marshal(msgObject)
	res.WriteHeader(status)
	fmt.Fprintf(res, "%s\n", msgJson)
}
