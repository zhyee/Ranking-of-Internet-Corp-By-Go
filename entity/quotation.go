package entity

import (
	"sync"
)

type NameApiResp struct {
	Code string `json:"code"`
	Message string `json:"message"`
	Data struct{
		Stock []struct{
			StockName string `json:"stock_name"`
			MarketType int `json:"market_type"`
			Market string `json:"market"`
		} `json:"stock"`
	} `json:"data"`
}

type QuotationCode struct {
	QuotationCodeTable struct {
		Data []struct{
			Code,
			Name,
			SecurityTypeName,
			MktNum,
			QuoteID string
		}
	}
}
type ExchangeRateCache struct {
	Rates  map[string]float64
	Locker *sync.RWMutex
}

type TickerSymbol string

type StockSymbol struct {
	CorpName string
	Symbol   TickerSymbol
}

type MarketCapApiResp struct {
	Data struct{
		Name string `json:"f58"`
		Value float64 `json:"f116"`
		Unit string	`json:"f172"`
	} `json:"data"`
}

type MarketCap struct {
	Cooperation string
	USD,
	HKD,
	RMB int64
}

func (mc *MarketCap) GetValue() int64 {
	return mc.USD
}

func (mc *MarketCap) Compare(another Comparable) int {

	if mc.GetValue() > another.GetValue() {
		return -1
	} else if mc.GetValue() == another.GetValue() {
		return 0
	}

	return 1
}