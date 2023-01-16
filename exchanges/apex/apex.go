package apex

import (
	"bytes"
	"context"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/thrasher-corp/gocryptotrader/common"
	"github.com/thrasher-corp/gocryptotrader/common/crypto"
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
	apexNonce = "/generate-nonce"
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

// GenerateNonce generate and obtain nonce.
func (ap *Apex) GenerateNonce(ctx context.Context, ethAddress, starkKey, chainID string) (*NonceData, error) {
	resp := struct {
		Data NonceData `json:"data"`
		Error
	}{}

	params := url.Values{}
	if ethAddress == "" {
		return nil, errETHAddressMissing
	}
	params.Set("ethAddress", ethAddress)
	if starkKey == "" {
		return nil, errStarkKeyMissing
	}
	params.Set("starkKey", starkKey)
	if chainID == "" {
		return nil, errChainIDMissing
	}
	params.Set("chainId", chainID)
	return &resp.Data, ap.SendAuthHTTPRequest(ctx, exchange.RestSpot, http.MethodPost, apexNonce, params, &resp, publicSpotRate)
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

// SendAuthHTTPRequest sends an authenticated HTTP request
func (ap *Apex) SendAuthHTTPRequest(ctx context.Context, ePath exchange.URL, method, path string, params url.Values, result UnmarshalTo, f request.EndpointLimit) error {
	// creds, err := ap.GetCredentials(ctx)
	// if err != nil {
	// 	return err
	// }

	if result == nil {
		result = &Error{}
	}
	endpointPath, err := ap.API.Endpoints.GetURL(ePath)
	if err != nil {
		return err
	}
	if params == nil {
		params = url.Values{}
	}

	err = ap.SendPayload(ctx, f, func() (*request.Item, error) {
		var (
			payload []byte
			//hmacSignedStr string
		)
		headers := make(map[string]string)

		// timeStr := strconv.FormatInt(time.Now().UnixMilli(), 10)
		// message := timeStr + method + "/api/" + apexAPIVersion + path + params.Encode()
		// hmacSignedStr, err = getSign(message, creds.Secret)
		// if err != nil {
		// 	return nil, err
		// }
		// headers["APEX-SIGNATURE"] = hmacSignedStr
		// headers["APEX-TIMESTAMP"] = timeStr
		// headers["APEX-API-KEY"] = creds.Key
		// headers["APEX-PASSPHRASE"] = "UzmQ0kfonxwb_ZK6I4ue" // passphrase variable to be added

		switch method {
		case http.MethodPost:
			headers["Content-Type"] = "application/x-www-form-urlencoded"
		}
		payload = []byte(params.Encode())
		return &request.Item{
			Method:        method,
			Path:          endpointPath + apexAPIVersion + path,
			Headers:       headers,
			Body:          bytes.NewBuffer(payload),
			Result:        &result,
			AuthRequest:   true,
			Verbose:       ap.Verbose,
			HTTPDebugging: ap.HTTPDebugging,
			HTTPRecording: ap.HTTPRecording}, nil
	})
	if err != nil {
		return err
	}
	return result.GetError()
}

func getSign(msg, secret string) (string, error) {
	hmacSigned, err := crypto.GetHMAC(crypto.HashSHA256, []byte(msg), []byte(secret))
	if err != nil {
		return "", err
	}
	return crypto.HexEncodeToString(hmacSigned), nil
}
