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
