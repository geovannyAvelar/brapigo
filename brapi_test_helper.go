package brapigo

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
)

var stocksData *StockApiResponse

func runTestServer() *httptest.Server {
	data, err := loadAllAssets()

	if err != nil {
		panic("Cannot load test-data/list.json data file")
	}

	stocksData = data

	return httptest.NewServer(http.HandlerFunc(rootHandlerFunc))
}

type handler func(r *http.Request) ([]byte, error)

var endpoints map[string]handler = map[string]handler{
	"/api/quote/list": func(r *http.Request) ([]byte, error) {
		content, err := json.Marshal(stocksData)

		if err != nil {
			return nil, err
		}

		return content, nil
	},
	"/api/quote/PETR3": func(r *http.Request) ([]byte, error) {
		stocksFound, err := searchAssetsByTicker("PETR3")

		if err != nil {
			return nil, err
		}

		return json.Marshal(StockApiResponse{Stocks: stocksFound})
	},
	"/api/quote/PETR4,ITUB3": func(r *http.Request) ([]byte, error) {
		stocksFound, err := searchAssetsByTicker("PETR4,ITUB3")

		if err != nil {
			return nil, err
		}

		return json.Marshal(StockApiResponse{Stocks: stocksFound})
	},
	"/api/available": func(r *http.Request) ([]byte, error) {
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
	data, err := os.ReadFile("test-data/list.json")

	if err != nil {
		return nil, err
	}

	stocksData := StockApiResponse{}
	err = json.Unmarshal(data, &stocksData)

	return &stocksData, err
}
