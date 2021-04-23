package main

import (
	"Ranking-of-Internet-Corp-By-Go/entity"
	"Ranking-of-Internet-Corp-By-Go/util/config"
	"Ranking-of-Internet-Corp-By-Go/util/crawler"
	"Ranking-of-Internet-Corp-By-Go/util/sort"
	"fmt"
	"log"
	"strings"
	"sync"
)

func main() {

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

			if symbol, err := crawler.SearchStockSymbol(corpName); err != nil {
				log.Printf("查找公司 %s 的股票代码失败，请确认名称是否正确:%s\n", corpName, err.Error())
			} else {
				symbolsChan <- symbol
			}
			wg.Done()

		}(corpName)
	}


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
	}()


	finish := make(chan bool)
	go func() {
		arr := make([]sort.Comparable, 0, len(config.Cfg.Companies))

		for mc := range marketCapChan {
			arr = sort.Insert(arr, mc)
		}

		for i,mcap := range arr{

			if mc, ok := mcap.(*entity.MarketCap); ok {
				fmt.Printf("%2d. %s%s%.2f亿美元\t\t\t%.2f亿元\t\t\t%.2f亿港元\n",
					i+1,
					mc.Cooperation,
					strings.Repeat("\t", (32 - len([]rune(mc.Cooperation)) * 4)/ 8),
					float64(mc.USD) / 1e8,
					float64(mc.RMB) / 1e8,
					float64(mc.HKD) / 1e8,
				)
			}
		}

		close(finish)
	}()

	wg.Wait()
	close(symbolsChan)
	wg2.Wait()
	close(marketCapChan)

	<-finish

}
