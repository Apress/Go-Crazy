package main

import (
	"fmt"
	"math"
)

func SharpeRatio(returns []float64, riskFreeRate float64) float64 {

    // Sharpe ratio has been adapted to cater to leverage settings
    
	n := len(returns)
	if n == 0 {
		return 0
	}

	// Calculate the average return
	var avgReturn float64
	for _, r := range returns {
		avgReturn += r
	}
	avgReturn /= float64(n)

	// Calculate the standard deviation of returns
	var stdDev float64
	for _, r := range returns {
		diff := r - avgReturn
		stdDev += diff * diff
	}
	stdDev = math.Sqrt(stdDev / float64(n))

	// Calculate the Sharpe Ratio
	excessReturn := avgReturn - riskFreeRate
	sharpeRatio := excessReturn / stdDev
	return sharpeRatio
}

func MARRatio(returns []float64) float64 {

    // MAR ratio has been adapted to cater to leverage settings

	n := len(returns)
	if n == 0 {
		return 0
	}

	// Calculate the cumulative return
	cumulativeReturn := 0.0

	// Calculate the maximum drawdown
	maxDrawdown := 0.0000001
	peak := returns[0]
	trough := returns[0]
	for _, r := range returns[1:] {

        cumulativeReturn += r
        
		if r > peak {
			peak = r
			trough = r
		} else if r < trough {
			trough = r
		}

		drawdown := (peak - trough) / peak
		if drawdown > maxDrawdown {
			maxDrawdown = drawdown
		}

	}

    cumulativeReturn /= config.CashBalance
    fmt.Printf("cumulativeReturn: %.2f\n", cumulativeReturn)

	// Calculate the MAR Ratio
	marRatio := cumulativeReturn / maxDrawdown

    if marRatio < 0 {
        return 0.0
    }
    
	return marRatio
}

