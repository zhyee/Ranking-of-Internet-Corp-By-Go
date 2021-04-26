package crawler

import (
	"Ranking-of-Internet-Corp-By-Go/entity"
	"Ranking-of-Internet-Corp-By-Go/util/config"
	"Ranking-of-Internet-Corp-By-Go/util/http_util"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"
)

type CurrencyType uint8

const (
	USD CurrencyType = iota
	HKD
	RMB
)

var CurrencyUnits  = map[CurrencyType]string{
	USD:"美元",
	HKD:"港元",
	RMB:"元",
}

var CurrencyNames = map[CurrencyType]string{
	USD: "usd",
	HKD: "hkd",
	RMB: "cny",
}

var CurrencyNameTypeMap = map[string]CurrencyType {
	"USD":USD,
	"HKD":HKD,
	"RMB":RMB,
	"CNY":RMB,
}

var exchangeRateCache = &entity.ExchangeRateCache{
	Rates:  make(map[string]float64),
	Locker: &sync.RWMutex{},
}

func getExchangeRate(from, to CurrencyType) (float64, error) {
	if from == to {
		return 1, nil
	}

	key := fmt.Sprintf("%d->%d", from, to)
	exchangeRateCache.Locker.RLock()
	if rate, ok := exchangeRateCache.Rates[key]; ok {
		exchangeRateCache.Locker.RUnlock()
		return rate, nil
	}
	exchangeRateCache.Locker.RUnlock()

	if err := config.ParseCfgFile(); err != nil {
		log.Println("配置文件解析出错：")
		panic(err)
	}

	url := fmt.Sprintf(
		config.Cfg.ApiUrl.ExchangeRate,
		time.Now().UnixNano() / 1e6,
		CurrencyNames[from],
		CurrencyNames[to],
	)

	resp, err := http_util.GetWithRetry(url, 3)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	if body, err = http_util.Convert2UTF8(body, http_util.GetCharset(resp.Header)); err != nil {
		return 0, err
	}

	utf8Body := string(body)

	utf8Body = strings.ReplaceAll(utf8Body, " ", "")

	if idx := strings.Index(utf8Body, "=\""); idx > -1 {
		utf8Body = utf8Body[idx+2:]
	} else {
		return 0, errors.New("无法在接口返回中找到特征字符：=\"")
	}

	if idx := strings.Index(utf8Body, "\""); idx > -1 {
		utf8Body = utf8Body[:idx]
	} else {
		return 0, errors.New("无法在接口返回中找到特征字符：\"")
	}

	pieces := strings.Split(utf8Body, ",")

	rate, err := strconv.ParseFloat(pieces[2], 64)
	if err != nil {
		return 0, err
	}
	exchangeRateCache.Locker.Lock()
	exchangeRateCache.Rates[key] = rate
	exchangeRateCache.Locker.Unlock()
	return rate, nil
}

/**
币种之间的兑换
 */
func Exchange(from, to CurrencyType, amount float64) (float64,error) {
	if from == to {
		return amount, nil
	}
	rate, err := getExchangeRate(from, to)
	if err != nil {
		return 0, err
	}

	return amount * rate, nil
}
