package handlers

import (
	"encoding/json"
	"net/http"
)

type RatingHandler struct {
	commonHandler CommonHandler
}

func NewRatingHandler(commonHandler CommonHandler) *RatingHandler {
	return &RatingHandler{
		commonHandler: commonHandler,
	}
}

func (handler *RatingHandler) GetByBookId(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode("{}")

}
