package serializer

import (
	"encoding/json"
	"io"
)

func Decode(body io.ReadCloser, v any) error {
	return json.NewDecoder(body).Decode(&v)
}

func Encode(body io.Writer, v any) error {
	return json.NewEncoder(body).Encode(&v)
}
