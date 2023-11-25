package brapigo

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
)

var stocksData *StockApiResponse
var cryptoData []CryptoCoin

func runTestServer() *httptest.Server {
	stocks, err := loadAllAssets()

	if err != nil {
		panic("Cannot load testdata/list.json stocksData file")
	}

	stocksData = stocks

	crypto, err := loadCryptocoins()

	if err != nil {
		panic("Cannot load testdata/cryptocoins.json cryptocoins file")
	}

	cryptoData = crypto

	return httptest.NewServer(http.HandlerFunc(rootHandlerFunc))
}

type handler func(r *http.Request) ([]byte, error)

var endpoints map[string]handler = map[string]handler{
	"/quote/list": func(r *http.Request) ([]byte, error) {
		content, err := json.Marshal(stocksData)

		if err != nil {
			return nil, err
		}

		return content, nil
	},
	"/quote/PETR3": func(r *http.Request) ([]byte, error) {
		stocksFound, err := searchAssetsByTicker("PETR3")

		if err != nil {
			return nil, err
		}

		return json.Marshal(StockApiResponse{Stocks: stocksFound})
	},
	"/quote/PETR4,ITUB3": func(r *http.Request) ([]byte, error) {
		quotes, err := loadQuoteData()

		if err != nil {
			return nil, err
		}

		return json.Marshal(quotes)
	},
	"/available": func(r *http.Request) ([]byte, error) {
		keyword := r.URL.Query().Get("search")

		assets, err := searchAssetsByKeyword(keyword)

		if err != nil {
			return nil, err
		}

		response := TickerApiResponse{Stocks: []string{}}

		for _, asset := range assets {
			response.Stocks = append(response.Stocks, asset.Stock)
		}

		return json.Marshal(response)
	},
	"/v2/crypto/available": func(r *http.Request) ([]byte, error) {
		return json.Marshal(CoinsApiResponse{Coins: []string{"BTC"}})
	},
	"/v2/crypto": func(r *http.Request) ([]byte, error) {
		q := r.URL.Query()
		coins := q.Get("coin")
		currency := q.Get("currency")

		cryptos := searchCryptoCoins(strings.Split(coins, ","), currency)

		data := struct {
			Coins []CryptoCoin `json:"coins"`
		}{Coins: cryptos}

		return json.Marshal(data)
	},
}

var rootHandlerFunc = func(w http.ResponseWriter, r *http.Request) {
	if endpointFunc, ok := endpoints[r.URL.Path]; ok {
		data, err := endpointFunc(r)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write(data)
	} else {
		http.Error(w, "Not Found", http.StatusNotFound)
	}
}

func searchAssetsByTicker(tickers string) ([]Stock, error) {
	tickersSplit := strings.Split(tickers, ",")

	stocksFound := []Stock{}

	for _, ticker := range stocksData.Stocks {
		for _, tickerSearch := range tickersSplit {
			if ticker.Stock == tickerSearch {
				stocksFound = append(stocksFound, ticker)
			}
		}
	}

	return stocksFound, nil
}

func searchAssetsByKeyword(keyword string) ([]Stock, error) {
	stocksFound := []Stock{}

	for _, ticker := range stocksData.Stocks {
		if strings.Contains(ticker.Stock, keyword) {
			stocksFound = append(stocksFound, ticker)
		}
	}

	return stocksFound, nil
}

func loadAllAssets() (*StockApiResponse, error) {
	data, err := os.ReadFile("testdata/list.json")

	if err != nil {
		return nil, err
	}

	stocksData := StockApiResponse{}
	err = json.Unmarshal(data, &stocksData)

	return &stocksData, err
}

func loadQuoteData() (*QuoteApiResponse, error) {
	data, err := os.ReadFile("testdata/quote_data.json")

	if err != nil {
		return nil, err
	}

	quoteData := QuoteApiResponse{}
	err = json.Unmarshal(data, &quoteData)

	return &quoteData, err
}

func loadCryptocoins() ([]CryptoCoin, error) {
	data, err := os.ReadFile("testdata/cryptocoins.json")

	if err != nil {
		return nil, err
	}

	coinsRes := struct {
		Coins []CryptoCoin `json:"coins"`
	}{}

	err = json.Unmarshal(data, &coinsRes)

	return coinsRes.Coins, err
}

func searchCryptoCoins(coins []string, currency string) (cryptos []CryptoCoin) {
	for _, c := range cryptoData {
		if c.Currency == currency {
			for _, ticker := range coins {
				if ticker == c.Coin {
					cryptos = append(cryptos, c)
				}
			}
		}
	}

	return cryptos
}
