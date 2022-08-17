package kucoin

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"github.com/thrasher-corp/gocryptotrader/config"
	"github.com/thrasher-corp/gocryptotrader/currency"
	exchange "github.com/thrasher-corp/gocryptotrader/exchanges"
)

// Please supply your own keys here to do authenticated endpoint testing
const (
	apiKey                  = ""
	apiSecret               = ""
	passPhrase              = ""
	canManipulateRealOrders = false
)

var k Kucoin

func TestMain(m *testing.M) {
	k.SetDefaults()
	cfg := config.GetConfig()
	err := cfg.LoadConfig("../../testdata/configtest.json", true)
	if err != nil {
		log.Fatal(err)
	}

	exchCfg, err := cfg.GetExchangeConfig("Kucoin")
	if err != nil {
		log.Fatal(err)
	}

	exchCfg.API.AuthenticatedSupport = true
	exchCfg.API.AuthenticatedWebsocketSupport = true
	exchCfg.API.Credentials.Key = apiKey
	exchCfg.API.Credentials.Secret = apiSecret
	exchCfg.API.Credentials.OTPSecret = passPhrase
	err = k.Setup(exchCfg)
	if err != nil {
		log.Fatal(err)
	}

	os.Exit(m.Run())
}

// Ensures that this exchange package is compatible with IBotExchange
func TestInterface(t *testing.T) {
	var e exchange.IBotExchange
	if e = new(Kucoin); e == nil {
		t.Fatal("unable to allocate exchange")
	}
}

func areTestAPIKeysSet() bool {
	return k.ValidateAPICredentials(k.GetDefaultCredentials()) == nil
}

// Spot asset test cases starts from here
func TestGetSymbols(t *testing.T) {
	t.Parallel()

	_, err := k.GetSymbols(context.Background(), "")
	if err != nil {
		t.Error("GetSymbols() error", err)
	}

	_, err = k.GetSymbols(context.Background(), currency.BTC.String())
	if err != nil {
		t.Error("GetSymbols() error", err)
	}
}

func TestGetTicker(t *testing.T) {
	t.Parallel()

	_, err := k.GetTicker(context.Background(), "BTC-USDT")
	if err != nil {
		t.Error("GetTicker() error", err)
	}
}

func TestGetAllTickers(t *testing.T) {
	t.Parallel()

	_, err := k.GetAllTickers(context.Background())
	if err != nil {
		t.Error("GetAllTickers() error", err)
	}
}

func TestGet24hrStats(t *testing.T) {
	t.Parallel()

	_, err := k.Get24hrStats(context.Background(), "BTC-USDT")
	if err != nil {
		t.Error("Get24hrStats() error", err)
	}
}

func TestGetMarketList(t *testing.T) {
	t.Parallel()

	_, err := k.GetMarketList(context.Background())
	if err != nil {
		t.Error("GetMarketList() error", err)
	}
}

func TestGetPartOrderbook20(t *testing.T) {
	t.Parallel()

	_, err := k.GetPartOrderbook20(context.Background(), "BTC-USDT")
	if err != nil {
		t.Error("GetPartOrderbook20() error", err)
	}
}

func TestGetPartOrderbook100(t *testing.T) {
	t.Parallel()

	_, err := k.GetPartOrderbook100(context.Background(), "BTC-USDT")
	if err != nil {
		t.Error("GetPartOrderbook100() error", err)
	}
}

func TestGetOrderbook(t *testing.T) {
	t.Parallel()
	if !areTestAPIKeysSet() {
		t.Skip("skipping test: api keys not set")
	}

	_, err := k.GetOrderbook(context.Background(), "BTC-USDT")
	if err != nil {
		t.Error("GetOrderbook() error", err)
	}
}

func TestGetTradeHistory(t *testing.T) {
	t.Parallel()

	_, err := k.GetTradeHistory(context.Background(), "BTC-USDT")
	if err != nil {
		t.Error("GetTradeHistory() error", err)
	}
}

func TestGetKlines(t *testing.T) {
	t.Parallel()

	_, err := k.GetKlines(context.Background(), "BTC-USDT", "1week", time.Time{}, time.Time{})
	if err != nil {
		t.Error("GetKlines() error", err)
	}

	_, err = k.GetKlines(context.Background(), "BTC-USDT", "5min", time.Now().Add(-time.Hour*1), time.Now())
	if err != nil {
		t.Error("GetKlines() error", err)
	}
}

func TestGetCurrencies(t *testing.T) {
	t.Parallel()

	_, err := k.GetCurrencies(context.Background())
	if err != nil {
		t.Error("GetCurrencies() error", err)
	}
}

func TestGetCurrency(t *testing.T) {
	t.Parallel()

	_, err := k.GetCurrency(context.Background(), "BTC", "")
	if err != nil {
		t.Error("GetCurrency() error", err)
	}

	_, err = k.GetCurrency(context.Background(), "BTC", "ETH")
	if err != nil {
		t.Error("GetCurrency() error", err)
	}
}

func TestGetFiatPrice(t *testing.T) {
	t.Parallel()

	_, err := k.GetFiatPrice(context.Background(), "", "")
	if err != nil {
		t.Error("GetFiatPrice() error", err)
	}

	_, err = k.GetFiatPrice(context.Background(), "EUR", "ETH,BTC")
	if err != nil {
		t.Error("GetFiatPrice() error", err)
	}
}

func TestGetMarkPrice(t *testing.T) {
	t.Parallel()

	_, err := k.GetMarkPrice(context.Background(), "USDT-BTC")
	if err != nil {
		t.Error("GetMarkPrice() error", err)
	}
}

func TestGetMarginConfiguration(t *testing.T) {
	t.Parallel()

	_, err := k.GetMarginConfiguration(context.Background())
	if err != nil {
		t.Error("GetMarginConfiguration() error", err)
	}
}

func TestGetMarginAccount(t *testing.T) {
	t.Parallel()
	if !areTestAPIKeysSet() {
		t.Skip("skipping test: api keys not set")
	}

	_, err := k.GetMarginAccount(context.Background())
	if err != nil {
		t.Error("GetMarginAccount() error", err)
	}
}

func TestGetMarginRiskLimit(t *testing.T) {
	t.Parallel()
	if !areTestAPIKeysSet() {
		t.Skip("skipping test: api keys not set")
	}

	_, err := k.GetMarginRiskLimit(context.Background(), "cross")
	if err != nil {
		t.Error("GetMarginRiskLimit() error", err)
	}

	_, err = k.GetMarginRiskLimit(context.Background(), "isolated")
	if err != nil {
		t.Error("GetMarginRiskLimit() error", err)
	}
}

func TestPostBorrowOrder(t *testing.T) {
	t.Parallel()
	if !areTestAPIKeysSet() || !canManipulateRealOrders {
		t.Skip("skipping test: api keys not set or canManipulateRealOrders set to false")
	}

	_, err := k.PostBorrowOrder(context.Background(), "USDT", "FOK", "", 10, 0)
	if err != nil {
		t.Error("PostBorrowOrder() error", err)
	}

	_, err = k.PostBorrowOrder(context.Background(), "USDT", "IOC", "7,14,28", 10, 10)
	if err != nil {
		t.Error("PostBorrowOrder() error", err)
	}
}

func TestGetBorrowOrder(t *testing.T) {
	t.Parallel()
	if !areTestAPIKeysSet() {
		t.Skip("skipping test: api keys not set")
	}

	_, err := k.GetBorrowOrder(context.Background(), "orderID")
	if err != nil && err.Error() != "Not Found" {
		t.Error("GetBorrowOrder() error", err)
	}
}

func TestGetOutstandingRecord(t *testing.T) {
	t.Parallel()
	if !areTestAPIKeysSet() {
		t.Skip("skipping test: api keys not set")
	}

	_, err := k.GetOutstandingRecord(context.Background(), "BTC")
	if err != nil {
		t.Error("GetOutstandingRecord() error", err)
	}
}

func TestGetRepaidRecord(t *testing.T) {
	t.Parallel()
	if !areTestAPIKeysSet() {
		t.Skip("skipping test: api keys not set")
	}

	_, err := k.GetRepaidRecord(context.Background(), "BTC")
	if err != nil {
		t.Error("GetRepaidRecord() error", err)
	}
}

func TestOneClickRepayment(t *testing.T) {
	t.Parallel()
	if !areTestAPIKeysSet() || !canManipulateRealOrders {
		t.Skip("skipping test: api keys not set or canManipulateRealOrders set to false")
	}

	err := k.OneClickRepayment(context.Background(), "BTC", "RECENTLY_EXPIRE_FIRST", 2.5)
	if err != nil && err.Error() != "Balance insufficient" {
		t.Error("OneClickRepayment() error", err)
	}
}

func TestSingleOrderRepayment(t *testing.T) {
	t.Parallel()
	if !areTestAPIKeysSet() || !canManipulateRealOrders {
		t.Skip("skipping test: api keys not set or canManipulateRealOrders set to false")
	}

	err := k.SingleOrderRepayment(context.Background(), "BTC", "fa3e34c980062c10dad74016", 2.5)
	if err != nil && err.Error() != "Balance insufficient" {
		t.Error("SingleOrderRepayment() error", err)
	}
}

func TestPostLendOrder(t *testing.T) {
	t.Parallel()
	if !areTestAPIKeysSet() || !canManipulateRealOrders {
		t.Skip("skipping test: api keys not set or canManipulateRealOrders set to false")
	}

	_, err := k.PostLendOrder(context.Background(), "BTC", 0.001, 5, 7)
	if err != nil && err.Error() != "Balance insufficient" {
		t.Error("PostLendOrder() error", err)
	}
}

func TestCancelLendOrder(t *testing.T) {
	t.Parallel()
	if !areTestAPIKeysSet() || !canManipulateRealOrders {
		t.Skip("skipping test: api keys not set or canManipulateRealOrders set to false")
	}

	err := k.CancelLendOrder(context.Background(), "OrderID")
	if err != nil && err.Error() != "order not exist" {
		t.Error("CancelLendOrder() error", err)
	}
}

func TestSetAutoLend(t *testing.T) {
	t.Parallel()
	if !areTestAPIKeysSet() || !canManipulateRealOrders {
		t.Skip("skipping test: api keys not set or canManipulateRealOrders set to false")
	}

	err := k.SetAutoLend(context.Background(), "BTC", 0.002, 0.005, 7, true)
	if err != nil {
		t.Error("SetAutoLend() error", err)
	}
}

func TestGetActiveOrder(t *testing.T) {
	t.Parallel()
	if !areTestAPIKeysSet() {
		t.Skip("skipping test: api keys not set")
	}

	_, err := k.GetActiveOrder(context.Background(), "")
	if err != nil {
		t.Error("GetActiveOrder() error", err)
	}

	_, err = k.GetActiveOrder(context.Background(), "BTC")
	if err != nil {
		t.Error("GetActiveOrder() error", err)
	}
}

func TestGetLendHistory(t *testing.T) {
	t.Parallel()
	if !areTestAPIKeysSet() {
		t.Skip("skipping test: api keys not set")
	}

	_, err := k.GetLendHistory(context.Background(), "")
	if err != nil {
		t.Error("GetLendHistory() error", err)
	}

	_, err = k.GetLendHistory(context.Background(), "BTC")
	if err != nil {
		t.Error("GetLendHistory() error", err)
	}
}

func TestGetUnsettleLendOrder(t *testing.T) {
	t.Parallel()
	if !areTestAPIKeysSet() {
		t.Skip("skipping test: api keys not set")
	}

	_, err := k.GetUnsettleLendOrder(context.Background(), "")
	if err != nil {
		t.Error("GetUnsettleLendOrder() error", err)
	}

	_, err = k.GetUnsettleLendOrder(context.Background(), "BTC")
	if err != nil {
		t.Error("GetUnsettleLendOrder() error", err)
	}
}

func TestGetSettleLendOrder(t *testing.T) {
	t.Parallel()
	if !areTestAPIKeysSet() {
		t.Skip("skipping test: api keys not set")
	}

	_, err := k.GetSettleLendOrder(context.Background(), "")
	if err != nil {
		t.Error("GetSettleLendOrder() error", err)
	}

	_, err = k.GetSettleLendOrder(context.Background(), "BTC")
	if err != nil {
		t.Error("GetSettleLendOrder() error", err)
	}
}

func TestGetAccountLendRecord(t *testing.T) {
	t.Parallel()
	if !areTestAPIKeysSet() {
		t.Skip("skipping test: api keys not set")
	}

	_, err := k.GetAccountLendRecord(context.Background(), "")
	if err != nil {
		t.Error("GetAccountLendRecord() error", err)
	}

	_, err = k.GetAccountLendRecord(context.Background(), "BTC")
	if err != nil {
		t.Error("GetAccountLendRecord() error", err)
	}
}

func TestGetLendingMarketData(t *testing.T) {
	t.Parallel()
	if !areTestAPIKeysSet() {
		t.Skip("skipping test: api keys not set")
	}

	_, err := k.GetLendingMarketData(context.Background(), "BTC", 0)
	if err != nil {
		t.Error("GetLendingMarketData() error", err)
	}

	_, err = k.GetLendingMarketData(context.Background(), "BTC", 7)
	if err != nil {
		t.Error("GetLendingMarketData() error", err)
	}
}

func TestGetMarginTradeData(t *testing.T) {
	t.Parallel()
	if !areTestAPIKeysSet() {
		t.Skip("skipping test: api keys not set")
	}

	_, err := k.GetMarginTradeData(context.Background(), "BTC")
	if err != nil {
		t.Error("GetMarginTradeData() error", err)
	}
}

func TestGetIsolatedMarginPairConfig(t *testing.T) {
	t.Parallel()
	if !areTestAPIKeysSet() {
		t.Skip("skipping test: api keys not set")
	}

	_, err := k.GetIsolatedMarginPairConfig(context.Background())
	if err != nil {
		t.Error("GetIsolatedMarginPairConfig() error", err)
	}
}

func TestGetIsolatedMarginAccountInfo(t *testing.T) {
	t.Parallel()
	if !areTestAPIKeysSet() {
		t.Skip("skipping test: api keys not set")
	}

	_, err := k.GetIsolatedMarginAccountInfo(context.Background(), "")
	if err != nil {
		t.Error("GetIsolatedMarginAccountInfo() error", err)
	}

	_, err = k.GetIsolatedMarginAccountInfo(context.Background(), "USDT")
	if err != nil {
		t.Error("GetIsolatedMarginAccountInfo() error", err)
	}
}

func TestGetSingleIsolatedMarginAccountInfo(t *testing.T) {
	t.Parallel()
	if !areTestAPIKeysSet() {
		t.Skip("skipping test: api keys not set")
	}

	_, err := k.GetSingleIsolatedMarginAccountInfo(context.Background(), "BTC-USDT")
	if err != nil {
		t.Error("GetSingleIsolatedMarginAccountInfo() error", err)
	}
}

func TestInitiateIsolateMarginBorrowing(t *testing.T) {
	t.Parallel()
	if !areTestAPIKeysSet() || !canManipulateRealOrders {
		t.Skip("skipping test: api keys not set or canManipulateRealOrders set to false")
	}

	_, _, _, err := k.InitiateIsolateMarginBorrowing(context.Background(), "BTC-USDT", "USDT", "FOK", "", 10, 0)
	if err != nil {
		t.Error("InitiateIsolateMarginBorrowing() error", err)
	}
}

func TestGetIsolatedOutstandingRepaymentRecords(t *testing.T) {
	t.Parallel()
	if !areTestAPIKeysSet() {
		t.Skip("skipping test: api keys not set")
	}

	_, err := k.GetIsolatedOutstandingRepaymentRecords(context.Background(), "", "", 0, 0)
	if err != nil {
		t.Error("GetIsolatedOutstandingRepaymentRecords() error", err)
	}

	_, err = k.GetIsolatedOutstandingRepaymentRecords(context.Background(), "BTC-USDT", "USDT", 0, 0)
	if err != nil {
		t.Error("GetIsolatedOutstandingRepaymentRecords() error", err)
	}
}

func TestGetIsolatedMarginRepaymentRecords(t *testing.T) {
	t.Parallel()
	if !areTestAPIKeysSet() {
		t.Skip("skipping test: api keys not set")
	}

	_, err := k.GetIsolatedMarginRepaymentRecords(context.Background(), "", "", 0, 0)
	if err != nil {
		t.Error("GetIsolatedMarginRepaymentRecords() error", err)
	}

	_, err = k.GetIsolatedMarginRepaymentRecords(context.Background(), "BTC-USDT", "USDT", 0, 0)
	if err != nil {
		t.Error("GetIsolatedMarginRepaymentRecords() error", err)
	}
}

func TestInitiateIsolatedMarginQuickRepayment(t *testing.T) {
	t.Parallel()
	if !areTestAPIKeysSet() || !canManipulateRealOrders {
		t.Skip("skipping test: api keys not set or canManipulateRealOrders set to false")
	}

	err := k.InitiateIsolatedMarginQuickRepayment(context.Background(), "BTC-USDT", "USDT", "RECENTLY_EXPIRE_FIRST", 10)
	if err != nil {
		t.Error("InitiateIsolatedMarginQuickRepayment() error", err)
	}
}

func TestInitiateIsolatedMarginSingleRepayment(t *testing.T) {
	t.Parallel()
	if !areTestAPIKeysSet() || !canManipulateRealOrders {
		t.Skip("skipping test: api keys not set or canManipulateRealOrders set to false")
	}

	err := k.InitiateIsolatedMarginSingleRepayment(context.Background(), "BTC-USDT", "USDT", "628c570f7818320001d52b69", 10)
	if err != nil {
		t.Error("InitiateIsolatedMarginSingleRepayment() error", err)
	}
}

func TestGetCurrentServerTime(t *testing.T) {
	t.Parallel()

	_, err := k.GetCurrentServerTime(context.Background())
	if err != nil {
		t.Error("GetCurrentServerTime() error", err)
	}
}

func TestGetServiceStatus(t *testing.T) {
	t.Parallel()

	_, _, err := k.GetServiceStatus(context.Background())
	if err != nil {
		t.Error("GetServiceStatus() error", err)
	}
}

func TestPostOrder(t *testing.T) {
	t.Parallel()
	if !areTestAPIKeysSet() || !canManipulateRealOrders {
		t.Skip("skipping test: api keys not set or canManipulateRealOrders set to false")
	}

	// default order type is limit
	_, err := k.PostOrder(context.Background(), "5bd6e9286d99522a52e458de", "buy", "BTC-USDT", "", "", "", "10000", "", 0.1, 0, 0, 0, true, false, false)
	if err != nil && err.Error() != "Balance insufficient!" {
		t.Error("PostOrder() error", err)
	}

	// market order
	_, err = k.PostOrder(context.Background(), "5bd6e9286d99522a52e458de", "buy", "BTC-USDT", "market", "remark", "", "", "", 0.1, 0, 0, 0, true, false, false)
	if err != nil && err.Error() != "Balance insufficient!" {
		t.Error("PostOrder() error", err)
	}
}

func TestPostMarginOrder(t *testing.T) {
	t.Parallel()
	if !areTestAPIKeysSet() || !canManipulateRealOrders {
		t.Skip("skipping test: api keys not set or canManipulateRealOrders set to false")
	}

	// default order type is limit and margin mode is cross
	_, err := k.PostMarginOrder(context.Background(), "5bd6e9286d99522a52e458de", "buy", "BTC-USDT", "", "", "", "", "10000", "", 0.1, 0, 0, 0, true, false, false, false)
	if err != nil && err.Error() != "Balance insufficient!" {
		t.Error("PostMarginOrder() error", err)
	}

	// market isolated order
	_, err = k.PostMarginOrder(context.Background(), "5bd6e9286d99522a52e458de", "buy", "BTC-USDT", "market", "remark", "", "isolated", "", "", 0.1, 0, 0, 5, true, false, false, true)
	if err != nil && err.Error() != "Balance insufficient!" {
		t.Error("PostMarginOrder() error", err)
	}
}

func TestPostBulkOrder(t *testing.T) {
	t.Parallel()
	if !areTestAPIKeysSet() || !canManipulateRealOrders {
		t.Skip("skipping test: api keys not set or canManipulateRealOrders set to false")
	}

	req := []OrderRequest{
		{
			ClientOID: "3d07008668054da6b3cb12e432c2b13a",
			Side:      "buy",
			Type:      "limit",
			Price:     "1000",
			Size:      "0.01",
		},
		{
			ClientOID: "37245dbe6e134b5c97732bfb36cd4a9d",
			Side:      "buy",
			Type:      "limit",
			Price:     "1000",
			Size:      "0.01",
		},
	}

	_, err := k.PostBulkOrder(context.Background(), "BTC-USDT", req)
	if err != nil {
		t.Error("PostBulkOrder() error", err)
	}
}

func TestCancelSingleOrder(t *testing.T) {
	t.Parallel()
	if !areTestAPIKeysSet() || !canManipulateRealOrders {
		t.Skip("skipping test: api keys not set or canManipulateRealOrders set to false")
	}

	_, err := k.CancelSingleOrder(context.Background(), "5bd6e9286d99522a52e458de")
	if err != nil && err.Error() != "order_not_exist_or_not_allow_to_cancel" {
		t.Error("CancelSingleOrder() error", err)
	}
}

func TestCancelOrderByClientOID(t *testing.T) {
	t.Parallel()
	if !areTestAPIKeysSet() || !canManipulateRealOrders {
		t.Skip("skipping test: api keys not set or canManipulateRealOrders set to false")
	}

	_, _, err := k.CancelOrderByClientOID(context.Background(), "5bd6e9286d99522a52e458de")
	if err != nil && err.Error() != "order_not_exist_or_not_allow_to_cancel" {
		t.Error("CancelOrderByClientOID() error", err)
	}
}

func TestCancelAllOpenOrders(t *testing.T) {
	t.Parallel()
	if !areTestAPIKeysSet() || !canManipulateRealOrders {
		t.Skip("skipping test: api keys not set or canManipulateRealOrders set to false")
	}

	_, err := k.CancelAllOpenOrders(context.Background(), "", "")
	if err != nil {
		t.Error("CancelAllOpenOrders() error", err)
	}
}

func TestGetOrders(t *testing.T) {
	t.Parallel()
	if !areTestAPIKeysSet() {
		t.Skip("skipping test: api keys not set")
	}

	_, err := k.GetOrders(context.Background(), "", "", "", "", "", time.Time{}, time.Time{})
	if err != nil {
		t.Error("GetOrders() error", err)
	}
}

// TODO: ambiguity in doc. and API response
func TestGetRecentOrders(t *testing.T) {
	t.Parallel()
	if !areTestAPIKeysSet() {
		t.Skip("skipping test: api keys not set")
	}

	_, err := k.GetRecentOrders(context.Background())
	if err != nil {
		t.Error("GetRecentOrders() error", err)
	}
}

// TODO: not sure of response after looking at doc.
func TestGetOrderByID(t *testing.T) {
	t.Parallel()
	if !areTestAPIKeysSet() {
		t.Skip("skipping test: api keys not set")
	}

	_, err := k.GetOrderByID(context.Background(), "5c35c02703aa673ceec2a168")
	if err != nil && err.Error() != "order not exist." {
		t.Error("GetOrderByID() error", err)
	}
}

func TestGetOrderByClientOID(t *testing.T) {
	t.Parallel()
	if !areTestAPIKeysSet() {
		t.Skip("skipping test: api keys not set")
	}

	_, err := k.GetOrderByClientOID(context.Background(), "6d539dc614db312")
	if err != nil && err.Error() != "order not exist." {
		t.Error("GetOrderByClientOID() error", err)
	}
}

func TestGetFills(t *testing.T) {
	t.Parallel()
	if !areTestAPIKeysSet() {
		t.Skip("skipping test: api keys not set")
	}

	_, err := k.GetFills(context.Background(), "", "", "", "", "", time.Time{}, time.Time{})
	if err != nil {
		t.Error("GetFills() error", err)
	}

	_, err = k.GetFills(context.Background(), "5c35c02703aa673ceec2a168", "BTC-USDT", "buy", "limit", "TRADE", time.Now().Add(-time.Hour*12), time.Now())
	if err != nil {
		t.Error("GetFills() error", err)
	}
}

func TestGetRecentFills(t *testing.T) {
	t.Parallel()
	if !areTestAPIKeysSet() {
		t.Skip("skipping test: api keys not set")
	}

	_, err := k.GetRecentFills(context.Background())
	if err != nil {
		t.Error("GetRecentFills() error", err)
	}
}

func TestPostStopOrder(t *testing.T) {
	t.Parallel()
	if !areTestAPIKeysSet() || !canManipulateRealOrders {
		t.Skip("skipping test: api keys not set or canManipulateRealOrders set to false")
	}

	_, err := k.PostStopOrder(context.Background(), "5bd6e9286d99522a52e458de", "buy", "BTC-USDT", "", "", "entry", "10000", "11000", "", "", "", 0.1, 0, 0, 0, true, false, false)
	if err != nil {
		t.Error("PostStopOrder() error", err)
	}
}

func TestCancelStopOrder(t *testing.T) {
	t.Parallel()
	if !areTestAPIKeysSet() || !canManipulateRealOrders {
		t.Skip("skipping test: api keys not set or canManipulateRealOrders set to false")
	}

	_, err := k.CancelStopOrder(context.Background(), "5bd6e9286d99522a52e458de")
	if err != nil && err.Error() != "order_not_exist_or_not_allow_to_cancel" {
		t.Error("CancelStopOrder() error", err)
	}
}

func TestCancelAllStopOrder(t *testing.T) {
	t.Parallel()
	if !areTestAPIKeysSet() || !canManipulateRealOrders {
		t.Skip("skipping test: api keys not set or canManipulateRealOrders set to false")
	}

	_, err := k.CancelAllStopOrder(context.Background(), "", "", "")
	if err != nil {
		t.Error("CancelAllStopOrder() error", err)
	}
}

func TestGetStopOrder(t *testing.T) {
	t.Parallel()
	if !areTestAPIKeysSet() {
		t.Skip("skipping test: api keys not set")
	}

	_, err := k.GetStopOrder(context.Background(), "5bd6e9286d99522a52e458de")
	if err != nil {
		t.Error("GetStopOrder() error", err)
	}
}

func TestGetAllStopOrder(t *testing.T) {
	t.Parallel()
	if !areTestAPIKeysSet() {
		t.Skip("skipping test: api keys not set")
	}

	_, err := k.GetAllStopOrder(context.Background(), "", "", "", "", "", time.Time{}, time.Time{}, 0, 0)
	if err != nil {
		t.Error("GetAllStopOrder() error", err)
	}
}

func TestGetStopOrderByClientID(t *testing.T) {
	t.Parallel()
	if !areTestAPIKeysSet() {
		t.Skip("skipping test: api keys not set")
	}

	_, err := k.GetStopOrderByClientID(context.Background(), "", "5bd6e9286d99522a52e458de")
	if err != nil {
		t.Error("GetStopOrderByClientID() error", err)
	}
}

func TestCancelStopOrderByClientID(t *testing.T) {
	t.Parallel()
	if !areTestAPIKeysSet() {
		t.Skip("skipping test: api keys not set")
	}

	_, _, err := k.CancelStopOrderByClientID(context.Background(), "", "5bd6e9286d99522a52e458de")
	if err != nil && err.Error() != "order_not_exist_or_not_allow_to_cancel" {
		t.Error("CancelStopOrderByClientID() error", err)
	}
}

func TestCreateAccount(t *testing.T) {
	t.Parallel()
	if !areTestAPIKeysSet() || !canManipulateRealOrders {
		t.Skip("skipping test: api keys not set or canManipulateRealOrders set to false")
	}

	_, err := k.CreateAccount(context.Background(), "BTC", "main")
	if err != nil {
		t.Error("CreateAccount() error", err)
	}
}

func TestGetAllAccounts(t *testing.T) {
	t.Parallel()
	if !areTestAPIKeysSet() {
		t.Skip("skipping test: api keys not set")
	}

	_, err := k.GetAllAccounts(context.Background(), "", "")
	if err != nil {
		t.Error("GetAllAccounts() error", err)
	}
}
