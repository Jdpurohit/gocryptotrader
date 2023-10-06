package apex

import (
	"fmt"
	"log"
	"math/big"
	"os"
	"testing"

	"github.com/thrasher-corp/gocryptotrader/config"
	exchange "github.com/thrasher-corp/gocryptotrader/exchanges"
)

// Please supply your own keys here to do authenticated endpoint testing
const (
	apiKey                  = "a22282a1-1c9d-3c9b-defa-5f7d359dc72f"
	apiSecret               = "rh8NEXm6rbP3xn2yMOUrGGvFip0uC79mQ63vXMP-"
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
	exchCfg.Verbose = true
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
/*
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

*/
func TestGenerateNonce(t *testing.T) {
	t.Parallel()
	if !areTestAPIKeysSet() {
		t.Skip("skipping test: api keys not set")
	}
	// _, err := ap.GenerateNonce(context.Background(), "0xbe7b1BE F4b9ce7A1C0A24A243F3b559C4f7Bd084", "0x7cbeead71aeec75d2bbb869a9baa1c1241a8d3a90dde8ea304f97f3a594672d", "1")
	// if err != nil {
	// 	t.Fatal(err)
	// }

	// x, y, err := deriveStarkKey()
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// fmt.Printf("deriveStarkKey: \n\nX: %s \nY: %s\n", x, y)

	// a, _ := new(big.Int).SetString("1959202139686649971135379195946240302532837600025877174746757219065680349137", 10)
	// b, _ := new(big.Int).SetString("2964034215857088791696663278430423156886956162064077437109373246614297849484", 10)
	// c, _ := new(big.Int).SetString("1", 10)
	// d, _ := new(big.Int).SetString("3618502788666131213697322783095070105623107215331596699973092056135872020481", 10)
	// point := ecDouble([2]*big.Int{a, b}, c, d)
	// fmt.Printf("\nX: %s Y: %s\n", point[0].String(), point[1].String())
	// X: 1858070200429287878750589068775421330844900027307286517151728242584957236426
	// Y: 1227026973562978579584701011068063615621182824056995865242742798017265498284

	a, _ := new(big.Int).SetString("233685396261341212670245012040713484489473139098736409389704419900786264492", 10)
	b, _ := new(big.Int).SetString("2448737574461136843506220627791380048993560238917253908122263523840710031760", 10)
	c, _ := new(big.Int).SetString("3136030469135674343172465880817263454880219855664441593466904169223571314065", 10)
	d, _ := new(big.Int).SetString("3230850854683103635133032411878658931556916918508772276704988424959453909526", 10)
	e, _ := new(big.Int).SetString("3618502788666131213697322783095070105623107215331596699973092056135872020481", 10)
	point := ecAdd([2]*big.Int{a, b}, [2]*big.Int{c, d}, e)
	fmt.Printf("\nX: %s \nY: %s\n", point[0].String(), point[1].String())

	// m := new(big.Int)
	// a, b := new(big.Int).DivMod(x, y, m)
	// fmt.Printf("A: %d B: %d", a, b)
	// // fmt.Println("Encoded Hex String: ", encodedString)
	// fmt.Println("1: ", hex.EncodeToString(getHashString("ApeX(string action,string onlySignOn,string nonce)")))
	// fmt.Println("2: ", hex.EncodeToString(getHashString("ApeX Onboarding")))
	// fmt.Println("3: ", hex.EncodeToString(getHashString("https://pro.apex.exchange")))
	// fmt.Println("4: ", hex.EncodeToString(getHashString("1517566403329392640")))

	// structHash := solsha3.SoliditySHA3(
	// 	[]string{"bytes32", "bytes32", "bytes32", "bytes32"},
	// 	[]interface{}{
	// 		getHashString("ApeX(string action,string onlySignOn,string nonce)"),
	// 		getHashString("ApeX Onboarding"),
	// 		getHashString("https://pro.apex.exchange"),
	// 		getHashString("1517566403329392640"),
	// 	},
	// )

	// // fmt.Println("5: ", hex.EncodeToString(getDomainHash()))

	// fmt.Println("5: ", hex.EncodeToString(solsha3.SoliditySHA3(
	// 	[]string{"bytes2", "bytes32", "bytes32"},
	// 	[]interface{}{
	// 		"0x1901",
	// 		getDomainHash(),
	// 		structHash,
	// 	},
	// )))
}

/*
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
*/
