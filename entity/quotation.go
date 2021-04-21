package entity

import "sync"

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
	Rates map[string]float64
	Lock sync.Locker
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
	USD,
	HKD,
	RMB int64
}