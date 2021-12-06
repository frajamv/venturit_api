package utils

import (
	"encoding/json"
	"net/http"
)

func SendErrorResponse(err error, res http.ResponseWriter) {
	res.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(res).Encode(err.Error())
}

func SendSuccessResponse(data interface{}, res http.ResponseWriter) {
	res.WriteHeader(http.StatusOK)
	json.NewEncoder(res).Encode(data)
}

func SendCustomResponse(data interface{}, status_code int, res http.ResponseWriter) {
	res.WriteHeader(status_code)
	json.NewEncoder(res).Encode(data)
}
