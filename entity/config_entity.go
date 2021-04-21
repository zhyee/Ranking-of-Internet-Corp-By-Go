package entity

type Config struct {
	Companies []string `yaml:"companies"`
	ApiUrl struct{
		ExchangeRate string `yaml:"exchange_rate"`
		MarketCode string `yaml:"market_code"`
		MarketValue string `yaml:"market_value"`
	} `yaml:"apiUrl"`
}
