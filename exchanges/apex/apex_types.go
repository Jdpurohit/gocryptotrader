package apex

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/thrasher-corp/gocryptotrader/exchanges/orderbook"
)

var validIntervals = []string{
	"1", "3", "5", "15", "30", "60", "120", "240", "360", "720",
	"D", "M", "W",
}

//1 3 5 15 30 + 60 120 240 360 720 + D + W + M
var (
	errInvalidInterval = errors.New("invalid interval")
)

// UnmarshalTo acts as interface to exchange API response
type UnmarshalTo interface {
	GetError() error
}

// Error defines all error information for each request
type Error struct {
	ReturnCode int64  `json:"code"`
	ReturnMsg  string `json:"msg"`
}

// GetError checks and returns an error if it is supplied.
func (e Error) GetError() error {
	if e.ReturnCode != 0 && e.ReturnMsg != "" {
		return errors.New(e.ReturnMsg)
	}
	return nil
}

// apexTimeMilliSec provides an internal conversion helper
type apexTimeMilliSec time.Time

// UnmarshalJSON is custom type json unmarshaller for apexTimeMilliSec
func (ap *apexTimeMilliSec) UnmarshalJSON(data []byte) error {
	var timestamp int64
	err := json.Unmarshal(data, &timestamp)
	if err != nil {
		return err
	}
	*ap = apexTimeMilliSec(time.UnixMilli(timestamp))
	return nil
}

// Time returns a time.Time object
func (b apexTimeMilliSec) Time() time.Time {
	return time.Time(b)
}

type Config struct {
	Currency          []CurrencyConfig          `json:"currency"`
	Global            GlobalConfig              `json:"global"`
	PerpetualContract []PerpetualContractConfig `json:"perpetualContract"`
	MultiChain        struct {
		Chains      []ChainConfig `json:"chains"`
		Currency    string        `json:"currency"`
		MaxWithdraw int64         `json:"maxWithdraw,string"`
		MinDeposit  int64         `json:"minDeposit,string"`
		MinWithdraw int64         `json:"minWithdraw,string"`
	} `json:"multiChain"`
}

type GlobalConfig struct {
	FeeAccountID                    int64  `json:"feeAccountId,string"`
	FeeAccountL2Key                 string `json:"feeAccountL2Key"`
	StarkExCollateralCurrencyID     string `json:"starkExCollateralCurrencyId"`
	StarkExFundingValidityPeriod    int64  `json:"starkExFundingValidityPeriod"`
	StarkExMaxFundingRate           int64  `json:"starkExMaxFundingRate,string"`
	StarkExOrdersTreeHeight         int64  `json:"starkExOrdersTreeHeight"`
	StarkExPositionsTreeHeight      int64  `json:"starkExPositionsTreeHeight"`
	StarkExPriceValidityPeriod      int64  `json:"starkExPriceValidityPeriod"`
	StarkExContractAddress          string `json:"starkExContractAddress"`
	RegisterEnvId                   int64  `json:"registerEnvId"`
	CrossChainAccountId             int64  `json:"crossChainAccountId,string"`
	CrossChainL2Key                 string `json:"crossChainL2Key"`
	FastWithdrawAccountId           int64  `json:"fastWithdrawAccountId,string"`
	FastWithdrawFactRegisterAddress string `json:"fastWithdrawFactRegisterAddress"`
	FastWithdrawL2Key               string `json:"fastWithdrawL2Key"`
	FastWithdrawMaxAmount           int64  `json:"fastWithdrawMaxAmount,string"`
}

type PerpetualContractConfig struct {
	BaselinePositionValue            float64 `json:"baselinePositionValue,string"`
	CrossID                          int64   `json:"crossId"`
	CrossSymbolID                    int64   `json:"crossSymbolId"`
	CrossSymbolName                  string  `json:"crossSymbolName"`
	DigitMerge                       string  `json:"digitMerge"`
	DisplayMaxLeverage               int64   `json:"displayMaxLeverage,string"`
	DisplayMinLeverage               int64   `json:"displayMinLeverage,string"`
	EnableDisplay                    bool    `json:"enableDisplay"`
	EnableOpenPosition               bool    `json:"enableOpenPosition"`
	EnableTrade                      bool    `json:"enableTrade"`
	FundingImpactMarginNotional      int64   `json:"fundingImpactMarginNotional,string"`
	FundingInterestRate              float64 `json:"fundingInterestRate,string"`
	IncrementalInitialMarginRate     float64 `json:"incrementalInitialMarginRate,string"`
	IncrementalMaintenanceMarginRate float64 `json:"incrementalMaintenanceMarginRate,string"`
	IncrementalPositionValue         float64 `json:"incrementalPositionValue,string"`
	InitialMarginRate                float64 `json:"initialMarginRate,string"`
	MaintenanceMarginRate            float64 `json:"maintenanceMarginRate,string"`
	MaxOrderSize                     float64 `json:"maxOrderSize,string"`
	MaxPositionSize                  float64 `json:"maxPositionSize,string"`
	MinOrderSize                     float64 `json:"minOrderSize,string"`
	MaxMarketPriceRange              float64 `json:"maxMarketPriceRange,string"`
	SettleCurrencyID                 string  `json:"settleCurrencyId"`
	StarkExOraclePriceQuorum         int64   `json:"starkExOraclePriceQuorum,string"`
	StarkExResolution                int64   `json:"starkExResolution,string"`
	StarkExRiskFactor                int64   `json:"starkExRiskFactor,string"`
	StarkExSyntheticAssetID          string  `json:"starkExSyntheticAssetId"`
	StepSize                         float64 `json:"stepSize,string"`
	Symbol                           string  `json:"symbol"`
	SymbolDisplayName                string  `json:"symbolDisplayName"`
	TickSize                         float64 `json:"tickSize,string"`
	UnderlyingCurrencyID             string  `json:"underlyingCurrencyId"`
	MaxMaintenanceMarginRate         float64 `json:"maxMaintenanceMarginRate,string"`
	MaxPositionValue                 float64 `json:"maxPositionValue,string"`
}

type CurrencyConfig struct {
	ID                string `json:"id"`
	StarkExAssetID    string `json:"starkExAssetId"`
	StarkExResolution string `json:"starkExResolution"`
	StepSize          string `json:"stepSize"`
	ShowStep          string `json:"showStep"`
	IconURL           string `json:"iconUrl"`
}

type ChainConfig struct {
	Chain              string        `json:"chain"`
	ChainID            int64         `json:"chainId"`
	ChainIconUrl       string        `json:"chainIconUrl"`
	ContractAddress    string        `json:"contractAddress"`
	DepositGasFeeLess  bool          `json:"depositGasFeeLess"`
	FeeLess            bool          `json:"feeLess"`
	FeeRate            float64       `json:"feeRate,string"`
	GasLess            bool          `json:"gasLess"`
	GasToken           string        `json:"gasToken"`
	MinFee             int64         `json:"minFee,string"`
	RPCUrl             string        `json:"rpcUrl"`
	WebTxUrl           string        `json:"webTxUrl"`
	TxConfirm          int64         `json:"txConfirm"`
	Tokens             []TokenConfig `json:"tokens"`
	WithdrawGasFeeLess bool          `json:"withdrawGasFeeLess"`
}

type TokenConfig struct {
	Decimals     int64  `json:"decimals"`
	IconUrl      string `json:"iconUrl"`
	Token        string `json:"token"`
	TokenAddress string `json:"tokenAddress"`
	PullOff      bool   `json:"pullOff"`
}

type orderbookResponse struct {
	Data struct {
		Asks   [][2]string `json:"a"`
		Bids   [][2]string `json:"b"`
		Symbol string      `json:"s"`
	} `json:"data"`
	Error
}

// Orderbook stores the orderbook data
type Orderbook struct {
	Bids   []orderbook.Item
	Asks   []orderbook.Item
	Symbol string
	Time   time.Time
}

type TradeData struct {
	Side   string           `json:"S"`
	Size   float64          `json:"v,string"`
	Price  float64          `json:"p,string"`
	Symbol string           `json:"s"`
	Time   apexTimeMilliSec `json:"T"`
}

type KlineData struct {
	Start    apexTimeMilliSec `json:"start"`
	Symbol   string           `json:"symbol"`
	Interval string           `json:"interval"`
	Low      float64          `json:"low,string"`
	High     float64          `json:"high,string"`
	Open     float64          `json:"open,string"`
	Close    float64          `json:"close,string"`
	Volume   float64          `json:"volume,string"`
	Turnover float64          `json:"turnover,string"`
}
