package entity

import (
	"Ranking-of-Internet-Corp-By-Go/util/sort"
	"sync"
)

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

func (mc *MarketCap) Compare(another sort.Comparable) int {

	if mc.GetValue() > another.GetValue() {
		return -1
	} else if mc.GetValue() == another.GetValue() {
		return 0
	}

	return 1
}