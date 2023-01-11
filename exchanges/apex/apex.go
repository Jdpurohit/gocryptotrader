package apex

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/thrasher-corp/gocryptotrader/common"
	exchange "github.com/thrasher-corp/gocryptotrader/exchanges"
	"github.com/thrasher-corp/gocryptotrader/exchanges/orderbook"
	"github.com/thrasher-corp/gocryptotrader/exchanges/request"
)

// Apex is the overarching type across this package
type Apex struct {
	exchange.Base
}

const (
	apexAPIURL     = "https://pro.apex.exchange/api/"
	apexAPIVersion = "v1"

	// Public endpoints
	apexSystemTime         = "/time"
	apexAllConfigData      = "/symbols"
	apexMarketDepth        = "/depth"
	apexLatestTrades       = "/trades"
	apexKlines             = "/klines"
	apexTicker             = "/ticker"
	apexFundingRateHistory = "/history-funding"
	apexCheckUserExists    = "/check-user-exist"

	// Authenticated endpoints
)

// GetSystemTime gets system time
func (ap *Apex) GetSystemTime(ctx context.Context) (time.Time, error) {
	resp := struct {
		Data struct {
			Time apexTimeMilliSec `json:"time"`
		} `json:"data"`
		Error
	}{}
	return resp.Data.Time.Time(), ap.SendHTTPRequest(ctx, exchange.RestSpot, apexSystemTime, publicSpotRate, &resp)
}

// GetAllConfigData gets all config data
func (ap *Apex) GetAllConfig(ctx context.Context) (Config, error) {
	resp := struct {
		Data     Config `json:"data"`
		TimeCost int64  `json:"timeCost"`
		Error
	}{}
	return resp.Data, ap.SendHTTPRequest(ctx, exchange.RestSpot, apexAllConfigData, publicSpotRate, &resp)
}

// GetMarketDepthData retrieve all active orderbook for one symbol, inclue all bids and asks.
func (ap *Apex) GetMarketDepth(ctx context.Context, symbol string, limit int64) (*Orderbook, error) {
	var o orderbookResponse

	params := url.Values{}
	params.Set("symbol", symbol)
	if limit != 0 {
		params.Set("limit", strconv.FormatInt(limit, 10))
	}
	path := common.EncodeURLValues(apexMarketDepth, params)
	err := ap.SendHTTPRequest(ctx, exchange.RestSpot, path, publicSpotRate, &o)
	if err != nil {
		return nil, err
	}
	return constructOrderbook(&o)
}

func processOB(ob [][2]string) ([]orderbook.Item, error) {
	o := make([]orderbook.Item, len(ob))
	for x := range ob {
		var price, amount float64
		amount, err := strconv.ParseFloat(ob[x][1], 64)
		if err != nil {
			return nil, err
		}
		price, err = strconv.ParseFloat(ob[x][0], 64)
		if err != nil {
			return nil, err
		}
		o[x] = orderbook.Item{
			Price:  price,
			Amount: amount,
		}
	}
	return o, nil
}

func constructOrderbook(o *orderbookResponse) (*Orderbook, error) {
	var (
		s   Orderbook
		err error
	)
	s.Bids, err = processOB(o.Data.Bids)
	if err != nil {
		return nil, err
	}
	s.Asks, err = processOB(o.Data.Asks)
	if err != nil {
		return nil, err
	}
	return &s, err
}

// GetLatestTrades retrieves all latest trading data.
func (ap *Apex) GetLatestTrades(ctx context.Context, symbol string, limit, from int64) ([]TradeData, error) {
	resp := struct {
		Data []TradeData `json:"data"`
		Error
	}{}

	params := url.Values{}
	params.Set("symbol", symbol)
	if limit != 0 {
		params.Set("limit", strconv.FormatInt(limit, 10))
	}
	if from != 0 {
		params.Set("from", strconv.FormatInt(from, 10))
	}
	path := common.EncodeURLValues(apexLatestTrades, params)
	return resp.Data, ap.SendHTTPRequest(ctx, exchange.RestSpot, path, publicSpotRate, &resp)
}

// GetCandlestickChart retrieves all candlestick chart data.
// Note: API response was not as per the v1 documentation so made changes as per the response.
func (ap *Apex) GetCandlestickChart(ctx context.Context, interval, symbol string, startTime, endTime time.Time, limit int64) ([]KlineData, error) {
	resp := struct {
		Data map[string][]KlineData `json:"data"`
		Error
	}{}

	params := url.Values{}
	if interval != "" {
		if !common.StringDataCompare(validIntervals, interval) {
			return nil, errInvalidInterval
		}
		params.Set("interval", interval)
	}
	if symbol == "" {
		return nil, errSymbolMissing
	}
	params.Set("symbol", symbol)

	if !startTime.IsZero() {
		params.Add("start", strconv.FormatInt(startTime.UnixMilli(), 10))
	}
	if !endTime.IsZero() {
		params.Add("end", strconv.FormatInt(endTime.UnixMilli(), 10))
	}
	if limit != 0 {
		params.Set("limit", strconv.FormatInt(limit, 10))
	}
	path := common.EncodeURLValues(apexKlines, params)
	return resp.Data[symbol], ap.SendHTTPRequest(ctx, exchange.RestSpot, path, publicSpotRate, &resp)
}

// GetTicker retrieves the latest data on symbol tickers.
func (ap *Apex) GetTicker(ctx context.Context, symbol string) ([]TickerData, error) {
	resp := struct {
		Data []TickerData `json:"data"`
		Error
	}{}

	params := url.Values{}
	if symbol == "" {
		return nil, errSymbolMissing
	}
	params.Set("symbol", symbol)
	path := common.EncodeURLValues(apexTicker, params)
	return resp.Data, ap.SendHTTPRequest(ctx, exchange.RestSpot, path, publicSpotRate, &resp)
}

// GetFundingRateHistory retrieves the funding rate history.
func (ap *Apex) GetFundingRateHistory(ctx context.Context, symbol string, startTime, endTime time.Time, limit, page int64) ([]FundData, error) {
	resp := struct {
		Data struct {
			FundingHistory []FundData `json:"historyFunds"`
		} `json:"data"`
		Error
	}{}

	params := url.Values{}
	if symbol == "" {
		return nil, errSymbolMissing
	}
	params.Set("symbol", symbol)
	if !startTime.IsZero() {
		params.Add("beginTimeInclusive", strconv.FormatInt(startTime.UnixMilli(), 10))
	}
	if !endTime.IsZero() {
		params.Add("endTimeExclusive", strconv.FormatInt(endTime.UnixMilli(), 10))
	}
	if limit != 0 {
		params.Set("limit", strconv.FormatInt(limit, 10))
	}
	if page >= 0 {
		params.Set("page", strconv.FormatInt(page, 10))
	}
	path := common.EncodeURLValues(apexFundingRateHistory, params)
	return resp.Data.FundingHistory, ap.SendHTTPRequest(ctx, exchange.RestSpot, path, publicSpotRate, &resp)
}

// CheckIfUserExists validates if user exists in the system.
func (ap *Apex) CheckIfUserExists(ctx context.Context, ethAddress string) (bool, error) {
	resp := struct {
		Data bool `json:"data"`
		Error
	}{}

	params := url.Values{}
	if ethAddress != "" {
		params.Set("ethAddress", ethAddress)
	}

	return resp.Data, ap.SendHTTPRequest(ctx, exchange.RestSpot, apexCheckUserExists, publicSpotRate, &resp)
}

// SendHTTPRequest sends an unauthenticated request
func (ap *Apex) SendHTTPRequest(ctx context.Context, ePath exchange.URL, path string, f request.EndpointLimit, result UnmarshalTo) error {
	endpointPath, err := ap.API.Endpoints.GetURL(ePath)
	if err != nil {
		return err
	}

	err = ap.SendPayload(ctx, f, func() (*request.Item, error) {
		return &request.Item{
			Method:        http.MethodGet,
			Path:          endpointPath + apexAPIVersion + path,
			Result:        result,
			Verbose:       ap.Verbose,
			HTTPDebugging: ap.HTTPDebugging,
			HTTPRecording: ap.HTTPRecording}, nil
	})
	if err != nil {
		return err
	}
	return result.GetError()
}
