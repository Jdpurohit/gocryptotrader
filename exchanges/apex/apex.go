package apex

import (
	"bytes"
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/thrasher-corp/gocryptotrader/common"
	exchange "github.com/thrasher-corp/gocryptotrader/exchanges"
	"github.com/thrasher-corp/gocryptotrader/exchanges/orderbook"
	"github.com/thrasher-corp/gocryptotrader/exchanges/request"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	solsha3 "github.com/miguelmota/go-solidity-sha3"
)

// Apex is the overarching type across this package
type Apex struct {
	exchange.Base
}

const (
	apexAPIURL        = "https://pro.apex.exchange/api/"
	apexAPIVersion    = "v1"
	accountPrivateKey = "4107c052723f1a92e6a6f6fd81d6b20d75578637584a4c72808f1d44be6c473e"
	accountETHAddress = "0x4315c720e1c256A800B93c1742a6525fF40aB7C5"

	// Public endpoints
	apexSystemTime         = "/time"
	apexAllConfigData      = "/symbols"
	apexMarketDepth        = "/depth"
	apexLatestTrades       = "/trades"
	apexKlines             = "/klines"
	apexTicker             = "/ticker"
	apexFundingRateHistory = "/history-funding"
	apexCheckUserExists    = "/check-user-exist"

	// Authenticated endpoints
	apexNonce        = "/generate-nonce"
	apexRegistration = "/onboarding"
)

// GetSystemTime gets system time
func (ap *Apex) GetSystemTime(ctx context.Context) (time.Time, error) {
	resp := struct {
		Data struct {
			Time apexTimeMilliSec `json:"time"`
		} `json:"data"`
		Error
	}{}
	return resp.Data.Time.Time(), ap.SendHTTPRequest(ctx, exchange.RestSpot, http.MethodGet, apexSystemTime, publicSpotRate, &resp)
}

// GetAllConfigData gets all config data
func (ap *Apex) GetAllConfig(ctx context.Context) (Config, error) {
	resp := struct {
		Data     Config `json:"data"`
		TimeCost int64  `json:"timeCost"`
		Error
	}{}
	return resp.Data, ap.SendHTTPRequest(ctx, exchange.RestSpot, http.MethodGet, apexAllConfigData, publicSpotRate, &resp)
}

// GetMarketDepthData retrieve all active orderbook for one symbol, inclue all bids and asks.
func (ap *Apex) GetMarketDepth(ctx context.Context, symbol string, limit int64) (*Orderbook, error) {
	var o orderbookResponse

	params := url.Values{}
	params.Set("symbol", symbol)
	if limit != 0 {
		params.Set("limit", strconv.FormatInt(limit, 10))
	}
	path := common.EncodeURLValues(apexMarketDepth, params)
	err := ap.SendHTTPRequest(ctx, exchange.RestSpot, http.MethodGet, path, publicSpotRate, &o)
	if err != nil {
		return nil, err
	}
	return constructOrderbook(&o)
}

func processOB(ob [][2]string) ([]orderbook.Item, error) {
	o := make([]orderbook.Item, len(ob))
	for x := range ob {
		var price, amount float64
		amount, err := strconv.ParseFloat(ob[x][1], 64)
		if err != nil {
			return nil, err
		}
		price, err = strconv.ParseFloat(ob[x][0], 64)
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
	return &s, err
}

// GetLatestTrades retrieves all latest trading data.
func (ap *Apex) GetLatestTrades(ctx context.Context, symbol string, limit, from int64) ([]TradeData, error) {
	resp := struct {
		Data []TradeData `json:"data"`
		Error
	}{}

	params := url.Values{}
	params.Set("symbol", symbol)
	if limit != 0 {
		params.Set("limit", strconv.FormatInt(limit, 10))
	}
	if from != 0 {
		params.Set("from", strconv.FormatInt(from, 10))
	}
	path := common.EncodeURLValues(apexLatestTrades, params)
	return resp.Data, ap.SendHTTPRequest(ctx, exchange.RestSpot, http.MethodGet, path, publicSpotRate, &resp)
}

// GetCandlestickChart retrieves all candlestick chart data.
// Note: API response was not as per the v1 documentation so made changes as per the response.
func (ap *Apex) GetCandlestickChart(ctx context.Context, interval, symbol string, startTime, endTime time.Time, limit int64) ([]KlineData, error) {
	resp := struct {
		Data map[string][]KlineData `json:"data"`
		Error
	}{}

	params := url.Values{}
	if interval != "" {
		if !common.StringDataCompare(validIntervals, interval) {
			return nil, errInvalidInterval
		}
		params.Set("interval", interval)
	}
	if symbol == "" {
		return nil, errSymbolMissing
	}
	params.Set("symbol", symbol)

	if !startTime.IsZero() {
		params.Add("start", strconv.FormatInt(startTime.UnixMilli(), 10))
	}
	if !endTime.IsZero() {
		params.Add("end", strconv.FormatInt(endTime.UnixMilli(), 10))
	}
	if limit != 0 {
		params.Set("limit", strconv.FormatInt(limit, 10))
	}
	path := common.EncodeURLValues(apexKlines, params)
	return resp.Data[symbol], ap.SendHTTPRequest(ctx, exchange.RestSpot, http.MethodGet, path, publicSpotRate, &resp)
}

// GetTicker retrieves the latest data on symbol tickers.
func (ap *Apex) GetTicker(ctx context.Context, symbol string) ([]TickerData, error) {
	resp := struct {
		Data []TickerData `json:"data"`
		Error
	}{}

	params := url.Values{}
	if symbol == "" {
		return nil, errSymbolMissing
	}
	params.Set("symbol", symbol)
	path := common.EncodeURLValues(apexTicker, params)
	return resp.Data, ap.SendHTTPRequest(ctx, exchange.RestSpot, http.MethodGet, path, publicSpotRate, &resp)
}

// GetFundingRateHistory retrieves the funding rate history.
func (ap *Apex) GetFundingRateHistory(ctx context.Context, symbol string, startTime, endTime time.Time, limit, page int64) ([]FundData, error) {
	resp := struct {
		Data struct {
			FundingHistory []FundData `json:"historyFunds"`
		} `json:"data"`
		Error
	}{}

	params := url.Values{}
	if symbol == "" {
		return nil, errSymbolMissing
	}
	params.Set("symbol", symbol)
	if !startTime.IsZero() {
		params.Add("beginTimeInclusive", strconv.FormatInt(startTime.UnixMilli(), 10))
	}
	if !endTime.IsZero() {
		params.Add("endTimeExclusive", strconv.FormatInt(endTime.UnixMilli(), 10))
	}
	if limit != 0 {
		params.Set("limit", strconv.FormatInt(limit, 10))
	}
	if page >= 0 {
		params.Set("page", strconv.FormatInt(page, 10))
	}
	path := common.EncodeURLValues(apexFundingRateHistory, params)
	return resp.Data.FundingHistory, ap.SendHTTPRequest(ctx, exchange.RestSpot, http.MethodGet, path, publicSpotRate, &resp)
}

// CheckIfUserExists validates if user exists in the system.
func (ap *Apex) CheckIfUserExists(ctx context.Context, ethAddress string) (bool, error) {
	resp := struct {
		Data bool `json:"data"`
		Error
	}{}

	params := url.Values{}
	if ethAddress != "" {
		params.Set("ethAddress", ethAddress)
	}
	path := common.EncodeURLValues(apexCheckUserExists, params)
	return resp.Data, ap.SendHTTPRequest(ctx, exchange.RestSpot, http.MethodGet, path, publicSpotRate, &resp)
}

// GenerateNonce generate and obtain nonce.
func (ap *Apex) GenerateNonce(ctx context.Context, ethAddress, starkKey, chainID string) (NonceData, error) {
	resp := struct {
		Data NonceData `json:"data"`
		Error
	}{}

	params := url.Values{}
	if ethAddress == "" {
		return resp.Data, errETHAddressMissing
	}
	params.Set("ethAddress", ethAddress)
	if starkKey == "" {
		return resp.Data, errStarkKeyMissing
	}
	params.Set("starkKey", starkKey)
	if chainID == "" {
		return resp.Data, errChainIDMissing
	}
	params.Set("chainId", chainID)
	params.Set("category", "CATEGORY_API")
	path := common.EncodeURLValues(apexNonce, params)
	return resp.Data, ap.SendHTTPRequest(ctx, exchange.RestSpot, http.MethodPost, path, publicSpotRate, &resp)
}

// Registration and User Onboarding
func (ap *Apex) Registration(ctx context.Context, starkKey, starkKeyYCoordinate, ethAddress, referredByAffiliateLink, country, chainID string) (interface{}, error) {
	resp := struct {
		Data interface{} `json:"data"`
		Error
	}{}

	params := url.Values{}
	if starkKey == "" {
		return nil, errStarkKeyMissing
	}
	params.Set("starkKey", starkKey)
	if starkKeyYCoordinate == "" {
		return nil, errStarkKeyYCoordinateMisssing
	}
	params.Set("starkKeyYCoordinate", starkKeyYCoordinate)
	if ethAddress == "" {
		return nil, errETHAddressMissing
	}
	params.Set("ethereumAddress", ethAddress)
	if referredByAffiliateLink != "" {
		params.Set("referredByAffiliateLink", referredByAffiliateLink)
	}
	if country != "" {
		params.Set("country", country)
	}
	params.Set("category", "CATEGORY_API")

	// generate new nonce and use it
	nonce, err := ap.GenerateNonce(ctx, ethAddress, starkKey, chainID)
	if err != nil {
		return nil, err
	}

	sign, err := getSign(nonce.Nonce)
	if err != nil {
		return nil, err
	}

	headers := make(map[string]string)
	headers["APEX-SIGNATURE"] = sign
	headers["APEX-ETHEREUM-ADDRESS"] = accountETHAddress
	return &resp.Data, ap.SendAuthHTTPRequest(ctx, exchange.RestSpot, http.MethodPost, apexRegistration, params, headers, &resp, publicSpotRate, true)
}

// SendHTTPRequest sends an unauthenticated request
func (ap *Apex) SendHTTPRequest(ctx context.Context, ePath exchange.URL, method, path string, f request.EndpointLimit, result UnmarshalTo) error {
	endpointPath, err := ap.API.Endpoints.GetURL(ePath)
	if err != nil {
		return err
	}

	err = ap.SendPayload(ctx, f, func() (*request.Item, error) {
		return &request.Item{
			Method:        method,
			Path:          endpointPath + apexAPIVersion + path,
			Result:        result,
			Verbose:       ap.Verbose,
			HTTPDebugging: ap.HTTPDebugging,
			HTTPRecording: ap.HTTPRecording}, nil
	})
	if err != nil {
		return err
	}
	return result.GetError()
}

// SendAuthHTTPRequest sends an authenticated HTTP request
// TODO: remove jsonPayload if non of the request requires it
func (ap *Apex) SendAuthHTTPRequest(ctx context.Context, ePath exchange.URL, method, path string, params url.Values, headers map[string]string, result UnmarshalTo, f request.EndpointLimit, isRegisterAPI bool) error {
	if headers == nil {
		headers = make(map[string]string)
	}

	if !isRegisterAPI {
		creds, err := ap.GetCredentials(ctx)
		if err != nil {
			return err
		}
		headers["APEX-API-KEY"] = creds.Key
		headers["APEX-PASSPHRASE"] = creds.Secret
	}

	if result == nil {
		result = &Error{}
	}
	endpointPath, err := ap.API.Endpoints.GetURL(ePath)
	if err != nil {
		return err
	}
	if params == nil {
		params = url.Values{}
	}

	// TODO:
	// call generate stark key once and store all keys
	// call generate nonce everytime a private request is made
	// sign the request

	err = ap.SendPayload(ctx, f, func() (*request.Item, error) {
		switch method {
		case http.MethodPost:
			headers["Content-Type"] = "application/x-www-form-urlencoded" // required for all private API
		}
		return &request.Item{
			Method:        method,
			Path:          endpointPath + apexAPIVersion + path,
			Headers:       headers,
			Body:          bytes.NewBuffer([]byte(params.Encode())),
			Result:        &result,
			AuthRequest:   true,
			Verbose:       ap.Verbose,
			HTTPDebugging: ap.HTTPDebugging,
			HTTPRecording: ap.HTTPRecording}, nil
	})
	if err != nil {
		return err
	}
	return result.GetError()
}

func deriveStarkKey() (string, string, error) {
	msgStr := "name: ApeX\nversion: 1.0\nenvId: 1\naction: L2 Key\nonlySignOn: https://pro.apex.exchange"
	bHash := solsha3.SoliditySHA3WithPrefix([]byte(msgStr))
	fmt.Println("Bytes Hash: ", hex.EncodeToString(bHash))
	sign, err := signMethod(bHash, SignatureTypePersonal)
	if err != nil {
		return "", "", err
	}
	fmt.Println("Sign: ", sign)

	bHashSign := solsha3.SoliditySHA3(sign)
	fmt.Println("bHashSign: ", hex.EncodeToString(bHashSign))

	m := new(big.Int)
	m.SetString(hex.EncodeToString(bHashSign), 16)
	fmt.Println("bHashSign Int: ", m.String())
	// privKey := n.Text(16)

	// // fmt.Println("privKey INT: ", n.String())
	// // fmt.Println("privKey Hex: ", privKey)

	ecPoint := privKeyToECPointOnStarkCurve(new(big.Int).Rsh(m, 5))
	return ecPoint[0].Text(16), ecPoint[1].Text(16), nil
}

func privKeyToECPointOnStarkCurve(privKeyInt *big.Int) [2]*big.Int {
	a, _ := new(big.Int).SetString("874739451078007766457464989774322083649278607533249481151382481072868806602", 10)
	b, _ := new(big.Int).SetString("152666792071518830868575557812948353041420400780739481342941381225525861407", 10)

	ecGenerator := [2]*big.Int{a, b}
	alpha, _ := new(big.Int).SetString("1", 10)
	fieldPrime, _ := new(big.Int).SetString("3618502788666131213697322783095070105623107215331596699973092056135872020481", 10)
	return ecMult(privKeyInt, ecGenerator, alpha, fieldPrime)
}

// Multiplies by m a point on the elliptic curve with equation y^2 = x^3 + alpha*x + beta mod p.
// Assumes the point is given in affine form (x, y) and that 0 < m < order(point).
func ecMult(privKeyInt *big.Int, ecGenPair [2]*big.Int, alpha, fieldPrime *big.Int) [2]*big.Int {
	if privKeyInt.Cmp(big.NewInt(1)) == 0 {
		return ecGenPair
	}
	if big.NewInt(0).Cmp(new(big.Int).Mod(privKeyInt, big.NewInt(2))) == 0 {
		return ecMult(new(big.Int).Div(privKeyInt, big.NewInt(2)), ecDouble(ecGenPair, alpha, fieldPrime), alpha, fieldPrime)
	}
	return ecAdd(ecMult(new(big.Int).Sub(privKeyInt, big.NewInt(1)), ecGenPair, alpha, fieldPrime), ecGenPair, fieldPrime)
}

// Gets two points on an elliptic curve mod p and returns their sum.
// Assumes the points are given in affine form (x, y) and have different x coordinates.
func ecAdd(ecPoint1, ecPoint2 [2]*big.Int, p *big.Int) [2]*big.Int {
	m := divMod(new(big.Int).Sub(ecPoint1[1], ecPoint2[1]), new(big.Int).Sub(ecPoint1[0], ecPoint2[0]), p)
	var ecPoint [2]*big.Int
	ecPoint[0] = new(big.Int).Mod(new(big.Int).Sub(new(big.Int).Sub(new(big.Int).Mul(m, m), ecPoint1[0]), ecPoint2[0]), p)
	ecPoint[1] = new(big.Int).Mod(new(big.Int).Sub(new(big.Int).Mul(m, new(big.Int).Sub(ecPoint1[0], ecPoint[0])), ecPoint1[1]), p)
	return ecPoint
}

// Doubles a point on an elliptic curve with the equation y^2 = x^3 + alpha*x + beta mod p.
// Assumes the point is given in affine form (x, y) and has y != 0.
func ecDouble(point [2]*big.Int, alpha, fieldPrime *big.Int) [2]*big.Int {
	var rPoint [2]*big.Int
	m := divMod(new(big.Int).Add(new(big.Int).Mul(big.NewInt(3), new(big.Int).Exp(point[0], big.NewInt(2), nil)), alpha), new(big.Int).Mul(big.NewInt(2), point[1]), fieldPrime)

	rPoint[0] = new(big.Int).Mod(new(big.Int).Sub(new(big.Int).Exp(m, big.NewInt(2), nil), new(big.Int).Mul(big.NewInt(2), point[0])), fieldPrime)
	rPoint[1] = new(big.Int).Mod(new(big.Int).Sub(new(big.Int).Mul(m, new(big.Int).Sub(point[0], rPoint[0])), point[1]), fieldPrime)
	return rPoint
}

// Finds a nonnegative integer 0 <= x < b such that (a * x) % b == c
func divMod(n, m, p *big.Int) *big.Int {
	a := big.NewInt(0)
	b := big.NewInt(0)
	gcd := new(big.Int).GCD(a, b, m, p)
	if gcd.Cmp(big.NewInt(1)) != 0 {
		// TODO: remove this panic
		panic("divMod: GCD not equal to 1")
	}
	return new(big.Int).Mod(new(big.Int).Mul(n, a), p)
}

func getSign(nonce string) (string, error) {
	sign, err := signMethod(getHash(nonce), SignatureTypeNOPrepend)
	if err != nil {
		return "", err
	}
	return sign, nil
}

func getEIP712Message(nonce string) string {
	return fmt.Sprintf(`{
		'types': {
			'EIP712Domain': [
				{
					'name': 'name',
					'type': 'string'
				},
				{
					'name': 'version',
					'type': 'string'
				},
				{
					'name': 'chainId',
					'type': 'uint256'
				}
			],
			'ApeX': [
				{
					'type': 'string', 
					'name': 'action'
				},
				{	
					'type': 'string', 
					'name': 'onlySignOn'
				},
				{
					'type': 'string', 
					'name': 'nonce'
				}
			]
		},
		'domain': {
			'name': 'ApeX',
			'version': '1.0',
			'chainId': 1
		},
		'primaryType': 'ApeX',
		'message': {
			'action': 'ApeX Onboarding',
			'nonce': '%s',
			'onlySignOn': 'https://pro.apex.exchange'
		}
	}`, nonce)
}

func getHash(nonce string) []byte {
	return solsha3.SoliditySHA3(
		[]string{"bytes2", "bytes32", "bytes32"},
		[]interface{}{
			"0x1901",
			getDomainHash(),
			getStructHash(nonce),
		},
	)
}

func getDomainHash() []byte {
	return solsha3.SoliditySHA3(
		[]string{"bytes32", "bytes32", "bytes32", "uint256"},
		[]interface{}{
			getHashString("EIP712Domain(string name,string version,uint256 chainId)"),
			getHashString("ApeX"),
			getHashString("1.0"),
			"1",
		},
	)
}

func getStructHash(nonce string) []byte {
	return solsha3.SoliditySHA3(
		[]string{"bytes32", "bytes32", "bytes32", "bytes32"},
		[]interface{}{
			getHashString("ApeX(string action,string onlySignOn,string nonce)"),
			getHashString("ApeX Onboarding"),
			getHashString("https://pro.apex.exchange"),
			getHashString(nonce),
		},
	)
}

func getHashString(str string) []byte {
	return solsha3.SoliditySHA3(
		[]string{"string"},
		[]interface{}{str},
	)
}

func signMethod(payload []byte, signType string) (string, error) {
	privateKey, err := crypto.HexToECDSA(accountPrivateKey)
	if err != nil {
		return "", err
	}

	signature, err := crypto.Sign(payload, privateKey)
	if err != nil {
		return "", err
	}

	tSign, err := createTypedSign(hexutil.Encode(signature), signType)
	if err != nil {
		return "", err
	}
	return tSign, nil
}

func createTypedSign(sign, signatureType string) (string, error) {
	if strings.HasPrefix(sign, "0x") {
		sign = sign[2:]
	}

	if len(sign) != 130 {
		return "", errors.New("invalid raw signature")
	}

	rs := sign[:128]
	v := sign[128:130]

	if v == "00" {
		return "0x" + rs + "1b" + "0" + signatureType, nil
	}
	if v == "01" {
		return "0x" + rs + "1c" + "0" + signatureType, nil
	}
	if v == "1b" || v == "1c" {
		return "0x" + sign + "0" + signatureType, nil
	}
	return "", fmt.Errorf("invalid value of v: %s", v)
}
