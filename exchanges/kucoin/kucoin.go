package kucoin

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/thrasher-corp/gocryptotrader/common"
	"github.com/thrasher-corp/gocryptotrader/common/crypto"
	exchange "github.com/thrasher-corp/gocryptotrader/exchanges"
	"github.com/thrasher-corp/gocryptotrader/exchanges/orderbook"
	"github.com/thrasher-corp/gocryptotrader/exchanges/request"
)

// TODO:
// run linter
// handle rate limit for all API

// Kucoin is the overarching type across this package
type Kucoin struct {
	exchange.Base
}

const (
	kucoinAPIURL        = "https://api.kucoin.com"
	kucoinAPIVersion    = "1"
	kucoinAPIKeyVersion = "2"

	// Public endpoints
	kucoinGetSymbols             = "/api/v1/symbols"
	kucoinGetTicker              = "/api/v1/market/orderbook/level1"
	kucoinGetAllTickers          = "/api/v1/market/allTickers"
	kucoinGet24hrStats           = "/api/v1/market/stats"
	kucoinGetMarketList          = "/api/v1/markets"
	kucoinGetPartOrderbook20     = "/api/v1/market/orderbook/level2_20"
	kucoinGetPartOrderbook100    = "/api/v1/market/orderbook/level2_100"
	kucoinGetTradeHistory        = "/api/v1/market/histories"
	kucoinGetKlines              = "/api/v1/market/candles"
	kucoinGetCurrencies          = "/api/v1/currencies"
	kucoinGetCurrency            = "/api/v2/currencies/"
	kucoinGetFiatPrice           = "/api/v1/prices"
	kucoinGetMarkPrice           = "/api/v1/mark-price/%s/current"
	kucoinGetMarginConfiguration = "/api/v1/margin/config"
	kucoinGetServerTime          = "/api/v1/timestamp"
	kucoinGetServiceStatus       = "/api/v1/status"

	// Authenticated endpoints
	kucoinGetOrderbook         = "/api/v3/market/orderbook/level2"
	kucoinGetMarginAccount     = "/api/v1/margin/account"
	kucoinGetMarginRiskLimit   = "/api/v1/risk/limit/strategy"
	kucoinBorrowOrder          = "/api/v1/margin/borrow"
	kucoinGetOutstandingRecord = "/api/v1/margin/borrow/outstanding"
	kucoinGetRepaidRecord      = "/api/v1/margin/borrow/repaid"
	kucoinOneClickRepayment    = "/api/v1/margin/repay/all"
	kucoinRepaySingleOrder     = "/api/v1/margin/repay/single"
	kucoinPostLendOrder        = "/api/v1/margin/lend"
	kucoinCancelLendOrder      = "/api/v1/margin/lend/%s"
	kucoinSetAutoLend          = "/api/v1/margin/toggle-auto-lend"
	kucoinGetActiveOrder       = "/api/v1/margin/lend/active"
	kucoinGetLendHistory       = "/api/v1/margin/lend/done"
	kucoinGetUnsettleLendOrder = "/api/v1/margin/lend/trade/unsettled"
	kucoinGetSettleLendOrder   = "/api/v1/margin/lend/trade/settled"
	kucoinGetAccountLendRecord = "/api/v1/margin/lend/assets"
	kucoinGetLendingMarketData = "/api/v1/margin/market"
	kucoinGetMarginTradeData   = "/api/v1/margin/trade/last"

	kucoinGetIsolatedMarginPairConfig            = "/api/v1/isolated/symbols"
	kucoinGetIsolatedMarginAccountInfo           = "/api/v1/isolated/accounts"
	kucoinGetSingleIsolatedMarginAccountInfo     = "/api/v1/isolated/account/%s"
	kucoinInitiateIsolatedMarginBorrowing        = "/api/v1/isolated/borrow"
	kucoinGetIsolatedOutstandingRepaymentRecords = "/api/v1/isolated/borrow/outstanding"
	kucoinGetIsolatedMarginRepaymentRecords      = "/api/v1/isolated/borrow/repaid"
	kucoinInitiateIsolatedMarginQuickRepayment   = "/api/v1/isolated/repay/all"
	kucoinInitiateIsolatedMarginSingleRepayment  = "/api/v1/isolated/repay/single"

	kucoinPostOrder        = "/api/v1/orders"
	kucoinPostMarginOrder  = "/api/v1/margin/order"
	kucoinPostBulkOrder    = "/api/v1/orders/multi"
	kucoinOrderByID        = "/api/v1/orders/%s"             // used by CancelSingleOrder and GetOrderByID
	kucoinOrderByClientOID = "/api/v1/order/client-order/%s" // used by CancelOrderByClientOID and GetOrderByClientOID
	kucoinOrders           = "/api/v1/orders"                // used by CancelAllOpenOrders and GetOrders
	kucoinGetRecentOrders  = "/api/v1/limit/orders"

	kucoinGetFills       = "/api/v1/fills"
	kucoinGetRecentFills = "/api/v1/limit/fills"

	kucoinStopOrder                 = "/api/v1/stop-order"
	kucoinStopOrderByID             = "/api/v1/stop-order/%s"
	kucoinCancelAllStopOrder        = "/api/v1/stop-order/cancel"
	kucoinGetStopOrderByClientID    = "/api/v1/stop-order/queryOrderByClientOid"
	kucoinCancelStopOrderByClientID = "/api/v1/stop-order/cancelOrderByClientOid"

	// account
	kucoinAccount                        = "/api/v1/accounts"
	kucoinGetAccount                     = "/api/v1/accounts/%s"
	kucoinGetAccountLedgers              = "/api/v1/accounts/ledgers"
	kucoinGetSubAccountBalance           = "/api/v1/sub-accounts/%s"
	kucoinGetAggregatedSubAccountBalance = "/api/v1/sub-accounts"
	kucoinGetTransferableBalance         = "/api/v1/accounts/transferable"
	kucoinTransferMainToSubAccount       = "/api/v2/accounts/sub-transfer"
	kucoinInnerTransfer                  = "/api/v2/accounts/inner-transfer"

	// deposit
	kucoinCreateDepositAddress     = "/api/v1/deposit-addresses"
	kucoinGetDepositAddressV2      = "/api/v2/deposit-addresses"
	kucoinGetDepositAddressV1      = "/api/v1/deposit-addresses"
	kucoinGetDepositList           = "/api/v1/deposits"
	kucoinGetHistoricalDepositList = "/api/v1/hist-deposits"

	// withdrawal
	kucoinWithdrawal                  = "/api/v1/withdrawals"
	kucoinGetHistoricalWithdrawalList = "/api/v1/hist-withdrawals"
	kucoinGetWithdrawalQuotas         = "/api/v1/withdrawals/quotas"
	kucoinCancelWithdrawal            = "/api/v1/withdrawals/%s"

	kucoinBasicFee   = "/api/v1/base-fee"
	kucoinTradingFee = "/api/v1/trade-fees"
)

// GetSymbols gets pairs details on the exchange
func (k *Kucoin) GetSymbols(ctx context.Context, currency string) ([]SymbolInfo, error) {
	resp := struct {
		Data []SymbolInfo `json:"data"`
		Error
	}{}

	params := url.Values{}
	if currency != "" {
		params.Set("market", currency)
	}
	return resp.Data, k.SendHTTPRequest(ctx, exchange.RestSpot, common.EncodeURLValues(kucoinGetSymbols, params), publicSpotRate, &resp)
}

// GetTicker gets pair ticker information
func (k *Kucoin) GetTicker(ctx context.Context, pair string) (Ticker, error) {
	resp := struct {
		Data Ticker `json:"data"`
		Error
	}{}

	params := url.Values{}
	if pair == "" {
		return Ticker{}, errors.New("pair can't be empty") // TODO: error as constant
	}
	params.Set("symbol", pair)
	return resp.Data, k.SendHTTPRequest(ctx, exchange.RestSpot, common.EncodeURLValues(kucoinGetTicker, params), publicSpotRate, &resp)
}

// GetAllTickers gets all trading pair ticker information including 24h volume
func (k *Kucoin) GetAllTickers(ctx context.Context) ([]TickerInfo, error) {
	resp := struct {
		Data struct {
			Time    uint64       `json:"time"` // TODO: find a way to convert it to time.Time
			Tickers []TickerInfo `json:"ticker"`
		} `json:"data"`
		Error
	}{}

	return resp.Data.Tickers, k.SendHTTPRequest(ctx, exchange.RestSpot, kucoinGetAllTickers, publicSpotRate, &resp)
}

// Get24hrStats get the statistics of the specified pair in the last 24 hours
func (k *Kucoin) Get24hrStats(ctx context.Context, pair string) (Stats24hrs, error) {
	resp := struct {
		Data Stats24hrs `json:"data"`
		Error
	}{}

	params := url.Values{}
	if pair == "" {
		return Stats24hrs{}, errors.New("pair can't be empty")
	}
	params.Set("symbol", pair)
	return resp.Data, k.SendHTTPRequest(ctx, exchange.RestSpot, common.EncodeURLValues(kucoinGet24hrStats, params), publicSpotRate, &resp)
}

// GetMarketList get the transaction currency for the entire trading market
func (k *Kucoin) GetMarketList(ctx context.Context) ([]string, error) {
	resp := struct {
		Data []string `json:"data"`
		Error
	}{}

	return resp.Data, k.SendHTTPRequest(ctx, exchange.RestSpot, kucoinGetMarketList, publicSpotRate, &resp)
}

func processOB(ob [][2]string) ([]orderbook.Item, error) {
	o := make([]orderbook.Item, len(ob))
	for x := range ob {
		amount, err := strconv.ParseFloat(ob[x][1], 64)
		if err != nil {
			return nil, err
		}
		price, err := strconv.ParseFloat(ob[x][0], 64)
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
	s.Time = o.Data.Time.Time()
	return &s, err
}

// GetPartOrderbook20 gets orderbook for a specified pair with depth 20
func (k *Kucoin) GetPartOrderbook20(ctx context.Context, pair string) (*Orderbook, error) {
	var o orderbookResponse
	params := url.Values{}
	if pair == "" {
		return nil, errors.New("pair can't be empty")
	}
	params.Set("symbol", pair)
	err := k.SendHTTPRequest(ctx, exchange.RestSpot, common.EncodeURLValues(kucoinGetPartOrderbook20, params), publicSpotRate, &o)
	if err != nil {
		return nil, err
	}
	return constructOrderbook(&o)
}

// GetPartOrderbook100 gets orderbook for a specified pair with depth 100
func (k *Kucoin) GetPartOrderbook100(ctx context.Context, pair string) (*Orderbook, error) {
	var o orderbookResponse
	params := url.Values{}
	if pair == "" {
		return nil, errors.New("pair can't be empty")
	}
	params.Set("symbol", pair)
	err := k.SendHTTPRequest(ctx, exchange.RestSpot, common.EncodeURLValues(kucoinGetPartOrderbook100, params), publicSpotRate, &o)
	if err != nil {
		return nil, err
	}
	return constructOrderbook(&o)
}

// GetOrderbook gets full orderbook for a specified pair
func (k *Kucoin) GetOrderbook(ctx context.Context, pair string) (*Orderbook, error) {
	var o orderbookResponse
	params := url.Values{}
	if pair == "" {
		return nil, errors.New("pair can't be empty")
	}
	params.Set("symbol", pair)
	err := k.SendAuthHTTPRequest(ctx, exchange.RestSpot, http.MethodGet, common.EncodeURLValues(kucoinGetOrderbook, params), nil, publicSpotRate, &o)
	if err != nil {
		return nil, err
	}
	return constructOrderbook(&o)
}

// GetTradeHistory gets trade history of the specified pair
func (k *Kucoin) GetTradeHistory(ctx context.Context, pair string) ([]Trade, error) {
	resp := struct {
		Data []Trade `json:"data"`
		Error
	}{}

	params := url.Values{}
	if pair == "" {
		return nil, errors.New("pair can't be empty")
	}
	params.Set("symbol", pair)
	return resp.Data, k.SendHTTPRequest(ctx, exchange.RestSpot, common.EncodeURLValues(kucoinGetTradeHistory, params), publicSpotRate, &resp)
}

// GetKlines gets kline of the specified pair
func (k *Kucoin) GetKlines(ctx context.Context, pair, period string, start, end time.Time) ([]Kline, error) {
	resp := struct {
		Data [][7]string `json:"data"`
		Error
	}{}

	params := url.Values{}
	if pair == "" {
		return nil, errors.New("pair can't be empty")
	}
	params.Set("symbol", pair)

	if period == "" {
		return nil, errors.New("period can't be empty")
	}
	if !common.StringDataContains(validPeriods, period) {
		return nil, errors.New("invalid period")
	}
	params.Set("type", period)

	if !start.IsZero() {
		params.Set("startAt", strconv.FormatInt(start.Unix(), 10))
	}

	if !end.IsZero() {
		params.Set("endAt", strconv.FormatInt(end.Unix(), 10))
	}
	err := k.SendHTTPRequest(ctx, exchange.RestSpot, common.EncodeURLValues(kucoinGetKlines, params), publicSpotRate, &resp)
	if err != nil {
		return nil, err
	}

	klines := make([]Kline, len(resp.Data))
	for i := range resp.Data {
		t, err := strconv.ParseInt(resp.Data[i][0], 10, 64)
		if err != nil {
			return nil, err
		}

		open, err := strconv.ParseFloat(resp.Data[i][1], 64)
		if err != nil {
			return nil, err
		}

		close, err := strconv.ParseFloat(resp.Data[i][2], 64)
		if err != nil {
			return nil, err
		}

		high, err := strconv.ParseFloat(resp.Data[i][3], 64)
		if err != nil {
			return nil, err
		}

		low, err := strconv.ParseFloat(resp.Data[i][4], 64)
		if err != nil {
			return nil, err
		}

		volume, err := strconv.ParseFloat(resp.Data[i][5], 64)
		if err != nil {
			return nil, err
		}

		amount, err := strconv.ParseFloat(resp.Data[i][6], 64)
		if err != nil {
			return nil, err
		}

		klines[i] = Kline{
			StartTime: time.Unix(t, 0),
			Open:      open,
			Close:     close,
			High:      high,
			Low:       low,
			Volume:    volume,
			Amount:    amount,
		}
	}
	return klines, nil
}

// GetCurrencies gets list of currencies
func (k *Kucoin) GetCurrencies(ctx context.Context) ([]Currency, error) {
	resp := struct {
		Data []Currency `json:"data"`
		Error
	}{}

	return resp.Data, k.SendHTTPRequest(ctx, exchange.RestSpot, kucoinGetCurrencies, publicSpotRate, &resp)
}

// GetCurrencies gets list of currencies
func (k *Kucoin) GetCurrency(ctx context.Context, currency, chain string) (CurrencyDetail, error) {
	resp := struct {
		Data CurrencyDetail `json:"data"`
		Error
	}{}

	if currency == "" {
		return CurrencyDetail{}, errors.New("currency can't be empty")
	}
	params := url.Values{}
	if chain != "" {
		params.Set("chain", chain)
	}
	return resp.Data, k.SendHTTPRequest(ctx, exchange.RestSpot, common.EncodeURLValues(kucoinGetCurrency+strings.ToUpper(currency), params), publicSpotRate, &resp)
}

// GetFiatPrice gets fiat prices of currencies, default base currency is USD
func (k *Kucoin) GetFiatPrice(ctx context.Context, base, currencies string) (map[string]string, error) {
	resp := struct {
		Data map[string]string `json:"data"`
		Error
	}{}

	params := url.Values{}
	if base != "" {
		params.Set("base", base)
	}
	if currencies != "" {
		params.Set("currencies", currencies)
	}
	return resp.Data, k.SendHTTPRequest(ctx, exchange.RestSpot, common.EncodeURLValues(kucoinGetFiatPrice, params), publicSpotRate, &resp)
}

// GetMarkPrice gets index price of the specified pair
func (k *Kucoin) GetMarkPrice(ctx context.Context, pair string) (MarkPrice, error) {
	resp := struct {
		Data MarkPrice `json:"data"`
		Error
	}{}

	if pair == "" {
		return MarkPrice{}, errors.New("pair can't be empty")
	}
	return resp.Data, k.SendHTTPRequest(ctx, exchange.RestSpot, fmt.Sprintf(kucoinGetMarkPrice, pair), publicSpotRate, &resp)
}

// GetMarginConfiguration gets configure info of the margin
func (k *Kucoin) GetMarginConfiguration(ctx context.Context) (MarginConfiguration, error) {
	resp := struct {
		Data MarginConfiguration `json:"data"`
		Error
	}{}

	return resp.Data, k.SendHTTPRequest(ctx, exchange.RestSpot, kucoinGetMarginConfiguration, publicSpotRate, &resp)
}

// GetMarginAccount gets configure info of the margin
func (k *Kucoin) GetMarginAccount(ctx context.Context) (MarginAccounts, error) {
	resp := struct {
		Data MarginAccounts `json:"data"`
		Error
	}{}

	return resp.Data, k.SendAuthHTTPRequest(ctx, exchange.RestSpot, http.MethodGet, kucoinGetMarginAccount, nil, publicSpotRate, &resp)
}

// GetMarginRiskLimit gets cross/isolated margin risk limit, default model is cross margin
func (k *Kucoin) GetMarginRiskLimit(ctx context.Context, marginModel string) ([]MarginRiskLimit, error) {
	resp := struct {
		Data []MarginRiskLimit `json:"data"`
		Error
	}{}

	params := url.Values{}
	if marginModel != "" {
		params.Set("marginModel", marginModel)
	}
	return resp.Data, k.SendAuthHTTPRequest(ctx, exchange.RestSpot, http.MethodGet, common.EncodeURLValues(kucoinGetMarginRiskLimit, params), nil, publicSpotRate, &resp)
}

// PostBorrowOrder used to post borrow order
func (k *Kucoin) PostBorrowOrder(ctx context.Context, currency, orderType, term string, size, maxRate float64) (PostBorrowOrderResp, error) {
	resp := struct {
		Data PostBorrowOrderResp `json:"data"`
		Error
	}{}

	params := make(map[string]interface{})
	if currency == "" {
		return PostBorrowOrderResp{}, errors.New("currency can't be empty")
	}
	params["currency"] = currency
	if orderType == "" {
		return PostBorrowOrderResp{}, errors.New("orderType can't be empty")
	}
	params["type"] = orderType
	if size == 0 {
		return PostBorrowOrderResp{}, errors.New("size can't be zero")
	}
	params["size"] = strconv.FormatFloat(size, 'f', -1, 64)

	if maxRate != 0 {
		params["maxRate"] = strconv.FormatFloat(maxRate, 'f', -1, 64)
	}
	if term != "" {
		params["term"] = term
	}
	return resp.Data, k.SendAuthHTTPRequest(ctx, exchange.RestSpot, http.MethodPost, kucoinBorrowOrder, params, publicSpotRate, &resp)
}

// GetBorrowOrder gets borrow order information
func (k *Kucoin) GetBorrowOrder(ctx context.Context, orderID string) (BorrowOrder, error) {
	resp := struct {
		Data BorrowOrder `json:"data"`
		Error
	}{}

	params := url.Values{}

	if orderID == "" {
		return resp.Data, errors.New("empty orderID")
	}
	params.Set("orderId", orderID)
	return resp.Data, k.SendAuthHTTPRequest(ctx, exchange.RestSpot, http.MethodGet, common.EncodeURLValues(kucoinBorrowOrder, params), nil, publicSpotRate, &resp)
}

// GetOutstandingRecord gets outstanding record information
func (k *Kucoin) GetOutstandingRecord(ctx context.Context, currency string) ([]OutstandingRecord, error) {
	resp := struct {
		Data []OutstandingRecord `json:"items"`
		Error
	}{}

	params := url.Values{}
	if currency != "" {
		params.Set("currency", currency)
	}
	return resp.Data, k.SendAuthHTTPRequest(ctx, exchange.RestSpot, http.MethodGet, common.EncodeURLValues(kucoinGetOutstandingRecord, params), nil, publicSpotRate, &resp)
}

// GetRepaidRecord gets repaid record information
func (k *Kucoin) GetRepaidRecord(ctx context.Context, currency string) ([]RepaidRecord, error) {
	resp := struct {
		Data []RepaidRecord `json:"items"`
		Error
	}{}

	params := url.Values{}
	if currency != "" {
		params.Set("currency", currency)
	}
	return resp.Data, k.SendAuthHTTPRequest(ctx, exchange.RestSpot, http.MethodGet, common.EncodeURLValues(kucoinGetRepaidRecord, params), nil, publicSpotRate, &resp)
}

// OneClickRepayment used to compplete repayment in single go
func (k *Kucoin) OneClickRepayment(ctx context.Context, currency, sequence string, size float64) error {
	resp := struct {
		Error
	}{}

	params := make(map[string]interface{})
	if currency == "" {
		return errors.New("currency can't be empty")
	}
	params["currency"] = currency
	if sequence == "" {
		return errors.New("sequence can't be empty")
	}
	params["sequence"] = sequence
	if size == 0 {
		return errors.New("size can't be zero")
	}
	params["size"] = strconv.FormatFloat(size, 'f', -1, 64)
	return k.SendAuthHTTPRequest(ctx, exchange.RestSpot, http.MethodPost, kucoinOneClickRepayment, params, publicSpotRate, &resp)
}

// SingleOrderRepayment used to repay single order
func (k *Kucoin) SingleOrderRepayment(ctx context.Context, currency, tradeID string, size float64) error {
	resp := struct {
		Error
	}{}

	params := make(map[string]interface{})
	if currency == "" {
		return errors.New("currency can't be empty")
	}
	params["currency"] = currency
	if tradeID == "" {
		return errors.New("tradeId can't be empty")
	}
	params["tradeId"] = tradeID
	if size == 0 {
		return errors.New("size can't be zero")
	}
	params["size"] = strconv.FormatFloat(size, 'f', -1, 64)
	return k.SendAuthHTTPRequest(ctx, exchange.RestSpot, http.MethodPost, kucoinRepaySingleOrder, params, publicSpotRate, &resp)
}

// PostLendOrder used to create lend order
func (k *Kucoin) PostLendOrder(ctx context.Context, currency string, dailyIntRate, size float64, term int64) (string, error) {
	resp := struct {
		OrderID string `json:"orderId"`
		Error
	}{}

	params := make(map[string]interface{})
	if currency == "" {
		return "", errors.New("currency can't be empty")
	}
	params["currency"] = currency
	if dailyIntRate == 0 {
		return "", errors.New("dailyIntRate can't be zero")
	}
	params["dailyIntRate"] = strconv.FormatFloat(dailyIntRate, 'f', -1, 64)
	if size == 0 {
		return "", errors.New("size can't be zero")
	}
	params["size"] = strconv.FormatFloat(size, 'f', -1, 64)
	if term == 0 {
		return "", errors.New("term can't be zero")
	}
	params["term"] = strconv.FormatInt(term, 10)
	return resp.OrderID, k.SendAuthHTTPRequest(ctx, exchange.RestSpot, http.MethodPost, kucoinPostLendOrder, params, publicSpotRate, &resp)
}

// CancelLendOrder used to cancel lend order
func (k *Kucoin) CancelLendOrder(ctx context.Context, orderID string) error {
	resp := struct {
		Error
	}{}

	return k.SendAuthHTTPRequest(ctx, exchange.RestSpot, http.MethodDelete, fmt.Sprintf(kucoinCancelLendOrder, orderID), nil, publicSpotRate, &resp)
}

// SetAutoLend used to set up the automatic lending for a specified currency
func (k *Kucoin) SetAutoLend(ctx context.Context, currency string, dailyIntRate, retainSize float64, term int64, isEnable bool) error {
	resp := struct {
		Error
	}{}

	params := make(map[string]interface{})
	if currency == "" {
		return errors.New("currency can't be empty")
	}
	params["currency"] = currency
	if dailyIntRate == 0 {
		return errors.New("dailyIntRate can't be zero")
	}
	params["dailyIntRate"] = strconv.FormatFloat(dailyIntRate, 'f', -1, 64)
	if retainSize == 0 {
		return errors.New("retainSize can't be zero")
	}
	params["retainSize"] = strconv.FormatFloat(retainSize, 'f', -1, 64)
	if term == 0 {
		return errors.New("term can't be zero")
	}
	params["term"] = strconv.FormatInt(term, 10)
	params["isEnable"] = isEnable
	return k.SendAuthHTTPRequest(ctx, exchange.RestSpot, http.MethodPost, kucoinSetAutoLend, params, publicSpotRate, &resp)
}

// GetActiveOrder gets active lend orders
func (k *Kucoin) GetActiveOrder(ctx context.Context, currency string) ([]LendOrder, error) {
	resp := struct {
		Data []LendOrder `json:"items"`
		Error
	}{}

	params := url.Values{}
	if currency != "" {
		params.Set("currency", currency)
	}
	return resp.Data, k.SendAuthHTTPRequest(ctx, exchange.RestSpot, http.MethodGet, common.EncodeURLValues(kucoinGetActiveOrder, params), nil, publicSpotRate, &resp)
}

// GetLendHistory gets lend orders
func (k *Kucoin) GetLendHistory(ctx context.Context, currency string) ([]LendOrderHistory, error) {
	resp := struct {
		Data []LendOrderHistory `json:"items"`
		Error
	}{}

	params := url.Values{}
	if currency != "" {
		params.Set("currency", currency)
	}
	return resp.Data, k.SendAuthHTTPRequest(ctx, exchange.RestSpot, http.MethodGet, common.EncodeURLValues(kucoinGetLendHistory, params), nil, publicSpotRate, &resp)
}

// GetUnsettleLendOrder gets outstanding lend order list
func (k *Kucoin) GetUnsettleLendOrder(ctx context.Context, currency string) ([]UnsettleLendOrder, error) {
	resp := struct {
		Data []UnsettleLendOrder `json:"items"`
		Error
	}{}

	params := url.Values{}
	if currency != "" {
		params.Set("currency", currency)
	}
	return resp.Data, k.SendAuthHTTPRequest(ctx, exchange.RestSpot, http.MethodGet, common.EncodeURLValues(kucoinGetUnsettleLendOrder, params), nil, publicSpotRate, &resp)
}

// GetSettleLendOrder gets settle lend orders
func (k *Kucoin) GetSettleLendOrder(ctx context.Context, currency string) ([]SettleLendOrder, error) {
	resp := struct {
		Data []SettleLendOrder `json:"items"`
		Error
	}{}

	params := url.Values{}
	if currency != "" {
		params.Set("currency", currency)
	}
	return resp.Data, k.SendAuthHTTPRequest(ctx, exchange.RestSpot, http.MethodGet, common.EncodeURLValues(kucoinGetSettleLendOrder, params), nil, publicSpotRate, &resp)
}

// GetAccountLendRecord get the lending history of the main account
func (k *Kucoin) GetAccountLendRecord(ctx context.Context, currency string) ([]LendRecord, error) {
	resp := struct {
		Data []LendRecord `json:"data"`
		Error
	}{}

	params := url.Values{}
	if currency != "" {
		params.Set("currency", currency)
	}
	return resp.Data, k.SendAuthHTTPRequest(ctx, exchange.RestSpot, http.MethodGet, common.EncodeURLValues(kucoinGetAccountLendRecord, params), nil, publicSpotRate, &resp)
}

// GetLendingMarketData get the lending market data
func (k *Kucoin) GetLendingMarketData(ctx context.Context, currency string, term int64) ([]LendMarketData, error) {
	resp := struct {
		Data []LendMarketData `json:"data"`
		Error
	}{}

	params := url.Values{}
	if currency == "" {
		return resp.Data, errors.New("currency can't be empty")
	}
	params.Set("currency", currency)
	if term != 0 {
		params.Set("term", strconv.FormatInt(term, 10))
	}
	return resp.Data, k.SendAuthHTTPRequest(ctx, exchange.RestSpot, http.MethodGet, common.EncodeURLValues(kucoinGetLendingMarketData, params), nil, publicSpotRate, &resp)
}

// GetMarginTradeData get the last 300 fills in the lending and borrowing market
func (k *Kucoin) GetMarginTradeData(ctx context.Context, currency string) ([]MarginTradeData, error) {
	resp := struct {
		Data []MarginTradeData `json:"data"`
		Error
	}{}

	params := url.Values{}
	if currency == "" {
		return resp.Data, errors.New("currency can't be empty")
	}
	params.Set("currency", currency)
	return resp.Data, k.SendAuthHTTPRequest(ctx, exchange.RestSpot, http.MethodGet, common.EncodeURLValues(kucoinGetMarginTradeData, params), nil, publicSpotRate, &resp)
}

// GetIsolatedMarginPairConfig get the current isolated margin trading pair configuration
func (k *Kucoin) GetIsolatedMarginPairConfig(ctx context.Context) ([]IsolatedMarginPairConfig, error) {
	resp := struct {
		Data []IsolatedMarginPairConfig `json:"data"`
		Error
	}{}

	return resp.Data, k.SendAuthHTTPRequest(ctx, exchange.RestSpot, http.MethodGet, kucoinGetIsolatedMarginPairConfig, nil, publicSpotRate, &resp)
}

// GetIsolatedMarginAccountInfo get all isolated margin accounts of the current user
func (k *Kucoin) GetIsolatedMarginAccountInfo(ctx context.Context, balanceCurrency string) (IsolatedMarginAccountInfo, error) {
	resp := struct {
		Data IsolatedMarginAccountInfo `json:"data"`
		Error
	}{}

	params := url.Values{}
	if balanceCurrency != "" {
		params.Set("balanceCurrency", balanceCurrency)
	}
	return resp.Data, k.SendAuthHTTPRequest(ctx, exchange.RestSpot, http.MethodGet, common.EncodeURLValues(kucoinGetIsolatedMarginAccountInfo, params), nil, publicSpotRate, &resp)
}

// GetSingleIsolatedMarginAccountInfo get single isolated margin accounts of the current user
func (k *Kucoin) GetSingleIsolatedMarginAccountInfo(ctx context.Context, symbol string) (AssetInfo, error) {
	resp := struct {
		Data AssetInfo `json:"data"`
		Error
	}{}

	if symbol == "" {
		return resp.Data, errors.New("symbol can't be empty")
	}
	return resp.Data, k.SendAuthHTTPRequest(ctx, exchange.RestSpot, http.MethodGet, fmt.Sprintf(kucoinGetSingleIsolatedMarginAccountInfo, symbol), nil, publicSpotRate, &resp)
}

// InitiateIsolateMarginBorrowing initiates isolated margin borrowing
func (k *Kucoin) InitiateIsolateMarginBorrowing(ctx context.Context, symbol, currency, borrowStrategy, period string, size, maxRate int64) (string, string, float64, error) {
	resp := struct {
		Data struct {
			OrderID    string  `json:"orderId"`
			Currency   string  `json:"currency"`
			ActualSize float64 `json:"actualSize,string"`
		} `json:"data"`
		Error
	}{}
	params := make(map[string]interface{})
	if symbol == "" {
		return resp.Data.OrderID, resp.Data.Currency, resp.Data.ActualSize, errors.New("symbol can't be empty")
	}
	params["symbol"] = symbol
	if currency == "" {
		return resp.Data.OrderID, resp.Data.Currency, resp.Data.ActualSize, errors.New("currency can't be empty")
	}
	params["currency"] = currency
	if borrowStrategy == "" {
		return resp.Data.OrderID, resp.Data.Currency, resp.Data.ActualSize, errors.New("borrowStrategy can't be empty")
	}
	params["borrowStrategy"] = borrowStrategy
	if size == 0 {
		return resp.Data.OrderID, resp.Data.Currency, resp.Data.ActualSize, errors.New("size can't be zero")
	}
	params["size"] = strconv.FormatInt(size, 10)

	if period != "" {
		params["period"] = period
	}
	if maxRate == 0 {
		params["maxRate"] = strconv.FormatInt(maxRate, 10)
	}
	return resp.Data.OrderID, resp.Data.Currency, resp.Data.ActualSize, k.SendAuthHTTPRequest(ctx, exchange.RestSpot, http.MethodPost, kucoinInitiateIsolatedMarginBorrowing, params, publicSpotRate, &resp)
}

// GetIsolatedOutstandingRepaymentRecords get the outstanding repayment records of isolated margin positions
func (k *Kucoin) GetIsolatedOutstandingRepaymentRecords(ctx context.Context, symbol, currency string, pageSize, currentPage int64) ([]OutstandingRepaymentRecord, error) {
	resp := struct {
		Data struct {
			CurrentPage int64                        `json:"currentPage"`
			PageSize    int64                        `json:"pageSize"`
			TotalNum    int64                        `json:"totalNum"`
			TotalPage   int64                        `json:"totalPage"`
			Items       []OutstandingRepaymentRecord `json:"items"`
		} `json:"data"`
		Error
	}{}

	params := url.Values{}
	if symbol != "" {
		params.Set("symbol", symbol)
	}
	if currency != "" {
		params.Set("currency", currency)
	}
	if pageSize != 0 {
		params.Set("pageSize", strconv.FormatInt(pageSize, 10))
	}
	if currentPage != 0 {
		params.Set("currentPage", strconv.FormatInt(currentPage, 10))
	}
	return resp.Data.Items, k.SendAuthHTTPRequest(ctx, exchange.RestSpot, http.MethodGet, common.EncodeURLValues(kucoinGetIsolatedOutstandingRepaymentRecords, params), nil, publicSpotRate, &resp)
}

// GetIsolatedMarginRepaymentRecords get the repayment records of isolated margin positions
func (k *Kucoin) GetIsolatedMarginRepaymentRecords(ctx context.Context, symbol, currency string, pageSize, currentPage int64) ([]CompletedRepaymentRecord, error) {
	resp := struct {
		Data struct {
			CurrentPage int64                      `json:"currentPage"`
			PageSize    int64                      `json:"pageSize"`
			TotalNum    int64                      `json:"totalNum"`
			TotalPage   int64                      `json:"totalPage"`
			Items       []CompletedRepaymentRecord `json:"items"`
		} `json:"data"`
		Error
	}{}

	params := url.Values{}
	if symbol != "" {
		params.Set("symbol", symbol)
	}
	if currency != "" {
		params.Set("currency", currency)
	}
	if pageSize != 0 {
		params.Set("pageSize", strconv.FormatInt(pageSize, 10))
	}
	if currentPage != 0 {
		params.Set("currentPage", strconv.FormatInt(currentPage, 10))
	}
	return resp.Data.Items, k.SendAuthHTTPRequest(ctx, exchange.RestSpot, http.MethodGet, common.EncodeURLValues(kucoinGetIsolatedMarginRepaymentRecords, params), nil, publicSpotRate, &resp)
}

// InitiateIsolatedMarginQuickRepayment is used to initiate quick repayment for isolated margin accounts
func (k *Kucoin) InitiateIsolatedMarginQuickRepayment(ctx context.Context, symbol, currency, seqStrategy string, size int64) error {
	resp := struct {
		Error
	}{}
	params := make(map[string]interface{})
	if symbol == "" {
		return errors.New("symbol can't be empty")
	}
	params["symbol"] = symbol
	if currency == "" {
		return errors.New("currency can't be empty")
	}
	params["currency"] = currency
	if seqStrategy == "" {
		return errors.New("seqStrategy can't be empty")
	}
	params["seqStrategy"] = seqStrategy
	if size == 0 {
		return errors.New("size can't be zero")
	}
	params["size"] = strconv.FormatInt(size, 10)
	return k.SendAuthHTTPRequest(ctx, exchange.RestSpot, http.MethodPost, kucoinInitiateIsolatedMarginQuickRepayment, params, publicSpotRate, &resp)
}

// InitiateIsolatedMarginSingleRepayment is used to initiate quick repayment for single margin accounts
func (k *Kucoin) InitiateIsolatedMarginSingleRepayment(ctx context.Context, symbol, currency, loanId string, size int64) error {
	resp := struct {
		Error
	}{}
	params := make(map[string]interface{})
	if symbol == "" {
		return errors.New("symbol can't be empty")
	}
	params["symbol"] = symbol
	if currency == "" {
		return errors.New("currency can't be empty")
	}
	params["currency"] = currency
	if loanId == "" {
		return errors.New("loanId can't be empty")
	}
	params["loanId"] = loanId
	if size == 0 {
		return errors.New("size can't be zero")
	}
	params["size"] = strconv.FormatInt(size, 10)
	return k.SendAuthHTTPRequest(ctx, exchange.RestSpot, http.MethodPost, kucoinInitiateIsolatedMarginSingleRepayment, params, publicSpotRate, &resp)
}

// GetCurrentServerTime gets the server time
func (k *Kucoin) GetCurrentServerTime(ctx context.Context) (time.Time, error) {
	resp := struct {
		Timestamp kucoinTimeMilliSec `json:"data"`
		Error
	}{}
	return resp.Timestamp.Time(), k.SendHTTPRequest(ctx, exchange.RestSpot, kucoinGetServerTime, publicSpotRate, &resp)
}

// GetServiceStatus gets the service status
func (k *Kucoin) GetServiceStatus(ctx context.Context) (string, string, error) {
	resp := struct {
		Data struct {
			Status  string `json:"status"`
			Message string `json:"msg"`
		} `json:"data"`
		Error
	}{}
	return resp.Data.Status, resp.Data.Message, k.SendHTTPRequest(ctx, exchange.RestSpot, kucoinGetServiceStatus, publicSpotRate, &resp)
}

// PostOrder used to place two types of orders: limit and market
// Note: use this only for SPOT trades
func (k *Kucoin) PostOrder(ctx context.Context, clientOID, side, symbol, orderType, remark, stop, price, timeInForce string, size, cancelAfter, visibleSize, funds float64, postOnly, hidden, iceberg bool) (string, error) {
	resp := struct {
		OrderID string `json:"orderId"`
		Error
	}{}

	params := make(map[string]interface{})
	if clientOID == "" {
		return resp.OrderID, errors.New("clientOid can't be empty")
	}
	params["clientOid"] = clientOID
	if side == "" {
		return resp.OrderID, errors.New("side can't be empty")
	}
	params["side"] = side
	if symbol == "" {
		return resp.OrderID, errors.New("symbol can't be empty")
	}
	params["symbol"] = symbol
	if remark != "" {
		params["remark"] = remark
	}
	if stop != "" {
		params["stp"] = stop
	}
	if orderType == "limit" || orderType == "" {
		if price == "" {
			return resp.OrderID, errors.New("price can't be empty")
		}
		params["price"] = price
		if size <= 0 {
			return resp.OrderID, errors.New("size can't be zero or negative")
		}
		params["size"] = strconv.FormatFloat(size, 'f', -1, 64)
		if timeInForce != "" {
			params["timeInForce"] = timeInForce
		}
		if cancelAfter > 0 && timeInForce == "GTT" {
			params["cancelAfter"] = strconv.FormatFloat(cancelAfter, 'f', -1, 64)
		}
		params["postOnly"] = postOnly
		params["hidden"] = hidden
		params["iceberg"] = iceberg
		if visibleSize > 0 {
			params["visibleSize"] = strconv.FormatFloat(visibleSize, 'f', -1, 64)
		}
	} else if orderType == "market" {
		if size > 0 {
			params["size"] = strconv.FormatFloat(size, 'f', -1, 64)
		} else if funds > 0 {
			params["funds"] = strconv.FormatFloat(funds, 'f', -1, 64)
		} else {
			return resp.OrderID, errors.New("atleast one required among size and funds")
		}
	} else {
		return resp.OrderID, errors.New("invalid orderType")
	}

	if orderType != "" {
		params["type"] = orderType
	}
	return resp.OrderID, k.SendAuthHTTPRequest(ctx, exchange.RestSpot, http.MethodPost, kucoinPostOrder, params, publicSpotRate, &resp)
}

// PostMarginOrder used to place two types of margin orders: limit and market
func (k *Kucoin) PostMarginOrder(ctx context.Context, clientOID, side, symbol, orderType, remark, stop, marginMode, price, timeInForce string, size, cancelAfter, visibleSize, funds float64, postOnly, hidden, iceberg, autoBorrow bool) (PostMarginOrderResp, error) {
	resp := struct {
		PostMarginOrderResp
		Error
	}{}

	params := make(map[string]interface{})
	if clientOID == "" {
		return resp.PostMarginOrderResp, errors.New("clientOid can't be empty")
	}
	params["clientOid"] = clientOID
	if side == "" {
		return resp.PostMarginOrderResp, errors.New("side can't be empty")
	}
	params["side"] = side
	if symbol == "" {
		return resp.PostMarginOrderResp, errors.New("symbol can't be empty")
	}
	params["symbol"] = symbol
	if remark != "" {
		params["remark"] = remark
	}
	if stop != "" {
		params["stp"] = stop
	}
	if marginMode != "" {
		params["marginMode"] = marginMode
	}
	params["autoBorrow"] = autoBorrow
	if orderType == "limit" || orderType == "" {
		if price == "" {
			return resp.PostMarginOrderResp, errors.New("price can't be empty")
		}
		params["price"] = price
		if size <= 0 {
			return resp.PostMarginOrderResp, errors.New("size can't be zero or negative")
		}
		params["size"] = strconv.FormatFloat(size, 'f', -1, 64)
		if timeInForce != "" {
			params["timeInForce"] = timeInForce
		}
		if cancelAfter > 0 && timeInForce == "GTT" {
			params["cancelAfter"] = strconv.FormatFloat(cancelAfter, 'f', -1, 64)
		}
		params["postOnly"] = postOnly
		params["hidden"] = hidden
		params["iceberg"] = iceberg
		if visibleSize > 0 {
			params["visibleSize"] = strconv.FormatFloat(visibleSize, 'f', -1, 64)
		}
	} else if orderType == "market" {
		if size > 0 {
			params["size"] = strconv.FormatFloat(size, 'f', -1, 64)
		} else if funds > 0 {
			params["funds"] = strconv.FormatFloat(funds, 'f', -1, 64)
		} else {
			return resp.PostMarginOrderResp, errors.New("atleast one required among size and funds")
		}
	} else {
		return resp.PostMarginOrderResp, errors.New("invalid orderType")
	}

	if orderType != "" {
		params["type"] = orderType
	}
	return resp.PostMarginOrderResp, k.SendAuthHTTPRequest(ctx, exchange.RestSpot, http.MethodPost, kucoinPostMarginOrder, params, publicSpotRate, &resp)
}

// PostBulkOrder used to place 5 orders at the same time. The order type must be a limit order of the same symbol
// Note: it supports only SPOT trades
// Note: To check if order was posted successfully, check status field in response
func (k *Kucoin) PostBulkOrder(ctx context.Context, symbol string, orderList []OrderRequest) ([]PostBulkOrderResp, error) {
	resp := struct {
		Data struct {
			Data []PostBulkOrderResp `json:"data"`
		} `json:"data"`
		Error
	}{}

	if symbol == "" {
		return resp.Data.Data, errors.New("symbol can't be empty")
	}
	for i := range orderList {
		if orderList[i].ClientOID == "" {
			return resp.Data.Data, errors.New("clientOid can't be empty")
		}
		if orderList[i].Side == "" {
			return resp.Data.Data, errors.New("side can't be empty")
		}
		if orderList[i].Price == "" {
			return resp.Data.Data, errors.New("price can't be empty")
		}
		if orderList[i].Size == "" {
			return resp.Data.Data, errors.New("size can't be empty")
		}
	}
	params := make(map[string]interface{})
	params["symbol"] = symbol
	params["orderList"] = orderList
	return resp.Data.Data, k.SendAuthHTTPRequest(ctx, exchange.RestSpot, http.MethodPost, kucoinPostBulkOrder, params, publicSpotRate, &resp)
}

// CancelSingleOrder used to cancel single order previously placed
func (k *Kucoin) CancelSingleOrder(ctx context.Context, orderID string) ([]string, error) {
	resp := struct {
		CancelledOrderIDs []string `json:"cancelledOrderIds"`
		Error
	}{}

	if orderID == "" {
		return resp.CancelledOrderIDs, errors.New("orderID can't be empty")
	}
	return resp.CancelledOrderIDs, k.SendAuthHTTPRequest(ctx, exchange.RestSpot, http.MethodDelete, fmt.Sprintf(kucoinOrderByID, orderID), nil, publicSpotRate, &resp)
}

// CancelOrderByClientOID used to cancel order via the clientOid
func (k *Kucoin) CancelOrderByClientOID(ctx context.Context, orderID string) (string, string, error) {
	resp := struct {
		CancelledOrderID string `json:"cancelledOrderId"`
		ClientOID        string `json:"clientOid"`
		Error
	}{}

	return resp.CancelledOrderID, resp.ClientOID, k.SendAuthHTTPRequest(ctx, exchange.RestSpot, http.MethodDelete, fmt.Sprintf(kucoinOrderByClientOID, orderID), nil, publicSpotRate, &resp)
}

// CancelAllOpenOrders used to cancel all order based upon the parameters passed
func (k *Kucoin) CancelAllOpenOrders(ctx context.Context, symbol, tradeType string) ([]string, error) {
	resp := struct {
		CancelledOrderIDs []string `json:"cancelledOrderIds"`
		Error
	}{}

	params := url.Values{}
	if symbol != "" {
		params.Set("symbol", symbol)
	}
	if tradeType != "" {
		params.Set("tradeType", tradeType)
	}
	return resp.CancelledOrderIDs, k.SendAuthHTTPRequest(ctx, exchange.RestSpot, http.MethodDelete, common.EncodeURLValues(kucoinOrders, params), nil, publicSpotRate, &resp)
}

// GetOrders gets the user order list
func (k *Kucoin) GetOrders(ctx context.Context, status, symbol, side, orderType, tradeType string, startAt, endAt time.Time) ([]OrderDetail, error) {
	resp := struct {
		Data struct {
			CurrentPage int64         `json:"currentPage"`
			PageSize    int64         `json:"pageSize"`
			TotalNum    int64         `json:"totalNum"`
			TotalPage   int64         `json:"totalPage"`
			Items       []OrderDetail `json:"items"`
		} `json:"data"`
		Error
	}{}

	params := url.Values{}
	if status != "" {
		params.Set("status", status)
	}
	if symbol != "" {
		params.Set("symbol", symbol)
	}
	if side != "" {
		params.Set("side", side)
	}
	if orderType != "" {
		params.Set("type", orderType)
	}
	if tradeType != "" {
		params.Set("tradeType", tradeType)
	}
	if !startAt.IsZero() {
		params.Set("startAt", strconv.FormatInt(startAt.UnixMilli(), 10))
	}
	if !endAt.IsZero() {
		params.Set("startAt", strconv.FormatInt(endAt.UnixMilli(), 10))
	}
	return resp.Data.Items, k.SendAuthHTTPRequest(ctx, exchange.RestSpot, http.MethodGet, common.EncodeURLValues(kucoinOrders, params), nil, publicSpotRate, &resp)
}

// GetRecentOrders get orders in the last 24 hours.
func (k *Kucoin) GetRecentOrders(ctx context.Context) ([]OrderDetail, error) {
	resp := struct {
		Data struct {
			CurrentPage int64         `json:"currentPage"`
			PageSize    int64         `json:"pageSize"`
			TotalNum    int64         `json:"totalNum"`
			TotalPage   int64         `json:"totalPage"`
			Items       []OrderDetail `json:"items"`
		} `json:"data"`
		Error
	}{}

	return resp.Data.Items, k.SendAuthHTTPRequest(ctx, exchange.RestSpot, http.MethodGet, kucoinGetRecentOrders, nil, publicSpotRate, &resp)
}

// GetOrderByID get a single order info by order ID
func (k *Kucoin) GetOrderByID(ctx context.Context, orderID string) (OrderDetail, error) {
	resp := struct {
		Data OrderDetail `json:"data"`
		Error
	}{}

	if orderID == "" {
		return resp.Data, errors.New("orderID can't be empty")
	}
	return resp.Data, k.SendAuthHTTPRequest(ctx, exchange.RestSpot, http.MethodGet, fmt.Sprintf(kucoinOrderByID, orderID), nil, publicSpotRate, &resp)
}

// GetOrderByClientOID get a single order info by client order ID
func (k *Kucoin) GetOrderByClientOID(ctx context.Context, clientOID string) (OrderDetail, error) {
	resp := struct {
		Data OrderDetail `json:"data"`
		Error
	}{}

	if clientOID == "" {
		return resp.Data, errors.New("clientOID can't be empty")
	}
	return resp.Data, k.SendAuthHTTPRequest(ctx, exchange.RestSpot, http.MethodGet, fmt.Sprintf(kucoinOrderByClientOID, clientOID), nil, publicSpotRate, &resp)
}

// GetFills get fills
func (k *Kucoin) GetFills(ctx context.Context, orderID, symbol, side, orderType, tradeType string, startAt, endAt time.Time) ([]Fill, error) {
	resp := struct {
		Data []Fill `json:"items"`
		Error
	}{}

	params := url.Values{}
	if orderID != "" {
		params.Set("orderId", orderID)
	}
	if symbol != "" {
		params.Set("symbol", symbol)
	}
	if side != "" {
		params.Set("side", side)
	}
	if orderType != "" {
		params.Set("type", orderType)
	}
	if !startAt.IsZero() {
		params.Set("startAt", strconv.FormatInt(startAt.UnixMilli(), 10))
	}
	if !endAt.IsZero() {
		params.Set("endAt", strconv.FormatInt(endAt.UnixMilli(), 10))
	}
	if tradeType != "" {
		params.Set("tradeType", tradeType)
	}
	return resp.Data, k.SendAuthHTTPRequest(ctx, exchange.RestSpot, http.MethodGet, common.EncodeURLValues(kucoinGetFills, params), nil, publicSpotRate, &resp)
}

// GetRecentFills get a list of 1000 fills in last 24 hours
func (k *Kucoin) GetRecentFills(ctx context.Context) ([]Fill, error) {
	resp := struct {
		Data []Fill `json:"data"`
		Error
	}{}
	return resp.Data, k.SendAuthHTTPRequest(ctx, exchange.RestSpot, http.MethodGet, kucoinGetRecentFills, nil, publicSpotRate, &resp)
}

// PostStopOrder used to place two types of stop orders: limit and market
func (k *Kucoin) PostStopOrder(ctx context.Context, clientOID, side, symbol, orderType, remark, stop, price, stopPrice, stp, tradeType, timeInForce string, size, cancelAfter, visibleSize, funds float64, postOnly, hidden, iceberg bool) (string, error) {
	resp := struct {
		OrderID string `json:"orderId"`
		Error
	}{}

	params := make(map[string]interface{})
	if clientOID == "" {
		return resp.OrderID, errors.New("clientOid can't be empty")
	}
	params["clientOid"] = clientOID
	if side == "" {
		return resp.OrderID, errors.New("side can't be empty")
	}
	params["side"] = side
	if symbol == "" {
		return resp.OrderID, errors.New("symbol can't be empty")
	}
	params["symbol"] = symbol
	if remark != "" {
		params["remark"] = remark
	}
	if stop != "" {
		params["stop"] = stop
		if stopPrice == "" {
			return resp.OrderID, errors.New("stopPrice can't be empty when stop is set")
		}
		params["stopPrice"] = stopPrice
	}
	if stp != "" {
		params["stp"] = stp
	}
	if tradeType != "" {
		params["tradeType"] = tradeType
	}
	if orderType == "limit" || orderType == "" {
		if price == "" {
			return resp.OrderID, errors.New("price can't be empty")
		}
		params["price"] = price
		if size <= 0 {
			return resp.OrderID, errors.New("size can't be zero or negative")
		}
		params["size"] = strconv.FormatFloat(size, 'f', -1, 64)
		if timeInForce != "" {
			params["timeInForce"] = timeInForce
		}
		if cancelAfter > 0 && timeInForce == "GTT" {
			params["cancelAfter"] = strconv.FormatFloat(cancelAfter, 'f', -1, 64)
		}
		params["postOnly"] = postOnly
		params["hidden"] = hidden
		params["iceberg"] = iceberg
		if visibleSize > 0 {
			params["visibleSize"] = strconv.FormatFloat(visibleSize, 'f', -1, 64)
		}
	} else if orderType == "market" {
		if size > 0 {
			params["size"] = strconv.FormatFloat(size, 'f', -1, 64)
		} else if funds > 0 {
			params["funds"] = strconv.FormatFloat(funds, 'f', -1, 64)
		} else {
			return resp.OrderID, errors.New("atleast one required among size and funds")
		}
	} else {
		return resp.OrderID, errors.New("invalid orderType")
	}

	if orderType != "" {
		params["type"] = orderType
	}
	return resp.OrderID, k.SendAuthHTTPRequest(ctx, exchange.RestSpot, http.MethodPost, kucoinStopOrder, params, publicSpotRate, &resp)
}

// CancelStopOrder used to cancel single stop order previously placed
func (k *Kucoin) CancelStopOrder(ctx context.Context, orderID string) ([]string, error) {
	resp := struct {
		Data []string `json:"cancelledOrderIds"`
		Error
	}{}

	if orderID == "" {
		return resp.Data, errors.New("orderID can't be empty")
	}
	return resp.Data, k.SendAuthHTTPRequest(ctx, exchange.RestSpot, http.MethodDelete, fmt.Sprintf(kucoinStopOrderByID, orderID), nil, publicSpotRate, &resp)
}

// CancelAllStopOrder used to cancel all order based upon the parameters passed
func (k *Kucoin) CancelAllStopOrder(ctx context.Context, symbol, tradeType, orderIDs string) ([]string, error) {
	resp := struct {
		CancelledOrderIDs []string `json:"cancelledOrderIds"`
		Error
	}{}

	params := url.Values{}
	if symbol != "" {
		params.Set("symbol", symbol)
	}
	if tradeType != "" {
		params.Set("tradeType", tradeType)
	}
	if orderIDs != "" {
		params.Set("orderIds", orderIDs)
	}
	return resp.CancelledOrderIDs, k.SendAuthHTTPRequest(ctx, exchange.RestSpot, http.MethodDelete, common.EncodeURLValues(kucoinCancelAllStopOrder, params), nil, publicSpotRate, &resp)
}

// GetStopOrder used to cancel single stop order previously placed
func (k *Kucoin) GetStopOrder(ctx context.Context, orderID string) (StopOrder, error) {
	resp := struct {
		StopOrder
		Error
	}{}

	if orderID == "" {
		return resp.StopOrder, errors.New("orderID can't be empty")
	}
	return resp.StopOrder, k.SendAuthHTTPRequest(ctx, exchange.RestSpot, http.MethodGet, fmt.Sprintf(kucoinStopOrderByID, orderID), nil, publicSpotRate, &resp)
}

// GetAllStopOrder get all current untriggered stop orders
func (k *Kucoin) GetAllStopOrder(ctx context.Context, symbol, side, orderType, tradeType, orderIDs string, startAt, endAt time.Time, currentPage, pageSize int64) ([]StopOrder, error) {
	resp := struct {
		Data struct {
			CurrentPage int64       `json:"currentPage"`
			PageSize    int64       `json:"pageSize"`
			TotalNum    int64       `json:"totalNum"`
			TotalPage   int64       `json:"totalPage"`
			Items       []StopOrder `json:"items"`
		} `json:"data"`
		Error
	}{}

	params := url.Values{}
	if symbol != "" {
		params.Set("symbol", symbol)
	}
	if side != "" {
		params.Set("side", side)
	}
	if orderType != "" {
		params.Set("type", orderType)
	}
	if tradeType != "" {
		params.Set("tradeType", tradeType)
	}
	if orderIDs != "" {
		params.Set("orderIds", orderIDs)
	}
	if !startAt.IsZero() {
		params.Set("startAt", strconv.FormatInt(startAt.Unix(), 10))
	}
	if !endAt.IsZero() {
		params.Set("endAt", strconv.FormatInt(endAt.Unix(), 10))
	}
	if currentPage != 0 {
		params.Set("currentPage", strconv.FormatInt(currentPage, 10))
	}
	if pageSize != 0 {
		params.Set("pageSize", strconv.FormatInt(pageSize, 10))
	}
	return resp.Data.Items, k.SendAuthHTTPRequest(ctx, exchange.RestSpot, http.MethodGet, common.EncodeURLValues(kucoinStopOrder, params), nil, publicSpotRate, &resp)
}

// GetStopOrderByClientID get a stop order information via the clientOID
func (k *Kucoin) GetStopOrderByClientID(ctx context.Context, symbol, clientOID string) ([]StopOrder, error) {
	resp := struct {
		Data []StopOrder `json:"data"`
		Error
	}{}
	//TODO: verify response

	params := url.Values{}
	if clientOID == "" {
		return resp.Data, errors.New("clientOID can't be empty")
	}
	params.Set("clientOid", clientOID)
	if symbol != "" {
		params.Set("symbol", symbol)
	}
	return resp.Data, k.SendAuthHTTPRequest(ctx, exchange.RestSpot, http.MethodGet, common.EncodeURLValues(kucoinGetStopOrderByClientID, params), nil, publicSpotRate, &resp)
}

// CancelStopOrderByClientID used to cancel a stop order via the clientOID.
func (k *Kucoin) CancelStopOrderByClientID(ctx context.Context, symbol, clientOID string) (string, string, error) {
	resp := struct {
		CancelledOrderID string `json:"cancelledOrderId"`
		ClientOID        string `json:"clientOid"`
		Error
	}{}

	params := url.Values{}
	if clientOID == "" {
		return resp.CancelledOrderID, resp.ClientOID, errors.New("clientOID can't be empty")
	}
	params.Set("clientOid", clientOID)
	if symbol != "" {
		params.Set("symbol", symbol)
	}
	return resp.CancelledOrderID, resp.ClientOID, k.SendAuthHTTPRequest(ctx, exchange.RestSpot, http.MethodDelete, common.EncodeURLValues(kucoinCancelStopOrderByClientID, params), nil, publicSpotRate, &resp)
}

// CreateAccount creates a account
func (k *Kucoin) CreateAccount(ctx context.Context, currency, accountType string) (string, error) {
	resp := struct {
		ID string `json:"id"`
		Error
	}{}

	params := make(map[string]interface{})
	if accountType == "" {
		return resp.ID, errors.New("accountType can't be empty")
	}
	params["type"] = accountType
	if currency == "" {
		return resp.ID, errors.New("currency can't be empty")
	}
	params["currency"] = currency
	return resp.ID, k.SendAuthHTTPRequest(ctx, exchange.RestSpot, http.MethodPost, kucoinAccount, params, publicSpotRate, &resp)
}

// GetAllAccounts get all accounts
func (k *Kucoin) GetAllAccounts(ctx context.Context, currency, accountType string) ([]AccountInfo, error) {
	resp := struct {
		Data []AccountInfo `json:"data"`
		Error
	}{}

	params := url.Values{}
	if currency != "" {
		params.Set("currency", currency)
	}
	if accountType != "" {
		params.Set("type", accountType)
	}
	return resp.Data, k.SendAuthHTTPRequest(ctx, exchange.RestSpot, http.MethodGet, common.EncodeURLValues(kucoinAccount, params), nil, publicSpotRate, &resp)
}

// GetAccount get information of single account
func (k *Kucoin) GetAccount(ctx context.Context, accountID string) (AccountInfo, error) {
	resp := struct {
		Data AccountInfo `json:"data"`
		Error
	}{}

	return resp.Data, k.SendAuthHTTPRequest(ctx, exchange.RestSpot, http.MethodGet, fmt.Sprintf(kucoinGetAccount, accountID), nil, publicSpotRate, &resp)
}

// GetAccountLedgers get the history of deposit/withdrawal of all accounts, supporting inquiry of various currencies
func (k *Kucoin) GetAccountLedgers(ctx context.Context, currency, direction, bizType string, startAt, endAt time.Time) ([]LedgerInfo, error) {
	resp := struct {
		Data struct {
			CurrentPage int64        `json:"currentPage"`
			PageSize    int64        `json:"pageSize"`
			TotalNum    int64        `json:"totalNum"`
			TotalPage   int64        `json:"totalPage"`
			Items       []LedgerInfo `json:"items"`
		} `json:"data"`
		Error
	}{}

	params := url.Values{}
	if currency != "" {
		params.Set("currency", currency)
	}
	if direction != "" {
		params.Set("direction", direction)
	}
	if bizType != "" {
		params.Set("bizType", bizType)
	}
	if !startAt.IsZero() {
		params.Set("startAt", strconv.FormatInt(startAt.UnixMilli(), 10))
	}
	if !endAt.IsZero() {
		params.Set("endAt", strconv.FormatInt(endAt.UnixMilli(), 10))
	}
	return resp.Data.Items, k.SendAuthHTTPRequest(ctx, exchange.RestSpot, http.MethodGet, common.EncodeURLValues(kucoinGetAccountLedgers, params), nil, publicSpotRate, &resp)
}

// GetSubAccountBalance get account info of a sub-user specified by the subUserID
func (k *Kucoin) GetSubAccountBalance(ctx context.Context, subUserID string) (SubAccountInfo, error) {
	resp := struct {
		Data SubAccountInfo `json:"data"`
		Error
	}{}

	return resp.Data, k.SendAuthHTTPRequest(ctx, exchange.RestSpot, http.MethodGet, fmt.Sprintf(kucoinGetSubAccountBalance, subUserID), nil, publicSpotRate, &resp)
}

// GetAggregatedSubAccountBalance get the account info of all sub-users
func (k *Kucoin) GetAggregatedSubAccountBalance(ctx context.Context) ([]SubAccountInfo, error) {
	resp := struct {
		Data []SubAccountInfo `json:"data"`
		Error
	}{}

	return resp.Data, k.SendAuthHTTPRequest(ctx, exchange.RestSpot, http.MethodGet, kucoinGetAggregatedSubAccountBalance, nil, publicSpotRate, &resp)
}

// GetTransferableBalance get the transferable balance of a specified account
func (k *Kucoin) GetTransferableBalance(ctx context.Context, currency, accountType, tag string) (TransferableBalanceInfo, error) {
	resp := struct {
		Data TransferableBalanceInfo `json:"data"`
		Error
	}{}

	params := url.Values{}
	if currency == "" {
		return resp.Data, errors.New("currency can't be empty")
	}
	params.Set("currency", currency)
	if accountType == "" {
		return resp.Data, errors.New("accountType can't be empty")
	}
	params.Set("type", accountType)
	if tag != "" {
		params.Set("tag", tag)
	}
	return resp.Data, k.SendAuthHTTPRequest(ctx, exchange.RestSpot, http.MethodGet, common.EncodeURLValues(kucoinGetTransferableBalance, params), nil, publicSpotRate, &resp)
}

// TransferMainToSubAccount used to transfer funds from main account to sub-account
func (k *Kucoin) TransferMainToSubAccount(ctx context.Context, clientOID, currency, amount, direction, accountType, subAccountType, subUserID string) (string, error) {
	resp := struct {
		Data struct {
			OrderID string `json:"orderId"`
		} `json:"data"`
		Error
	}{}

	params := make(map[string]interface{})
	if clientOID == "" {
		return resp.Data.OrderID, errors.New("clientOID can't be empty")
	}
	params["clientOid"] = clientOID
	if currency == "" {
		return resp.Data.OrderID, errors.New("currency can't be empty")
	}
	params["currency"] = currency
	if amount == "" {
		return resp.Data.OrderID, errors.New("amount can't be empty")
	}
	params["amount"] = amount
	if direction == "" {
		return resp.Data.OrderID, errors.New("direction can't be empty")
	}
	params["direction"] = direction
	if accountType != "" {
		params["accountType"] = accountType
	}
	if subAccountType != "" {
		params["subAccountType"] = subAccountType
	}
	if subUserID == "" {
		return resp.Data.OrderID, errors.New("subUserID can't be empty")
	}
	params["subUserId"] = subUserID
	return resp.Data.OrderID, k.SendAuthHTTPRequest(ctx, exchange.RestSpot, http.MethodPost, kucoinTransferMainToSubAccount, params, publicSpotRate, &resp)
}

// MakeInnerTransfer used to transfer funds between accounts internally
func (k *Kucoin) MakeInnerTransfer(ctx context.Context, clientOID, currency, from, to, amount, fromTag, toTag string) (string, error) {
	resp := struct {
		Data struct {
			OrderID string `json:"orderId"`
		} `json:"data"`
		Error
	}{}

	params := make(map[string]interface{})
	if clientOID == "" {
		return resp.Data.OrderID, errors.New("clientOID can't be empty")
	}
	params["clientOid"] = clientOID
	if currency == "" {
		return resp.Data.OrderID, errors.New("currency can't be empty")
	}
	params["currency"] = currency
	if amount == "" {
		return resp.Data.OrderID, errors.New("amount can't be empty")
	}
	params["amount"] = amount
	if from == "" {
		return resp.Data.OrderID, errors.New("from can't be empty")
	}
	params["from"] = from
	if to == "" {
		return resp.Data.OrderID, errors.New("to can't be empty")
	}
	params["to"] = to
	if fromTag != "" {
		params["fromTag"] = fromTag
	}
	if toTag != "" {
		params["toTag"] = toTag
	}
	return resp.Data.OrderID, k.SendAuthHTTPRequest(ctx, exchange.RestSpot, http.MethodPost, kucoinInnerTransfer, params, publicSpotRate, &resp)
}

// CreateDepositAddress create a deposit address for a currency you intend to deposit
func (k *Kucoin) CreateDepositAddress(ctx context.Context, currency, chain string) (DepositAddress, error) {
	resp := struct {
		Data DepositAddress `json:"data"`
		Error
	}{}

	params := make(map[string]interface{})
	if currency == "" {
		return resp.Data, errors.New("currency can't be empty")
	}
	params["currency"] = currency
	if chain != "" {
		params["chain"] = chain
	}
	return resp.Data, k.SendAuthHTTPRequest(ctx, exchange.RestSpot, http.MethodPost, kucoinCreateDepositAddress, params, publicSpotRate, &resp)
}

// GetDepositAddressV2 get all deposit addresses for the currency you intend to deposit
func (k *Kucoin) GetDepositAddressV2(ctx context.Context, currency string) ([]DepositAddress, error) {
	resp := struct {
		Data []DepositAddress `json:"data"`
		Error
	}{}

	params := url.Values{}
	if currency == "" {
		return resp.Data, errors.New("currency can't be empty")
	}
	params.Set("currency", currency)
	return resp.Data, k.SendAuthHTTPRequest(ctx, exchange.RestSpot, http.MethodGet, common.EncodeURLValues(kucoinGetDepositAddressV2, params), nil, publicSpotRate, &resp)
}

// GetDepositAddressV1 get a deposit address for the currency you intend to deposit
func (k *Kucoin) GetDepositAddressV1(ctx context.Context, currency, chain string) (DepositAddress, error) {
	resp := struct {
		Data DepositAddress `json:"data"`
		Error
	}{}

	params := url.Values{}
	if currency == "" {
		return resp.Data, errors.New("currency can't be empty")
	}
	params.Set("currency", currency)
	if chain != "" {
		params.Set("chain", chain)
	}
	return resp.Data, k.SendAuthHTTPRequest(ctx, exchange.RestSpot, http.MethodGet, common.EncodeURLValues(kucoinGetDepositAddressV1, params), nil, publicSpotRate, &resp)
}

// GetDepositList get deposit list items and sorted to show the latest first
func (k *Kucoin) GetDepositList(ctx context.Context, currency, status string, startAt, endAt time.Time) ([]Deposit, error) {
	resp := struct {
		Data struct {
			CurrentPage int64     `json:"currentPage"`
			PageSize    int64     `json:"pageSize"`
			TotalNum    int64     `json:"totalNum"`
			TotalPage   int64     `json:"totalPage"`
			Items       []Deposit `json:"items"`
		} `json:"data"`
		Error
	}{}
	params := url.Values{}
	if currency != "" {
		params.Set("currency", currency)
	}
	if status != "" {
		params.Set("status", status)
	}
	if !startAt.IsZero() {
		params.Set("startAt", strconv.FormatInt(startAt.UnixMilli(), 10))
	}
	if !endAt.IsZero() {
		params.Set("endAt", strconv.FormatInt(endAt.UnixMilli(), 10))
	}
	return resp.Data.Items, k.SendAuthHTTPRequest(ctx, exchange.RestSpot, http.MethodGet, common.EncodeURLValues(kucoinGetDepositList, params), nil, publicSpotRate, &resp)
}

// GetHistoricalDepositList get historical deposit list items
func (k *Kucoin) GetHistoricalDepositList(ctx context.Context, currency, status string, startAt, endAt time.Time) ([]HistoricalDepositWithdrawal, error) {
	resp := struct {
		Data struct {
			CurrentPage int64                         `json:"currentPage"`
			PageSize    int64                         `json:"pageSize"`
			TotalNum    int64                         `json:"totalNum"`
			TotalPage   int64                         `json:"totalPage"`
			Items       []HistoricalDepositWithdrawal `json:"items"`
		} `json:"data"`
		Error
	}{}
	params := url.Values{}
	if currency != "" {
		params.Set("currency", currency)
	}
	if status != "" {
		params.Set("status", status)
	}
	if !startAt.IsZero() {
		params.Set("startAt", strconv.FormatInt(startAt.UnixMilli(), 10))
	}
	if !endAt.IsZero() {
		params.Set("endAt", strconv.FormatInt(endAt.UnixMilli(), 10))
	}
	return resp.Data.Items, k.SendAuthHTTPRequest(ctx, exchange.RestSpot, http.MethodGet, common.EncodeURLValues(kucoinGetHistoricalDepositList, params), nil, publicSpotRate, &resp)
}

// GetWithdrawalList get withdrawal list items
func (k *Kucoin) GetWithdrawalList(ctx context.Context, currency, status string, startAt, endAt time.Time) ([]Withdrawal, error) {
	resp := struct {
		Data struct {
			CurrentPage int64        `json:"currentPage"`
			PageSize    int64        `json:"pageSize"`
			TotalNum    int64        `json:"totalNum"`
			TotalPage   int64        `json:"totalPage"`
			Items       []Withdrawal `json:"items"`
		} `json:"data"`
		Error
	}{}
	params := url.Values{}
	if currency != "" {
		params.Set("currency", currency)
	}
	if status != "" {
		params.Set("status", status)
	}
	if !startAt.IsZero() {
		params.Set("startAt", strconv.FormatInt(startAt.UnixMilli(), 10))
	}
	if !endAt.IsZero() {
		params.Set("endAt", strconv.FormatInt(endAt.UnixMilli(), 10))
	}
	return resp.Data.Items, k.SendAuthHTTPRequest(ctx, exchange.RestSpot, http.MethodGet, common.EncodeURLValues(kucoinWithdrawal, params), nil, publicSpotRate, &resp)
}

// GetHistoricalWithdrawalList get historical withdrawal list items
func (k *Kucoin) GetHistoricalWithdrawalList(ctx context.Context, currency, status string, startAt, endAt time.Time, currentPage, pageSize int64) ([]HistoricalDepositWithdrawal, error) {
	resp := struct {
		Data struct {
			CurrentPage int64                         `json:"currentPage"`
			PageSize    int64                         `json:"pageSize"`
			TotalNum    int64                         `json:"totalNum"`
			TotalPage   int64                         `json:"totalPage"`
			Items       []HistoricalDepositWithdrawal `json:"items"`
		} `json:"data"`
		Error
	}{}
	params := url.Values{}
	if currency != "" {
		params.Set("currency", currency)
	}
	if status != "" {
		params.Set("status", status)
	}
	if !startAt.IsZero() {
		params.Set("startAt", strconv.FormatInt(startAt.UnixMilli(), 10))
	}
	if !endAt.IsZero() {
		params.Set("endAt", strconv.FormatInt(endAt.UnixMilli(), 10))
	}
	if currentPage != 0 {
		params.Set("currentPage", strconv.FormatInt(currentPage, 10))
	}
	if pageSize != 0 {
		params.Set("pageSize", strconv.FormatInt(pageSize, 10))
	}
	return resp.Data.Items, k.SendAuthHTTPRequest(ctx, exchange.RestSpot, http.MethodGet, common.EncodeURLValues(kucoinGetHistoricalWithdrawalList, params), nil, publicSpotRate, &resp)
}

// GetWithdrawalQuotas get withdrawal quota details
func (k *Kucoin) GetWithdrawalQuotas(ctx context.Context, currency, chain string) (WithdrawalQuota, error) {
	resp := struct {
		Data WithdrawalQuota `json:"data"`
		Error
	}{}
	params := url.Values{}
	if currency == "" {
		return resp.Data, errors.New("currency can't be empty")
	}
	params.Set("currency", currency)
	if chain != "" {
		params.Set("chain", chain)
	}
	return resp.Data, k.SendAuthHTTPRequest(ctx, exchange.RestSpot, http.MethodGet, common.EncodeURLValues(kucoinGetWithdrawalQuotas, params), nil, publicSpotRate, &resp)
}

// ApplyWithdrawal create a withdrawal request
func (k *Kucoin) ApplyWithdrawal(ctx context.Context, currency, address, memo, remark, chain, feeDeductType string, isInner bool, amount float64) (string, error) {
	resp := struct {
		WithdrawalID string `json:"withdrawalId"`
		Error
	}{}

	params := make(map[string]interface{})
	if currency == "" {
		return resp.WithdrawalID, errors.New("currency can't be empty")
	}
	params["currency"] = currency
	if address == "" {
		return resp.WithdrawalID, errors.New("address can't be empty")
	}
	params["address"] = address
	if amount == 0 {
		return resp.WithdrawalID, errors.New("amount can't be empty")
	}
	params["amount"] = amount
	if memo != "" {
		params["memo"] = memo
	}
	params["isInner"] = isInner
	if remark != "" {
		params["remark"] = remark
	}
	if chain != "" {
		params["chain"] = chain
	}
	if feeDeductType != "" {
		params["feeDeductType"] = feeDeductType
	}
	return resp.WithdrawalID, k.SendAuthHTTPRequest(ctx, exchange.RestSpot, http.MethodPost, kucoinWithdrawal, params, publicSpotRate, &resp)
}

// CancelWithdrawal used to cancel a withdrawal request
func (k *Kucoin) CancelWithdrawal(ctx context.Context, withdrawalID string) error {
	resp := struct {
		Error
	}{}

	return k.SendAuthHTTPRequest(ctx, exchange.RestSpot, http.MethodDelete, fmt.Sprintf(kucoinCancelWithdrawal, withdrawalID), nil, publicSpotRate, &resp)
}

// GetBasicFee get basic fee rate of users
func (k *Kucoin) GetBasicFee(ctx context.Context, currencyType string) (Fees, error) {
	resp := struct {
		Data Fees `json:"data"`
		Error
	}{}

	params := url.Values{}
	if currencyType != "" {
		params.Set("currencyType", currencyType)
	}
	return resp.Data, k.SendAuthHTTPRequest(ctx, exchange.RestSpot, http.MethodGet, common.EncodeURLValues(kucoinBasicFee, params), nil, publicSpotRate, &resp)
}

// GetTradingFee get fee rate of trading pairs
func (k *Kucoin) GetTradingFee(ctx context.Context, symbols string) ([]Fees, error) {
	resp := struct {
		Data []Fees `json:"data"`
		Error
	}{}

	params := url.Values{}
	if symbols != "" {
		params.Set("symbols", symbols)
	}
	return resp.Data, k.SendAuthHTTPRequest(ctx, exchange.RestSpot, http.MethodGet, common.EncodeURLValues(kucoinTradingFee, params), nil, publicSpotRate, &resp)
}

// SendHTTPRequest sends an unauthenticated HTTP request
func (k *Kucoin) SendHTTPRequest(ctx context.Context, ePath exchange.URL, path string, f request.EndpointLimit, result UnmarshalTo) error {
	endpointPath, err := k.API.Endpoints.GetURL(ePath)
	if err != nil {
		return err
	}

	err = k.SendPayload(ctx, f, func() (*request.Item, error) {
		return &request.Item{
			Method:        http.MethodGet,
			Path:          endpointPath + path,
			Result:        result,
			Verbose:       k.Verbose,
			HTTPDebugging: k.HTTPDebugging,
			HTTPRecording: k.HTTPRecording}, nil
	})
	if err != nil {
		return err
	}
	return result.GetError()
}

// SendAuthHTTPRequest sends an authenticated HTTP request
// Request parameters are added to path variable for GET and DELETE request and for other requests its passed in params variable
func (k *Kucoin) SendAuthHTTPRequest(ctx context.Context, ePath exchange.URL, method, path string, params map[string]interface{}, f request.EndpointLimit, result UnmarshalTo) error {
	creds, err := k.GetCredentials(ctx)
	if err != nil {
		return err
	}

	endpointPath, err := k.API.Endpoints.GetURL(ePath)
	if err != nil {
		return err
	}

	if result == nil {
		result = &Error{}
	}

	err = k.SendPayload(ctx, f, func() (*request.Item, error) {
		var (
			body    io.Reader
			payload []byte
			err     error
		)
		if params != nil && len(params) != 0 {
			payload, err = json.Marshal(params)
			if err != nil {
				return nil, err
			}
			body = bytes.NewBuffer(payload)
		}
		timeStamp := strconv.FormatInt(time.Now().UnixMilli(), 10)
		var signHash, passPhraseHash []byte
		signHash, err = crypto.GetHMAC(crypto.HashSHA256, []byte(timeStamp+method+path+string(payload)), []byte(creds.Secret))
		if err != nil {
			return nil, err
		}

		passPhraseHash, err = crypto.GetHMAC(crypto.HashSHA256, []byte(creds.OneTimePassword), []byte(creds.Secret))
		if err != nil {
			return nil, err
		}
		headers := map[string]string{
			"KC-API-KEY":         creds.Key,
			"KC-API-SIGN":        crypto.Base64Encode(signHash),
			"KC-API-TIMESTAMP":   timeStamp,
			"KC-API-PASSPHRASE":  crypto.Base64Encode(passPhraseHash), // TODO: need pass phrase here!,
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
	if err != nil {
		return err
	}
	return result.GetError()
}