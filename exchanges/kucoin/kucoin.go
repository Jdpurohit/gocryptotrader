package kucoin

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/thrasher-corp/gocryptotrader/common/crypto"
	exchange "github.com/thrasher-corp/gocryptotrader/exchanges"
	"github.com/thrasher-corp/gocryptotrader/exchanges/request"
)

// Kucoin is the overarching type across this package
type Kucoin struct {
	exchange.Base
}

const (
	kucoinAPIURL        = "https://api.kucoin.com"
	kucoinAPIVersion    = "1"
	kucoinAPIKeyVersion = "2"

	// Public endpoints
	kucoinGetSymbols = "/api/v1/symbols"

	// Authenticated endpoints
)

// Start implementing public and private exchange API funcs below

// GetAllSpotPairs gets all pairs on the exchange
func (k *Kucoin) GetSymbols(ctx context.Context) ([]SymbolInfo, error) {
	resp := struct {
		Data []SymbolInfo `json:"data"`
		Code string       `json:"code"`
	}{}
	return resp.Data, k.SendHTTPRequest(ctx, exchange.RestSpot, kucoinGetSymbols, publicSpotRate, &resp)
}

// SendHTTPRequest sends an unauthenticated HTTP request
func (k *Kucoin) SendHTTPRequest(ctx context.Context, ePath exchange.URL, path string, f request.EndpointLimit, result interface{}) error {
	endpointPath, err := k.API.Endpoints.GetURL(ePath)
	if err != nil {
		return err
	}

	return k.SendPayload(ctx, f, func() (*request.Item, error) {
		return &request.Item{
			Method:        http.MethodGet,
			Path:          endpointPath + path,
			Result:        result,
			Verbose:       k.Verbose,
			HTTPDebugging: k.HTTPDebugging,
			HTTPRecording: k.HTTPRecording}, nil
	})
}

// SendAuthHTTPRequest sends an authenticated HTTP request
func (k *Kucoin) SendAuthHTTPRequest(ctx context.Context, ePath exchange.URL, method, path string, params map[string]interface{}, f request.EndpointLimit, result interface{}) error {
	creds, err := k.GetCredentials(ctx)
	if err != nil {
		return err
	}

	endpointPath, err := k.API.Endpoints.GetURL(ePath)
	if err != nil {
		return err
	}

	err = k.SendPayload(ctx, f, func() (*request.Item, error) {
		var body io.Reader
		var payload []byte
		if len(params) != 0 {
			payload, err = json.Marshal(params)
			if err != nil {
				return nil, err
			}
			body = bytes.NewBuffer(payload)
		}
		signature := strconv.FormatInt(time.Now().UnixMilli(), 10) + method + path + string(payload)
		var hmacSigned []byte
		hmacSigned, err = crypto.GetHMAC(crypto.HashSHA256, []byte(signature), []byte(creds.Secret))
		if err != nil {
			return nil, err
		}

		headers := map[string]string{
			"KC-API-KEY":         creds.Key,
			"KC-API-SIGN":        crypto.HexEncodeToString(hmacSigned),
			"KC-API-TIMESTAMP":   strconv.FormatInt(time.Now().UnixMilli(), 10),
			"KC-API-PASSPHRASE":  creds.OneTimePassword, // TODO: need pass phrase here!,
			"KC-API-KEY-VERSION": kucoinAPIKeyVersion,
			"Content-Type":       "application/json",
		}

		return &request.Item{
			Method:        method,
			Path:          endpointPath + path,
			Headers:       headers,
			Body:          body,
			Result:        &result,
			AuthRequest:   true,
			Verbose:       k.Verbose,
			HTTPDebugging: k.HTTPDebugging,
			HTTPRecording: k.HTTPRecording}, nil
	})
	return err
}
