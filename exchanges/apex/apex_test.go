package apex

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"github.com/thrasher-corp/gocryptotrader/config"
	exchange "github.com/thrasher-corp/gocryptotrader/exchanges"
)

// Please supply your own keys here to do authenticated endpoint testing
const (
	apiKey                  = ""
	apiSecret               = ""
	canManipulateRealOrders = false
)

var ap Apex

func TestMain(m *testing.M) {
	ap.SetDefaults()
	cfg := config.GetConfig()
	err := cfg.LoadConfig("../../testdata/configtest.json", true)
	if err != nil {
		log.Fatal(err)
	}

	exchCfg, err := cfg.GetExchangeConfig("Apex")
	if err != nil {
		log.Fatal(err)
	}

	exchCfg.API.AuthenticatedSupport = true
	exchCfg.API.AuthenticatedWebsocketSupport = true
	exchCfg.API.Credentials.Key = apiKey
	exchCfg.API.Credentials.Secret = apiSecret

	err = ap.Setup(exchCfg)
	if err != nil {
		log.Fatal(err)
	}

	os.Exit(m.Run())
}

// Ensures that this exchange package is compatible with IBotExchange
func TestInterface(t *testing.T) {
	var e exchange.IBotExchange
	if e = new(Apex); e == nil {
		t.Fatal("unable to allocate exchange")
	}
}

func areTestAPIKeysSet() bool {
	return ap.ValidateAPICredentials(ap.GetDefaultCredentials()) == nil
}

// Implement tests for API endpoints below

func TestGetSystemTime(t *testing.T) {
	t.Parallel()
	_, err := ap.GetSystemTime(context.Background())
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetAllConfigData(t *testing.T) {
	t.Parallel()
	_, err := ap.GetAllConfig(context.Background())
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetMarketDepthData(t *testing.T) {
	t.Parallel()
	_, err := ap.GetMarketDepth(context.Background(), "BTCUSDC", 0)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetLatestTrades(t *testing.T) {
	t.Parallel()
	_, err := ap.GetLatestTrades(context.Background(), "BTCUSDC", 0, 0)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetCandlestickChart(t *testing.T) {
	t.Parallel()
	_, err := ap.GetCandlestickChart(context.Background(), "D", "BTCUSDC", time.Time{}, time.Time{}, 0)
	if err != nil {
		t.Fatal(err)
	}

	_, err = ap.GetCandlestickChart(context.Background(), "", "BTCUSDC", time.Now().Add(-time.Hour*24), time.Now(), 0)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetTicker(t *testing.T) {
	t.Parallel()
	_, err := ap.GetTicker(context.Background(), "BTCUSDC")
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetFundingRateHistory(t *testing.T) {
	t.Parallel()
	_, err := ap.GetFundingRateHistory(context.Background(), "BTCUSDC", time.Time{}, time.Time{}, 0, -1)
	if err != nil {
		t.Fatal(err)
	}
}

func TestCheckIfUserExists(t *testing.T) {
	t.Parallel()
	_, err := ap.CheckIfUserExists(context.Background(), "")
	if err != nil {
		t.Fatal(err)
	}
}

func TestGenerateNonce(t *testing.T) {
	t.Parallel()
	if !areTestAPIKeysSet() {
		t.Skip("skipping test: api keys not set")
	}
	_, err := ap.GenerateNonce(context.Background(), "0x4315c720e1c256A800B93c1742a6525fF40aB7C5", "0x7c9fec5834aaa1e30143544ee0d8ed91025d1336bb188d57592d5e64e5b7c5f", "1")
	if err != nil {
		t.Fatal(err)
	}
}

func TestRegistration(t *testing.T) {
	t.Parallel()
	if !areTestAPIKeysSet() {
		t.Skip("skipping test: api keys not set")
	}
	_, err := ap.Registration(context.Background(), "0x7c9fec5834aaa1e30143544ee0d8ed91025d1336bb188d57592d5e64e5b7c5f", "0x2d99d8e5060171bac631b7efd7d97464fc98b0efcee87aef3a9eca1f965b569", "0x4315c720e1c256A800B93c1742a6525fF40aB7C5", "", "")
	if err != nil {
		t.Fatal(err)
	}
}
