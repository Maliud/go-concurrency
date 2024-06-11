package main

import (
	"fmt"
	"time"

	"github.com/Maliud/go-concurrency/internal/currency"
)

func runCurrencyWorker(
	worderId int,
	currenctChan <- chan currency.Currency,
	resultChan chan<- currency.Currency,
) {
	fmt.Printf(" %d başlanan\n", worderId)
	for c := range currenctChan{
		rates, err := currency.FetchCurrencyRates(c.Code)
		if err != nil {
			panic(err)
		}
		c.Rates = rates
		resultChan <- c
	}

	fmt.Printf("Duran", worderId)
}

func main() {
	ce := &currency.MyCurrencyExchange{
		Currencies: make(map[string]currency.Currency),
	}
	err := ce.FetchAllCurrencies()
	if err != nil{
		panic(err)
	}
	currencyChan := make(chan currency.Currency, len(ce.Currencies))
	resultChan := make(chan currency.Currency, len(ce.Currencies))

	for i := 0; i < 5; i++ {
		go runCurrencyWorker(i, currencyChan, resultChan)
	}
	startTime := time.Now()

	resultCount := 0

	for _, curr := range ce.Currencies {
		currencyChan <- curr
	}

	for {
		if resultCount == len(ce.Currencies) {
			fmt.Println("resultchan kapat")
			close(currencyChan)
			break
		}
		select{
		case c:= <- resultChan:
			ce.Currencies[c.Code] = c
			resultCount++
		case <- time.After(3 *time.Second):
			fmt.Println("Timeout")
			return
		}
	}
	endTime := time.Now()

	fmt.Println("============= Sonuclar =================")
	for _, curr := range ce.Currencies {
		fmt.Printf("%s (%s): %d oran:", curr.Name, curr.Code, len(curr.Rates))
	}

	fmt.Println("=====================================")
	fmt.Println("geçen süre", endTime.Sub(startTime))


}