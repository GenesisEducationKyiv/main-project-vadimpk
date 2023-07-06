package entity

type CryptoCurrency string

const (
	CryptoCurrencyBTC CryptoCurrency = "BTC"
	CryptoCurrencyETH CryptoCurrency = "ETH"
)

func (c CryptoCurrency) String() string {
	return string(c)
}

var cryptoCurrencies = map[CryptoCurrency]struct{}{
	CryptoCurrencyBTC: {}, CryptoCurrencyETH: {},
}

func (c CryptoCurrency) IsValid() bool {
	_, ok := cryptoCurrencies[c]
	return ok
}

type FiatCurrency string

const (
	FiatCurrencyUSD FiatCurrency = "USD"
	FiatCurrencyUAH FiatCurrency = "UAH"
)

func (f FiatCurrency) String() string {
	return string(f)
}

var fiatCurrencies = map[FiatCurrency]struct{}{
	FiatCurrencyUSD: {}, FiatCurrencyUAH: {},
}

func (f FiatCurrency) IsValid() bool {
	_, ok := fiatCurrencies[f]
	return ok
}
