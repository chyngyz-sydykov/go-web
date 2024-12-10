package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	my_error "github.com/chyngyz-sydykov/go-web/error"
	"github.com/chyngyz-sydykov/go-web/internal/rating"
)

type RatingHandler struct {
	commonHandler CommonHandler
	service       rating.RatingServiceInterface
}

func NewRatingHandler(service rating.RatingServiceInterface, commonHandler CommonHandler) *RatingHandler {
	return &RatingHandler{
		service:       service,
		commonHandler: commonHandler,
	}
}

func (handler *RatingHandler) SaveRating(w http.ResponseWriter, r *http.Request) {
	var ratingDTO rating.RatingDTO
	err := json.NewDecoder(r.Body).Decode(&ratingDTO)
	if err != nil {
		handler.commonHandler.HandleError(w, err, http.StatusBadRequest,
			ErrorResponse{
				Error: ErrorDetail{
					Code:    INVALID_REQUEST,
					Message: "provided rating payload is invalid.",
				},
			})
		return
	}

	if err := handler.service.Create(&ratingDTO); err != nil {

		if errors.Is(err, my_error.ErrInvalidArgument) {
			handler.commonHandler.HandleError(w, err, http.StatusBadRequest,
				ErrorResponse{
					Error: ErrorDetail{
						Code:    INVALID_REQUEST,
						Message: "provided rating payload is invalid.",
					},
				})
			return
		} else {
			handler.commonHandler.HandleError(w, err, http.StatusInternalServerError,
				ErrorResponse{
					Error: ErrorDetail{
						Code:    SERVER_ERROR,
						Message: "Couldn't create payload.",
					},
				})
		}
		return
	}
	w.Header().Set("Location", "api/v1/rating/"+ratingDTO.RatingID)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(ratingDTO)

}
