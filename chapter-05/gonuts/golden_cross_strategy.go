package main

import (
    "fmt"
)

type GoldenCrossStrategy interface {

	SetShortLookback(shortLookback int)
	SetLongLookback(longLookback int)
	ShouldEnterGoldenCrossMarket(data []MarketData, me_index int) bool
	ShouldExitGoldenCrossMarket(data []MarketData, me_index int) bool
}

type GoldenCrossMaStrategy struct {
	shortPeriod int
	longPeriod  int
}

func (s *GoldenCrossMaStrategy) SetShortLookback(shortLookback int) {
    s.shortPeriod = shortLookback
}
func (s *GoldenCrossMaStrategy) SetLongLookback(longLookback int) {
    s.longPeriod = longLookback
}

func BacktestGoldenCross(strategy GoldenCrossStrategy, indData []MarketData, leverage float64, trade_size float64, cash_balance float64) (gcprofits []float64, gctrades []Trade, gcprofit float64) {

	positionOpen := false
	entryPrice   := 0.0
	counter      := 0
    quantity     := 0.0
    poslen       := 0

	for i := 0; i< len(indData); i++ {

		if !positionOpen && strategy.ShouldEnterGoldenCrossMarket(indData, counter) {
			positionOpen = true
			entryPrice = indData[i].Price

            if cash_balance < trade_size {
                fmt.Println("Not enough capital left... Exiting.")
                break
            }

            // leverage
            quantity = trade_size / entryPrice * leverage
            
            gctrade := Trade{
                Datetime:  indData[i].Date,
                Indicator: "buy",
                Price:     indData[i].Price,
                Quantity:  quantity, // Adjust the quantity as needed
                Poslen:    0,
                PnL:       0,
            }
            gctrades = append(gctrades, gctrade)
            
			entryPrice = indData[i].Price

		} else if positionOpen && strategy.ShouldExitGoldenCrossMarket(indData, counter) {

            _trade_pl_pct := (indData[i].Price / entryPrice) - 1

            _trade_pl_usd := _trade_pl_pct * (quantity * indData[i].Price)

            if _trade_pl_usd < (-1 * trade_size) {
                _trade_pl_usd = (-1 * trade_size)
            }
            
			gcprofit += _trade_pl_usd
            cash_balance += _trade_pl_usd
            
			positionOpen = false

			gcprofits = append(gcprofits, gcprofit)
            
            gctrade := Trade{
                Datetime:  indData[i].Date,
                Indicator: "sell",
                Price:     indData[i].Price,
                Quantity:  quantity, // Adjust the quantity as needed
                Poslen:    poslen,
                PnL:       _trade_pl_usd,
            }
            gctrades = append(gctrades, gctrade)

			gcprofits = append(gcprofits, gcprofit)
		}
		counter++

        // upcount poslen
        if poslen > 0 {
            poslen++
        }
        
	}

	return gcprofits, gctrades, gcprofit
}

func (s *GoldenCrossMaStrategy) ShouldEnterGoldenCrossMarket(data []MarketData, i int) bool {

    if i < 200 {
        return false
    }

	// Check for Golden Cross
	if data[i].Sma50 > data[i].Sma200 && data[i-1].Sma50 <= data[i-1].Sma200 {
		return true
	}

	return false
}

func (s *GoldenCrossMaStrategy) ShouldExitGoldenCrossMarket(data []MarketData, i int) bool {

    if i < 200 {
        return false
    }

	// Check for Death Cross
	if data[i].Sma50 < data[i].Sma200 && data[i-1].Sma50 >= data[i-1].Sma200 {
		return true
	}

	return false
}
