package response

import (
	"currency-converter/internal/models"
	"currency-converter/pkg/serializer"
	"errors"
	"net/http"
)

type errorResponse struct {
	Error apiError `json:"error"`
}

type apiError struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func findErr(err error, errs ...error) bool {
	for _, e := range errs {
		if errors.Is(err, e) {
			return true
		}
	}

	return false
}

func Error(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")
	var errResp errorResponse
	errResp.Error.Message = err.Error()

	switch {
	case findErr(err,
		models.ErrPasswordTooLong,
		models.ErrLastNameTooLong,
		models.ErrFirstNameTooLong,
		models.ErrInvalidInput,
		models.ErrInvalidPass,
		models.ErrEmailAlreadyExists,
		models.ErrInvalidEmail):

		errResp.Error.Code = http.StatusBadRequest
		w.WriteHeader(http.StatusBadRequest)
	case findErr(err, models.ErrUserNotFound):
		errResp.Error.Code = http.StatusNotFound
		w.WriteHeader(http.StatusNotFound)
	case findErr(err, models.ErrInvalidCredentials):
		errResp.Error.Code = http.StatusUnauthorized
		w.WriteHeader(http.StatusUnauthorized)
	default:
		errResp.Error.Code = http.StatusInternalServerError
		w.WriteHeader(http.StatusInternalServerError)
	}

	if err := serializer.Encode(w, errResp); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
