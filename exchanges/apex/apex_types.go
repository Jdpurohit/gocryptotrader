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

type SignatureType int

const (
	SignatureTypeNOPrepend   = "0"
	SignatureTypeDecimal     = "1"
	SignatureTypeHexaDecimal = "2"
	SignatureTypePersonal    = "3"

	StarkAlpha      = 1
	StarkFieldPrime = "3618502788666131213697322783095070105623107215331596699973092056135872020481"
)

var (
	errInvalidInterval             = errors.New("invalid interval")
	errSymbolMissing               = errors.New("symbol missing")
	errETHAddressMissing           = errors.New("ethAddress missing")
	errStarkKeyMissing             = errors.New("starkKey missing")
	errChainIDMissing              = errors.New("chainId missing")
	errStarkKeyYCoordinateMisssing = errors.New("starkKeyYCoordinate missing")
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
	Start    apexTimeMilliSec `json:"t"`
	Symbol   string           `json:"s"`
	Interval string           `json:"i"`
	Low      float64          `json:"l,string"`
	High     float64          `json:"h,string"`
	Open     float64          `json:"o,string"`
	Close    float64          `json:"c,string"`
	Volume   float64          `json:"v,string"`
	Turnover float64          `json:"tr,string"`
}

type TickerData struct {
	Symbol               string  `json:"symbol"`
	Price24hPcnt         float64 `json:"price24hPcnt,string"`
	LastPrice            float64 `json:"lastPrice,string"`
	HighPrice24h         float64 `json:"highPrice24h,string"`
	LowPrice24h          float64 `json:"lowPrice24h,string"`
	OraclePrice          float64 `json:"oraclePrice,string"`
	IndexPrice           float64 `json:"indexPrice,string"`
	OpenInterest         float64 `json:"openInterest,string"`
	Turnover24h          float64 `json:"turnover24h,string"`
	Volume24h            float64 `json:"volume24h,string"`
	FundingRate          float64 `json:"fundingRate,string"`
	PredictedFundingRate float64 `json:"predictedFundingRate,string"`
	NextFundingTime      string  `json:"nextFundingTime"`
	TradeCount           string  `json:"tradeCount"`
}

type FundData struct {
	Symbol           string  `json:"symbol"`
	Rate             float64 `json:"rate"`
	Price            float64 `json:"price"`
	FundingTime      int64   `json:"fundingTime"`
	FundingTimestamp int64   `json:"fundingTimestamp"`
}

type VersionData struct {
	PlatformName       string `json:"platformName"`
	LatestVersion      string `json:"latestVersion"`
	ForceUpdateVersion string `json:"forceUpdateVersion"`
	DownloadUrl        string `json:"downloadUrl"`
	MarketUrl          string `json:"marketUrl"`
	Description        string `json:"description"`
	Title              string `json:"title"`
}

type NonceData struct {
	Nonce       string           `json:"nonce"`
	NonceExpiry apexTimeMilliSec `json:"nonceExpired"`
}

type UserInfo struct {
	EthereumAddress          string      `json:"ethereumAddress"`
	IsRegistered             bool        `json:"isRegistered"`
	Email                    string      `json:"email"`
	Username                 string      `json:"username"`
	UserData                 interface{} `json:"userData"`
	IsEmailVerified          bool        `json:"isEmailVerified"`
	EmailNotifyGeneralEnable bool        `json:"emailNotifyGeneralEnable"`
	EmailNotifyTradingEnable bool        `json:"emailNotifyTradingEnable"`
	EmailNotifyAccountEnable bool        `json:"emailNotifyAccountEnable"`
	PopupNotifyTradingEnable bool        `json:"popupNotifyTradingEnable"`
}
