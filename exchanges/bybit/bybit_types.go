package bybit

import (
	"errors"
	"time"
)

var (
	errTypeAssert = errors.New("type assertion failed")
	errStrParsing = errors.New("parsing string failed")
)

// PairData stores pair data
type PairData struct {
	Name              string  `json:"name"`
	Alias             string  `json:"alias"`
	BaseCurrency      string  `json:"baseCurrency"`
	QuoteCurrency     string  `json:"quoteCurrency"`
	BasePrecision     float64 `json:"basePrecision,string"`
	QuotePrecision    float64 `json:"quotePrecision,string"`
	MinTradeQuantity  float64 `json:"minTradeQuantity,string"`
	MinTradeAmount    float64 `json:"minTradeAmount,string"`
	MinPricePrecision float64 `json:"minPricePrecision,string"`
	MaxTradeQuantity  float64 `json:"maxTradeQuantity,string"`
	MaxTradeAmount    float64 `json:"maxTradeAmount,string"`
	Category          int64   `json:"category"`
}

// OrderbookItem stores an individual orderbook item
type OrderbookItem struct {
	Price  float64
	Amount float64
}

// Orderbook stores the orderbook data
type Orderbook struct {
	Bids   []OrderbookItem
	Asks   []OrderbookItem
	Symbol string
	Time   time.Time
}

// TradeItem stores a single trade
type TradeItem struct {
	CurrencyPair string
	Price        float64
	Side         string
	Volume       float64
	TradeTime    time.Time
}

// KlineItem stores an individual kline data item
type KlineItem struct {
	StartTime        time.Time
	EndTime          time.Time
	Open             float64
	Close            float64
	High             float64
	Low              float64
	Volume           float64
	QuoteAssetVolume float64
	TakerBaseVolume  float64
	TakerQuoteVolume float64
	TradesCount      int64
}

// PriceChangeStats contains statistics for the last 24 hours trade
type PriceChangeStats struct {
	Time         time.Time
	Symbol       string
	BestBidPrice float64
	BestAskPrice float64
	LastPrice    float64
	OpenPrice    float64
	HighPrice    float64
	LowPrice     float64
	Volume       float64
	QuoteVolume  float64
}

// LastTradePrice contains price for last trade
type LastTradePrice struct {
	Symbol string  `json:"symbol"`
	Price  float64 `json:"price,string"`
}

// TickerData stores ticker data
type TickerData struct {
	Symbol      string
	BidPrice    float64
	BidQuantity float64
	AskPrice    float64
	AskQuantity float64
	Time        time.Time
}
