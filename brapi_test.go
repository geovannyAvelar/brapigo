package brapigo

import (
	"testing"
)

func TestFindAssetByTicker(t *testing.T) {
	testApiServer := runTestServer()

	brApi := NewBrApiWithCustomBaseUrl(testApiServer.URL)

	stocks, err := brApi.FindAssetByTicker("PETR3")

	if err != nil {
		t.Errorf("Error in /api/quote/PETR3 request. Cause %s", err.Error())
	}

	var petr3Stock *Stock

	for _, stock := range stocks {
		if stock.Stock == "PETR3" {
			petr3Stock = &stock
		}
	}

	if petr3Stock == nil {
		t.Errorf("Error in /api/quote/PETR3 request. Could not find PETR3")
	}

	stocks, err = brApi.FindAssetByTicker("PETR4,ITUB3")

	var petr4Stock *Stock
	var itub3Stock *Stock

	for _, stock := range stocks {
		if stock.Stock == "PETR4" {
			petr4Stock = &stock
		}

		if stock.Stock == "ITUB3" {
			itub3Stock = &stock
		}
	}

	if petr4Stock == nil || itub3Stock == nil {
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

func TestListAllTickers(t *testing.T) {
	testApiServer := runTestServer()

	brApi := NewBrApiWithCustomBaseUrl(testApiServer.URL)

	tickers, err := brApi.ListTickers()

	if err != nil {
		t.Errorf("Error in /api/quote/list request. Cause %s", err.Error())
	}

	if len(tickers) != 1585 {
		t.Errorf("Error is request. Received %d but expected %d", len(tickers), 1585)
	}
}
