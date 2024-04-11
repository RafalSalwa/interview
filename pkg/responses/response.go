package responses

import (
	"encoding/json"
	"net/http"

	"github.com/RafalSalwa/interview-app-srv/pkg/models"
)

type data struct {
	Success bool    `json:"success"`
	Message *string `json:"message"`
}

type UserResponse struct {
	*models.UserResponse `json:"user"`
}

func InternalServerError(w http.ResponseWriter) {
	NewErrorBuilder().
		SetResponseCode(http.StatusInternalServerError).
		SetReason("Internal server error").
		SetWriter(w).
		Respond()
}
func NotFound(w http.ResponseWriter) {
	NewErrorBuilder().
		SetResponseCode(http.StatusNotFound).
		SetReason("Not found").
		SetWriter(w).
		Respond()
}

func RespondNotFound(w http.ResponseWriter) {
	response := NewNotFoundResponse()
	responseBody := marshalErrorResponse(response)
	Respond(w, http.StatusNotFound, responseBody)
}

func RespondNotAuthorized(w http.ResponseWriter, msg string) {
	w.Header().Set("WWW-Authenticate", "Basic realm=Restricted")

	errorResponse := NewUnauthorizedErrorResponse(msg)
	responseBody := marshalErrorResponse(errorResponse)

	Respond(w, http.StatusUnauthorized, responseBody)
}

func RespondConflict(w http.ResponseWriter, msg string) {
	resp := NewConflictResponse(msg)
	responseBody := marshalErrorResponse(resp)
	Respond(w, http.StatusConflict, responseBody)
}

func RespondBadRequest(w http.ResponseWriter, msg string) {
	errorResponse := NewBadRequestErrorResponse(msg)

	responseBody := marshalErrorResponse(errorResponse)
	Respond(w, http.StatusBadRequest, responseBody)
}

func Respond(w http.ResponseWriter, statusCode int, responseBody []byte) {
	setHTTPHeaders(w, statusCode)
	_, _ = w.Write(responseBody)
}

func RespondOk(w http.ResponseWriter) {
	setHTTPHeaders(w, http.StatusOK)
	_, err := w.Write([]byte("{\"status\":\"ok\"}"))

	if err != nil {
		InternalServerError(w)
	}
}

func RespondCreated(w http.ResponseWriter) {
	setHTTPHeaders(w, http.StatusCreated)
	_, err := w.Write([]byte("{\"status\":\"created\"}"))

	if err != nil {
		InternalServerError(w)
	}
}

func User(w http.ResponseWriter, u models.UserResponse) {
	if u.LastLogin != nil && u.LastLogin.Unix() == 0 {
		u.LastLogin = nil
	}
	response := &UserResponse{&u}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	js, err := json.MarshalIndent(response, "", "   ")
	if err != nil {
		InternalServerError(w)
	}

	Respond(w, http.StatusOK, js)
}

func NewUserResponse(u *models.UserResponse, w http.ResponseWriter) {
	response := &UserResponse{u}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	js, err := json.MarshalIndent(response, "", "   ")
	if err != nil {
		InternalServerError(w)
	}

	Respond(w, http.StatusOK, js)
}

func setHTTPHeaders(w http.ResponseWriter, statusCode int) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(statusCode)
}

func marshalErrorResponse(err interface{}) []byte {
	body, _ := json.Marshal(err)

	return body
}
