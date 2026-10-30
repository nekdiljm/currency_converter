package handler

import (
	"currency-converter/internal/converter"
	"net/http"
)

type Handler interface {
	Convert(w http.ResponseWriter, r *http.Request)
}

type handler struct {
	currConverter converter.Converter
}

func New(currConverter converter.Converter) Handler {
	return &handler{
		currConverter: currConverter,
	}
}

func (h *handler) Convert(w http.ResponseWriter, r *http.Request) {
	h.currConverter.Convert()
}
