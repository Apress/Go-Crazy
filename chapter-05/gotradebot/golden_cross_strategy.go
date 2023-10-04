package main

type GoldenCrossStrategy interface {

	SetShortLookback(shortLookback int)
	SetLongLookback(longLookback int)
	ShouldEnterMarket(data []MarketData, me_index int) bool
	ShouldExitMarket(data []MarketData, me_index int) bool
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

func (s *GoldenCrossMaStrategy) ShouldEnterMarket(data []MarketData, i int) bool {

    if i < 200 {
        return false
    }

	// Check for Golden Cross
	if data[i].Sma50 > data[i].Sma200 && data[i-1].Sma50 <= data[i-1].Sma200 {
		return true
	}

	return false
}

func (s *GoldenCrossMaStrategy) ShouldExitMarket(data []MarketData, i int) bool {

    if i < 200 {
        return false
    }

	// Check for Death Cross
	if data[i].Sma50 < data[i].Sma200 && data[i-1].Sma50 >= data[i-1].Sma200 {
		return true
	}

	return false
}
