package converter

import (
	"context"
	"net/http"
	"time"
)

type Converter interface {
	Convert(ctx context.Context, from, to string) (any, error)
}

type currConverter struct {
	cl *http.Client
	url string
}

func New(url string) Converter {
	return &currConverter {
		url: url,
		cl: &http.Client{
			Timeout: time.Second * 5,
		},
	}
}

func (c *currConverter) Convert(ctx context.Context, from, to string) (any, error)  {
	//c.cl.Get(c.url)

	return nil, nil
}