package brapigo

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"time"
)

const DEFAULT_BASE_URL = "https://brapi.dev/api"

func NewBrApi() BrApi {
	return BrApi{baseUrl: DEFAULT_BASE_URL, Client: &http.Client{
		Timeout: 10 * time.Second,
	}}
}

func NewBrApiWithCustomBaseUrl(baseUrl string) BrApi {
	return BrApi{baseUrl: baseUrl, Client: &http.Client{
		Timeout: 10 * time.Second,
	}}
}

type BrApi struct {
	baseUrl string
	Token   string
	Client  *http.Client
}

type StockApiResponse struct {
	Stocks []Stock `json:"stocks"`
}

type QuoteApiResponse struct {
	Results []Quote `json:"results"`
}

type TickerApiResponse struct {
	Stocks []string `json:"stocks"`
}

type Stock struct {
	Stock     string  `json:"stock"`
	Name      string  `json:"name"`
	Close     float64 `json:"close"`
	Change    float64 `json:"change"`
	Volume    int64   `json:"volume"`
	MarketCap float64 `json:"market_cap"`
	Logo      string  `json:"logo"`
	Sector    string  `json:"sector"`
}

type Quote struct {
	Symbol             string  `json:"symbol"`
	ShortName          string  `json:"shortName"`
	LongName           string  `json:"LongName"`
	Currency           string  `json:"Currency"`
	RegularMarketPrice float64 `json:"RegularMarketPrice"`
}

type CoinsApiResponse struct {
	Coins []string `json:"coins"`
}

type CryptoCoin struct {
	Currency                   string  `json:"currency"`
	CurrencyRateFromUSD        float64 `json:"currencyRateFromUSD"`
	CoinName                   string  `json:"coinName"`
	Coin                       string  `json:"coin"`
	RegularMarketChange        float64 `json:"regularMarketChange"`
	RegularMarketPrice         float64 `json:"regularMarketPrice"`
	RegularMarketChangePercent float64 `json:"regularMarketChangePercent"`
	RegularMarketDayLow        float64 `json:"regularMarketDayLow"`
	RegularMarketDayHigh       float64 `json:"regularMarketDayHigh"`
	RegularMarketDayRange      string  `json:"regularMarketDayRange"`
	RegularMarketVolume        float64 `json:"regularMarketVolume"`
	MarketCap                  float64 `json:"marketCap"`
	RegularMarketTime          int     `json:"regularMarketTime"`
	CoinImageUrl               string  `json:"coinImageUrl"`
}

func (a BrApi) FindAssetByTicker(tickers ...string) ([]Quote, error) {
	tickersParam := strings.Join(tickers, ",")

	req, err := http.NewRequest("GET", a.baseUrl+"/quote/"+tickersParam, nil)

	q := req.URL.Query()

	if a.Token != "" {
		q.Add("token", a.Token)
	}

	req.URL.RawQuery = q.Encode()

	resp, err := a.Client.Do(req)

	return parseQuoteResponse(resp, err)
}

func (a BrApi) SearchTickets(keyword string) ([]string, error) {
	req, err := http.NewRequest("GET", a.baseUrl+"/available", nil)

	q := req.URL.Query()
	q.Add("search", keyword)

	if a.Token != "" {
		q.Add("token", a.Token)
	}

	req.URL.RawQuery = q.Encode()

	if err != nil {
		return nil, err
	}

	resp, err := a.Client.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	responseData, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	tickerResponse := TickerApiResponse{}
	err = json.Unmarshal(responseData, &tickerResponse)

	if err != nil {
		return nil, err
	}

	return tickerResponse.Stocks, nil
}

func (a BrApi) ListStocks() ([]Stock, error) {
	req, err := http.NewRequest("GET", a.baseUrl+"/quote/list", nil)

	q := req.URL.Query()

	if a.Token != "" {
		q.Add("token", a.Token)
	}

	req.URL.RawQuery = q.Encode()

	if err != nil {
		return nil, err
	}

	resp, err := a.Client.Do(req)

	return parseStockResponse(resp, err)
}

func (a BrApi) ListCryptoCoins() ([]string, error) {
	req, err := http.NewRequest("GET", a.baseUrl+"/v2/crypto/available", nil)

	q := req.URL.Query()

	if a.Token != "" {
		q.Add("token", a.Token)
	}

	req.URL.RawQuery = q.Encode()

	if err != nil {
		return nil, err
	}

	resp, err := a.Client.Do(req)

	responseData, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	coinsResponse := CoinsApiResponse{}
	err = json.Unmarshal(responseData, &coinsResponse)

	if err != nil {
		return nil, err
	}

	return coinsResponse.Coins, nil
}

func (a BrApi) FindCryptoCoin(coins []string, currency string) ([]CryptoCoin, error) {
	if len(coins) == 0 {
		return nil, errors.New("it is necessary to inform at least one cryptocoin")
	}

	req, err := http.NewRequest("GET", a.baseUrl+"/v2/crypto", nil)

	q := req.URL.Query()

	if a.Token != "" {
		q.Add("token", a.Token)
	}

	q.Add("coin", strings.Join(coins, ","))

	if currency != "" {
		q.Add("currency", currency)
	}

	req.URL.RawQuery = q.Encode()

	if err != nil {
		return nil, err
	}

	resp, err := a.Client.Do(req)

	responseData, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	coinsRes := struct {
		Coins []CryptoCoin `json:"coins"`
	}{}
	err = json.Unmarshal(responseData, &coinsRes)

	if err != nil {
		return nil, err
	}

	return coinsRes.Coins, nil
}

func parseStockResponse(resp *http.Response, err error) ([]Stock, error) {
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	responseData, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	stocksData := StockApiResponse{}
	err = json.Unmarshal(responseData, &stocksData)

	if err != nil {
		return nil, err
	}

	return stocksData.Stocks, nil
}

func parseQuoteResponse(resp *http.Response, err error) ([]Quote, error) {
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	responseData, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	quoteData := QuoteApiResponse{}
	err = json.Unmarshal(responseData, &quoteData)

	if err != nil {
		return nil, err
	}

	return quoteData.Results, nil
}
