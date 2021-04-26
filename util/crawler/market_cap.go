package crawler

import (
	"Ranking-of-Internet-Corp-By-Go/entity"
	"Ranking-of-Internet-Corp-By-Go/util/config"
	"Ranking-of-Internet-Corp-By-Go/util/http_util"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
	"time"
)

var MarketNames = map[string]int{
	"美股":1,
	"港股":2,
	"深A":3,
	"沪A":4,
	"科创板":5,
}

/**
查询股票代号
 */
func SearchStockSymbol(corpName string) (*entity.StockSymbol, error) {
	if err := config.ParseCfgFile(); err != nil {
		return nil, err
	}

	url := fmt.Sprintf(config.Cfg.ApiUrl.MarketCode, corpName, time.Now().UnixNano() / 1e6)

	resp, err := http_util.GetWithRetry(url, 3)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if body, err = http_util.Convert2UTF8(body, http_util.GetCharset(resp.Header)); err != nil {
		return nil, err
	}

	code := entity.QuotationCode{}

	if err := json.Unmarshal(body, &code); err != nil {
		return nil, nil
	}

	for _, val := range code.QuotationCodeTable.Data {
		if _, ok := MarketNames[val.SecurityTypeName]; ok {
			return &entity.StockSymbol{CorpName: corpName, Symbol:entity.TickerSymbol(val.QuoteID)}, nil
		}
	}

	return nil, fmt.Errorf("quotation code for %s not found", corpName)
}

func GetMarketCap(symbol *entity.StockSymbol) (*entity.MarketCap, error) {

	if err := config.ParseCfgFile(); err != nil {
		return nil, err
	}

	url := fmt.Sprintf(config.Cfg.ApiUrl.MarketValue, symbol.Symbol)

	resp, err := http_util.GetWithRetry(url, 3)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if body, err = http_util.Convert2UTF8(body, http_util.GetCharset(resp.Header)); err != nil {
		return nil, err
	}

	apiResp := entity.MarketCapApiResp{}

	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, err
	}

	marketCap, err := parseValue(apiResp)
	if err != nil {
		return nil, err
	}
	marketCap.Cooperation = symbol.CorpName

	return marketCap, nil
}

func getCurrencyType(unit string) CurrencyType {
	unit = strings.ToUpper(unit)
	if MoneyType, ok := CurrencyNameTypeMap[unit]; ok {
		return MoneyType
	}
	return RMB
}

func parseValue(apiResp entity.MarketCapApiResp) (*entity.MarketCap, error) {
	moneyType := getCurrencyType(apiResp.Data.Unit)
	amount := apiResp.Data.Value
	amountRMB, err := Exchange(moneyType, RMB, amount)
	if err != nil {
		return nil, err
	}
	amountUSD, err := Exchange(moneyType, USD, amount)
	if err != nil {
		return nil, err
	}
	amountHKD, err := Exchange(moneyType, HKD, amount)
	if err != nil {
		return nil, err
	}

	return &entity.MarketCap{
		Cooperation:apiResp.Data.Name,
		USD:int64(amountUSD),
		HKD:int64(amountHKD),
		RMB:int64(amountRMB),
	}, nil
}
