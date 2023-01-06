package apex

import (
	"context"
	"net/http"
	"time"

	exchange "github.com/thrasher-corp/gocryptotrader/exchanges"
	"github.com/thrasher-corp/gocryptotrader/exchanges/request"
)

// Apex is the overarching type across this package
type Apex struct {
	exchange.Base
}

const (
	apexAPIURL     = "https://pro.apex.exchange/api/"
	apexAPIVersion = "/v1"

	// Public endpoints
	apexSystemTime         = "/time"
	apexAllConfigData      = "/symbols"
	apexMarketDepth        = "/depth"
	apexLatestTrades       = "/trades"
	apexKlines             = "/klines"
	apexTicker             = "/ticker"
	apexFundingRateHistory = "/history-funding"
	apexCheckVersion       = "/check-version"
	apexCheckUserExists    = "/check-user-exist"
	apexAdvertisementData  = "/ads-banner"

	// Authenticated endpoints
)

// GetSystemTime gets system time
func (a *Apex) GetSystemTime(ctx context.Context) (time.Time, error) {
	resp := struct {
		Data struct {
			Time int64 `json:"time"`
		} `json:"result"`
		Error
	}{}
	return time.Unix(resp.Data.Time, 0), a.SendHTTPRequest(ctx, exchange.RestSpot, apexSystemTime, publicSpotRate, &resp)
}

// SendHTTPRequest sends an unauthenticated request
func (a *Apex) SendHTTPRequest(ctx context.Context, ePath exchange.URL, path string, f request.EndpointLimit, result UnmarshalTo) error {
	endpointPath, err := a.API.Endpoints.GetURL(ePath)
	if err != nil {
		return err
	}

	err = a.SendPayload(ctx, f, func() (*request.Item, error) {
		return &request.Item{
			Method:        http.MethodGet,
			Path:          apexAPIVersion + endpointPath + path,
			Result:        result,
			Verbose:       a.Verbose,
			HTTPDebugging: a.HTTPDebugging,
			HTTPRecording: a.HTTPRecording}, nil
	})
	if err != nil {
		return err
	}
	return result.GetError()
}
