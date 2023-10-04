package main

import (
    "fmt"
)

type TemaStrategy interface {
	ShouldEnterMarket(data []MarketData, me_index int) bool
	ShouldExitMarket(data []MarketData, me_index int) bool
	SetShortLookback(shortLookback int)
	SetLongLookback(longLookback int)
}

type MovingAverageCrossoverStrategy struct {
	shortLookback int
	longLookback  int
}

func (s *MovingAverageCrossoverStrategy) SetShortLookback(shortLookback int) {
	s.shortLookback = shortLookback
}

func (s *MovingAverageCrossoverStrategy) SetLongLookback(longLookback int) {
	s.longLookback = longLookback
}

func (s *MovingAverageCrossoverStrategy) ShouldEnterMarket(data []MarketData, i int) bool {

    // probably too early to trade
    // indicators are not yet accurate    
    if i < 10 {
        return false
    }
    
	// Check for Ema crossing
	if data[i].EmaShort > data[i].EmaLong && data[i-1].EmaShort <= data[i-1].EmaLong {
		return true
	}

	return false
}

func (s *MovingAverageCrossoverStrategy) ShouldExitMarket(data []MarketData, i int) bool {

    // probably too early to trade
    // indicators are not yet accurate
    if i < 10 {
        return false
    }

	if data[i].EmaShort < data[i].EmaLong && data[i-1].EmaShort >= data[i-1].EmaLong {
		return true
	}

    return false
}

func Backtest(strategy TemaStrategy, indData []MarketData, leverage float64, trade_size float64, cash_balance float64) (profits []float64, trades []Trade, ema_profit float64) {

	positionOpen := false
	entryPrice   := 0.0
	counter      := 0
    poslen       := 0
    quantity     := 0.0

    for i := 0; i < len(indData); i++ {

		if !positionOpen && strategy.ShouldEnterMarket(indData, i) {

			positionOpen = true
			entryPrice = indData[i].Price

            if cash_balance < trade_size {
                fmt.Println("Not enough capital left... Exiting.")
                break
            }

            // leverage
            quantity   = trade_size / entryPrice * leverage
            
			//fmt.Println("Entry price: ", price)

            trade := Trade{
                Datetime:  indData[i].Date,
                Indicator: "buy",
                Price:     indData[i].Price,
                Quantity:  quantity, // Adjust the quantity as needed
                Poslen:    0,
                PnL:       0,
            }
            trades = append(trades, trade)
            
		} else if positionOpen && strategy.ShouldExitMarket(indData, i) {

            _trade_pl_pct := (indData[i].Price / entryPrice) - 1

            _trade_pl_usd := _trade_pl_pct * (quantity * indData[i].Price)

            if _trade_pl_usd < (-1 * trade_size) {
                _trade_pl_usd = (-1 * trade_size)
            }
            
			ema_profit += _trade_pl_usd
            cash_balance += _trade_pl_usd
            
			positionOpen = false

			profits = append(profits, ema_profit)
            trade := Trade{
                Datetime:  indData[i].Date,
                Indicator: "sell",
                Price:     indData[i].Price,
                Quantity:  quantity, // Adjust the quantity as needed
                Poslen:    poslen,
                PnL:       _trade_pl_usd,
            }
            trades = append(trades, trade)
            poslen = 0
		}
		counter++

        // upcount poslen
        if poslen > 0 {
            poslen++
        }
        
	}

	return profits, trades, ema_profit
}
