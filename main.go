package main

import (
	"Ranking-of-Internet-Corp-By-Go/entity"
	"Ranking-of-Internet-Corp-By-Go/util/config"
	"Ranking-of-Internet-Corp-By-Go/util/crawler"
	"fmt"
	"log"
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

	for _, corpName := range config.Cfg.Companies {
		go func(corpName string) {

			fmt.Println(corpName)

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
					fmt.Println(symbol.CorpName, mcap.USD, mcap.RMB, mcap.HKD)

					fmt.Println()
				}
				wg2.Done()
			}(symbol)
		}
	}()

	wg.Wait()
	close(symbolsChan)
	wg2.Wait()

}
