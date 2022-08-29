package main

import (
	"github.com/zhyee/Ranking-of-Internet-Corp-By-Go/entity"
	"github.com/zhyee/Ranking-of-Internet-Corp-By-Go/util/config"
	"github.com/zhyee/Ranking-of-Internet-Corp-By-Go/util/crawler"
	"github.com/zhyee/Ranking-of-Internet-Corp-By-Go/util/sort"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

func main() {
	// 设置http请求的超时时间
	http.DefaultClient.Timeout = time.Second * 5

	if err := config.ParseCfgFile(); err != nil {
		panic(err)
	}

	wg := &sync.WaitGroup{}
	wg.Add(len(config.Cfg.Companies))
	wg2 := &sync.WaitGroup{}

	symbolsChan := make(chan *entity.StockSymbol, len(config.Cfg.Companies))
	marketCapChan := make(chan *entity.MarketCap, len(config.Cfg.Companies))

	for _, corpName := range config.Cfg.Companies {
		go func(corpName string) {
			symbol, err := crawler.SearchStockSymbol(corpName)
			if err != nil {
				fmt.Println(err)
				var name string
				if name, err = crawler.QueryEnterpriseLatestName(corpName); name != "" && err == nil {
					symbol, err = crawler.SearchStockSymbol(name)
				}
			}
			if err != nil {
				log.Printf("查找公司 %s 的股票代码失败，请确认名称是否正确:%s\n", corpName, err.Error())
			} else {
				symbolsChan <- symbol
			}
			wg.Done()

		}(corpName)
	}


	symbolsChanEmpty := make(chan bool)
	go func() {
		for symbol := range symbolsChan {
			wg2.Add(1)
			go func(symbol *entity.StockSymbol) {
				mcap, err := crawler.GetMarketCap(symbol)
				if err != nil {
					log.Printf("查询公司 %s 市值失败： %s\n", symbol.CorpName, err.Error())
				} else {
					marketCapChan <- mcap
				}
				wg2.Done()
			}(symbol)
		}
		close(symbolsChanEmpty)
	}()


	finish := make(chan bool)
	go func() {
		arr := make([]entity.Comparable, 0, len(config.Cfg.Companies))

		for mc := range marketCapChan {
			arr = sort.Insert(arr, mc)
		}

		for i,mcap := range arr{

			if mc, ok := mcap.(*entity.MarketCap); ok {
				name := fmt.Sprintf("%2d. %s", i+1, mc.Cooperation)
				usd := fmt.Sprintf("%.2f亿美元", float64(mc.USD) / 1e8)
				rmb := fmt.Sprintf("%.2f亿元", float64(mc.RMB) / 1e8)
				hkd := fmt.Sprintf("%.2f亿港元", float64(mc.HKD) / 1e8)

				fmt.Printf("%s%s%s%s%s%s%s\n",
					name,
					"\t\t",
					usd,
					"\t\t",
					rmb,
					"\t\t",
					hkd,
				)
			}
		}

		close(finish)
	}()

	wg.Wait()
	close(symbolsChan)
	<-symbolsChanEmpty
	wg2.Wait()
	close(marketCapChan)

	<-finish
}
