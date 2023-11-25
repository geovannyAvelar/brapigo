package brapigo

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestFindAssetByTicker(t *testing.T) {
	testApiServer := runTestServer()

	brApi := NewBrApiWithCustomBaseUrl(testApiServer.URL)

	quotes, err := brApi.FindAssetByTicker("PETR4", "ITUB3")

	if err != nil {
		t.Errorf("Error during /api/quote/PETR4,ITUB3. Cause: %s", err.Error())
	}

	var petr4Quote *Quote
	var itub3Quote *Quote

	for _, quote := range quotes {
		if quote.Symbol == "PETR4" {
			petr4Quote = &quote
		}

		if quote.Symbol == "ITUB3" {
			itub3Quote = &quote
		}
	}

	if petr4Quote == nil || itub3Quote == nil {
		t.Errorf("Error in /api/quote/PETR4,ITUB3 request. Could not find PETR4 or ITUB3")
	}
}

func TestSearchTickets(t *testing.T) {
	testApiServer := runTestServer()

	brApi := NewBrApiWithCustomBaseUrl(testApiServer.URL)

	keyword := "PETR"
	tickers, err := brApi.SearchTickets(keyword)

	if err != nil {
		t.Errorf("Error in /api/available request. Cause %s", err.Error())
	}

	if len(tickers) == 0 {
		t.Errorf("Could not find a ticker with %s keyword", keyword)
	}
}

func TestListStocks(t *testing.T) {
	testApiServer := runTestServer()

	brApi := NewBrApiWithCustomBaseUrl(testApiServer.URL)

	tickers, err := brApi.ListStocks()

	if err != nil {
		t.Errorf("Error in /api/quote/list request. Cause %s", err.Error())
	}

	if len(tickers) != 1585 {
		t.Errorf("Error is request. Received %d but expected %d", len(tickers), 1585)
	}
}

func TestListCryptoCoins(t *testing.T) {
	t.Parallel()

	testApiServer := runTestServer()

	brApi := NewBrApiWithCustomBaseUrl(testApiServer.URL)

	coin, err := brApi.ListCryptoCoins()

	if err != nil {
		t.Errorf("Error in /v2/crypto/available request. Cause: %s", err)
	}

	if len(coin) == 0 {
		t.Errorf("Coins slice is empty")
	}
}

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
