package currencyratesprovider

import (
	"context"
	"currency-converter/pkg/serializer"
	"fmt"
	"net/http"
	"time"
)

type CurrencyRatesProvider interface {
	GetCurrencyRatesPair(ctx context.Context, from, to string) (CurrencyRate, error)
	GetCurrencyRatesByDate(ctx context.Context, from, to, baseCurrency string) ([]CurrencyRate, error)
	Close()
}

type client struct {
	host       string
	httpClient *http.Client
}

func NewClient(host string, timeout time.Duration) CurrencyRatesProvider {
	return &client{
		host: host,
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}
}

func (c *client) GetCurrencyRatesPair(ctx context.Context, from, to string) (CurrencyRate, error) {
	url := fmt.Sprintf("%s/rate/%s/%s", c.host, from, to)
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)

	response, err := c.httpClient.Do(request)
	if err != nil {
		return CurrencyRate{}, fmt.Errorf("failed to get currency rates: %w", err)
	}
	defer func() {
		_ = response.Body.Close()
	}()

	if response.StatusCode != http.StatusOK {
		return CurrencyRate{}, fmt.Errorf("unexpected status code: %v", response.StatusCode)
	}

	var currencyRate CurrencyRate
	if err := serializer.Decode(response.Body, &currencyRate); err != nil {
		return CurrencyRate{}, fmt.Errorf("failed to decode response body: %w", err)
	}

	return currencyRate, nil
}

func (c *client) GetCurrencyRatesByDate(ctx context.Context, from, to, baseCurrency string) ([]CurrencyRate, error) {
	url := fmt.Sprintf("%s/rates?from=%s&to=%s&base=%s", c.host, from, to, baseCurrency)
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)

	response, err := c.httpClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("failed to get currency rates: %w", err)
	}
	defer func() {
		_ = response.Body.Close()
	}()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %v", response.StatusCode)
	}

	currencyRateByDate := make([]CurrencyRate, 0)
	if err := serializer.Decode(response.Body, &currencyRateByDate); err != nil {
		return nil, fmt.Errorf("failed to decode response body: %w", err)
	}

	return currencyRateByDate, nil
}

func (c *client) Close() {
	c.httpClient.CloseIdleConnections()
}
