package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/chyngyz-sydykov/go-web/logger"
)

type AppHandlerInterface interface {
}

type ErrorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error ErrorDetail `json:"error"`
}

type CommonHandler struct {
	logger *logger.Logger
}

const INVALID_REQUEST string = "INVALID_REQUEST"
const RESOURCE_NOT_FOUND string = "RESOURCE_NOT_FOUND"
const SERVER_ERROR string = "SERVER_ERROR"

func NewCommonHandler(logger *logger.Logger) *CommonHandler {

	return &CommonHandler{logger: logger}
}

func (c *CommonHandler) HandleError(w http.ResponseWriter, err error, statusCode int, errorResponse ErrorResponse) {
	c.logger.LogError(statusCode, err)
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(errorResponse)
}
