package utils

import (
	"encoding/json"
	"net/http"
)

// Create and send response object for REST client using statuscode 200 (OK).
func SendErrorResponse(err error, res http.ResponseWriter) {
	res.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(res).Encode(err.Error())
}

// Creante and send response object for REST client using statuscode 500 (Internal Server Error).
func SendSuccessResponse(data interface{}, res http.ResponseWriter) {
	res.WriteHeader(http.StatusOK)
	json.NewEncoder(res).Encode(data)
}
