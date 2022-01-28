package bybit

import (
	"time"

	"github.com/thrasher-corp/gocryptotrader/exchanges/request"
	"golang.org/x/time/rate"
)

const (
	spotInterval          = time.Second
	spotRequestRate       = 70
	futuresPublicInterval = time.Second
	futuresRequestRate    = 50

	spotPrivateRequestRate  = 20
	futuresInterval         = time.Minute
	futuresDefaultRate      = 100
	futuresOrderRate        = 100
	futuresOrderListRate    = 600
	futuresExecutionRate    = 120
	futuresPositionRate     = 75
	futuresPositionListRate = 120
	futuresFundingRate      = 120
	futuresWalletRate       = 120
	futuresAccountRate      = 600
)

const (
	publicSpotRate request.EndpointLimit = iota
	publicFuturesRate
	privateSpotRate

	cFuturesDefaultRate

	cFuturesCancelActiveOrderRate
	cFuturesCancelAllActiveOrderRate
	cFuturesCreateConditionalOrderRate
	cFuturesCancelConditionalOrderRate
	cFuturesReplaceActiveOrderRate
	cFuturesReplaceConditionalOrderRate
	cFuturesCreateOrderRate
	cFuturesCancelAllConditionalOrderRate

	cFuturesGetActiveOrderRate
	cFuturesGetConditionalOrderRate
	cFuturesGetRealtimeOrderRate

	cFuturesTradeRate

	cFuturesSetLeverageRate
	cFuturesUpdateMarginRate
	cFuturesSetTradingRate
	cFuturesSwitchPositionRate

	cFuturesPositionRate
	cFuturesWalletBalanceRate

	cFuturesLastFundingFeeRate
	cFuturesPredictFundingRate

	cFuturesWalletFundRecordRate
	cFuturesWalletWithdrawalRate

	cFuturesAPIKeyInfoRate

	uFuturesDefaultRate

	uFuturesCreateOrderRate
	uFuturesCancelOrderRate
	uFuturesCancelAllOrderRate
	uFuturesCreateConditionalOrderRate
	uFuturesCancelConditionalOrderRate
	uFuturesCancelAllConditionalOrderRate

	uFuturesSetLeverageRate
	uFuturesSwitchMargin
	uFuturesSwitchPosition
	uFuturesSetMarginRate
	uFuturesSetTradingStopRate
	uFuturesUpdateMarginRate

	uFuturesPositionRate
	uFuturesGetClosedTradesRate
	uFuturesGetTradesRate

	uFuturesGetActiveOrderRate
	uFuturesGetActiveRealtimeOrderRate
	uFuturesGetCondtionalOrderRate
	uFuturesGetCondtionalRealtimeOrderRate

	uFuturesGetMyLastFundingFeeRate
	uFuturesPredictFundingRate

	FuturesDefaultRate

	FuturesCancelOrderRate
	FuturesCreateOrderRate
	FuturesReplaceOrderRate
	FuturesCancelAllOrderRate
	FuturesCancelAllCondtionalOrderRate
	FuturesReplaceConditionalOrderRate
	FuturesCancelCondtionalOrderRate
	FuturesCreateConditionalOrderRate

	FuturesGetActiveOrderRate
	FuturesGetConditionalOrderRate
	FuturesGetActiveRealtimeOrderRate
	FuturesGetConditionalRealtimeOrderRate

	FuturesGetTradeRate

	FuturesSetLeverateRate
	FuturesUpdateMarginRate
	FuturesSetTradingStopRate
	FuturesSwitchPositionModeRate
	FuturesSwitchMarginRate
	FuturesSwitchPositionRate

	FuturesPositionRate
)

// RateLimit implements the request.Limiter interface
type RateLimit struct {
	SpotRate                  *rate.Limiter
	FuturesRate               *rate.Limiter
	PrivateSpotRate           *rate.Limiter
	CMFuturesDefaultRate      *rate.Limiter
	CMFuturesOrderRate        *rate.Limiter
	CMFuturesOrderListRate    *rate.Limiter
	CMFuturesExecutionRate    *rate.Limiter
	CMFuturesPositionRate     *rate.Limiter
	CMFuturesPositionListRate *rate.Limiter
	CMFuturesFundingRate      *rate.Limiter
	CMFuturesWalletRate       *rate.Limiter
	CMFuturesAccountRate      *rate.Limiter
	UFuturesDefaultRate       *rate.Limiter
	UFuturesOrderRate         *rate.Limiter
	UFuturesPositionRate      *rate.Limiter
	UFuturesPositionListRate  *rate.Limiter
	UFuturesOrderListRate     *rate.Limiter
	UFuturesFundingRate       *rate.Limiter
	FuturesDefaultRate        *rate.Limiter
	FuturesOrderRate          *rate.Limiter
	FuturesOrderListRate      *rate.Limiter
	FuturesExecutionRate      *rate.Limiter
	FuturesPositionRate       *rate.Limiter
	FuturesPositionListRate   *rate.Limiter
}

// Limit executes rate limiting functionality for Binance
func (r *RateLimit) Limit(f request.EndpointLimit) error {
	var limiter *rate.Limiter
	var tokens int
	switch f {
	case publicSpotRate:
		limiter, tokens = r.SpotRate, 1
	case privateSpotRate:
		limiter, tokens = r.PrivateSpotRate, 1

	case cFuturesDefaultRate:
		limiter, tokens = r.CMFuturesDefaultRate, 1

	case cFuturesCancelActiveOrderRate, cFuturesCreateConditionalOrderRate, cFuturesCancelConditionalOrderRate, cFuturesReplaceActiveOrderRate,
		cFuturesReplaceConditionalOrderRate, cFuturesCreateOrderRate:
		limiter, tokens = r.CMFuturesOrderRate, 1
	case cFuturesCancelAllActiveOrderRate, cFuturesCancelAllConditionalOrderRate:
		limiter, tokens = r.CMFuturesOrderRate, 10

	case cFuturesGetActiveOrderRate, cFuturesGetConditionalOrderRate, cFuturesGetRealtimeOrderRate:
		limiter, tokens = r.CMFuturesOrderListRate, 1

	case cFuturesTradeRate:
		limiter, tokens = r.CMFuturesExecutionRate, 1

	case cFuturesSetLeverageRate, cFuturesUpdateMarginRate, cFuturesSetTradingRate, cFuturesSwitchPositionRate:
		limiter, tokens = r.CMFuturesPositionRate, 1

	case cFuturesPositionRate, cFuturesWalletBalanceRate:
		limiter, tokens = r.CMFuturesPositionListRate, 1

	case cFuturesLastFundingFeeRate, cFuturesPredictFundingRate:
		limiter, tokens = r.CMFuturesFundingRate, 1

	case cFuturesWalletFundRecordRate, cFuturesWalletWithdrawalRate:
		limiter, tokens = r.CMFuturesWalletRate, 1

	case cFuturesAPIKeyInfoRate:
		limiter, tokens = r.CMFuturesAccountRate, 1

	case uFuturesDefaultRate:
		limiter, tokens = r.UFuturesDefaultRate, 1

	case uFuturesCreateOrderRate, uFuturesCancelOrderRate, uFuturesCreateConditionalOrderRate, uFuturesCancelConditionalOrderRate:
		limiter, tokens = r.UFuturesOrderRate, 1

	case uFuturesCancelAllOrderRate, uFuturesCancelAllConditionalOrderRate:
		limiter, tokens = r.UFuturesOrderRate, 10

	case uFuturesSetLeverageRate, uFuturesSwitchMargin, uFuturesSwitchPosition, uFuturesSetMarginRate, uFuturesSetTradingStopRate, uFuturesUpdateMarginRate:
		limiter, tokens = r.UFuturesPositionRate, 1

	case uFuturesPositionRate, uFuturesGetClosedTradesRate, uFuturesGetTradesRate:
		limiter, tokens = r.UFuturesPositionListRate, 1

	case uFuturesGetActiveOrderRate, uFuturesGetActiveRealtimeOrderRate, uFuturesGetCondtionalOrderRate, uFuturesGetCondtionalRealtimeOrderRate:
		limiter, tokens = r.UFuturesOrderListRate, 1

	case uFuturesGetMyLastFundingFeeRate, uFuturesPredictFundingRate:
		limiter, tokens = r.UFuturesFundingRate, 1

	case FuturesDefaultRate:
		limiter, tokens = r.FuturesDefaultRate, 1

	case FuturesCancelOrderRate, FuturesCreateOrderRate, FuturesReplaceOrderRate, FuturesReplaceConditionalOrderRate, FuturesCancelCondtionalOrderRate,
		FuturesCreateConditionalOrderRate:
		limiter, tokens = r.FuturesOrderRate, 1

	case FuturesCancelAllOrderRate, FuturesCancelAllCondtionalOrderRate:
		limiter, tokens = r.FuturesOrderRate, 10

	case FuturesGetActiveOrderRate, FuturesGetConditionalOrderRate, FuturesGetActiveRealtimeOrderRate, FuturesGetConditionalRealtimeOrderRate:
		limiter, tokens = r.FuturesOrderListRate, 1

	case FuturesGetTradeRate:
		limiter, tokens = r.FuturesExecutionRate, 1

	case FuturesSetLeverateRate, FuturesUpdateMarginRate, FuturesSetTradingStopRate, FuturesSwitchPositionModeRate, FuturesSwitchMarginRate, FuturesSwitchPositionRate:
		limiter, tokens = r.FuturesPositionRate, 1

	case FuturesPositionRate:
		limiter, tokens = r.FuturesPositionListRate, 1

	default:
		limiter, tokens = r.SpotRate, 1
	}

	var finalDelay time.Duration
	for i := 0; i < tokens; i++ {
		// Consume tokens 1 at a time as this avoids needing burst capacity in the limiter,
		// which would otherwise allow the rate limit to be exceeded over short periods
		finalDelay = limiter.Reserve().Delay()
	}
	time.Sleep(finalDelay)
	return nil
}

// SetRateLimit returns the rate limit for the exchange
func SetRateLimit() *RateLimit {
	return &RateLimit{
		SpotRate:                  request.NewRateLimit(spotInterval, spotRequestRate),
		FuturesRate:               request.NewRateLimit(futuresPublicInterval, futuresRequestRate),
		PrivateSpotRate:           request.NewRateLimit(spotInterval, spotPrivateRequestRate),
		CMFuturesDefaultRate:      request.NewRateLimit(futuresInterval, futuresDefaultRate),
		CMFuturesOrderRate:        request.NewRateLimit(futuresInterval, futuresOrderRate),
		CMFuturesOrderListRate:    request.NewRateLimit(futuresInterval, futuresOrderListRate),
		CMFuturesExecutionRate:    request.NewRateLimit(futuresInterval, futuresExecutionRate),
		CMFuturesPositionRate:     request.NewRateLimit(futuresInterval, futuresPositionRate),
		CMFuturesPositionListRate: request.NewRateLimit(futuresInterval, futuresPositionListRate),
		CMFuturesFundingRate:      request.NewRateLimit(futuresInterval, futuresFundingRate),
		CMFuturesWalletRate:       request.NewRateLimit(futuresInterval, futuresWalletRate),
		CMFuturesAccountRate:      request.NewRateLimit(futuresInterval, futuresAccountRate),
		UFuturesDefaultRate:       request.NewRateLimit(futuresInterval, futuresDefaultRate),
		UFuturesOrderRate:         request.NewRateLimit(futuresInterval, futuresOrderRate),
		UFuturesPositionRate:      request.NewRateLimit(futuresInterval, futuresPositionRate),
		UFuturesPositionListRate:  request.NewRateLimit(futuresInterval, futuresPositionListRate),
		UFuturesOrderListRate:     request.NewRateLimit(futuresInterval, futuresOrderListRate),
		UFuturesFundingRate:       request.NewRateLimit(futuresInterval, futuresFundingRate),
		FuturesDefaultRate:        request.NewRateLimit(futuresInterval, futuresDefaultRate),
		FuturesOrderRate:          request.NewRateLimit(futuresInterval, futuresOrderRate),
		FuturesOrderListRate:      request.NewRateLimit(futuresInterval, futuresOrderListRate),
		FuturesExecutionRate:      request.NewRateLimit(futuresInterval, futuresExecutionRate),
		FuturesPositionRate:       request.NewRateLimit(futuresInterval, futuresPositionRate),
		FuturesPositionListRate:   request.NewRateLimit(futuresInterval, futuresPositionListRate),
	}
}