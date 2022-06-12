package bybit

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/thrasher-corp/gocryptotrader/common"
	"github.com/thrasher-corp/gocryptotrader/currency"
	exchange "github.com/thrasher-corp/gocryptotrader/exchanges"
	"github.com/thrasher-corp/gocryptotrader/exchanges/asset"
	"github.com/thrasher-corp/gocryptotrader/exchanges/orderbook"
)

const (

	// public endpoint
	usdcfuturesGetOrderbook          = "/perpetual/usdc/openapi/public/v1/order-book"
	usdcfuturesGetContracts          = "/perpetual/usdc/openapi/public/v1/symbols"
	usdcfuturesGetSymbols            = "/perpetual/usdc/openapi/public/v1/tick"
	usdcfuturesGetKlines             = "/perpetual/usdc/openapi/public/v1/kline/list"
	usdcfuturesGetMarkPriceKlines    = "/perpetual/usdc/openapi/public/v1/mark-price-kline"
	usdcfuturesGetIndexPriceKlines   = "/perpetual/usdc/openapi/public/v1/index-price-kline"
	usdcfuturesGetPremiumIndexKlines = "/perpetual/usdc/openapi/public/v1/premium-index-kline"
	usdcfuturesGetOpenInterest       = "/perpetual/usdc/openapi/public/v1/open-interest"
	usdcfuturesGetLargeOrders        = "/perpetual/usdc/openapi/public/v1/big-deal"
	usdcfuturesGetAccountRatio       = "/perpetual/usdc/openapi/public/v1/account-ratio"
	usdcfuturesGetLatestTrades       = "/option/usdc/openapi/public/v1/query-trade-latest"

	// auth endpoint
	usdcfuturesPlaceOrder           = "/perpetual/usdc/openapi/private/v1/place-order"
	usdcfuturesModifyOrder          = "/perpetual/usdc/openapi/private/v1/replace-order"
	usdcfuturesCancelOrder          = "/perpetual/usdc/openapi/private/v1/cancel-order"
	usdcfuturesCancelAllActiveOrder = "/perpetual/usdc/openapi/private/v1/cancel-all"
	usdcfuturesGetActiveOrder       = "/option/usdc/openapi/private/v1/query-active-orders"
	usdcfuturesGetOrderHistory      = "/option/usdc/openapi/private/v1/query-order-history"
	usdcfuturesGetTradeHistory      = "/option/usdc/openapi/private/v1/execution-list"
)

// GetUSDCFuturesOrderbook gets orderbook data for USDCMarginedFutures.
func (by *Bybit) GetUSDCFuturesOrderbook(ctx context.Context, symbol currency.Pair) (Orderbook, error) {
	var resp Orderbook
	data := struct {
		Result []USDCOrderbookData `json:"result"`
		Error
	}{}

	params := url.Values{}
	symbolValue, err := by.FormatSymbol(symbol, asset.USDCMarginedFutures)
	if err != nil {
		return resp, err
	}
	params.Set("symbol", symbolValue)

	err = by.SendHTTPRequest(ctx, exchange.RestUSDCMargined, common.EncodeURLValues(usdcfuturesGetOrderbook, params), publicFuturesRate, &data)
	if err != nil {
		return resp, err
	}

	for x := range data.Result {
		switch data.Result[x].Side {
		case sideBuy:
			resp.Bids = append(resp.Bids, orderbook.Item{
				Price:  data.Result[x].Price,
				Amount: data.Result[x].Size,
			})
		case sideSell:
			resp.Asks = append(resp.Asks, orderbook.Item{
				Price:  data.Result[x].Price,
				Amount: data.Result[x].Size,
			})
		default:
			return resp, errInvalidSide
		}
	}
	return resp, nil
}

// GetUSDCContracts gets all contract information for USDCMarginedFutures.
func (by *Bybit) GetUSDCContracts(ctx context.Context, symbol currency.Pair, direction string, limit int64) ([]USDCContract, error) {
	resp := struct {
		Data []USDCContract `json:"result"`
		Error
	}{}

	params := url.Values{}
	if !symbol.IsEmpty() {
		symbolValue, err := by.FormatSymbol(symbol, asset.USDCMarginedFutures)
		if err != nil {
			return resp.Data, err
		}
		params.Set("symbol", symbolValue)
	}

	if direction != "" {
		params.Set("direction", direction)
	}
	if limit > 0 && limit <= 200 {
		params.Set("limit", strconv.FormatInt(limit, 10))
	}

	return resp.Data, by.SendHTTPRequest(ctx, exchange.RestUSDCMargined, common.EncodeURLValues(usdcfuturesGetContracts, params), publicFuturesRate, &resp)
}

// GetUSDCSymbols gets all symbol information for USDCMarginedFutures.
func (by *Bybit) GetUSDCSymbols(ctx context.Context, symbol currency.Pair) (USDCSymbol, error) {
	resp := struct {
		Data USDCSymbol `json:"result"`
		Error
	}{}

	params := url.Values{}
	if !symbol.IsEmpty() {
		symbolValue, err := by.FormatSymbol(symbol, asset.USDCMarginedFutures)
		if err != nil {
			return resp.Data, err
		}
		params.Set("symbol", symbolValue)
	} else {
		return USDCSymbol{}, errSymbolMissing
	}

	return resp.Data, by.SendHTTPRequest(ctx, exchange.RestUSDCMargined, common.EncodeURLValues(usdcfuturesGetSymbols, params), publicFuturesRate, &resp)
}

// GetUSDCKlines gets kline of symbol for USDCMarginedFutures.
func (by *Bybit) GetUSDCKlines(ctx context.Context, symbol currency.Pair, period string, startTime time.Time, limit int64) ([]USDCKline, error) {
	resp := struct {
		Data []USDCKline `json:"result"`
		Error
	}{}

	params := url.Values{}
	if !symbol.IsEmpty() {
		symbolValue, err := by.FormatSymbol(symbol, asset.USDCMarginedFutures)
		if err != nil {
			return resp.Data, err
		}
		params.Set("symbol", symbolValue)
	} else {
		return nil, errSymbolMissing
	}

	if !common.StringDataCompare(validFuturesIntervals, period) {
		return resp.Data, errInvalidPeriod
	}
	params.Set("period", period)

	if startTime.IsZero() {
		return nil, errInvalidStartTime
	} else {
		params.Set("startTime", strconv.FormatInt(startTime.Unix(), 10))
	}

	if limit > 0 && limit <= 200 {
		params.Set("limit", strconv.FormatInt(limit, 10))
	}
	return resp.Data, by.SendHTTPRequest(ctx, exchange.RestUSDCMargined, common.EncodeURLValues(usdcfuturesGetKlines, params), publicFuturesRate, &resp)
}

// GetUSDCMarkPriceKlines gets mark price kline of symbol for USDCMarginedFutures.
func (by *Bybit) GetUSDCMarkPriceKlines(ctx context.Context, symbol currency.Pair, period string, startTime time.Time, limit int64) ([]USDCKlineBase, error) {
	resp := struct {
		Data []USDCKlineBase `json:"result"`
		Error
	}{}

	params := url.Values{}
	if !symbol.IsEmpty() {
		symbolValue, err := by.FormatSymbol(symbol, asset.USDCMarginedFutures)
		if err != nil {
			return resp.Data, err
		}
		params.Set("symbol", symbolValue)
	} else {
		return nil, errSymbolMissing
	}

	if !common.StringDataCompare(validFuturesIntervals, period) {
		return resp.Data, errInvalidPeriod
	}
	params.Set("period", period)

	if startTime.IsZero() {
		return nil, errInvalidStartTime
	} else {
		params.Set("startTime", strconv.FormatInt(startTime.Unix(), 10))
	}

	if limit > 0 && limit <= 200 {
		params.Set("limit", strconv.FormatInt(limit, 10))
	}
	return resp.Data, by.SendHTTPRequest(ctx, exchange.RestUSDCMargined, common.EncodeURLValues(usdcfuturesGetMarkPriceKlines, params), publicFuturesRate, &resp)
}

// GetUSDCIndexPriceKlines gets index price kline of symbol for USDCMarginedFutures.
func (by *Bybit) GetUSDCIndexPriceKlines(ctx context.Context, symbol currency.Pair, period string, startTime time.Time, limit int64) ([]USDCKlineBase, error) {
	resp := struct {
		Data []USDCKlineBase `json:"result"`
		Error
	}{}

	params := url.Values{}
	if !symbol.IsEmpty() {
		symbolValue, err := by.FormatSymbol(symbol, asset.USDCMarginedFutures)
		if err != nil {
			return resp.Data, err
		}
		params.Set("symbol", symbolValue)
	} else {
		return nil, errSymbolMissing
	}

	if !common.StringDataCompare(validFuturesIntervals, period) {
		return resp.Data, errInvalidPeriod
	}
	params.Set("period", period)

	if startTime.IsZero() {
		return nil, errInvalidStartTime
	} else {
		params.Set("startTime", strconv.FormatInt(startTime.Unix(), 10))
	}

	if limit > 0 && limit <= 200 {
		params.Set("limit", strconv.FormatInt(limit, 10))
	}
	return resp.Data, by.SendHTTPRequest(ctx, exchange.RestUSDCMargined, common.EncodeURLValues(usdcfuturesGetIndexPriceKlines, params), publicFuturesRate, &resp)
}

// GetUSDCPremiumIndexKlines gets premium index kline of symbol for USDCMarginedFutures.
func (by *Bybit) GetUSDCPremiumIndexKlines(ctx context.Context, symbol currency.Pair, period string, startTime time.Time, limit int64) ([]USDCKlineBase, error) {
	resp := struct {
		Data []USDCKlineBase `json:"result"`
		Error
	}{}

	params := url.Values{}
	if !symbol.IsEmpty() {
		symbolValue, err := by.FormatSymbol(symbol, asset.USDCMarginedFutures)
		if err != nil {
			return resp.Data, err
		}
		params.Set("symbol", symbolValue)
	} else {
		return nil, errSymbolMissing
	}

	if !common.StringDataCompare(validFuturesIntervals, period) {
		return resp.Data, errInvalidPeriod
	}
	params.Set("period", period)

	if startTime.IsZero() {
		return nil, errInvalidStartTime
	} else {
		params.Set("startTime", strconv.FormatInt(startTime.Unix(), 10))
	}

	if limit > 0 && limit <= 200 {
		params.Set("limit", strconv.FormatInt(limit, 10))
	}
	return resp.Data, by.SendHTTPRequest(ctx, exchange.RestUSDCMargined, common.EncodeURLValues(usdcfuturesGetPremiumIndexKlines, params), publicFuturesRate, &resp)
}

// GetUSDCOpenInterest gets open interest of symbol for USDCMarginedFutures.
func (by *Bybit) GetUSDCOpenInterest(ctx context.Context, symbol currency.Pair, period string, limit int64) ([]USDCOpenInterest, error) {
	resp := struct {
		Data []USDCOpenInterest `json:"result"`
		Error
	}{}

	params := url.Values{}
	if !symbol.IsEmpty() {
		symbolValue, err := by.FormatSymbol(symbol, asset.USDCMarginedFutures)
		if err != nil {
			return resp.Data, err
		}
		params.Set("symbol", symbolValue)
	} else {
		return nil, errSymbolMissing
	}

	if !common.StringDataCompare(validFuturesPeriods, period) {
		return resp.Data, errInvalidPeriod
	}
	params.Set("period", period)

	if limit > 0 && limit <= 200 {
		params.Set("limit", strconv.FormatInt(limit, 10))
	}
	return resp.Data, by.SendHTTPRequest(ctx, exchange.RestUSDCMargined, common.EncodeURLValues(usdcfuturesGetOpenInterest, params), publicFuturesRate, &resp)
}

// GetUSDCLargeOrders gets large order of symbol for USDCMarginedFutures.
func (by *Bybit) GetUSDCLargeOrders(ctx context.Context, symbol currency.Pair, limit int64) ([]USDCLargeOrder, error) {
	resp := struct {
		Data []USDCLargeOrder `json:"result"`
		Error
	}{}

	params := url.Values{}
	if !symbol.IsEmpty() {
		symbolValue, err := by.FormatSymbol(symbol, asset.USDCMarginedFutures)
		if err != nil {
			return resp.Data, err
		}
		params.Set("symbol", symbolValue)
	} else {
		return nil, errSymbolMissing
	}

	if limit > 0 && limit <= 100 {
		params.Set("limit", strconv.FormatInt(limit, 10))
	}
	return resp.Data, by.SendHTTPRequest(ctx, exchange.RestUSDCMargined, common.EncodeURLValues(usdcfuturesGetLargeOrders, params), publicFuturesRate, &resp)
}

// GetUSDCAccountRatio gets account long short ratio of symbol for USDCMarginedFutures.
func (by *Bybit) GetUSDCAccountRatio(ctx context.Context, symbol currency.Pair, period string, limit int64) ([]USDCAccountRatio, error) {
	resp := struct {
		Data []USDCAccountRatio `json:"result"`
		Error
	}{}

	params := url.Values{}
	if !symbol.IsEmpty() {
		symbolValue, err := by.FormatSymbol(symbol, asset.USDCMarginedFutures)
		if err != nil {
			return resp.Data, err
		}
		params.Set("symbol", symbolValue)
	} else {
		return nil, errSymbolMissing
	}

	if !common.StringDataCompare(validFuturesPeriods, period) {
		return resp.Data, errInvalidPeriod
	}
	params.Set("period", period)

	if limit > 0 && limit <= 500 {
		params.Set("limit", strconv.FormatInt(limit, 10))
	}
	return resp.Data, by.SendHTTPRequest(ctx, exchange.RestUSDCMargined, common.EncodeURLValues(usdcfuturesGetAccountRatio, params), publicFuturesRate, &resp)
}

// GetUSDCLatestTrades gets lastest 500 trades for USDCMarginedFutures.
func (by *Bybit) GetUSDCLatestTrades(ctx context.Context, symbol currency.Pair, category string, limit int64) ([]USDCTrade, error) {
	resp := struct {
		Result struct {
			ResultSize int64       `json:"resultTotalSize"`
			Cursor     string      `json:"cursor"`
			Data       []USDCTrade `json:"dataList"`
		} `json:"result"`
		Error
	}{}

	params := url.Values{}
	if category != "" {
		params.Set("category", category)
	} else {
		return nil, errors.New("invalid category")
	}

	if !symbol.IsEmpty() {
		symbolValue, err := by.FormatSymbol(symbol, asset.USDCMarginedFutures)
		if err != nil {
			return resp.Result.Data, err
		}
		params.Set("symbol", symbolValue)
	}

	if limit > 0 && limit <= 500 {
		params.Set("limit", strconv.FormatInt(limit, 10))
	}
	return resp.Result.Data, by.SendHTTPRequest(ctx, exchange.RestUSDCMargined, common.EncodeURLValues(usdcfuturesGetLatestTrades, params), publicFuturesRate, &resp)
}

// PlaceUSDCOrder create new USDC derivatives order.
func (by *Bybit) PlaceUSDCOrder(ctx context.Context, symbol currency.Pair, orderType, orderFilter, side, timeInForce, orderLinkID string, orderPrice, orderQty, takeProfit, stopLoss, tptriggerby, slTriggerBy, triggerPrice, triggerBy float64, reduceOnly, closeOnTrigger, mmp bool) (USDCCreateOrderResp, error) {
	resp := struct {
		Result USDCCreateOrderResp `json:"result"`
		Error
	}{}

	req := make(map[string]interface{})
	if !symbol.IsEmpty() {
		symbolValue, err := by.FormatSymbol(symbol, asset.USDCMarginedFutures)
		if err != nil {
			return resp.Result, err
		}
		req["symbol"] = symbolValue
	} else {
		return USDCCreateOrderResp{}, errSymbolMissing
	}

	if orderType != "" {
		req["orderType"] = orderType
	} else {
		return USDCCreateOrderResp{}, errInvalidOrderType
	}

	if orderFilter != "" {
		req["orderFilter"] = orderFilter
	} else {
		return USDCCreateOrderResp{}, errInvalidOrderFilter
	}

	if side != "" {
		req["side"] = side
	} else {
		return USDCCreateOrderResp{}, errInvalidSide
	}

	if orderQty != 0 {
		req["orderQty"] = strconv.FormatFloat(orderQty, 'f', -1, 64)
	} else {
		return USDCCreateOrderResp{}, errInvalidQuantity
	}

	if orderPrice != 0 {
		req["orderPrice"] = strconv.FormatFloat(orderPrice, 'f', -1, 64)
	}

	if timeInForce != "" {
		req["timeInForce"] = timeInForce
	}

	if orderLinkID != "" {
		req["orderLinkId"] = orderLinkID
	}

	if reduceOnly {
		req["reduceOnly"] = true
	} else {
		req["reduceOnly"] = false
	}

	if closeOnTrigger {
		req["closeOnTrigger"] = true
	} else {
		req["closeOnTrigger"] = false
	}

	if mmp {
		req["mmp"] = true
	} else {
		req["mmp"] = false
	}

	if takeProfit != 0 {
		req["takeProfit"] = strconv.FormatFloat(takeProfit, 'f', -1, 64)
	}

	if stopLoss != 0 {
		req["stopLoss"] = strconv.FormatFloat(stopLoss, 'f', -1, 64)
	}

	if tptriggerby != 0 {
		req["tptriggerby"] = tptriggerby
	}

	if slTriggerBy != 0 {
		req["slTriggerBy"] = strconv.FormatFloat(slTriggerBy, 'f', -1, 64)
	}

	if triggerPrice != 0 {
		req["triggerPrice"] = strconv.FormatFloat(triggerPrice, 'f', -1, 64)
	}

	if triggerBy != 0 {
		req["triggerBy"] = triggerBy
	}
	return resp.Result, by.SendAuthHTTPRequest(ctx, exchange.RestUSDCMargined, http.MethodPost, usdcfuturesPlaceOrder, nil, req, &resp, publicFuturesRate)
}

// ModifyUSDCOrder modifies USDC derivatives order.
func (by *Bybit) ModifyUSDCOrder(ctx context.Context, symbol currency.Pair, orderFilter, orderID, orderLinkID string, orderPrice, orderQty, takeProfit, stopLoss, tptriggerby, slTriggerBy, triggerPrice float64) (string, error) {
	resp := struct {
		Result struct {
			OrderID       string `json:"orderId"`
			OrderLinkedID string `json:"orderLinkId"`
		} `json:"result"`
		Error
	}{}

	params := url.Values{}
	if !symbol.IsEmpty() {
		symbolValue, err := by.FormatSymbol(symbol, asset.USDCMarginedFutures)
		if err != nil {
			return resp.Result.OrderID, err
		}
		params.Set("symbol", symbolValue)
	} else {
		return resp.Result.OrderID, errSymbolMissing
	}

	if orderFilter != "" {
		params.Set("orderFilter", orderFilter)
	} else {
		return resp.Result.OrderID, errInvalidOrderFilter
	}

	if orderID != "" {
		params.Set("orderId", orderID)
	}

	if orderLinkID != "" {
		params.Set("orderLinkId", orderLinkID)
	}

	if orderPrice != 0 {
		params.Set("orderPrice", strconv.FormatFloat(orderPrice, 'f', -1, 64))
	}

	if orderQty != 0 {
		params.Set("orderQty", strconv.FormatFloat(orderQty, 'f', -1, 64))
	}

	if takeProfit != 0 {
		params.Set("takeProfit", strconv.FormatFloat(takeProfit, 'f', -1, 64))
	}

	if stopLoss != 0 {
		params.Set("stopLoss", strconv.FormatFloat(stopLoss, 'f', -1, 64))
	}

	if tptriggerby != 0 {
		params.Set("tptriggerby", strconv.FormatFloat(tptriggerby, 'f', -1, 64))
	}

	if slTriggerBy != 0 {
		params.Set("slTriggerBy", strconv.FormatFloat(slTriggerBy, 'f', -1, 64))
	}

	if triggerPrice != 0 {
		params.Set("triggerPrice", strconv.FormatFloat(triggerPrice, 'f', -1, 64))
	}
	return resp.Result.OrderID, by.SendAuthHTTPRequest(ctx, exchange.RestUSDCMargined, http.MethodPost, usdcfuturesModifyOrder, params, nil, &resp, publicFuturesRate)
}

// CancelUSDCOrder cancels USDC derivatives order.
func (by *Bybit) CancelUSDCOrder(ctx context.Context, symbol currency.Pair, orderFilter, orderID, orderLinkID string) (string, error) {
	resp := struct {
		Result struct {
			OrderID string `json:"orderId"`
		} `json:"result"`
		Error
	}{}

	params := url.Values{}
	if !symbol.IsEmpty() {
		symbolValue, err := by.FormatSymbol(symbol, asset.USDCMarginedFutures)
		if err != nil {
			return resp.Result.OrderID, err
		}
		params.Set("symbol", symbolValue)
	} else {
		return resp.Result.OrderID, errSymbolMissing
	}

	if orderFilter != "" {
		params.Set("orderFilter", orderFilter)
	} else {
		return resp.Result.OrderID, errInvalidOrderFilter
	}

	if orderID != "" {
		params.Set("orderId", orderID)
	}

	if orderLinkID != "" {
		params.Set("orderLinkId", orderLinkID)
	}
	return resp.Result.OrderID, by.SendAuthHTTPRequest(ctx, exchange.RestUSDCMargined, http.MethodPost, usdcfuturesCancelOrder, params, nil, &resp, publicFuturesRate)
}

// CancelAllActiveUSDCOrder cancels all active USDC derivatives order.
func (by *Bybit) CancelAllActiveUSDCOrder(ctx context.Context, symbol currency.Pair, orderFilter string) error {
	resp := struct {
		Error
	}{}

	params := url.Values{}
	if !symbol.IsEmpty() {
		symbolValue, err := by.FormatSymbol(symbol, asset.USDCMarginedFutures)
		if err != nil {
			return err
		}
		params.Set("symbol", symbolValue)
	} else {
		return errSymbolMissing
	}

	if orderFilter != "" {
		params.Set("orderFilter", orderFilter)
	} else {
		return errInvalidOrderFilter
	}
	return by.SendAuthHTTPRequest(ctx, exchange.RestUSDCMargined, http.MethodPost, usdcfuturesCancelAllActiveOrder, params, nil, &resp, publicFuturesRate)
}

// GetActiveUSDCOrder gets all active USDC derivatives order.
func (by *Bybit) GetActiveUSDCOrder(ctx context.Context, symbol currency.Pair, category, orderID, orderLinkID, orderFilter, direction, cursor string, limit int64) ([]USDCOrder, error) {
	resp := struct {
		Result struct {
			Cursor          string      `json:"cursor"`
			ResultTotalSize int64       `json:"resultTotalSize"`
			Data            []USDCOrder `json:"dataList"`
		} `json:"result"`
		Error
	}{}

	params := url.Values{}
	if !symbol.IsEmpty() {
		symbolValue, err := by.FormatSymbol(symbol, asset.USDCMarginedFutures)
		if err != nil {
			return resp.Result.Data, err
		}
		params.Set("symbol", symbolValue)
	}

	if category != "" {
		params.Set("category", category)
	} else {
		return nil, errors.New("invalid category")
	}

	if orderID != "" {
		params.Set("orderId", orderID)
	}

	if orderLinkID != "" {
		params.Set("orderLinkId", orderLinkID)
	}

	if orderFilter != "" {
		params.Set("orderFilter", orderFilter)
	}

	if direction != "" {
		params.Set("direction", direction)
	}

	if limit != 0 {
		params.Set("limit", strconv.FormatInt(limit, 10))
	}

	if cursor != "" {
		params.Set("cursor", cursor)
	}
	return resp.Result.Data, by.SendAuthHTTPRequest(ctx, exchange.RestUSDCMargined, http.MethodPost, usdcfuturesGetActiveOrder, params, nil, &resp, publicFuturesRate)
}

// GetUSDCOrderHistory gets order history with support of last 30 days of USDC derivatives order.
func (by *Bybit) GetUSDCOrderHistory(ctx context.Context, symbol currency.Pair, category, orderID, orderLinkID, orderStatus, direction, cursor string, limit int64) ([]USDCOrderHistory, error) {
	resp := struct {
		Result struct {
			Cursor          string             `json:"cursor"`
			ResultTotalSize int64              `json:"resultTotalSize"`
			Data            []USDCOrderHistory `json:"dataList"`
		} `json:"result"`
		Error
	}{}

	params := url.Values{}
	if !symbol.IsEmpty() {
		symbolValue, err := by.FormatSymbol(symbol, asset.USDCMarginedFutures)
		if err != nil {
			return resp.Result.Data, err
		}
		params.Set("symbol", symbolValue)
	}

	if category != "" {
		params.Set("category", category)
	} else {
		return nil, errors.New("invalid category")
	}

	if orderID != "" {
		params.Set("orderId", orderID)
	}

	if orderLinkID != "" {
		params.Set("orderLinkId", orderLinkID)
	}

	if orderStatus != "" {
		params.Set("orderStatus", orderStatus)
	}

	if direction != "" {
		params.Set("direction", direction)
	}

	if limit != 0 {
		params.Set("limit", strconv.FormatInt(limit, 10))
	}

	if cursor != "" {
		params.Set("cursor", cursor)
	}
	return resp.Result.Data, by.SendAuthHTTPRequest(ctx, exchange.RestUSDCMargined, http.MethodPost, usdcfuturesGetOrderHistory, params, nil, &resp, publicFuturesRate)
}

// GetUSDCTradeHistory gets trade history with support of last 30 days of USDC derivatives trades.
func (by *Bybit) GetUSDCTradeHistory(ctx context.Context, symbol currency.Pair, category, orderID, orderLinkID, direction, cursor string, limit int64, startTime time.Time) ([]USDCTradeHistory, error) {
	resp := struct {
		Result struct {
			Cursor          string             `json:"cursor"`
			ResultTotalSize int64              `json:"resultTotalSize"`
			Data            []USDCTradeHistory `json:"dataList"`
		} `json:"result"`
		Error
	}{}

	params := url.Values{}
	if !symbol.IsEmpty() {
		symbolValue, err := by.FormatSymbol(symbol, asset.USDCMarginedFutures)
		if err != nil {
			return resp.Result.Data, err
		}
		params.Set("symbol", symbolValue)
	}

	if category != "" {
		params.Set("category", category)
	} else {
		return nil, errors.New("invalid category")
	}

	if orderID != "" {
		params.Set("orderId", orderID)
	}

	if orderLinkID != "" {
		params.Set("orderLinkId", orderLinkID)
	}

	if startTime.IsZero() {
		return nil, errInvalidStartTime
	} else {
		params.Set("startTime", strconv.FormatInt(startTime.Unix(), 10))
	}

	if direction != "" {
		params.Set("direction", direction)
	}

	if limit > 0 && limit <= 50 {
		params.Set("limit", strconv.FormatInt(limit, 10))
	}

	if cursor != "" {
		params.Set("cursor", cursor)
	}
	return resp.Result.Data, by.SendAuthHTTPRequest(ctx, exchange.RestUSDCMargined, http.MethodPost, usdcfuturesGetTradeHistory, params, nil, &resp, publicFuturesRate)
}
