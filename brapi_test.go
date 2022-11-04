package brapigo

import (
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
