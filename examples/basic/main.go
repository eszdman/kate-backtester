package main

import (
	"fmt"

	"github.com/eszdman/kate-backtester/kate"
)

// SimpleStrategy is a basic trading strategy that open long positions when prices rise
type SimpleStrategy struct {
	lastPrice    *kate.DataPoint
	currentPrice *kate.DataPoint
}

func main() {
	data, err := kate.PricesFromCSV("../../testdata/ETHUSD5.csv")

	if err != nil {
		panic("could`t load data." + err.Error())
	}

	kate.NewBacktester(&SimpleStrategy{}, data)
	backtester := kate.NewBacktester(&SimpleStrategy{}, data)
	stats := backtester.Run()
	fmt.Println(stats)
}

// PreProcessIndicators allows the pre processing of indicators
func (stg *SimpleStrategy) PreProcessIndicators(latestPrice kate.DataPoint) {
	//No indicators to process
	stg.lastPrice = stg.currentPrice
	stg.currentPrice = &latestPrice
}

// OpenNewPosition process the next data point and checks if a position should be opened
func (stg *SimpleStrategy) OpenNewPosition(latestPrice kate.DataPoint) *kate.OpenPositionEvt {
	if stg.lastPrice != nil && stg.currentPrice.Close() > stg.lastPrice.Close() {
		return &kate.OpenPositionEvt{Direction: kate.LONG, Leverage: 30}
	}
	return nil
}

// SetStoploss defines a stoploss for the current open position
func (stg *SimpleStrategy) SetStoploss(openPosition kate.Position) *kate.StoplossEvt {
	if openPosition.Direction == kate.LONG && openPosition.Stoploss <= 0 {
		return &kate.StoplossEvt{Price: openPosition.EntryPrice * 0.995}
	}
	return nil
}

// SetTakeProfit defines a takeprofit for the current open position
func (stg *SimpleStrategy) SetTakeProfit(openPosition kate.Position) *kate.TakeProfitEvt {
	if openPosition.Direction == kate.LONG && openPosition.TakeProfit <= 0 {
		return &kate.TakeProfitEvt{Price: openPosition.EntryPrice * 1.005}
	}
	return nil
}
